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
	err = c.CheckPassword("test")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Correct password found!")
	}
}
