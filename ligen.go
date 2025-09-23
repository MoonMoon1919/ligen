package ligen

import (
	"bytes"
	"errors"
	"strings"
	"time"
)

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

type Copyright struct {
	Holder    string
	StartYear int
	EndYear   int
}

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

func (c *Copyright) Validate() error {
	if c.EndYear == 0 {
		return nil
	}

	if c.EndYear < c.StartYear {
		return EndYearBeforeStartError
	}

	return nil
}

func (c *Copyright) SetStartYear(year int) error {
	if c.EndYear != 0 && year > c.EndYear {
		return EndYearBeforeStartError
	}

	c.StartYear = year

	return nil
}

func (c *Copyright) SetEndYear(year int) error {
	if year < c.StartYear {
		return EndYearBeforeStartError
	}

	c.EndYear = year

	return nil
}

// License Generators
func MITGenerator(projectName *string, cr *Copyright, dest *bytes.Buffer) ([]Writeable, error) {
	if err := MITTemplate.Execute(dest, cr); err != nil {
		return nil, err
	}

	writeableSlice := make([]Writeable, 1)
	writeableSlice[0] = Writeable{Content: dest.String(), Path: "LICENSE"}
	dest.Reset()

	return writeableSlice, nil
}

func BoostGenerator(projectName *string, cr *Copyright, dest *bytes.Buffer) ([]Writeable, error) {
	writeableSlice := make([]Writeable, 1)
	writeableSlice[0] = Writeable{Content: BoostBody, Path: "LICENSE"}
	dest.Reset()

	return writeableSlice, nil
}

func UnlicenseGenerator(projectName *string, cr *Copyright, dest *bytes.Buffer) ([]Writeable, error) {
	writeableSlice := make([]Writeable, 1)
	writeableSlice[0] = Writeable{Content: UnlicenseBody, Path: "UNLICENSE"}
	dest.Reset()

	return writeableSlice, nil
}

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

// License stuff

type LicenseType int

const (
	MIT LicenseType = iota + 1
	BOOST_1_0
	UNLICENSE
	APACHE_2_0
	MOZILLA_2_0
	GNU_LESSER_3_0
)

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

type Writeable struct {
	Content string
	Path    string
}

type WriteableGenerator func(projectName *string, cr *Copyright, dest *bytes.Buffer) ([]Writeable, error)

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

func (lt LicenseType) Compare(left string, comparisonFunc func(left, right string) float64) (float64, error) {
	tmp, err := lt.Template()
	if err != nil {
		return 0.0, err
	}

	return comparisonFunc(left, tmp), nil
}

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

func (lt LicenseType) RequiresNotice() bool {
	switch lt {
	case MOZILLA_2_0, GNU_LESSER_3_0, APACHE_2_0:
		return true
	default:
		return false
	}
}

func (lt LicenseType) RequiresCopyright() bool {
	switch lt {
	case UNLICENSE, BOOST_1_0:
		return false
	default:
		return true
	}
}

type License struct {
	projectName string
	copyright   Copyright
	licenseType LicenseType
}

func New(projectName string, holder string, startYear int, endYear int, licenseType LicenseType) (*License, error) {
	projectName = strings.TrimSpace(projectName)

	if len(projectName) == 0 {
		return &License{}, NameTooShortError
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

func (l *License) SetHolder(holder string) error {
	l.copyright.Holder = holder
	return nil
}

func (l *License) SetCopyrightEndYear(year int) error {
	return l.copyright.SetEndYear(year)
}

func (l *License) SetCopyrightStartYear(year int) error {
	return l.copyright.SetStartYear(year)
}

func (l *License) SetLicenseType(licenseType LicenseType) error {
	l.licenseType = licenseType
	return nil
}
