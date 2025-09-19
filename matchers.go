package ligen

import (
	"errors"
	"slices"
)

var (
	DetectionFailedError = errors.New("License detection failed")
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

func newBigrams(s string) biggrams {
	bg := make(biggrams)

	for i := 0; i < len(s)-1; i++ {
		bg[s[i:i+2]] = exists{}
	}

	return bg
}

func SorensonDiceCoefficient(left, right string) float64 {
	bigramsleft := newBigrams(left)
	bigramsright := newBigrams(right)

	intersection := bigramsleft.intersection(bigramsright)

	numerator := float64(2 * intersection)
	denominator := float64(bigramsleft.len() + bigramsright.len())

	if denominator == 0 {
		return 0.0
	}

	return numerator / denominator
}

type score struct {
	licenseType LicenseType
	distance    float64
}

func Match(content string, threshold float64) (LicenseType, error) {
	knownLicenseTypes := AllLicensesTypes()
	scores := make([]score, len(knownLicenseTypes))

	for idx, licenseType := range knownLicenseTypes {
		distance, err := licenseType.Compare(content, SorensonDiceCoefficient)

		if err != nil {
			return LicenseType(-1), DetectionFailedError
		}

		// If the coefficient is 1, it's an exect match
		// so don't bother iterating through the rest
		if distance == 1 {
			return licenseType, nil
		}

		scores[idx] = score{licenseType: licenseType, distance: distance}
	}

	slices.SortFunc(scores, func(a, b score) int {
		adst := a.distance
		bdst := b.distance

		if adst > bdst {
			return 1
		}

		if a == b {
			return 0
		}

		// a < b
		return -1
	})

	// The last item in the slice has the highest coefficient
	// thus is the most similar, so we select it to see if
	// the match exceeds our threshold
	bestMatch := scores[len(scores)-1]
	if bestMatch.distance < threshold {
		return LicenseType(-1), DetectionFailedError
	}

	return bestMatch.licenseType, nil
}
