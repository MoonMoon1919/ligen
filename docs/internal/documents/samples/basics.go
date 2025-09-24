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
