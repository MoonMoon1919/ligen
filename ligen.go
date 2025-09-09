package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/MoonMoon1919/ligen/pkg/licenses"
)

/*
Generation -
- Each license may have different inputs
- A license may have _no_ inputs

Checking -
- Answer "what license is in this repo?"
*/

// MIT
func MITLicense(holder string, year int) (string, error) {
	licenseContent := licenses.MIT

	copyright, err := Copyright(holder, year)
	if err != nil {
		return "", err
	}

	fullLicense := strings.Join([]string{copyright + "\n", licenseContent}, "\n")

	return fullLicense, nil
}

// General use copyright line
var (
	InvalidYearError = errors.New("Invalid year")
	EmptyNameError   = errors.New("Name must not be empty")
	NameTooLongError = errors.New("Name must be 128 chars")
)

const (
	// MAX_NAME_LENGTH is the maximum amount of chars the holder of a copyright can contain
	// 128 picked arbitrarily, seemed reasonable
	MAX_NAME_LENGTH = 128
	// MAX_YEARS_PAST is the maximum amount of time in years that a copyright can be backdated
	// 50 picked arbitrarily, seemed reasonable
	MAX_YEARS_PAST = 50
)

func Copyright(name string, year int) (string, error) {
	currentYear := time.Now().Year()
	fiftyYearsAgo := currentYear - MAX_YEARS_PAST

	if year > currentYear || year < fiftyYearsAgo {
		return "", InvalidYearError
	}

	strippedName := strings.TrimSpace(name)
	if len(strippedName) == 0 {
		return "", EmptyNameError
	}

	if len(name) > MAX_NAME_LENGTH {
		return "", NameTooLongError
	}

	return fmt.Sprintf("Copyright %d %s", year, strippedName), nil
}

// License stuff
type License struct {
	content string
}

func New(holder string, year int) (License, error) {
	content, err := MITLicense(holder, year)
	if err != nil {
		return License{}, err
	}

	return License{
		content: content,
	}, nil
}

func Render(license *License) (string, error) {
	return license.content, nil
}

// File management
type RenderOptions struct {
	TrailingNewline bool
}

func Write(writer io.Writer, licence *License, renderOpts *RenderOptions) error {
	content, err := Render(licence)
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(content))
	return err
}

type FileRepository struct {
	Path string
}

func (f FileRepository) Write(license *License) error {
	file, err := os.OpenFile("LICENSE", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	return Write(file, license, &RenderOptions{})
}
