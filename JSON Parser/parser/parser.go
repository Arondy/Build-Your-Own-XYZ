package parser

import "fmt"

type Parser struct {
	Tokens   []Token
	position int
}

func (p *Parser) Parse() error {
	if len(p.Tokens) == 0 {
		return fmt.Errorf("empty sequence provided!")
	}
	if err := p.parseToken(); err != nil {
		return err
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

	if err := p.safeIncrementPosition(); err != nil {
		return err
	}

	switch currentToken.Type {
	case LEFT_BRACE:
		return p.parseMap()
	case STRING:
		return p.parseObject()
	case COMMA:
		if err := p.safeIncrementPosition(); err != nil {
			return err
		}

		return p.parseObject()
	default:
		return fmt.Errorf("unexpected token type: %v", p.Tokens[p.position-1].Type)
	}
}

func (p *Parser) parseMap() error {
	for p.Tokens[p.position].Type != RIGHT_BRACE {
		if err := p.parseToken(); err != nil {
			return err
		}
	}

	p.position++

	if p.position != len(p.Tokens) {
		return fmt.Errorf("unexpected tokens after ending '}': '%v'", p.Tokens[p.position:len(p.Tokens)])
	}
	return nil
}

func (p *Parser) parseObject() error {
	currentToken := p.Tokens[p.position]

	if currentToken.Type != COLON {
		return fmt.Errorf("no colon after key, found '%v' instead", currentToken.Type)
	}

	if err := p.safeIncrementPosition(); err != nil {
		return err
	}
	currentToken = p.Tokens[p.position]

	if currentToken.Type != STRING {
		return fmt.Errorf("no value after colon, found '%v' instead", currentToken.Type)
	}

	if err := p.safeIncrementPosition(); err != nil {
		return err
	}
	return nil
}
