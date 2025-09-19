package ligen

type Repository interface {
	Load(path string, license *License) error
	Write(license *License) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) Create(projectName string, holder string, start, end int, licenseType LicenseType) error {
	license, err := New(projectName, holder, start, end, licenseType)
	if err != nil {
		return err
	}

	if err = s.repo.Write(license); err != nil {
		return err
	}

	return nil
}

type CopyrightYears struct {
	Start int
	End   int
}

func (s Service) GetYears() (CopyrightYears, error) {
	var license License
	err := s.repo.Load("LICENSE", &license)
	if err != nil {
		return CopyrightYears{}, err
	}

	return CopyrightYears{
		Start: license.copyright.StartYear,
		End:   license.copyright.EndYear,
	}, nil
}

func (s Service) GetLicenseType() (LicenseType, error) {
	var license License
	err := s.repo.Load("LICENSE", &license)
	if err != nil {
		return LicenseType(-1), err
	}

	return license.licenseType, nil
}

func (s Service) UpdateEndYear(year int) error {
	var license License
	err := s.repo.Load("LICENSE", &license)
	if err != nil {
		return err
	}

	license.SetCopyrightEndYear(year)

	return s.repo.Write(&license)
}
