package parser

import "fmt"

type TokenType int

const (
	LEFT_BRACE TokenType = iota
	RIGHT_BRACE
	WHITESPACE
)

type Token struct {
	Type TokenType
}

func Lex(contents []byte) (tokens []Token, err error) {
	jsonStr := string(contents)

	for i := 0; i < len(jsonStr); i++ {
		switch jsonStr[i] {
		case '{':
			tokens = append(tokens, Token{
				Type: LEFT_BRACE,
			})
		case '}':
			tokens = append(tokens, Token{
				Type: RIGHT_BRACE,
			})
		case ' ', '\t', '\n', '\r':

		default:
			return nil, fmt.Errorf("wrong token in position %d: '%c'", i, jsonStr[i])
		}
	}

	return tokens, nil
}
