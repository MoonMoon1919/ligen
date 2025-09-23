package documents

import "github.com/MoonMoon1919/doyoucompute"

func ReadMe() (doyoucompute.Document, error) {
	document, err := doyoucompute.NewDocument("LIGEN")
	if err != nil {
		return doyoucompute.Document{}, err
	}

	document.WriteIntro().
		Text("Go package for managing license files.")

	// Features
	featuresSection := document.CreateSection("Features")
	featuresList := featuresSection.CreateList(doyoucompute.BULLET)
	featuresList.Append("Create licenses for your projects")
	featuresList.Append("Detect and identify existing license types")
	featuresList.Append("Manage copyright years and holder information")
	featuresList.Append("Parse existing license files")
	featuresList.Append("Template-based license generation")

	// Supported Licenses
	licensesSection := featuresSection.CreateSection("Supported Licenses")
	licensesList := licensesSection.CreateList(doyoucompute.BULLET)
	licensesList.Append("MIT")
	licensesList.Append("Apache 2.0")
	licensesList.Append("Mozilla Public License 2.0")
	licensesList.Append("Boost Software License 1.0")
	licensesList.Append("The Unlicense")
	licensesList.Append("GNU Lesser General Public License 3.0")

	// Quickstart
	quickStartSection := document.CreateSection("Quick Start")
	installSection := quickStartSection.CreateSection("Installation")
	installSection.WriteCodeBlock("bash", []string{"go get github.com/MoonMoon1919/ligen"}, doyoucompute.Exec)

	// Contrib
	contributing := document.CreateSection("Contributing")
	contributing.WriteIntro().
		Text("See").
		Link("CONTRIBUTING", "./CONTRIBUTING.md").
		Text("for details.")

	// License
	licenseSection := document.CreateSection("License")
	licenseSection.WriteIntro().
		Text("MIT License - see").
		Link("LICENSE", "./LICENSE").
		Text("for details.")

	// Disclaimer
	disclaimerSection := document.CreateSection("Disclaimers")
	disclaimerSection.WriteIntro().
		Text("This work does not represent the interests or technologies of any employer, past or present.").
		Text("It is a personal project only.")

	return document, nil
}
