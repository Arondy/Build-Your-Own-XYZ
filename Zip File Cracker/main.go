package main

import (
	"zip_file_cracker/cracker"
)

func main() {
	c, err := cracker.NewCracker("test.zip")
	if err != nil {
		panic(err)
	}
	_, err = c.Bruteforce(1, 3)
	if err != nil {
		panic(err)
	}
}
