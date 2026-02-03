package wc

import (
	"bytes"
	"fmt"
	"os"
)

type WordCount struct {
}

func (wc WordCount) CountBytes(filepath string) (bytesNumber int, err error) {
	file, err := os.ReadFile(filepath)

	if err != nil {
		return 0, fmt.Errorf("failed to read file '%s': %w", filepath, err)
	}

	return len(file), nil
}

func (wc WordCount) CountLines(filepath string) (linesNumber int, err error) {
	file, err := os.ReadFile(filepath)

	if err != nil {
		return 0, fmt.Errorf("failed to read file '%s': %w", filepath, err)
	}

	return bytes.Count(file, []byte{'\n'}), nil
}
