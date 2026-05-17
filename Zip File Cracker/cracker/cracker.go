package cracker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

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
			return err
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

func (c Cracker) generatePasswords(pos int, password []byte) string {
	if pos == len(password) {
		err := c.CheckPassword(string(password))
		if err == nil {
			return string(password)
		}
		return ""
	}

	// for i := 97; i <= 122; i++ {
	for i := 33; i <= 126; i++ {
		password[pos] = byte(i)
		pass := c.generatePasswords(pos+1, password)
		if pass != "" {
			return pass
		}
	}

	return ""
}

func (c Cracker) Bruteforce(min, max int) (string, error) {
	s := time.Now()
	for lettersNum := min; lettersNum <= max; lettersNum++ {
		buffer := make([]byte, lettersNum)

		password := c.generatePasswords(0, buffer)
		if password != "" {
			fmt.Printf("Found correct password in %v: %s", time.Since(s), password)
			return password, nil
		}
	}

	return "", fmt.Errorf("no password with length [%d, %d] is correct, searched for %v", min, max, time.Since(s))
}
