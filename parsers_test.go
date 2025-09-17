package ligen

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParseDoc(t *testing.T) {
	buildInput := func(f WriteableGenerator, projectName string, holder string, startYear, endYear int, dest *bytes.Buffer) ([]Writeable, error) {
		cr, err := NewCopyright(holder, startYear, endYear)
		if err != nil {
			return nil, err
		}

		writeable, err := f(&projectName, &cr, dest)
		if err != nil {
			return nil, err
		}

		return writeable, nil
	}

	builder := func(t *testing.T, lt LicenseType, startYear, endYear int, holder string) []Writeable {
		var buf bytes.Buffer
		generatorFunc, _ := lt.GeneratorFunc()
		builtLicense, err := buildInput(generatorFunc, "Ligen", holder, startYear, endYear, &buf)
		if err != nil {
			t.FailNow()
			return nil
		}

		buf.Reset()

		return builtLicense
	}

	type input struct {
		holder      string
		startYear   int
		endYear     int
		licenseType LicenseType
	}

	commonInput := input{
		holder:    "Max Moon",
		startYear: 2024,
		endYear:   2025,
	}

	tests := []struct {
		name         string
		inputBuilder func(t *testing.T, startYear, endYear int, holder string) string
		input        input
		errorMessage string
	}{
		{
			name: "Passing-MIT",
			inputBuilder: func(t *testing.T, startYear, endYear int, holder string) string {
				docs := builder(t, MIT, startYear, endYear, holder)
				return docs[0].Content
			},
			input:        commonInput,
			errorMessage: "",
		},
		{
			name: "Passing-Apache_2_0",
			inputBuilder: func(t *testing.T, startYear, endYear int, holder string) string {
				docs := builder(t, APACHE_2_0, startYear, endYear, holder)
				return docs[0].Content
			},
			input:        commonInput,
			errorMessage: "",
		},
		{
			name: "Passing-Mozilla",
			inputBuilder: func(t *testing.T, startYear, endYear int, holder string) string {
				docs := builder(t, MOZILLA_2_0, startYear, endYear, holder)
				return docs[1].Content
			},
			input:        commonInput,
			errorMessage: "",
		},
		{
			name: "Passing-GNULesser",
			inputBuilder: func(t *testing.T, startYear, endYear int, holder string) string {
				docs := builder(t, GNU_LESSER_3_0, startYear, endYear, holder)
				return docs[1].Content
			},
			input:        commonInput,
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := tc.inputBuilder(t,
				tc.input.startYear,
				tc.input.endYear,
				tc.input.holder,
			)

			parsedCopyright, err := ParseDoc(input)

			checkError(tc.errorMessage, err, t)
			if tc.errorMessage != "" {
				return
			}

			expectedOutput := Copyright{
				Holder:    tc.input.holder,
				StartYear: tc.input.startYear,
				EndYear:   tc.input.endYear,
			}

			if !reflect.DeepEqual(expectedOutput, parsedCopyright) {
				t.Errorf("Expected %v, got %v", expectedOutput, parsedCopyright)
			}
		})
	}
}

func TestParseCopyright(t *testing.T) {
	tests := []struct {
		name           string
		inputBuilder   func() string
		expectedOutput Copyright
		errorMessage   string
	}{
		{
			name: "Passing",
			inputBuilder: func() string {
				return "Copyright (C) 2024-2025 Max Moon"
			},
			expectedOutput: Copyright{
				Holder:    "Max Moon",
				StartYear: 2024,
				EndYear:   2025,
			},
			errorMessage: "",
		},
		{
			name: "Failing",
			inputBuilder: func() string {
				return "Copyright (C) LALA-L000 Max Moon"
			},
			expectedOutput: Copyright{},
			errorMessage:   noMatchError.Error(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cr, err := ParseCopyright(tc.inputBuilder())

			checkError(tc.errorMessage, err, t)
			if tc.errorMessage != "" {
				return
			}

			if !reflect.DeepEqual(tc.expectedOutput, cr) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, cr)
			}
		})
	}
}
