package main

import (
	"compression_tool/args"
	"fmt"
)

func main() {
	err := args.Parse()
	if err != nil {
		fmt.Print(err)
	}
}
