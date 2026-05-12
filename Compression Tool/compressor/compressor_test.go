package compressor

import (
	"bufio"
	"io"
	"maps"
	"os"
	"testing"
)

const testFilename = "../test.txt"
const testHeaderFilename = "../test.header"
const testEncodedFilename = "../test.enc"
const testDecodedFilename = "../test.dec"

func TestFrequencies(t *testing.T) {
	file, err := os.Open(testFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	frequencies, err := getCharactersFrequency(reader)
	if err != nil {
		t.FailNow()
	}
	if !(frequencies[byte('t')] == 223000 && frequencies[byte('X')] == 333) {
		t.Fail()
	}
}

func TestSaveAndLoadFreqs(t *testing.T) {
	file, err := os.Open(testFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	freqs, err := getCharactersFrequency(reader)
	if err != nil {
		t.FailNow()
	}
	headerFile, err := os.OpenFile(testHeaderFilename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer os.Remove(testHeaderFilename)
	defer headerFile.Close()

	writer := bufio.NewWriter(headerFile)
	err = writeFreqs(freqs, writer)
	if err != nil {
		t.FailNow()
	}

	_, err = headerFile.Seek(0, io.SeekStart)
	if err != nil {
		t.FailNow()
	}

	reader2 := bufio.NewReader(headerFile)
	freqs2, _, err := loadFreqs(reader2)
	if err != nil {
		t.FailNow()
	}

	if !maps.Equal(freqs, freqs2) {
		t.Fail()
	}
}

func TestEncodeAndDecode(t *testing.T) {
	err := Encode(testFilename, testEncodedFilename)
	if err != nil {
		t.FailNow()
	}

	err = Decode(testEncodedFilename, testDecodedFilename)
	if err != nil {
		t.FailNow()
	}

	srcStats, err := os.Stat(testFilename)
	if err != nil {
		t.FailNow()
	}
	decodedStats, err := os.Stat(testDecodedFilename)
	if err != nil {
		t.FailNow()
	}
	if srcStats.Size() != decodedStats.Size() {
		t.Fail()
	}

	os.Remove(testEncodedFilename)
	os.Remove(testDecodedFilename)
}
