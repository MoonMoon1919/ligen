package ligen

import (
	"io"
	"reflect"
	"slices"
	"strings"
	"testing"
)

type FakeRepo struct {
	files map[string]string
}

func NewFakeRepo(files ...Writeable) FakeRepo {
	fileMap := make(map[string]string, len(files))

	for _, file := range files {
		fileMap[file.Path] = file.Content
	}

	return FakeRepo{
		files: fileMap,
	}
}

func (f *FakeRepo) Load(path string, license *License) error {
	ll := func() (io.Reader, error) {
		return strings.NewReader(f.files[path]), nil
	}

	nl := func() (io.Reader, error) {
		return strings.NewReader(f.files["NOTICE"]), nil
	}

	return Load(license, ll, nl)
}

func (f *FakeRepo) Write(license *License) error {
	files, err := license.Render()
	if err != nil {
		return err
	}

	for _, file := range files {
		f.files[file.Path] = file.Content
	}

	return nil
}

func (f *FakeRepo) reset() {
	f.files = make(map[string]string, 0)
}

func TestServiceCreate(t *testing.T) {
	type input struct {
		start       int
		end         int
		holder      string
		licenseType LicenseType
		projectName string
	}

	tests := []struct {
		name         string
		input        input
		fileToCheck  string
		errorMessage string
	}{
		{
			name: "Pass-MIT",
			input: input{
				start:       2025,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: MIT,
				projectName: "Ligen",
			},
			fileToCheck: "LICENSE",
		},
		{
			name: "Pass-Apache-2.0",
			input: input{
				start:       2025,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: APACHE_2_0,
				projectName: "Ligen",
			},
			fileToCheck: "LICENSE",
		},
		{
			name: "Pass-GNULesser-3.0",
			input: input{
				start:       2025,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: GNU_LESSER_3_0,
				projectName: "Ligen",
			},
			fileToCheck: "COPYING.LESSER",
		},
		{
			name: "Pass-Mozzila-2.0",
			input: input{
				start:       2025,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: MOZILLA_2_0,
				projectName: "Ligen",
			},
			fileToCheck: "LICENSE",
		},
		{
			name: "Pass-Boost-1.0",
			input: input{
				start:       2025,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: BOOST_1_0,
				projectName: "Ligen",
			},
			fileToCheck: "LICENSE",
		},
		{
			name: "Pass-Unlicense",
			input: input{
				start:       2025,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: UNLICENSE,
				projectName: "Ligen",
			},
			fileToCheck: "UNLICENSE",
		},
	}

	for _, tc := range tests {
		repo := NewFakeRepo()
		svc := NewService(&repo)

		t.Run(tc.name, func(t *testing.T) {
			// Given
			expected := License{
				copyright: Copyright{
					Holder:    tc.input.holder,
					StartYear: tc.input.start,
					EndYear:   tc.input.end,
				},
				licenseType: tc.input.licenseType,
			}

			if slices.Contains([]LicenseType{MOZILLA_2_0, GNU_LESSER_3_0, APACHE_2_0}, tc.input.licenseType) {
				expected.projectName = tc.input.projectName
			}

			// When
			err := svc.Create(
				tc.input.projectName,
				tc.input.holder,
				tc.input.start,
				tc.input.end,
				tc.input.licenseType,
			)

			// Then
			checkError(tc.errorMessage, err, t)

			var license License
			err = repo.Load(tc.fileToCheck, &license)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}

			if license.licenseType.RequiresCopyright() {
				if !reflect.DeepEqual(expected, license) {
					t.Errorf("Expected %v, got %v", expected, license)
				}
			}

			repo.reset()
		})
	}
}

func TestServiceGetYears(t *testing.T) {
	tests := []struct {
		name string
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {})
	}
}

func TestServiceGetLicenseType(t *testing.T) {
	tests := []struct {
		name string
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {})
	}
}

func TestServiceUpdateEndYear(t *testing.T) {
	tests := []struct {
		name string
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {})
	}
}
