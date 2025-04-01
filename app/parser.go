package main

import (
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) parse() (Expr, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.isAtEnd() {
		return nil, fmt.Errorf("unexpected token after expression")
	}

	return expr, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.term()
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(PLUS, MINUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(STAR, SLASH) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Unary{
			Operator: operator,
			Right:    right,
		}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(TRUE) {
		return &Literal{Value: true}, nil
	}
	if p.match(FALSE) {
		return &Literal{Value: false}, nil
	}
	if p.match(NIL) {
		return &Literal{Value: nil}, nil
	}
	if p.match(NUMBER) {
		return &Literal{Value: p.previous().Literal}, nil
	}
	if p.match(STRING) {
		return &Literal{Value: p.previous().Literal}, nil
	}
	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if p.match(RIGHT_PAREN) {
			return &Grouping{Expression: expr}, nil
		}
		return &Grouping{Expression: expr}, nil
	}

	return nil, fmt.Errorf("Expect expression.")
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}
