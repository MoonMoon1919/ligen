package ligen

// Repository provides an abstraction for loading and writing licenses from different storage backends
type Repository interface {
	Load(path string, license *License) error
	Write(license *License) error
}

// Service provides business logic operations for managing licenses.
type Service struct {
	repo Repository
}

// NewService creates a new Service with the given repository.
func NewService(repo Repository) Service {
	return Service{repo: repo}
}

// Create creates a new license with the given parameters and writes it via the repository.
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

// CopyrightYears contains the start and end years of a copyright.
type CopyrightYears struct {
	Start int
	End   int
}

func (s Service) load(path string) (*License, error) {
	var license License
	err := s.repo.Load(path, &license)

	return &license, err
}

// GetYears loads a license from the given path and returns its copyright years.
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

// GetLicenseType loads a license from the given path and returns its license type.
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

// UpdateProjectName loads a license from the given path, updates its project name, and writes it back.
func (s Service) UpdateProjectName(path string, name string) error {
	return s.loadSetFlush(path, func(license *License) error {
		return license.SetProjectName(name)
	})
}

// UpdateHolder loads a license from the given path, updates its copyright holder, and writes it back.
func (s Service) UpdateHolder(path string, holder string) error {
	return s.loadSetFlush(path, func(license *License) error {
		return license.SetHolder(holder)
	})
}

// UpdateStartYear loads a license from the given path, updates its copyright start year, and writes it back.
func (s Service) UpdateStartYear(path string, year int) error {
	return s.loadSetFlush(path, func(license *License) error {
		return license.SetCopyrightStartYear(year)
	})
}

// UpdateEndYear loads a license from the given path, updates its copyright end year, and writes it back.
func (s Service) UpdateEndYear(path string, year int) error {
	return s.loadSetFlush(path, func(license *License) error {
		return license.SetCopyrightEndYear(year)
	})
}
