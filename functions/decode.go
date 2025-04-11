package functions

import (
	"strconv"
	"strings"
)
func DecodeString(input string, multiline bool) string {
	if multiline {
		lines := strings.Split(input, "\n")
		var resultLines []string
		for _, line := range lines {
			result := decodeLine(line)
			if result == "Error\n" {
				return result
			}
			resultLines = append(resultLines, result)
		}
		return strings.Join(resultLines, "\n")
	}

	return decodeLine(input)
}
func decodeLine(input string) string {
	var result strings.Builder
	length := len(input)

	for i := 0; i < length; {
		if input[i] == '[' {
			end := strings.IndexByte(input[i:], ']')
			if end == -1 {
				return printError()
			}
			end += i
			content := input[i+1 : end]

			//making sure there is atleast one space.
			if !strings.Contains(content, " ") {
				return printError()
			}

			parts := strings.SplitN(content, " ", 2)
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" || strings.ContainsAny(parts[1], "[]") {
				return printError()
			}

			count, err := strconv.Atoi(parts[0])
			if err != nil {
				return printError()
			}

			result.WriteString(strings.Repeat(parts[1], count))
			i = end + 1
		} else if input[i] == ']' {
			return printError()
		} else {
			result.WriteByte(input[i])
			i++
		}
	}

	return result.String()
}

func printError() string {
	return "Error\n"
}