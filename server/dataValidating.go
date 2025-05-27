package server

// helper functions for some data validation.

import (
	"regexp"
	"strconv"
)

// validates the input to check if it would exceed the set len(string) limit if decoded.
func decodedExceedsLimit(input string, limit int) bool {
	pattern := regexp.MustCompile(`\[(\d+)\s([^\[\]]+)\]`)
	matches := pattern.FindAllStringSubmatchIndex(input, -1)

	decodedLen := 0
	prevEnd := 0

	for _, match := range matches {
		start, end := match[0], match[1]
		countStr := input[match[2]:match[3]]
		patternStr := input[match[4]:match[5]]

		decodedLen += len(input[prevEnd:start])

		count, err := strconv.Atoi(countStr)
		if err != nil {
			continue // skip invalid block
		}
		decodedLen += count * len(patternStr)

		// too long, return true.
		if decodedLen > limit {
			return true
		}

		prevEnd = end
	}

	//adding remaining literal characters after last match
	decodedLen += len(input[prevEnd:])

	return decodedLen > limit
}
// for other inputs this should be fine for now.
func inputExceedsLimit(input string, limit int) bool {
	return len(input) > limit
}
