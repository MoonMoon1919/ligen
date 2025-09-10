package main

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

func TestMITLicense(t *testing.T) {
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
				holder: "Peanut Butter",
				year:   2025,
			},
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// When
			var rendered bytes.Buffer
			err := MITLicense(&Copyright{Holder: tc.input.holder, Year: tc.input.year}, &rendered)
			checkError(tc.errorMessage, err, t)
			if tc.errorMessage != "" {
				return
			}

			// Then
			var expected bytes.Buffer
			MITTemplate.Execute(&expected, Copyright{Year: tc.input.year, Holder: strings.TrimSpace(tc.input.holder)})

			if rendered.String() != expected.String() {
				t.Errorf("Expected %s, got %s", expected.String(), rendered.String())
			}
		})
	}
}

func TestLicenseRender(t *testing.T) {
	type input struct {
		year   int
		holder string
	}

	tests := []struct {
		name         string
		input        input
		errorMessage string
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			license, _ := New(tc.input.holder, tc.input.year, MIT)

			rendered, err := license.Render()
			checkError(tc.errorMessage, err, t)
			if tc.errorMessage != "" {
				return
			}

			var expected bytes.Buffer
			MITTemplate.Execute(&expected, Copyright{Year: tc.input.year, Holder: strings.TrimSpace(tc.input.holder)})

			if rendered != expected.String() {
				t.Errorf("Expected %s, got %s", expected.String(), rendered)
			}
		})
	}
}
