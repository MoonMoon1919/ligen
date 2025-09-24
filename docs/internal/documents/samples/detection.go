package samples

import (
	"fmt"

	"github.com/MoonMoon1919/ligen"
)

func Detect() {
	content := `MIT License

Copyright (c) 2025 J Doe
`
	// Match takes in a threshold
	// Set intentionally very low here to avoid copying an entire license into the example
	// Package uses 0.90 when matching content loaded from file
	licenseType, err := ligen.Match(content, 0.1)
	if err != nil {
		panic(fmt.Errorf("Error: %w", err))
	}

	fmt.Printf("Detected %s", licenseType)
}
