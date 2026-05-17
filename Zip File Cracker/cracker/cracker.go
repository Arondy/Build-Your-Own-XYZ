package cracker

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/yeka/zip"
)

var zipHeader = [4]byte{80, 75, 03, 04}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func isZipFile(filename string) bool {
	if !fileExists(filename) {
		return false
	}

	file, err := os.Open(filename)
	if err != nil {
		return false
	}

	reader := bufio.NewReader(file)
	header := make([]byte, 4)
	_, err = reader.Read(header)
	if err != nil {
		return false
	}

	for i := range header {
		if header[i] != zipHeader[i] {
			return false
		}
	}

	return true
}

type Cracker struct {
	reader *zip.ReadCloser
}

func NewCracker(filename string) (Cracker, error) {
	if !isZipFile(filename) {
		return Cracker{}, fmt.Errorf("'%s' isn't a valid ZIP file!", filename)
	}
	reader, err := zip.OpenReader(filename)
	if err != nil {
		return Cracker{}, err
	}
	return Cracker{reader}, nil
}

func (c Cracker) CheckPassword(password string) error {
	for _, file := range c.reader.File {
		if !file.IsEncrypted() {
			continue
		}

		file.SetPassword(password)
		r, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(io.Discard, r)
		r.Close()

		if err != nil {
			return fmt.Errorf("password is incorrect: %w", err)
		}
	}

	return nil
}

func (c Cracker) CheckWordlist(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", nil
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		password := scanner.Text()
		passErr := c.CheckPassword(password)
		if passErr == nil {
			return password, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("no password from the list is correct")
}
