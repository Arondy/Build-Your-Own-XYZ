package cracker

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/yeka/zip"
)

var zipHeader = [4]byte{80, 75, 03, 04}

// lowercase english characters
const asciiStart = 97
const asciiEnd = 123

// const asciiStart = 33
// const asciiEnd = 127

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
	Filename string
}

func NewCracker(filename string) (Cracker, error) {
	if !isZipFile(filename) {
		return Cracker{}, fmt.Errorf("'%s' isn't a valid ZIP file!", filename)
	}
	return Cracker{filename}, nil
}

func (c Cracker) createReaders(num int) ([]*zip.ReadCloser, error) {
	readers := make([]*zip.ReadCloser, 0, num)

	for range num {
		reader, err := zip.OpenReader(c.Filename)
		if err != nil {
			return nil, err
		}
		readers = append(readers, reader)
	}

	return readers, nil
}

func (c Cracker) CheckPassword(reader *zip.ReadCloser, password string) error {
	for _, file := range reader.File {
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

func (c Cracker) wordlistWorker(ctx context.Context, reader *zip.ReadCloser, passwords <-chan string, results chan<- string) {
	for {
		select {
		case <-ctx.Done():
			return
		case password, ok := <-passwords:
			if !ok {
				return
			}

			err := c.CheckPassword(reader, password)
			if err == nil {
				results <- password
			}
		}
	}
}

func (c Cracker) wordlistProducer(ctx context.Context, scanner *bufio.Scanner, passwords chan<- string) {
	defer close(passwords)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		case passwords <- scanner.Text():
		}
	}
}

func (c Cracker) WordlistAttack(wordlist string) (string, error) {
	s := time.Now()

	file, err := os.Open(wordlist)
	if err != nil {
		return "", err
	}
	defer file.Close()

	passwords := make(chan string, 8)
	results := make(chan string)
	scanner := bufio.NewScanner(file)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cores := runtime.NumCPU()
	wg := sync.WaitGroup{}

	wg.Go(func() {
		c.wordlistProducer(ctx, scanner, passwords)
	})

	for range cores - 1 {
		reader, err := zip.OpenReader(c.Filename)
		if err != nil {
			return "", err
		}
		defer reader.Close()

		wg.Go(func() {
			c.wordlistWorker(ctx, reader, passwords, results)
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	password, ok := <-results
	cancel() // for no race in scanner.Err() check

	if ok {
		fmt.Printf("Found correct password in %v: %s", time.Since(s), password)
		return password, nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("no password from the list '%s' is correct, searched for %s", wordlist, time.Since(s))
}

func (c Cracker) generatePasswords(ctx context.Context, results chan<- string, reader *zip.ReadCloser, pos, start, end int, password []byte) {
	if pos == len(password) {
		err := c.CheckPassword(reader, string(password))
		if err == nil {
			results <- string(password)
		}
		return
	}

	for i := start; i < end; i++ {
		if ctx.Err() != nil {
			return
		}

		password[pos] = byte(i)
		c.generatePasswords(ctx, results, reader, pos+1, asciiStart, asciiEnd, password)
	}
}

func (c Cracker) Bruteforce(min, max int) (string, error) {
	s := time.Now()

	cores := runtime.NumCPU()
	wg := sync.WaitGroup{}
	readers, err := c.createReaders(cores)
	if err != nil {
		return "", err
	}

	for lettersNum := min; lettersNum <= max; lettersNum++ {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		results := make(chan string)

		for i := range cores {
			start := asciiStart + (asciiEnd-asciiStart)*i/cores
			end := asciiStart + (asciiEnd-asciiStart)*(i+1)/cores
			buffer := make([]byte, lettersNum)

			wg.Go(func() {
				c.generatePasswords(ctx, results, readers[i], 0, start, end, buffer)
			})
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		password, ok := <-results

		if ok {
			fmt.Printf("Found correct password in %v: %s", time.Since(s), password)
			return password, nil
		}
	}

	return "", fmt.Errorf("no password with length [%d, %d] is correct, searched for %s", min, max, time.Since(s))
}
