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

// ParseProjectNameFromNotice extracts the project name from the first line of a NOTICE file.
// The project name must be the entire first line, trimmed of whitespace.
func ParseProjectNameFromNotice(document string) (string, error) {
	// Split document into lines
	lines := strings.Split(document, "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("empty document")
	}

	// Get the first line and trim whitespace and carriage returns
	firstLine := strings.TrimSpace(strings.TrimRight(lines[0], "\r"))

	// Check if first line is empty
	if firstLine == "" {
		return "", fmt.Errorf("first line is empty")
	}

	return firstLine, nil
}

// ParseDocForCopyright scans a document line by line and returns the first valid copyright it finds.
func ParseDocForCopyright(content string) (Copyright, error) {
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

// ParseCopyright parses a copyright line and extracts the holder name and year range.
// Expects format: "Copyright [©|©] YYYY[-YYYY] Holder Name"
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
