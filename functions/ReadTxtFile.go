package functions

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*This function reads from .txt file and returns string, when multiLine is false it will only
read one line from the .txt file, otherwise it will return entire content as string*/
func ReadTxtFile(filePath string, multiLine bool) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	//check if multiLine is enabled and read the entire file.
	if multiLine {
		var builder strings.Builder
		for scanner.Scan() {
			builder.WriteString(scanner.Text())
			builder.WriteString("\n")
		}

		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("error reading file: %w", err)
		}

		return builder.String(), nil
		//otherwise read only single line from .txt file.
	} else {
		if scanner.Scan() {
			return scanner.Text(), nil
		}
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("error reading file: %w", err)
		}
		return "", nil //file was empty
	}
}