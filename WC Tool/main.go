package main

import (
	"fmt"
	"wc_tool/args"
	"wc_tool/wc"
)

func main() {
	wordCount := wc.WordCount{}
	parser := args.Parser{
		WC: wordCount,
	}
	err := parser.Parse()
	if err != nil {
		fmt.Println(err)
	}
}
