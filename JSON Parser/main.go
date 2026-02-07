package main

import (
	"fmt"
	"json_parser/parser"
	"os"
)

func main() {
	json := []byte("{}")

	tokens, err := parser.Lex(json)
	p := parser.Parser{Tokens: tokens}

	err = p.Parse()
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	fmt.Print("Successful JSON parsing!")
}
