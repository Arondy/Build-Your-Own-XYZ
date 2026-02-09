package parser

import (
	"fmt"
	"slices"
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
		tokensStr := ""

		for _, token := range p.Tokens[p.position:len(p.Tokens)] {
			tokensStr += token.Type.String()
		}

		return fmt.Errorf("unexpected tokens after ending '}': '%s'", tokensStr)
	}

	return nil
}

func (p *Parser) safeIncrementPosition() error {
	p.position++

	if p.position >= len(p.Tokens) {
		return fmt.Errorf("reached unexpected tokens end (%d)", p.position)
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
		return fmt.Errorf("unexpected token type: '%s'", p.Tokens[p.position].Type)
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
		return fmt.Errorf("unexpected token: '%s' can't be value", currentToken.Type)
	}

	if currentToken.Type == LEFT_BRACE || currentToken.Type == LEFT_BRACKET {
		err := p.parseToken()
		if err != nil {
			return err
		}

		currentToken = p.Tokens[p.position]
		if !slices.Contains(canFollowValue, currentToken.Type) {
			return fmt.Errorf("unexpected token after value: '%s'", currentToken.Type)
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

	if p.Tokens[p.position-1].Type == COMMA {
		return fmt.Errorf("trailing ',' in array")
	}

	p.position++
	return nil
}

func (p *Parser) parsePair() error {
	currentToken := p.Tokens[p.position]

	if currentToken.Type != STRING {
		return fmt.Errorf("no key found, '%v' instead", currentToken.Type)
	}

	if err := p.safeIncrementPosition(); err != nil {
		return err
	}
	currentToken = p.Tokens[p.position]

	if currentToken.Type != COLON {
		return fmt.Errorf("no ':' after key, found '%v' instead", currentToken.Type)
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
