package parsers

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/MoonMoon1919/ligen"
)

func ParseCopyright(line string) (ligen.Copyright, error) {
	pattern := `^Copyright\s*(?:\(C\)\s*)?(\d{4})(?:-(\d{4}))?\s+(.+?)\s*$`

	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return ligen.Copyright{}, fmt.Errorf("line does not match copyright pattern: %s", line)
	}

	// Parse StartYear (always present)
	startYear, err := strconv.Atoi(matches[1])
	if err != nil {
		return ligen.Copyright{}, fmt.Errorf("invalid start year: %s", matches[1])
	}

	// Parse EndYear (optional)
	var endYear int
	if matches[2] != "" {
		endYear, err = strconv.Atoi(matches[2])
		if err != nil {
			return ligen.Copyright{}, fmt.Errorf("invalid end year: %s", matches[2])
		}
	}

	holder := matches[3]

	return ligen.Copyright{
		Holder:    holder,
		StartYear: startYear,
		EndYear:   endYear,
	}, nil
}
