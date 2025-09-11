package matchers

import (
	"bytes"
	"testing"

	"github.com/MoonMoon1919/ligen"
)

func buildInput(f ligen.WriteableGenerator, projectName string, holder string, startYear, endYear int, dest *bytes.Buffer) (string, error) {
	cr, err := ligen.NewCopyright(holder, startYear, endYear)
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

				left, err := buildInput(ligen.MITGenerator, "example", "Peanut butter", 2024, 2025, &dest)
				if err != nil {
					t.FailNow()
				}

				dest.Reset()

				right, err := buildInput(ligen.MITGenerator, "Ligen", "Max Moon", 2025, 2025, &dest)
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

				left, err := buildInput(ligen.MITGenerator, "example", "Peanut butter", 2024, 2025, &dest)
				if err != nil {
					t.FailNow()
				}

				dest.Reset()

				right, err := buildInput(ligen.ApacheGenerator, "Ligen", "Max Moon", 2025, 2025, &dest)
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
