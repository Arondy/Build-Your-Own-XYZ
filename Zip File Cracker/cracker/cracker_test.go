package cracker_test

import (
	"testing"
	"zip_file_cracker/cracker"

	"github.com/yeka/zip"
)

const testZipFilename = "../test.zip"
const testWordlist = "../wordlist.txt"
const testPassword = "test"

func TestCheckPassword(t *testing.T) {
	c, err := cracker.NewCracker(testZipFilename)
	if err != nil {
		t.Fatal(err)
	}

	reader, err := zip.OpenReader(c.Filename)
	if err != nil {
		t.Fatal(err)
	}

	err = c.CheckPassword(reader, testPassword)
	if err != nil {
		t.Fatalf("the correct password '%s' did not match: %s", testPassword, err)
	}
}

func TestWordlistAttack(t *testing.T) {
	c, err := cracker.NewCracker(testZipFilename)
	if err != nil {
		t.Fatal(err)
	}

	password, err := c.WordlistAttack(testWordlist)
	if err != nil {
		t.Fatal(err)
	}
	if password != testPassword {
		t.Fatalf("passwords don't match: expected %s, got %s", testPassword, password)
	}
}

func TestBruteforce(t *testing.T) {
	c, err := cracker.NewCracker(testZipFilename)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Bruteforce(1, 3)
	if err == nil {
		t.Fatal("password isn't in [1, 3] range")
	}

	password, err := c.Bruteforce(4, 4)
	if err != nil {
		t.Fatal(err)
	}

	if password != testPassword {
		t.Fatalf("passwords don't match: expected %s, got %s", testPassword, password)
	}
}
