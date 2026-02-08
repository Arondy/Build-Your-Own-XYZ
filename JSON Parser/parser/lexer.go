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
	TRUE
	FALSE
	NULL
	NUMBER
)

var tokenTypeNames = map[TokenType]string{
	LEFT_BRACE:  "{",
	RIGHT_BRACE: "}",
	WHITESPACE:  " ",
	STRING:      "string",
	COLON:       ":",
	COMMA:       ",",
	TRUE:        "true",
	FALSE:       "false",
	NULL:        "null",
	NUMBER:      "number",
}

func (t TokenType) String() string {
	return tokenTypeNames[t]
}

var charToType = map[byte]TokenType{
	'{': LEFT_BRACE,
	'}': RIGHT_BRACE,
	':': COLON,
	',': COMMA,
}

var charToString = map[byte]string{
	't': "true",
	'f': "false",
	'n': "null",
}

var stringToType = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
	"null":  NULL,
}

type Token struct {
	Type  TokenType
	Value string
}

func readWord(jsonStr string, position *int, tokens *[]Token) error {
	char := jsonStr[*position]
	expectedWord := charToString[char]

	if len(jsonStr) < *position+len(expectedWord) {
		return fmt.Errorf("wrong token in position %d: '%c'", position, jsonStr[*position])
	}

	actualWord := jsonStr[*position : *position+len(expectedWord)]

	if actualWord != expectedWord {
		return fmt.Errorf("wrong token sequence starting from position %d: '%v' ('%s' expected)", position, actualWord, expectedWord)
	}

	*tokens = append(*tokens, Token{Type: stringToType[expectedWord]})
	*position += len(expectedWord) - 1
	return nil
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isNumberPart(char byte) bool {
	return isDigit(char) || char == '-' || char == '+' || char == 'e' || char == 'E'
}

func incrementPosition(jsonStr string, position *int, value *[]byte) (byte, error) {
	*value = append(*value, jsonStr[*position])

	*position++
	if *position >= len(jsonStr) {
		return 0, fmt.Errorf("unexpected end of tokens")
	}

	return jsonStr[*position], nil
}

func readNumber(jsonStr string, position *int, tokens *[]Token) error {
	char := jsonStr[*position]
	value := make([]byte, 0)
	startsWithZero := char == '0'
	isNegative := char == '-'
	char, err := incrementPosition(jsonStr, position, &value)
	if err != nil {
		return err
	}

	// -?x
	if isNegative {
		startsWithZero = char == '0'

		if !isDigit(char) {
			return fmt.Errorf("no digit after '-' in the number in position %d, found '%v' instead", *position, char)
		}
		if char, err = incrementPosition(jsonStr, position, &value); err != nil {
			return err
		}
	}

	if !startsWithZero {
		for isDigit(char) {
			if char, err = incrementPosition(jsonStr, position, &value); err != nil {
				return err
			}
		}
	}

	// 0.x
	if startsWithZero && char != '.' {
		return fmt.Errorf("no '.' after starting '0' in the number in position %d, found '%v' instead", *position, char)
	}

	if char == '.' {
		if char, err = incrementPosition(jsonStr, position, &value); err != nil {
			return err
		}
		if !isDigit(char) {
			return fmt.Errorf("no digit after '.' in the number in position %d, found '%v' instead", *position, char)
		}
		for isDigit(char) {
			if char, err = incrementPosition(jsonStr, position, &value); err != nil {
				return err
			}
		}
	}

	if value[len(value)-1] == '.' {
		return fmt.Errorf("no digit after '.'' in position %d, found '%v' instead", *position, char)
	}

	// e(+-)?x
	if char == 'e' || char == 'E' {
		if char, err = incrementPosition(jsonStr, position, &value); err != nil {
			return err
		}

		if char == '-' || char == '+' {
			if char, err = incrementPosition(jsonStr, position, &value); err != nil {
				return err
			}
		}

		for isDigit(char) {
			if char, err = incrementPosition(jsonStr, position, &value); err != nil {
				return err
			}
		}
	}
	if !isDigit(value[len(value)-1]) {
		return fmt.Errorf("no digit after 'e' in position %d, found '%v' instead", *position, char)
	}
	*position--

	*tokens = append(*tokens, Token{Type: NUMBER, Value: string(value)})
	return nil
}

func Lex(contents []byte) (tokens []Token, err error) {
	jsonStr := string(contents)

	for i := 0; i < len(jsonStr); i++ {
		char := jsonStr[i]

		if isDigit(char) || char == '-' {
			if err := readNumber(jsonStr, &i, &tokens); err != nil {
				return nil, err
			}
			continue
		}

		switch char {
		case '{', '}', ':', ',':
			tokens = append(tokens, Token{
				Type: charToType[char],
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
		case 't', 'f', 'n':
			if err := readWord(jsonStr, &i, &tokens); err != nil {
				return nil, err
			}
		case ' ', '\t', '\n', '\r':

		default:
			return nil, fmt.Errorf("wrong token in position %d: '%c'", i, char)
		}
	}

	return tokens, nil
}
