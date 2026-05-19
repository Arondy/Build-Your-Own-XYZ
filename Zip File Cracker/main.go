package main

import (
	"fmt"
	"zip_file_cracker/args"
)

func main() {
	err := args.Parse()
	if err != nil {
		fmt.Println(err)
	}
}
