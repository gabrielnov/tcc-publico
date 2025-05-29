package service

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type FileManager interface {
	ReplaceLine(filename string, lineNumber int, newLine string) error
}

type fileManager struct{}

func NewFileManager() FileManager {
	return &fileManager{}
}

// ReplaceLine replaces a specific line in a file while maintaining indentation.
func (fm *fileManager) ReplaceLine(filename string, lineNumber int, newLine string) error {
	log.Printf("replacing line %d in file %s", lineNumber, filename)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	lineCount := 1

	newLine = strings.TrimSpace(newLine)

	for scanner.Scan() {
		line := scanner.Text()
		if lineCount == lineNumber {
			// Preserve indentation
			indentation := regexp.MustCompile(`^\s*`).FindString(line)
			lines = append(lines, indentation+newLine)
		} else {
			lines = append(lines, line)
		}
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Write back to the file
	file, err = os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}

	return writer.Flush()
}
