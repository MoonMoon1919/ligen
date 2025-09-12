package parsers

import (
	"reflect"
	"testing"

	"github.com/MoonMoon1919/ligen"
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

func TestParseCopyright(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput ligen.Copyright
		errorMessage   string
	}{
		{
			name:  "Passing",
			input: "Copyright (C) 2024-2025 Max Moon",
			expectedOutput: ligen.Copyright{
				Holder:    "Max Moon",
				StartYear: 2024,
				EndYear:   2025,
			},
			errorMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cr, err := ParseCopyright(tc.input)

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
