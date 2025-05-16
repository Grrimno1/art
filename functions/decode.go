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
		//looks for start of the pattern "["
		if input[i] == '[' {
			//finds the index of the closing "]" from current position forward.
			end := strings.IndexByte(input[i:], ']')
			if end == -1 { //no closing found == malformed input.
				return printError()
			}
			//adjusting end
			end += i
			content := input[i+1 : end] //extracts content inside the brackets

			//making sure there is atleast one space.
			if !strings.Contains(content, " ") {
				return printError()
			}

			//splits content into 2 parts i.e repetition count and pattern to repeat.
			parts := strings.SplitN(content, " ", 2)
			/*
				validating:
				-both parts must exist
				-neither part can be empty.
				-pattern must not contain [ or ], to avoid nested brackets.
			*/
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" || strings.ContainsAny(parts[1], "[]") {
				return printError()
			}

			//Trying to convert repetition count to an integer. 
			count, err := strconv.Atoi(parts[0])
			if err != nil {
				return printError()
			}
			//appends repeated pattern to the result.
			result.WriteString(strings.Repeat(parts[1], count))
			i = end + 1 //adjusting index.
		} else if input[i] == ']' {
			return printError()
		} else {
			//adding normal text outside of brackets to result.
			result.WriteByte(input[i])
			i++
		}
	}

	return result.String()
}

func printError() string {
	return "Error\n"
}