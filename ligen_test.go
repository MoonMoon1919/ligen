package ligen

import (
	"bytes"
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
			rendered, err := NewCopyright(tc.input.holder, tc.input.year)
			checkError(tc.errorMessage, err, t)
			if tc.errorMessage != "" {
				return
			}

			// Then
			expected := Copyright{
				Year:   tc.input.year,
				Holder: strings.TrimSpace(tc.input.holder),
			}

			if rendered != expected {
				t.Errorf("Expected %v, got %v", expected, rendered)
			}
		})
	}
}

func TestLicenseRender(t *testing.T) {
	type input struct {
		year        int
		holder      string
		licenseType LicenseType
	}

	tests := []struct {
		name            string
		input           input
		errorMessage    string
		expectedBuilder func(in input) (string, error)
	}{
		{
			name: "Pass-MIT",
			input: input{
				holder:      "Peanut Butter",
				year:        2025,
				licenseType: MIT,
			},
			errorMessage: "",
			expectedBuilder: func(in input) (string, error) {
				var expected bytes.Buffer

				if err := MITTemplate.Execute(&expected, Copyright{Year: in.year, Holder: strings.TrimSpace(in.holder)}); err != nil {
					return "", nil
				}

				return expected.String(), nil
			},
		},
		{
			name: "Pass-BSL-1.0",
			input: input{
				holder:      "Peanut Butter",
				year:        2025,
				licenseType: BOOST_1_0,
			},
			errorMessage: "",
			expectedBuilder: func(in input) (string, error) {
				return BoostBody, nil
			},
		},
		{
			name: "Pass-Unlicense",
			input: input{
				holder:      "Peanut Butter",
				year:        2025,
				licenseType: UNLICENSE,
			},
			errorMessage: "",
			expectedBuilder: func(in input) (string, error) {
				return UnlicenseBody, nil
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			license, _ := New(tc.input.holder, tc.input.year, tc.input.licenseType)

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

			if rendered[len(rendered)-1].content != expected {
				t.Errorf("Expected %s, got %s", expected, rendered)
			}
		})
	}
}
