package ligen

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func checkError(expected string, received error, t *testing.T) {
	var errMsg string
	if received != nil {
		errMsg = received.Error()
	}

	if expected != errMsg {
		t.Errorf("Expected error %s, got %s", expected, errMsg)
	}
}

func TestCopyrightRender(t *testing.T) {
	type input struct {
		year   int
		holder string
	}

	tests := []struct {
		name         string
		input        input
		errorMessage string
	}{
		{
			name: "Passing",
			input: input{
				year:   2025,
				holder: "Peanut Butter",
			},
			errorMessage: "",
		},
		{
			name: "Fail-InvalidName-Empty",
			input: input{
				year:   2025,
				holder: "",
			},
			errorMessage: EmptyNameError.Error(),
		},
		{
			name: "Fail-InvalidName-TooLong",
			input: input{
				year:   2025,
				holder: strings.Repeat("a", 129),
			},
			errorMessage: NameTooLongError.Error(),
		},
		{
			name: "Fail-InvalidYear-TooLongAgo",
			input: input{
				year:   1973,
				holder: "Jelly Sandwich",
			},
			errorMessage: InvalidYearError.Error(),
		},
		{
			name: "Fail-InvalidYear-InTheFuture",
			input: input{
				year:   2026,
				holder: "Peanut Butter",
			},
			errorMessage: InvalidYearError.Error(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// When
			rendered, err := NewCopyright(tc.input.holder, tc.input.year, 0)
			checkError(tc.errorMessage, err, t)
			if tc.errorMessage != "" {
				return
			}

			// Then
			expected := Copyright{
				StartYear: tc.input.year,
				Holder:    strings.TrimSpace(tc.input.holder),
			}

			if rendered != expected {
				t.Errorf("Expected %v, got %v", expected, rendered)
			}
		})
	}
}

func TestLicenseRender(t *testing.T) {
	type input struct {
		startYear   int
		endYear     int
		holder      string
		projectName string
		licenseType LicenseType
	}

	tests := []struct {
		name            string
		input           input
		errorMessage    string
		expectedBuilder func(in input) ([]string, error)
	}{
		{
			name: "Pass-MIT-NoEnd",
			input: input{
				holder:      "Peanut Butter",
				projectName: "Cool",
				startYear:   2025,
				licenseType: MIT,
			},
			errorMessage: "",
			expectedBuilder: func(in input) ([]string, error) {
				var expected bytes.Buffer

				if err := MITTemplate.Execute(&expected, Copyright{StartYear: in.startYear, Holder: strings.TrimSpace(in.holder)}); err != nil {
					return nil, nil
				}

				return []string{expected.String()}, nil
			},
		},
		{
			name: "Pass-MIT-WithEnd",
			input: input{
				holder:      "Peanut Butter",
				projectName: "Cool",
				startYear:   2025,
				endYear:     2026,
				licenseType: MIT,
			},
			errorMessage: "",
			expectedBuilder: func(in input) ([]string, error) {
				var expected bytes.Buffer

				if err := MITTemplate.Execute(&expected, Copyright{StartYear: in.startYear, EndYear: in.endYear, Holder: strings.TrimSpace(in.holder)}); err != nil {
					return nil, nil
				}

				return []string{expected.String()}, nil
			},
		},
		{
			name: "Pass-BSL-1.0",
			input: input{
				holder:      "Peanut Butter",
				projectName: "Cool",
				startYear:   2025,
				licenseType: BOOST_1_0,
			},
			errorMessage: "",
			expectedBuilder: func(in input) ([]string, error) {
				return []string{BoostBody}, nil
			},
		},
		{
			name: "Pass-Unlicense",
			input: input{
				holder:      "Peanut Butter",
				projectName: "Cool",
				startYear:   2025,
				licenseType: UNLICENSE,
			},
			errorMessage: "",
			expectedBuilder: func(in input) ([]string, error) {
				return []string{UnlicenseBody}, nil
			},
		},
		{
			name: "Pass-Apache",
			input: input{
				holder:      "Peanut Butter",
				projectName: "Cool",
				startYear:   2025,
				licenseType: APACHE_2_0,
			},
			errorMessage: "",
			expectedBuilder: func(in input) ([]string, error) {
				expected := make([]string, 2)

				var dest bytes.Buffer
				if err := ApacheTemplate.Execute(&dest, Copyright{StartYear: in.startYear, Holder: strings.TrimSpace(in.holder)}); err != nil {
					return nil, nil
				}

				expected[0] = dest.String()

				dest.Reset()
				if err := SimpleNoticeTemplate.Execute(&dest, &NoticeInput{ProjectName: in.projectName, StartYear: in.startYear, Holder: strings.TrimSpace(in.holder)}); err != nil {
					return nil, err
				}

				expected[1] = dest.String()

				return expected, nil
			},
		},
		{
			name: "Pass-Mozilla",
			input: input{
				holder:      "Peanut Butter",
				projectName: "Cool",
				startYear:   2025,
				licenseType: MOZILLA_2_0,
			},
			errorMessage: "",
			expectedBuilder: func(in input) ([]string, error) {
				expected := make([]string, 2)
				expected[0] = MozillaLicenseBody

				// Reset the buffer so we can re-use it
				var dest bytes.Buffer
				if err := SimpleNoticeTemplate.Execute(&dest, &NoticeInput{ProjectName: in.projectName, StartYear: in.startYear, Holder: in.holder}); err != nil {
					return nil, err
				}
				expected[1] = dest.String()

				return expected, nil
			},
		},
		{
			name: "Pass-GNU Lesser GPL 3.0",
			input: input{
				holder:      "Peanut Butter",
				projectName: "Cool",
				startYear:   2025,
				licenseType: GNU_LESSER_3_0,
			},
			errorMessage: "",
			expectedBuilder: func(in input) ([]string, error) {
				expected := make([]string, 2)
				expected[0] = GNULesserLicenseBody

				// Reset the buffer so we can re-use it
				var dest bytes.Buffer
				if err := GnuLesserNoticeTemplate.Execute(&dest, &NoticeInput{ProjectName: in.projectName, StartYear: in.startYear, Holder: in.holder}); err != nil {
					return nil, err
				}
				expected[1] = dest.String()

				return expected, nil
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			license, err := New(tc.input.projectName, tc.input.holder, tc.input.startYear, tc.input.endYear, tc.input.licenseType)
			if err != nil {
				t.Errorf("Unexpected error %s", err.Error())
				return
			}

			rendered, err := license.Render()
			checkError(tc.errorMessage, err, t)
			if tc.errorMessage != "" {
				return
			}

			expected, err := tc.expectedBuilder(tc.input)
			if err != nil {
				t.Fatalf("Got error generating expected output %s", err.Error())
				return
			}

			renderedConent := make([]string, len(rendered))
			for idx, render := range rendered {
				renderedConent[idx] = render.content
			}

			if !reflect.DeepEqual(expected, renderedConent) {
				t.Errorf("Expected %s, got %s", expected, rendered)
			}
		})
	}
}
