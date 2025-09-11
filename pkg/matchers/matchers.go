package matchers

import (
	"errors"

	"github.com/MoonMoon1919/ligen"
)

type exists struct{}

type biggrams map[string]exists

func (b biggrams) intersection(other biggrams) int {
	intersection := 0

	for bg := range b {
		if _, ok := other[bg]; ok {
			intersection++
		}
	}

	return intersection
}

func (b biggrams) len() int {
	return len(b)
}

func newBigram(s string) biggrams {
	bg := make(biggrams)

	for i := 0; i < len(s)-1; i++ {
		bg[s[i:i+2]] = exists{}
	}

	return bg
}

func SorensonDiceCoefficient(left, right string) float64 {
	bigramsleft := newBigram(left)
	bigramsright := newBigram(right)

	intersection := bigramsleft.intersection(bigramsright)

	numerator := float64(2 * intersection)
	denominator := float64(bigramsleft.len() + bigramsright.len())

	if denominator == 0 {
		return 0.0
	}

	return numerator / denominator
}

func Match(content string) (ligen.LicenseType, error) {
	return ligen.LicenseType(-1), errors.New("Detection failed")
}
