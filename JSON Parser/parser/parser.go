package parser

import (
	"fmt"
	"slices"
	"strings"
)

type Parser struct {
	Tokens   []Token
	position int
}

var canBeValue = []TokenType{STRING, NUMBER, TRUE, FALSE, NULL, LEFT_BRACE, LEFT_BRACKET}
var canFollowValue = []TokenType{COMMA, RIGHT_BRACE, RIGHT_BRACKET}

func (p *Parser) Parse() error {
	if len(p.Tokens) == 0 {
		return fmt.Errorf("empty sequence provided!")
	}
	if err := p.parseToken(); err != nil {
		return err
	}
	if p.position != len(p.Tokens) {
		tokensStr := strings.Builder{}

		for _, token := range p.Tokens[p.position:len(p.Tokens)] {
			tokensStr.WriteString(token.Type.String())
		}

		token := p.Tokens[p.position]
		return fmt.Errorf("unexpected tokens after ending '}' starting from line %d, position %d: '%s'", token.Line, token.RelativePosition, tokensStr.String())
	}

	return nil
}

func (p *Parser) safeIncrementPosition() error {
	p.position++

	if p.position >= len(p.Tokens) {
		token := p.Tokens[len(p.Tokens)-1]
		return fmt.Errorf("reached unexpected tokens end, last token: %s", token)
	}

	return nil
}

func (p *Parser) parseToken() error {
	currentToken := p.Tokens[p.position]

	switch currentToken.Type {
	case LEFT_BRACE:
		return p.parseObject()
	case LEFT_BRACKET:
		return p.parseArray()
	case STRING:
		return p.parseMultiplePairs()
	default:
		token := p.Tokens[p.position]
		return fmt.Errorf("unexpected token %s", token)
	}
}

func (p *Parser) parseObject() error {
	if err := p.safeIncrementPosition(); err != nil {
		return err
	}

	for p.Tokens[p.position].Type != RIGHT_BRACE {
		if err := p.parseToken(); err != nil {
			return err
		}
	}

	p.position++
	return nil
}

func (p *Parser) parseValue() error {
	currentToken := p.Tokens[p.position]

	if !slices.Contains(canBeValue, currentToken.Type) {
		token := p.Tokens[p.position]
		return fmt.Errorf("unexpected token can't be value: %s", token)
	}

	if currentToken.Type == LEFT_BRACE || currentToken.Type == LEFT_BRACKET {
		err := p.parseToken()
		if err != nil {
			return err
		}

		currentToken = p.Tokens[p.position]
		if !slices.Contains(canFollowValue, currentToken.Type) {
			token := p.Tokens[p.position]
			return fmt.Errorf("unexpected token after value: %s", token)
		}
	} else {
		if err := p.safeIncrementPosition(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) parseArray() error {
	if err := p.safeIncrementPosition(); err != nil {
		return err
	}

	for p.Tokens[p.position].Type != RIGHT_BRACKET {
		if err := p.parseValue(); err != nil {
			return err
		}

		if p.Tokens[p.position].Type == COMMA {
			if err := p.safeIncrementPosition(); err != nil {
				return err
			}
		} else {
			break
		}
	}

	token := p.Tokens[p.position]

	if token.Type != RIGHT_BRACKET {
		return fmt.Errorf("no closing ']': found %s instead", token)
	}

	token = p.Tokens[p.position-1]

	if token.Type == COMMA {
		return fmt.Errorf("trailing ',' in array: %s", token)
	}

	p.position++
	return nil
}

func (p *Parser) parsePair() error {
	token := p.Tokens[p.position]

	if token.Type != STRING {
		return fmt.Errorf("no key: found %s instead", token)
	}

	if err := p.safeIncrementPosition(); err != nil {
		return err
	}
	token = p.Tokens[p.position]

	if token.Type != COLON {
		return fmt.Errorf("no ':' after key: found %s instead", token)
	}

	if err := p.safeIncrementPosition(); err != nil {
		return err
	}
	if err := p.parseValue(); err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseMultiplePairs() error {
	for {
		if err := p.parsePair(); err != nil {
			return err
		}

		if p.Tokens[p.position].Type == COMMA {
			if err := p.safeIncrementPosition(); err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}
