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

func (s Service) load(path string) (*License, error) {
	var license License
	err := s.repo.Load(path, &license)

	return &license, err
}

func (s Service) GetYears(path string) (CopyrightYears, error) {
	license, err := s.load(path)
	if err != nil {
		return CopyrightYears{}, err
	}

	return CopyrightYears{
		Start: license.copyright.StartYear,
		End:   license.copyright.EndYear,
	}, nil
}

func (s Service) GetLicenseType(path string) (LicenseType, error) {
	license, err := s.load(path)
	if err != nil {
		return LicenseType(-1), err
	}

	return license.licenseType, nil
}

func (s Service) loadSetFlush(path string, op func(license *License) error) error {
	license, err := s.load(path)
	if err != nil {
		return err
	}

	if err = op(license); err != nil {
		return err
	}

	return s.repo.Write(license)
}

func (s Service) UpdateProjectName(path string, name string) error {
	return s.loadSetFlush(path, func(license *License) error {
		return license.SetProjectName(name)
	})
}

func (s Service) UpdateHolder(path string, holder string) error {
	return s.loadSetFlush(path, func(license *License) error {
		return license.SetHolder(holder)
	})
}

func (s Service) UpdateStartYear(path string, year int) error {
	return s.loadSetFlush(path, func(license *License) error {
		return license.SetCopyrightStartYear(year)
	})
}

func (s Service) UpdateEndYear(path string, year int) error {
	return s.loadSetFlush(path, func(license *License) error {
		return license.SetCopyrightEndYear(year)
	})
}
