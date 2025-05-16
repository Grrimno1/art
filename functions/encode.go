package functions

import (
	"fmt"
	"strings"
)

func EncodeString(input string, multiline bool) string {
	if strings.ContainsAny(input, "[]") {
		return "Error\n"
	}
	if multiline {
		lines := strings.Split(input, "\n")
		var resultLines []string 
		for _, line := range lines {
			resultLines = append(resultLines, encodeLine(line))
		}
		return strings.Join(resultLines, "\n")
	}
	
	return encodeLine(input)
}

func encodeLine(line string) string {
	var result strings.Builder
	n := len(line)
	i := 0

	for i < n {
		maxPatternLen := (n - i) / 2 //maximum pattern size to consider
		found := false

		for patternLen := 1; patternLen <= maxPatternLen; patternLen++ {
			pattern := line[i : i+patternLen]
			repeats := 1
			for j := i + patternLen; j+patternLen <= n && line[j:j+patternLen] == pattern; j += patternLen {
				repeats++
			}

			if repeats > 1 {
				result.WriteString(fmt.Sprintf("[%d %s]", repeats, pattern))
				i += repeats * patternLen
				found = true
				break
			}
		}

		if !found {
			//no repeating pattern found, encode single character.
			result.WriteString(fmt.Sprintf("[1 %c]", line[i]))
			i++
		}
	}

	return result.String()
}

/*
func encodeLine(line string) string {
	var result strings.Builder
	n := len(line)

	for i := 0; i < n; {
		ch := line[i]
		j := i + 1
		for j < n && line[j] == ch {
			j++
		}
		count := j - i
		result.WriteString(fmt.Sprintf("[%d %c]", count, ch))
		i = j
	}

	return result.String()
}
*/