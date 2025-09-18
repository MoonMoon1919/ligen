package ligen

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestFileLoad(t *testing.T) {
	type input struct {
		licenseType LicenseType
		projectName string
		holder      string
		startYear   int
		endYear     int
	}

	tests := []struct {
		name         string
		input        input
		errorMessage string
	}{
		{
			name: "Passing-MIT",
			input: input{
				licenseType: MIT,
				startYear:   2024,
				endYear:     2025,
				holder:      "Peanut Butter",
				projectName: "", // unused for MIT
			},
		},
		{
			name: "Passing-Apache-2.0",
			input: input{
				licenseType: APACHE_2_0,
				startYear:   2024,
				endYear:     2025,
				holder:      "Peanut Butter",
				projectName: "Ligen",
			},
		},
		{
			name: "Passing-Mozilla-2.0",
			input: input{
				licenseType: MOZILLA_2_0,
				startYear:   2024,
				endYear:     2025,
				holder:      "Peanut Butter",
				projectName: "Ligen",
			},
		},
		{
			name: "Passing-GNU-Lesser-3.0",
			input: input{
				licenseType: GNU_LESSER_3_0,
				startYear:   2024,
				endYear:     2025,
				holder:      "Peanut Butter",
				projectName: "Ligen",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			docs := builder(t, tc.input.licenseType, tc.input.startYear, tc.input.endYear, tc.input.holder, tc.input.projectName)

			licenseLoader := func() (io.Reader, error) {
				return strings.NewReader(docs[0].Content), nil
			}

			noticeLoader := func() (io.Reader, error) {
				if len(docs) > 1 {
					return strings.NewReader(docs[1].Content), nil
				}
				return nil, nil
			}

			expected := License{
				projectName: tc.input.projectName,
				copyright: Copyright{
					Holder:    tc.input.holder,
					EndYear:   tc.input.endYear,
					StartYear: tc.input.startYear,
				},
				licenseType: tc.input.licenseType,
			}

			// WHEN
			var license License
			err := Load(&license, licenseLoader, noticeLoader)

			// THEN
			checkError(tc.errorMessage, err, t)
			if tc.errorMessage != "" {
				return
			}

			if !reflect.DeepEqual(expected, license) {
				t.Errorf("Got copyright %v, expected %v", license, expected)
			}
		})
	}
}
