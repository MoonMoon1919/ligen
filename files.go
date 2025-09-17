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

type licenseLoadResult struct {
	licenseType LicenseType
	content     string
}

func loadLicense(reader io.Reader) (licenseLoadResult, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return licenseLoadResult{}, err
	}

	contentString := string(content)

	licenseType, err := Match(contentString, 0.90)
	if err != nil {
		return licenseLoadResult{}, err
	}

	return licenseLoadResult{
		licenseType: licenseType,
		content:     contentString,
	}, nil
}

type FileRepository struct{}

func (f FileRepository) Load(path string, license *License) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()
	result, err := loadLicense(file)
	if err != nil {
		return err
	}

	var copyright Copyright

	if result.licenseType.RequiresNotice() {
		noticeFile, err := os.OpenFile("NOTICE", os.O_RDONLY, 0644)
		if err != nil {
			return err
		}

		noticeContent, err := io.ReadAll(noticeFile)
		if err != nil {
			return err
		}
		copyright, err = ParseDoc(string(noticeContent))
		if err != nil {
			return err
		}
	} else {
		copyright, err = ParseDoc(result.content)
		if err != nil {
			return err
		}
	}

	// TODO: Set project name from notice, since that's the only place it is set
	license.projectName = ""
	license.copyright = copyright
	license.SetLicenseType(result.licenseType)

	return nil
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
