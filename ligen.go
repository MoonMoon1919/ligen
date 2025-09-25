package ligen

import (
	"bytes"
	"errors"
	"strings"
	"time"
)

// NoticeInput contains the information needed to generate a NOTICE file.
type NoticeInput struct {
	ProjectName string
	Holder      string
	StartYear   int
	EndYear     int
}

// General use copyright line
var (
	StartYearTooOldError        = errors.New("start year cannot be more than 50 years in the past")
	StartYearTooNewError        = errors.New("start year cannot be in the future")
	EndYearTooOldError          = errors.New("end year cannot be in the past")
	EndYearBeforeStartError     = errors.New("end year must be after start year")
	EmptyHolderError            = errors.New("holder must not be empty")
	HolderTooLongError          = errors.New("holder must be less than 128 chars")
	EmptyNameError              = errors.New("name must not be empty")
	NameTooLongError            = errors.New("name must be 128 chars")
	NameTooShortError           = errors.New("project name must have at least 1 character")
	InvalidLicenseType          = errors.New("invalid license type")
	UnsupportedLicenseTypeError = errors.New("unsupported license type")
	NoKnownTemplateError        = errors.New("no template found")
)

const (
	// MAX_NAME_LENGTH is the maximum amount of chars the holder of a copyright can contain
	// 128 picked arbitrarily, seemed reasonable
	MAX_NAME_LENGTH = 128
	// MAX_YEARS_PAST is the maximum amount of time in years that a copyright can be backdated
	// 50 picked arbitrarily, seemed reasonable
	MAX_YEARS_PAST = 50
)

// Copyright contains copyright information used to render license notices and files.
type Copyright struct {
	Holder    string
	StartYear int
	EndYear   int
}

// NewCopyright creates a new Copyright with the given holder name and year range.
// The startYear must be within the last 50 years and not in the future.
// If endYear is 0, only startYear is set. Otherwise, endYear must be after startYear and not in the past.
func NewCopyright(name string, startYear int, endYear int) (Copyright, error) {
	currentYear := time.Now().Year()
	fiftyYearsAgo := currentYear - MAX_YEARS_PAST

	if startYear > currentYear {
		return Copyright{}, StartYearTooNewError
	}

	if startYear < fiftyYearsAgo {
		return Copyright{}, StartYearTooOldError
	}

	strippedName := strings.TrimSpace(name)
	if len(strippedName) == 0 {
		return Copyright{}, EmptyNameError
	}

	if len(name) > MAX_NAME_LENGTH {
		return Copyright{}, NameTooLongError
	}

	if endYear != 0 {
		if endYear < startYear {
			return Copyright{}, EndYearBeforeStartError
		}

		if endYear < currentYear {
			return Copyright{}, EndYearTooOldError
		}

		return Copyright{Holder: name, StartYear: startYear, EndYear: endYear}, nil
	}

	return Copyright{Holder: name, StartYear: startYear}, nil
}

// Validate checks if the Copyright has a valid year range.
// Returns an error if EndYear is set and is before StartYear.
func (c *Copyright) Validate() error {
	if c.EndYear == 0 {
		return nil
	}

	if c.EndYear < c.StartYear {
		return EndYearBeforeStartError
	}

	return nil
}

// SetHolder updates the copyright holder name.
// The holder must be non-empty and less than 128 characters.
func (c *Copyright) SetHolder(holder string) error {
	if len(holder) == 0 {
		return EmptyHolderError
	}

	if len(holder) > 128 {
		return HolderTooLongError
	}

	c.Holder = holder

	return nil
}

// SetStartYear updates the start year of the copyright.
// The year must be non-zero and not after EndYear if EndYear is set.
func (c *Copyright) SetStartYear(year int) error {
	if year == 0 {
		return StartYearTooOldError
	}

	if c.EndYear != 0 && year > c.EndYear {
		return EndYearBeforeStartError
	}

	c.StartYear = year

	return nil
}

// SetEndYear updates the end year of the copyright.
// The year must be after or equal to StartYear.
func (c *Copyright) SetEndYear(year int) error {
	if year < c.StartYear {
		return EndYearBeforeStartError
	}

	c.EndYear = year

	return nil
}

// MITGenerator generates license files for the MIT license.
func MITGenerator(projectName *string, cr *Copyright, dest *bytes.Buffer) ([]Writeable, error) {
	if err := MITTemplate.Execute(dest, cr); err != nil {
		return nil, err
	}

	writeableSlice := make([]Writeable, 1)
	writeableSlice[0] = Writeable{Content: dest.String(), Path: "LICENSE"}
	dest.Reset()

	return writeableSlice, nil
}

// BoostGenerator generates license files for the Boost Software License 1.0.
func BoostGenerator(projectName *string, cr *Copyright, dest *bytes.Buffer) ([]Writeable, error) {
	writeableSlice := make([]Writeable, 1)
	writeableSlice[0] = Writeable{Content: BoostBody, Path: "LICENSE"}
	dest.Reset()

	return writeableSlice, nil
}

// UnlicenseGenerator generates license files for the Unlicense.
func UnlicenseGenerator(projectName *string, cr *Copyright, dest *bytes.Buffer) ([]Writeable, error) {
	writeableSlice := make([]Writeable, 1)
	writeableSlice[0] = Writeable{Content: UnlicenseBody, Path: "UNLICENSE"}
	dest.Reset()

	return writeableSlice, nil
}

// ApacheGenerator generates license files for the Apache License 2.0.
func ApacheGenerator(projectName *string, cr *Copyright, dest *bytes.Buffer) ([]Writeable, error) {
	if err := ApacheTemplate.Execute(dest, cr); err != nil {
		return nil, err
	}

	writeableSlice := make([]Writeable, 2)
	writeableSlice[0] = Writeable{Content: dest.String(), Path: "LICENSE"}

	// Reset the buffer so we can re-use it
	dest.Reset()
	if err := SimpleNoticeTemplate.Execute(dest, &NoticeInput{ProjectName: *projectName, StartYear: cr.StartYear, EndYear: cr.EndYear, Holder: cr.Holder}); err != nil {
		return nil, err
	}
	writeableSlice[1] = Writeable{Content: dest.String(), Path: "NOTICE"}
	dest.Reset()

	return writeableSlice, nil
}

// MozillaGenerator generates license files for the Mozilla Public License 2.0.
func MozillaGenerator(projectName *string, cr *Copyright, dest *bytes.Buffer) ([]Writeable, error) {
	writeableSlice := make([]Writeable, 2)
	writeableSlice[0] = Writeable{Content: MozillaLicenseBody, Path: "LICENSE"}

	// Reset the buffer so we can re-use it
	dest.Reset()
	if err := SimpleNoticeTemplate.Execute(dest, &NoticeInput{ProjectName: *projectName, StartYear: cr.StartYear, EndYear: cr.EndYear, Holder: cr.Holder}); err != nil {
		return nil, err
	}
	writeableSlice[1] = Writeable{Content: dest.String(), Path: "NOTICE"}
	dest.Reset()

	return writeableSlice, nil
}

// GNULesserGenerator generates license files for the GNU Lesser General Public License 3.0.
func GNULesserGenerator(projectName *string, cr *Copyright, dest *bytes.Buffer) ([]Writeable, error) {
	writeableSlice := make([]Writeable, 2)
	writeableSlice[0] = Writeable{Content: GNULesserLicenseBody, Path: "COPYING.LESSER"}

	dest.Reset()
	if err := GnuLesserNoticeTemplate.Execute(dest, &NoticeInput{ProjectName: *projectName, StartYear: cr.StartYear, EndYear: cr.EndYear, Holder: cr.Holder}); err != nil {
		return nil, err
	}
	writeableSlice[1] = Writeable{Content: dest.String(), Path: "NOTICE"}
	dest.Reset()

	return writeableSlice, nil
}

// LicenseType represents a supported open source license.
type LicenseType int

const (
	MIT LicenseType = iota + 1
	BOOST_1_0
	UNLICENSE
	APACHE_2_0
	MOZILLA_2_0
	GNU_LESSER_3_0
)

// AllLicensesTypes returns a slice of all supported license types.
func AllLicensesTypes() []LicenseType {
	return []LicenseType{
		MIT,
		BOOST_1_0,
		UNLICENSE,
		APACHE_2_0,
		MOZILLA_2_0,
		GNU_LESSER_3_0,
	}
}

// Writeable contains license file content and its destination path.
type Writeable struct {
	Content string
	Path    string
}

// WriteableGenerator is a function that generates license files for a given license type.
type WriteableGenerator func(projectName *string, cr *Copyright, dest *bytes.Buffer) ([]Writeable, error)

// Template returns the license text template for this license type.
func (lt LicenseType) Template() (string, error) {
	switch lt {
	case MIT:
		return MitTemplateBody, nil
	case BOOST_1_0:
		return BoostBody, nil
	case UNLICENSE:
		return UnlicenseBody, nil
	case APACHE_2_0:
		return ApacheTemplateBody, nil
	case MOZILLA_2_0:
		return MozillaLicenseBody, nil
	case GNU_LESSER_3_0:
		return GNULesserLicenseBody, nil
	default:
		return "", NoKnownTemplateError
	}
}

// String returns the string representation of the license type.
func (lt LicenseType) String() string {
	switch lt {
	case MIT:
		return "MIT"
	case BOOST_1_0:
		return "BOOST_1_0"
	case UNLICENSE:
		return "UNLICENSE"
	case APACHE_2_0:
		return "APACHE_2_0"
	case MOZILLA_2_0:
		return "MOZILLA_2_0"
	case GNU_LESSER_3_0:
		return "GNU_LESSER_3_0"
	default:
		return "UNKNOWN"
	}
}

// LicenseTypeFromString parses a license type from its string representation.
// The input is case-insensitive.
func LicenseTypeFromString(licenseType string) (LicenseType, error) {
	licenseType = strings.ToUpper(licenseType)

	switch licenseType {
	case "MIT":
		return MIT, nil
	case "BOOST":
		return BOOST_1_0, nil
	case "UNLICENSE":
		return UNLICENSE, nil
	case "APACHE":
		return APACHE_2_0, nil
	case "MOZILLA":
		return MOZILLA_2_0, nil
	case "GNU_LESSER":
		return GNU_LESSER_3_0, nil
	default:
		return LicenseType(-1), InvalidLicenseType
	}
}

// Compare compares the license template text with the provided text using the given comparison function.
// Returns the similarity score from the comparison function.
func (lt LicenseType) Compare(left string, comparisonFunc func(left, right string) float64) (float64, error) {
	tmp, err := lt.Template()
	if err != nil {
		return 0.0, err
	}

	return comparisonFunc(left, tmp), nil
}

// GeneratorFunc returns the generator function for this license type.
func (lt LicenseType) GeneratorFunc() (WriteableGenerator, error) {
	switch lt {
	case MIT:
		return MITGenerator, nil
	case BOOST_1_0:
		return BoostGenerator, nil
	case UNLICENSE:
		return UnlicenseGenerator, nil
	case APACHE_2_0:
		return ApacheGenerator, nil
	case MOZILLA_2_0:
		return MozillaGenerator, nil
	case GNU_LESSER_3_0:
		return GNULesserGenerator, nil
	default:
		return nil, UnsupportedLicenseTypeError
	}
}

// RequiresNotice returns true if this license type requires a NOTICE file.
func (lt LicenseType) RequiresNotice() bool {
	switch lt {
	case MOZILLA_2_0, GNU_LESSER_3_0, APACHE_2_0:
		return true
	default:
		return false
	}
}

// RequiresCopyright returns true if this license type requires copyright information.
func (lt LicenseType) RequiresCopyright() bool {
	switch lt {
	case UNLICENSE, BOOST_1_0:
		return false
	default:
		return true
	}
}

// License represents a complete license configuration with project name, copyright, and license type.
type License struct {
	projectName string
	copyright   Copyright
	licenseType LicenseType
}

func validateProjectName(name string) error {
	if len(name) == 0 {
		return NameTooShortError
	}

	if len(name) > 128 {
		return NameTooLongError
	}

	return nil
}

// New creates a new License with the given project name, copyright holder, year range, and license type.
// The project name must be 1-128 characters after trimming whitespace.
func New(projectName string, holder string, startYear int, endYear int, licenseType LicenseType) (*License, error) {
	projectName = strings.TrimSpace(projectName)

	if err := validateProjectName(projectName); err != nil {
		return &License{}, err
	}

	copyright, err := NewCopyright(holder, startYear, endYear)
	if err != nil {
		return &License{}, err
	}

	return &License{
		projectName: projectName,
		copyright:   copyright,
		licenseType: licenseType,
	}, nil
}

// Render generates the license files for this License.
// Returns a slice of Writeable containing the file content and paths where they should be written.
func (l *License) Render() ([]Writeable, error) {
	generatorFunc, err := l.licenseType.GeneratorFunc()
	if err != nil {
		return nil, err
	}

	var content bytes.Buffer

	writeable, err := generatorFunc(&l.projectName, &l.copyright, &content)

	if err != nil {
		return nil, err
	}

	return writeable, nil
}

// SetHolder updates the copyright holder name.
func (l *License) SetHolder(holder string) error {
	return l.copyright.SetHolder(holder)
}

// SetProjectName updates the project name.
// The name must be 1-128 characters.
func (l *License) SetProjectName(name string) error {
	if err := validateProjectName(name); err != nil {
		return err
	}

	l.projectName = name

	return nil
}

// SetCopyrightEndYear updates the copyright end year.
func (l *License) SetCopyrightEndYear(year int) error {
	return l.copyright.SetEndYear(year)
}

// SetCopyrightStartYear updates the copyright start year.
func (l *License) SetCopyrightStartYear(year int) error {
	return l.copyright.SetStartYear(year)
}

// SetLicenseType updates the license type.
func (l *License) SetLicenseType(licenseType LicenseType) error {
	l.licenseType = licenseType
	return nil
}
