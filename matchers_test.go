package ligen

import (
	"bytes"
	"testing"
)

func buildInput(f WriteableGenerator, projectName string, holder string, startYear, endYear int, dest *bytes.Buffer) (string, error) {
	cr, err := NewCopyright(holder, startYear, endYear)
	if err != nil {
		return "", err
	}

	writeable, err := f(&projectName, &cr, dest)
	if err != nil {
		return "", err
	}

	return writeable[0].Content, nil
}

func TestSorensonDiceCoefficient(t *testing.T) {
	tests := []struct {
		name          string
		inputBuilder  func(t *testing.T) (string, string)
		threshold     float64
		validatorFunc func(t *testing.T, coefficient, threshold float64)
	}{
		{
			name: "Passing",
			inputBuilder: func(t *testing.T) (string, string) {
				var dest bytes.Buffer

				left, err := buildInput(MITGenerator, "example", "Peanut butter", 2024, 2025, &dest)
				if err != nil {
					t.FailNow()
				}

				dest.Reset()

				right, err := buildInput(MITGenerator, "Ligen", "Max Moon", 2025, 2025, &dest)
				if err != nil {
					t.FailNow()
				}

				return left, right
			},
			threshold: 0.95,
			validatorFunc: func(t *testing.T, coefficient, threshold float64) {
				if coefficient < threshold {
					t.Errorf("Expected passing threshold to be greater than %f, got similarity of %f", threshold, coefficient)
				}
			},
		},
		{
			name: "Failing",
			inputBuilder: func(t *testing.T) (string, string) {
				var dest bytes.Buffer

				left, err := buildInput(MITGenerator, "example", "Peanut butter", 2024, 2025, &dest)
				if err != nil {
					t.FailNow()
				}

				dest.Reset()

				right, err := buildInput(ApacheGenerator, "Ligen", "Max Moon", 2025, 2025, &dest)
				if err != nil {
					t.FailNow()
				}

				return left, right
			},
			threshold: 0.60,
			validatorFunc: func(t *testing.T, coefficient, threshold float64) {
				if coefficient > threshold {
					t.Errorf("Coefficient match greater than %f, got similarity of %f", threshold, coefficient)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			left, right := tc.inputBuilder(t)

			coefficient := SorensonDiceCoefficient(left, right)

			tc.validatorFunc(t, coefficient, tc.threshold)
		})
	}
}

func TestMatch(t *testing.T) {
	passingThreshold := 0.90

	// Convenience method for building passing inputs to avoid duplicating the same
	// input builder func in every passing case
	inputBuilder := func(t *testing.T, lt LicenseType) string {
		var buf bytes.Buffer
		generatorFunc, _ := lt.GeneratorFunc()
		builtLicense, err := buildInput(generatorFunc, "Ligen", "Max Moon", 2025, 2025, &buf)
		if err != nil {
			t.FailNow()
		}

		buf.Reset()

		return builtLicense
	}

	tests := []struct {
		name         string
		inputBuilder func(t *testing.T, lt LicenseType) string
		expected     LicenseType
		threshold    float64
		errorMessage string
	}{
		{
			name:         "Pass-MatchFound-MIT",
			inputBuilder: inputBuilder,
			expected:     MIT,
			threshold:    passingThreshold,
			errorMessage: "",
		},
		{
			name:         "Pass-MatchFound-Apache",
			inputBuilder: inputBuilder,
			expected:     APACHE_2_0,
			threshold:    passingThreshold,
			errorMessage: "",
		},
		{
			name:         "Pass-MatchFound-Boost",
			inputBuilder: inputBuilder,
			expected:     BOOST_1_0,
			threshold:    passingThreshold,
			errorMessage: "",
		},
		{
			name:         "Pass-MatchFound-Unlicense",
			inputBuilder: inputBuilder,
			expected:     UNLICENSE,
			threshold:    passingThreshold,
			errorMessage: "",
		},
		{
			name:         "Pass-MatchFound-Mozilla",
			inputBuilder: inputBuilder,
			expected:     MOZILLA_2_0,
			threshold:    passingThreshold,
			errorMessage: "",
		},
		{
			name:         "Pass-MatchFound-GnuLesser",
			inputBuilder: inputBuilder,
			expected:     GNU_LESSER_3_0,
			threshold:    passingThreshold,
			errorMessage: "",
		},
		{
			name: "Fail-NoMatchFound",
			inputBuilder: func(t *testing.T, lt LicenseType) string {
				return "The dog likes to jump and play."
			},
			expected:     LicenseType(-1),
			threshold:    passingThreshold,
			errorMessage: DetectionFailedError.Error(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			discoveredType, err := Match(tc.inputBuilder(t, tc.expected), tc.threshold)

			checkError(tc.errorMessage, err, t)
			if tc.errorMessage != "" {
				return
			}

			if discoveredType != tc.expected {
				t.Errorf("Expected to find license type %d, found %d", tc.expected, discoveredType)
			}
		})
	}
}
