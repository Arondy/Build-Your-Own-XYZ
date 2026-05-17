package main

import (
	"fmt"
	"zip_file_cracker/cracker"
)

func main() {
	c, err := cracker.NewCracker("test.zip")
	if err != nil {
		panic(err)
	}
	password, err := c.CheckWordlist("wordlist.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(password)
}
