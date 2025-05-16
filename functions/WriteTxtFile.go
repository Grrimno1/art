package functions

import (
	"fmt"
	"os"
)
//writes the content from string to a .txt file
func WriteTxtFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}