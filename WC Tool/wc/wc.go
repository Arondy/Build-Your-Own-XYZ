package wc

import (
	"bytes"
	"fmt"
	"os"
	"slices"
)

type WordCount struct {
}

func (wc WordCount) GetFileContents(filepath string) (contents []byte, err error) {
	file, err := os.ReadFile(filepath)

	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s': %w", filepath, err)
	}

	return file, nil
}

func (wc WordCount) CountBytes(file []byte) int {
	return len(file)
}

func (wc WordCount) CountLines(file []byte) int {
	return bytes.Count(file, []byte{'\n'})
}

func (wc WordCount) CountWords(file []byte) int {
	file = bytes.TrimPrefix(file, []byte("\ufeff"))
	separators := []byte{' ', '\t', '\n', '\r', '\f'}
	inWord := false
	wordsNumber := 0

	for _, char := range file {
		if slices.Contains(separators, char) {
			inWord = false
		} else {
			if !inWord {
				wordsNumber++
				inWord = true
			}
		}
	}

	return wordsNumber
}

func (wc WordCount) CountCharacters(file []byte) int {
	charachterNumber := 0
	skipBytes := 0

	for _, char := range file {
		if skipBytes != 0 {
			skipBytes--
			continue
		}

		if char>>5 == 0b110 {
			skipBytes = 1
		} else if char>>4 == 0b1110 {
			skipBytes = 2
		} else if char>>3 == 0b11110 {
			skipBytes = 3
		}
		charachterNumber += 1
	}
	return charachterNumber
}
