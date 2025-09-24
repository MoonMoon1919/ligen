package samples

import (
	"fmt"

	"github.com/MoonMoon1919/ligen"
)

func Service() {
	repo := ligen.FileRepository{}
	service := ligen.NewService(repo)

	// Create and write license files
	err := service.Create("My Project", "J Doe", 2024, 0, ligen.MIT)
	if err != nil {
		panic(err)
	}

	// Read license information
	years, err := service.GetYears("LICENSE")
	if err != nil {
		panic(err)
	}
	fmt.Println(years)

	licenseType, err := service.GetLicenseType("LICENSE")
	if err != nil {
		panic(err)
	}
	fmt.Println(licenseType)

	// Update license
	err = service.UpdateEndYear("LICENSE", 2025)
	if err != nil {
		panic(err)
	}
}
