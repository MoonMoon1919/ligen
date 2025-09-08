package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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
func MITLicense(holder string, year int) (string, error) {
	licenseContent := `Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
	`

	copyright, err := Copyright(holder, year)
	if err != nil {
		return "", err
	}

	fullLicense := strings.Join([]string{copyright + "\n", licenseContent}, "\n")

	return fullLicense, nil
}

// General use copyright line
func Copyright(name string, year int) (string, error) {
	currentYear := time.Now().Year()
	fiftyYearsAgo := year - 50

	if year > currentYear || year < fiftyYearsAgo {
		return "", errors.New("Invalid year")
	}

	strippedName := strings.TrimSpace(name)
	if len(strippedName) == 0 {
		return "", errors.New("Name must not be empty")
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

// DOIT
func main() {
	repo := FileRepository{Path: "LICENSE"}

	license, err := New("Max Moon", 2025)
	if err != nil {
		panic(err)
	}

	if err = repo.Write(&license); err != nil {
		panic(err)
	}
}
