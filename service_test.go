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
	ll := func() (io.Reader, func() error, error) {
		return strings.NewReader(f.files[path]), func() error { return nil }, nil
	}

	nl := func() (io.Reader, func() error, error) {
		return strings.NewReader(f.files["NOTICE"]), func() error { return nil }, nil
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
		})
	}
}

func TestServiceGetYears(t *testing.T) {
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
	}

	for _, tc := range tests {
		repo := NewFakeRepo()
		svc := NewService(&repo)

		t.Run(tc.name, func(t *testing.T) {
			expected := CopyrightYears{Start: tc.input.start, End: tc.input.end}

			err := svc.Create(
				tc.input.projectName,
				tc.input.holder,
				tc.input.start,
				tc.input.end,
				tc.input.licenseType,
			)
			if err != nil {
				t.FailNow()
			}

			years, err := svc.GetYears(tc.fileToCheck)
			checkError(tc.errorMessage, err, t)
			if tc.errorMessage != "" {
				return
			}

			if !reflect.DeepEqual(years, expected) {
				t.Errorf("Expected %v, got %v", expected, years)
			}
		})
	}
}

func TestServiceGetLicenseType(t *testing.T) {
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
			err := svc.Create(
				tc.input.projectName,
				tc.input.holder,
				tc.input.start,
				tc.input.end,
				tc.input.licenseType,
			)
			if err != nil {
				t.FailNow()
			}

			lt, err := svc.GetLicenseType(tc.fileToCheck)
			checkError(tc.errorMessage, err, t)
			if tc.errorMessage != "" {
				return
			}

			if lt != tc.input.licenseType {
				t.Errorf("Expected %s, got %s", tc.input.licenseType.String(), lt.String())
			}
		})
	}
}

func TestServiceUpdateHolder(t *testing.T) {
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
		newHolder    string
		fileToCheck  string
		errorMessage string
	}{
		{
			name: "Pass-MIT",
			input: input{
				start:       2023,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: MIT,
				projectName: "Ligen",
			},
			fileToCheck: "LICENSE",
			newHolder:   "Jelly",
		},
		{
			name: "Pass-Apache-2.0",
			input: input{
				start:       2023,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: APACHE_2_0,
				projectName: "Ligen",
			},
			fileToCheck: "LICENSE",
			newHolder:   "Jelly",
		},
		{
			name: "Pass-GNULesser-3.0",
			input: input{
				start:       2023,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: GNU_LESSER_3_0,
				projectName: "Ligen",
			},
			fileToCheck: "COPYING.LESSER",
			newHolder:   "Jelly",
		},
		{
			name: "Pass-Mozzila-2.0",
			input: input{
				start:       2023,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: MOZILLA_2_0,
				projectName: "Ligen",
			},
			fileToCheck: "LICENSE",
			newHolder:   "Jelly",
		},
	}

	for _, tc := range tests {
		repo := NewFakeRepo()
		svc := NewService(&repo)

		t.Run(tc.name, func(t *testing.T) {
			err := svc.Create(
				tc.input.projectName,
				tc.input.holder,
				tc.input.start,
				tc.input.end,
				tc.input.licenseType,
			)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}

			err = svc.UpdateHolder(tc.fileToCheck, tc.newHolder)
			checkError(tc.errorMessage, err, t)
			if tc.errorMessage != "" {
				return
			}

			var license License
			err = repo.Load(tc.fileToCheck, &license)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}

			if license.copyright.Holder != tc.newHolder {
				t.Errorf("Expected %s, got %s", license.copyright.Holder, tc.newHolder)
			}
		})
	}
}

func TestServiceUpdateProjectName(t *testing.T) {
	type input struct {
		start       int
		end         int
		holder      string
		licenseType LicenseType
		projectName string
	}

	tests := []struct {
		name           string
		input          input
		newProjectName string
		fileToCheck    string
		errorMessage   string
	}{
		{
			name: "Pass-Apache-2.0",
			input: input{
				start:       2023,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: APACHE_2_0,
				projectName: "Ligen",
			},
			fileToCheck:    "LICENSE",
			newProjectName: "license-generator",
		},
		{
			name: "Pass-GNULesser-3.0",
			input: input{
				start:       2023,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: GNU_LESSER_3_0,
				projectName: "Ligen",
			},
			fileToCheck:    "COPYING.LESSER",
			newProjectName: "license-generator",
		},
		{
			name: "Pass-Mozzila-2.0",
			input: input{
				start:       2023,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: MOZILLA_2_0,
				projectName: "Ligen",
			},
			fileToCheck:    "LICENSE",
			newProjectName: "license-generator",
		},
	}

	for _, tc := range tests {
		repo := NewFakeRepo()
		svc := NewService(&repo)

		t.Run(tc.name, func(t *testing.T) {
			err := svc.Create(
				tc.input.projectName,
				tc.input.holder,
				tc.input.start,
				tc.input.end,
				tc.input.licenseType,
			)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}

			err = svc.UpdateProjectName(tc.fileToCheck, tc.newProjectName)
			checkError(tc.errorMessage, err, t)
			if tc.errorMessage != "" {
				return
			}

			var license License
			err = repo.Load(tc.fileToCheck, &license)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}

			if license.projectName != tc.newProjectName {
				t.Errorf("Expected %s, got %s", license.projectName, tc.newProjectName)
			}
		})
	}
}

func TestServiceUpdateStartYear(t *testing.T) {
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
		newStartYear int
		fileToCheck  string
		errorMessage string
	}{
		{
			name: "Pass-MIT",
			input: input{
				start:       2023,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: MIT,
				projectName: "Ligen",
			},
			fileToCheck:  "LICENSE",
			newStartYear: 2022,
		},
		{
			name: "Pass-Apache-2.0",
			input: input{
				start:       2023,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: APACHE_2_0,
				projectName: "Ligen",
			},
			fileToCheck:  "LICENSE",
			newStartYear: 2022,
		},
		{
			name: "Pass-GNULesser-3.0",
			input: input{
				start:       2023,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: GNU_LESSER_3_0,
				projectName: "Ligen",
			},
			fileToCheck:  "COPYING.LESSER",
			newStartYear: 2022,
		},
		{
			name: "Pass-Mozzila-2.0",
			input: input{
				start:       2023,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: MOZILLA_2_0,
				projectName: "Ligen",
			},
			fileToCheck:  "LICENSE",
			newStartYear: 2022,
		},
	}

	for _, tc := range tests {
		repo := NewFakeRepo()
		svc := NewService(&repo)

		t.Run(tc.name, func(t *testing.T) {
			err := svc.Create(
				tc.input.projectName,
				tc.input.holder,
				tc.input.start,
				tc.input.end,
				tc.input.licenseType,
			)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}

			err = svc.UpdateStartYear(tc.fileToCheck, tc.newStartYear)
			checkError(tc.errorMessage, err, t)
			if tc.errorMessage != "" {
				return
			}

			var license License
			err = repo.Load(tc.fileToCheck, &license)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}

			if license.copyright.StartYear != tc.newStartYear {
				t.Errorf("Expected %d, got %d", license.copyright.StartYear, tc.newStartYear)
			}
		})
	}
}

func TestServiceUpdateEndYear(t *testing.T) {
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
		newEndYear   int
		fileToCheck  string
		errorMessage string
	}{
		{
			name: "Pass-MIT",
			input: input{
				start:       2023,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: MIT,
				projectName: "Ligen",
			},
			fileToCheck: "LICENSE",
			newEndYear:  2026,
		},
		{
			name: "Pass-Apache-2.0",
			input: input{
				start:       2023,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: APACHE_2_0,
				projectName: "Ligen",
			},
			fileToCheck: "LICENSE",
			newEndYear:  2026,
		},
		{
			name: "Pass-GNULesser-3.0",
			input: input{
				start:       2023,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: GNU_LESSER_3_0,
				projectName: "Ligen",
			},
			fileToCheck: "COPYING.LESSER",
			newEndYear:  2026,
		},
		{
			name: "Pass-Mozzila-2.0",
			input: input{
				start:       2023,
				end:         2025,
				holder:      "Peanut Butter",
				licenseType: MOZILLA_2_0,
				projectName: "Ligen",
			},
			fileToCheck: "LICENSE",
			newEndYear:  2026,
		},
	}

	for _, tc := range tests {
		repo := NewFakeRepo()
		svc := NewService(&repo)

		t.Run(tc.name, func(t *testing.T) {
			err := svc.Create(
				tc.input.projectName,
				tc.input.holder,
				tc.input.start,
				tc.input.end,
				tc.input.licenseType,
			)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}

			err = svc.UpdateEndYear(tc.fileToCheck, tc.newEndYear)
			checkError(tc.errorMessage, err, t)
			if tc.errorMessage != "" {
				return
			}

			var license License
			err = repo.Load(tc.fileToCheck, &license)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}

			if license.copyright.EndYear != tc.newEndYear {
				t.Errorf("Expected %d, got %d", license.copyright.EndYear, tc.newEndYear)
			}
		})
	}
}
