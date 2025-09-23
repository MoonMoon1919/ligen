package ligen

import (
	"errors"
	"io"
	"os"
)

// File management
type RenderOptions struct {
	TrailingNewline bool
}

func Write(writer io.Writer, writeable *Writeable, renderOpts *RenderOptions) error {
	_, err := writer.Write([]byte(writeable.Content))

	return err
}

type LicenseLoadResult struct {
	licenseType LicenseType
	content     string
}

func loadLicense(reader io.Reader) (LicenseLoadResult, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return LicenseLoadResult{}, err
	}

	contentString := string(content)

	licenseType, err := Match(contentString, 0.90)
	if err != nil {
		return LicenseLoadResult{}, err
	}

	return LicenseLoadResult{
		licenseType: licenseType,
		content:     contentString,
	}, nil
}

func loadNotice(reader io.Reader) (string, error) {
	noticeContent, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(noticeContent), nil
}

type loader func() (io.Reader, func() error, error)

func Load(license *License, licenseLoader loader, noticeLoader loader) error {
	licenseReader, close, err := licenseLoader()
	if err != nil {
		return err
	}
	defer close()

	licenseResult, err := loadLicense(licenseReader)
	if err != nil {
		return err
	}

	var projectName string
	contentContainingCopyright := licenseResult.content

	if licenseResult.licenseType.RequiresNotice() {
		noticeReader, close, err := noticeLoader()
		if err != nil {
			return err
		}
		defer close()

		notice, err := loadNotice(noticeReader)
		if err != nil {
			return err
		}

		contentContainingCopyright = notice

		projectName, err = ParseProjectNameFromNotice(notice)
		if err != nil {
			return err
		}
	}

	var copyright Copyright
	if licenseResult.licenseType.RequiresCopyright() {
		copyright, err = ParseDocForCopyright(contentContainingCopyright)
		if err != nil {
			return err
		}
	}

	license.projectName = projectName
	license.copyright = copyright
	license.SetLicenseType(licenseResult.licenseType)

	return nil
}

func loadFile(path string) (io.Reader, func() error, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	close := func() error {
		return f.Close()
	}

	return f, close, nil
}

type FileRepository struct{}

func (f FileRepository) Load(path string, license *License) error {
	ll := func() (io.Reader, func() error, error) {
		return loadFile(path)
	}

	nl := func() (io.Reader, func() error, error) {
		return loadFile("NOTICE")
	}

	return Load(license, ll, nl)
}

func (f FileRepository) Write(license *License) error {
	writeables, err := license.Render()
	if err != nil {
		return err
	}

	renderOpts := &RenderOptions{}

	write := func(writeable *Writeable, render *RenderOptions) error {
		file, err := os.OpenFile(writeable.Path, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		return Write(file, writeable, render)
	}

	for _, writeable := range writeables {
		if err := write(&writeable, renderOpts); err != nil {
			return err
		}
	}

	return nil
}

func DiscoverLicenseFile() (string, error) {
	// Files without extensions first (standard convention)
	primaryCandidates := []string{
		"LICENSE",
		"UNLICENSE",
		"COPYING.LESSER",
	}

	// Check primary candidates first
	for _, candidate := range primaryCandidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	// Fallback to files with extensions if no standard files found
	fallbackCandidates := []string{
		"LICENSE.txt",
		"LICENSE.md",
	}

	for _, candidate := range fallbackCandidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", errors.New("no license file found in current directory")
}
