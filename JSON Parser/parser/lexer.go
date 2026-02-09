package parser

import (
	"fmt"
	"slices"
)

type TokenType int

const (
	LEFT_BRACE TokenType = iota
	RIGHT_BRACE
	LEFT_BRACKET
	RIGHT_BRACKET
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
	LEFT_BRACE:    "{",
	RIGHT_BRACE:   "}",
	LEFT_BRACKET:  "[",
	RIGHT_BRACKET: "]",
	WHITESPACE:    " ",
	STRING:        "string",
	COLON:         ":",
	COMMA:         ",",
	TRUE:          "true",
	FALSE:         "false",
	NULL:          "null",
	NUMBER:        "number",
}

func (t TokenType) String() string {
	return tokenTypeNames[t]
}

var canFollowValueChar = []byte{',', '}', ']'}
var canFollowBackslash = []byte{'"', '\\', '/', 'b', 'f', 'n', 'r', 't', 'u'}

var charToType = map[byte]TokenType{
	'{': LEFT_BRACE,
	'}': RIGHT_BRACE,
	'[': LEFT_BRACKET,
	']': RIGHT_BRACKET,
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

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isHex(char byte) bool {
	return (char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')
}

type Lexer struct {
	jsonStr           string
	position          int
	line              int
	lineStartPosition int
	tokens            []Token
}

func NewLexer(json []byte) Lexer {
	return Lexer{jsonStr: string(json), line: 1}
}

func (l *Lexer) getLineRelativePosition() int {
	return l.position - l.lineStartPosition + 1
}

func (l *Lexer) readString() error {
	token := Token{Type: STRING}
	l.position++

	for l.position < len(l.jsonStr) && l.jsonStr[l.position] != '"' {
		char := l.jsonStr[l.position]

		if char == '\t' || char == '\n' {
			return fmt.Errorf("'%c' isn't allowed in string in line %d, position %d", char, l.line, l.getLineRelativePosition())
		}

		token.Value += string(char)
		l.position++

		if char == '\\' {
			char = l.jsonStr[l.position]

			if !slices.Contains(canFollowBackslash, char) {
				return fmt.Errorf("wrong token '%c' after '\\' in line %d, position %d", char, l.line, l.getLineRelativePosition())
			}

			token.Value += string(char)
			l.position++

			if char == 'u' {
				for range 4 {
					char = l.jsonStr[l.position]

					if !isHex(char) {
						return fmt.Errorf("non-hex token '%c' after '\\u' in line %d, position %d", char, l.line, l.getLineRelativePosition())
					}

					l.position++
				}
			}
		}
	}

	if l.position == len(l.jsonStr) {
		return fmt.Errorf("no closing ' \" ' at the end of string in line %d, position %d: '%c' instead", l.line, l.getLineRelativePosition(), l.jsonStr[l.position-1])
	}

	l.tokens = append(l.tokens, token)
	return nil
}

func (l *Lexer) readWord() error {
	char := l.jsonStr[l.position]
	expectedWord := charToString[char]

	if len(l.jsonStr) < l.position+len(expectedWord) {
		return fmt.Errorf("wrong token '%c' in line %d, position %d", char, l.line, l.getLineRelativePosition())
	}

	actualWord := l.jsonStr[l.position : l.position+len(expectedWord)]

	if actualWord != expectedWord {
		return fmt.Errorf("wrong token sequence '%v' starting from line %d, position %d: '%s' expected", actualWord, l.line, l.getLineRelativePosition(), expectedWord)
	}

	l.tokens = append(l.tokens, Token{Type: stringToType[expectedWord]})
	l.position += len(expectedWord) - 1
	return nil
}

func (l *Lexer) appendCharThenIncrementPosition(value *[]byte) (byte, error) {
	*value = append(*value, l.jsonStr[l.position])

	l.position++
	if l.position >= len(l.jsonStr) {
		return 0, fmt.Errorf("unexpected end of tokens in line %d, position %d", l.line, l.getLineRelativePosition())
	}

	return l.jsonStr[l.position], nil
}

func (l *Lexer) readNumber() error {
	char := l.jsonStr[l.position]
	value := make([]byte, 0)
	startsWithZero := char == '0'
	isNegative := char == '-'
	char, err := l.appendCharThenIncrementPosition(&value)
	if err != nil {
		return err
	}

	// -?x
	if isNegative {
		startsWithZero = char == '0'

		if !isDigit(char) {
			return fmt.Errorf("no digit after '-' in the number in line %d, position %d: found '%v' instead", l.line, l.getLineRelativePosition(), char)
		}
		if char, err = l.appendCharThenIncrementPosition(&value); err != nil {
			return err
		}
	}

	if !startsWithZero {
		for isDigit(char) {
			if char, err = l.appendCharThenIncrementPosition(&value); err != nil {
				return err
			}
		}
	}

	if startsWithZero && slices.Contains(canFollowValueChar, char) {
		l.tokens = append(l.tokens, Token{Type: NUMBER, Value: string(value)})
		return nil
	}

	// 0.x
	if startsWithZero && !slices.Contains([]byte{'e', 'E', '.'}, char) {
		return fmt.Errorf("no 'e', 'E' or '.' after starting '0' in the number in line %d, position %d: found '%v' instead", l.line, l.getLineRelativePosition(), char)
	}

	if char == '.' {
		if char, err = l.appendCharThenIncrementPosition(&value); err != nil {
			return err
		}
		if !isDigit(char) {
			return fmt.Errorf("no digit after '.' in the number in line %d, position %d: found '%v' instead", l.line, l.getLineRelativePosition(), char)

		}
		for isDigit(char) {
			if char, err = l.appendCharThenIncrementPosition(&value); err != nil {
				return err
			}
		}
	}

	if value[len(value)-1] == '.' {
		return fmt.Errorf("no digit after '.' in the number in line %d, position %d: found '%v' instead", l.line, l.getLineRelativePosition(), char)
	}

	// e(+-)?x
	if char == 'e' || char == 'E' {
		if char, err = l.appendCharThenIncrementPosition(&value); err != nil {
			return err
		}

		if char == '-' || char == '+' {
			if char, err = l.appendCharThenIncrementPosition(&value); err != nil {
				return err
			}
		}

		for isDigit(char) {
			if char, err = l.appendCharThenIncrementPosition(&value); err != nil {
				return err
			}
		}
	}
	if !isDigit(value[len(value)-1]) {
		return fmt.Errorf("no digit after 'e' in the number in line %d, position %d: found '%v' instead", l.line, l.getLineRelativePosition(), char)
	}
	l.position--

	l.tokens = append(l.tokens, Token{Type: NUMBER, Value: string(value)})
	return nil
}

func (l *Lexer) Lex() (tokens []Token, err error) {
	for ; l.position < len(l.jsonStr); l.position++ {
		char := l.jsonStr[l.position]

		if isDigit(char) || char == '-' {
			if err := l.readNumber(); err != nil {
				return nil, err
			}
			continue
		}

		switch char {
		case '{', '}', '[', ']', ':', ',':
			l.tokens = append(l.tokens, Token{
				Type: charToType[char],
			})
		case '"':
			if err := l.readString(); err != nil {
				return nil, err
			}
		case 't', 'f', 'n':
			if err := l.readWord(); err != nil {
				return nil, err
			}
		case ' ', '\t', '\n', '\r':
			if char == '\n' {
				l.line++
				l.lineStartPosition = l.position
			}

		default:
			return nil, fmt.Errorf("wrong token '%c' in line %d, position %d", l.line, l.position-l.lineStartPosition, char)
		}
	}

	return l.tokens, nil
}
