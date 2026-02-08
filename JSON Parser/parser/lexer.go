package parser

import "fmt"

type TokenType int

const (
	LEFT_BRACE TokenType = iota
	RIGHT_BRACE
	WHITESPACE
	STRING
	COLON
	COMMA
)

var charToType = map[byte]TokenType{
	'{': LEFT_BRACE,
	'}': RIGHT_BRACE,
	':': COLON,
	',': COMMA,
}

type Token struct {
	Type  TokenType
	Value string
}

func Lex(contents []byte) (tokens []Token, err error) {
	jsonStr := string(contents)

	for i := 0; i < len(jsonStr); i++ {
		switch jsonStr[i] {
		case '{', '}', ':', ',':
			tokens = append(tokens, Token{
				Type: charToType[jsonStr[i]],
			})
		case '"':
			token := Token{Type: STRING}
			i++

			for i < len(jsonStr) && jsonStr[i] != '"' {
				token.Value += string(jsonStr[i])
				i++
			}

			if i == len(jsonStr) {
				return nil, fmt.Errorf("no closing \" in the end: '%c' instead", jsonStr[i-1])
			}

			tokens = append(tokens, token)
		case ' ', '\t', '\n', '\r':

		default:
			return nil, fmt.Errorf("wrong token in position %d: '%c'", i, jsonStr[i])
		}
	}

	return tokens, nil
}
