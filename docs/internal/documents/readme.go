package documents

import (
	"os"

	"github.com/MoonMoon1919/doyoucompute"
)

func codeBlockSectionFactory(name string, intro string, samplePath string) (doyoucompute.Section, error) {
	section := doyoucompute.NewSection(name)
	section.WriteIntro().
		Text(intro)

	sample, err := os.ReadFile(samplePath)
	if err != nil {
		return doyoucompute.Section{}, err
	}

	section.WriteCodeBlock("go", []string{string(sample)}, doyoucompute.Static)

	return section, nil
}

func quickStartSection() (doyoucompute.Section, error) {
	quickStartSection := doyoucompute.NewSection("Quick Start")
	installSection := quickStartSection.CreateSection("Installation")
	installSection.WriteCodeBlock("bash", []string{"go get github.com/MoonMoon1919/ligen"}, doyoucompute.Exec)

	// Usage
	usageSection := quickStartSection.CreateSection("Usage")
	usageSection.WriteIntro().
		Text("Ligen is flexible - you can define licenses and write your own file management or use service layer in the package.")

	// Basic
	basicUsage, err := codeBlockSectionFactory(
		"Basic Usage",
		"Create a license in just. a few lines of code:",
		"./docs/internal/documents/samples/basics.go",
	)
	if err != nil {
		return doyoucompute.Section{}, err
	}
	usageSection.AddSection(basicUsage)

	// Detection
	detectionUsage, err := codeBlockSectionFactory(
		"Detecting a license",
		"Determine the type of license a repository is using:",
		"./docs/internal/documents/samples/detection.go",
	)
	if err != nil {
		return doyoucompute.Section{}, err
	}
	usageSection.AddSection(detectionUsage)

	// Service layer
	serviceUsage, err := codeBlockSectionFactory(
		"Using the service",
		"The service handles loading files into a License and contains getters and setters:",
		"./docs/internal/documents/samples/service.go",
	)
	if err != nil {
		return doyoucompute.Section{}, err
	}
	usageSection.AddSection(serviceUsage)

	return quickStartSection, nil
}

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
	quickstart, err := quickStartSection()
	if err != nil {
		return document, err
	}
	document.AddSection(quickstart)

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
