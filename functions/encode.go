package functions

import (
	"strings"
	"fmt"
)

func EncodeString(input string, multiline bool) string {
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
	if line == "" {
		return ""
	}

	var result strings.Builder
	prev := line[0]
	count := 1

	for i := 1; i < len(line); i++ {
		if line[i] == prev {
			count++
		} else {
			result.WriteString(fmt.Sprintf("[%d %c]", count, prev))
			prev = line[i]
			count = 1
		}
	}

	result.WriteString(fmt.Sprintf("[%d %c]", count, prev))
	return result.String()
}