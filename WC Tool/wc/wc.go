package wc

import (
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
