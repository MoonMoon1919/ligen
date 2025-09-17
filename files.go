package ligen

import (
	"fmt"
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

func Load(reader io.Reader, license *License) error {
	content, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	contentString := string(content)

	licenseType, err := Match(contentString, 0.90)
	if err != nil {
		return err
	}

	if licenseType.RequiresNotice() {
		fmt.Print("License requires notice")
	}

	copyright, err := ParseDoc(contentString)
	if err != nil {
		return err
	}

	// Set it up
	license.SetCopyrightStartYear(copyright.StartYear)

	if copyright.EndYear != 0 {
		license.SetCopyrightEndYear(copyright.EndYear)
	}

	license.SetHolder(copyright.Holder)

	if err = license.SetLicenseType(licenseType); err != nil {
		return err
	}

	return nil
}

type FileRepository struct{}

func (f FileRepository) Load(path string, license *License) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()
	return Load(file, license)
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
