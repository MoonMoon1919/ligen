package documents

import "github.com/MoonMoon1919/doyoucompute"

func ReadMe() (doyoucompute.Document, error) {
	document, err := doyoucompute.NewDocument("LIGEN")
	if err != nil {
		return doyoucompute.Document{}, err
	}

	document.WriteIntro().
		Text("Go package for managing license files.")

	return document, nil
}
