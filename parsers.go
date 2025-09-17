package ligen

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	noMatchError           = errors.New("line does not match copyright pattern")
	copyrightNotFoundError = errors.New("no copyright line found")
)

func ParseDoc(content string) (Copyright, error) {
	reader := strings.NewReader(content)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		copyright, err := ParseCopyright(line)

		if err == nil {
			return copyright, nil
		}
	}

	return Copyright{}, copyrightNotFoundError
}

func ParseCopyright(line string) (Copyright, error) {
	line = strings.TrimSpace(line)

	pattern := `^Copyright\s*(?:\([Cc]\)\s*)?(\d{4})(?:-(\d{4}))?\s+(.+?)\s*$`

	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return Copyright{}, noMatchError
	}

	// Parse StartYear (should always be present)
	startYear, err := strconv.Atoi(matches[1])
	if err != nil {
		return Copyright{}, fmt.Errorf("invalid start year: %s", matches[1])
	}

	// Parse EndYear (optional)
	var endYear int
	if matches[2] != "" {
		endYear, err = strconv.Atoi(matches[2])
		if err != nil {
			return Copyright{}, fmt.Errorf("invalid end year: %s", matches[2])
		}
	}

	holder := matches[3]

	license := Copyright{
		Holder:    holder,
		StartYear: startYear,
		EndYear:   endYear,
	}

	if err := license.Validate(); err != nil {
		return Copyright{}, err
	}

	return license, nil
}
