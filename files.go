package ligen

import (
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

type loader func() (io.Reader, error)

func Load(license *License, licenseLoader loader, noticeLoader loader) error {
	licenseReader, err := licenseLoader()
	if err != nil {
		return err
	}

	licenseResult, err := loadLicense(licenseReader)
	if err != nil {
		return err
	}

	var projectName string
	contentContainingCopyright := licenseResult.content

	if licenseResult.licenseType.RequiresNotice() {
		noticeReader, err := noticeLoader()
		if err != nil {
			return err
		}

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

	copyright, err := ParseDocForCopyright(contentContainingCopyright)
	if err != nil {
		return err
	}

	// TODO: Set project name from notice, since that's the only place it is set
	license.projectName = projectName
	license.copyright = copyright
	license.SetLicenseType(licenseResult.licenseType)

	return nil
}

func loadFile(path string) (io.Reader, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return f, nil
}

type FileRepository struct{}

func (f FileRepository) Load(path string, license *License) error {
	ll := func() (io.Reader, error) {
		return loadFile(path)
	}

	nl := func() (io.Reader, error) {
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

		Write(file, writeable, render)

		return nil
	}

	for _, writeable := range writeables {
		write(&writeable, renderOpts)
	}

	return nil
}
