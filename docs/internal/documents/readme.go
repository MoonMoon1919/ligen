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
	featuresList.Append("🖊️ Create a license")
	featuresList.Append("🔎 Check what license is in the repo")
	featuresList.Append("✅ Check and update copyright years")

	return document, nil
}
