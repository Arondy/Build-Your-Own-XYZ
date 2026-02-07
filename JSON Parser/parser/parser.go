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

func (p *Parser) parseToken() error {
	if p.position >= len(p.Tokens) {
		return fmt.Errorf("reached tokens end (%d)", p.position)
	}

	currentToken := p.Tokens[p.position]
	p.position++

	if p.position >= len(p.Tokens) {
		return fmt.Errorf("reached tokens end (%d)", p.position)
	}

	switch currentToken.Type {
	case LEFT_BRACE:
		return p.parseMap()
	default:
		return fmt.Errorf("unexpected token type: %v", p.Tokens[p.position-1].Type)
	}
}

func (p *Parser) parseMap() error {
	currentToken := p.Tokens[p.position]

	if currentToken.Type == RIGHT_BRACE {
		p.position++
		return nil
	} else {
		return p.parseToken()
	}
}
