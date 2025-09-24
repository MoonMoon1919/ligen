# LIGEN

Go package for managing license files.

## Features

- Create licenses for your projects
- Detect and identify existing license types
- Manage copyright years and holder information
- Parse existing license files
- Template-based license generation


### Supported Licenses

- MIT
- Apache 2.0
- Mozilla Public License 2.0
- Boost Software License 1.0
- The Unlicense
- GNU Lesser General Public License 3.0


## Quick Start

### Installation

```bash
go get github.com/MoonMoon1919/ligen
```

### Usage

Ligen is flexible - you can define licenses and write your own file management or use service layer in the package.

#### Basic Usage

Create a license in just. a few lines of code:

```go
package samples

import (
	"fmt"

	"github.com/MoonMoon1919/ligen"
)

func Basic() {
	license, err := ligen.New(
		"Example",
		"J Doe",
		2025,
		0, // End year: 0 for ongoing
		ligen.MIT,
	)
	if err != nil {
		panic(err)
	}

	writeables, err := license.Render()
	if err != nil {
		panic(err)
	}

	// Do something more interesting than print
	fmt.Print(writeables)
}

```

#### Detecting a license

Determine the type of license a repository is using:

```go
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

```

#### Using the service

The service handles loading files into a License and contains getters and setters:

```go
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

```

## Contributing

See [CONTRIBUTING](./CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](./LICENSE) for details.

## Disclaimers

This work does not represent the interests or technologies of any employer, past or present. It is a personal project only.
