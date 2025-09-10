package ligen

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"text/template"
	"time"
)

/*
Generation -
- Each license may have different inputs
- A license may have _no_ inputs

Checking -
- Answer "what license is in this repo?"
*/

// MIT
// Body of text for an MIT License
const MitTemplateBody = `
Copyright {{.Year}} {{.Holder}}

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
`

var MITTemplate = template.Must(template.New("MIT").Parse(MitTemplateBody))

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

type Copyright struct {
	Holder string
	Year   int
}

func NewCopyright(name string, year int) (Copyright, error) {
	currentYear := time.Now().Year()
	fiftyYearsAgo := currentYear - MAX_YEARS_PAST

	if year > currentYear || year < fiftyYearsAgo {
		return Copyright{}, InvalidYearError
	}

	strippedName := strings.TrimSpace(name)
	if len(strippedName) == 0 {
		return Copyright{}, EmptyNameError
	}

	if len(name) > MAX_NAME_LENGTH {
		return Copyright{}, NameTooLongError
	}

	return Copyright{Holder: name, Year: year}, nil
}

// License stuff
type LicenseType int

const (
	MIT LicenseType = iota + 1
)

type Writeable struct {
	content string
	path    string
}

type writeableGenerator func(cr *Copyright, dest *bytes.Buffer) ([]Writeable, error)

func MITGenerator(cr *Copyright, dest *bytes.Buffer) ([]Writeable, error) {
	if err := MITTemplate.Execute(dest, cr); err != nil {
		return nil, err
	}

	writeableSlice := make([]Writeable, 1)
	writeableSlice[0] = Writeable{content: dest.String(), path: "LICENSE"}

	return writeableSlice, nil
}

func generatorFactory(licenseType LicenseType) (writeableGenerator, error) {
	switch licenseType {
	case MIT:
		return MITGenerator, nil
	default:
		return nil, errors.New("Unsupported license type")
	}
}

type License struct {
	copyright     Copyright
	generatorFunc writeableGenerator
}

func New(holder string, year int, licenseType LicenseType) (*License, error) {
	copyright, err := NewCopyright(holder, year)
	if err != nil {
		return &License{}, err
	}

	generatorFunc, err := generatorFactory(licenseType)
	if err != nil {
		return &License{}, err
	}

	return &License{
		copyright:     copyright,
		generatorFunc: generatorFunc,
	}, nil
}

func (l *License) Render() ([]Writeable, error) {
	var content bytes.Buffer

	writeable, err := l.generatorFunc(&l.copyright, &content)

	if err != nil {
		return nil, err
	}

	return writeable, nil
}

// File management
type RenderOptions struct {
	TrailingNewline bool
}

func Write(writer io.Writer, writeable *Writeable, renderOpts *RenderOptions) error {
	_, err := writer.Write([]byte(writeable.content))

	return err
}

type FileRepository struct{}

func (f FileRepository) Write(license *License) error {
	writeables, err := license.Render()
	if err != nil {
		return err
	}

	renderOpts := &RenderOptions{}

	write := func(writeable *Writeable, render *RenderOptions) error {
		file, err := os.OpenFile(writeable.path, os.O_CREATE|os.O_WRONLY, 0644)
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
