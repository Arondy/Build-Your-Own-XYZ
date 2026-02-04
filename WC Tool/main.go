package main

import (
	"fmt"
	"wc_tool/args"
)

func main() {
	parser := args.Parser{}
	err := parser.Parse()
	if err != nil {
		fmt.Println(err)
	}
}
