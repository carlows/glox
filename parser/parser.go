package parser

import (
	"glox/expr"
	"glox/scanner"
	"slices"
)

type Parser struct {
	tokens  []scanner.Token
	current int
	errFn   func(token scanner.Token, message string)
}

func NewParser(tokens []scanner.Token, errFn func(token scanner.Token, message string)) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
		errFn:   errFn,
	}
}

type parseError struct{}

func (p *Parser) error(token scanner.Token, message string) {
	p.errFn(token, message)
	panic(parseError{})
}

func (p *Parser) Parse() expr.Expr {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(parseError); !ok {
				panic(r) // re-panic if it's not a parse error
			}
			// parse error already reported via callback, just return
		}
	}()

	return p.expression()
}

func (p *Parser) expression() expr.Expr {
	return p.equality()
}

// 4 == 4 != 5
// Conects a syntax tree (right node connected to previous)
func (p *Parser) equality() expr.Expr {
	left := p.comparison()
	for p.match(scanner.EqualEqual, scanner.BangEqual) {
		op := p.previous()
		right := p.comparison()
		left = &expr.Binary{
			Left:  left,
			Op:    op,
			Right: right,
		}
	}
	return left
}

// 3 > 2 <= 3
func (p *Parser) comparison() expr.Expr {
	left := p.term()
	for p.match(scanner.Greater, scanner.GreaterEqual, scanner.Less, scanner.LessEqual) {
		op := p.previous()
		right := p.term()
		left = &expr.Binary{
			Left:  left,
			Op:    op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) term() expr.Expr {
	left := p.factor()
	for p.match(scanner.Slash, scanner.Star) {
		op := p.previous()
		right := p.factor()
		left = &expr.Binary{
			Left:  left,
			Op:    op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) factor() expr.Expr {
	left := p.unary()
	for p.match(scanner.Minus, scanner.Plus) {
		op := p.previous()
		right := p.unary()
		left = &expr.Binary{
			Left:  left,
			Op:    op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) unary() expr.Expr {
	if p.match(scanner.Bang, scanner.Minus) {
		op := p.previous()
		right := p.unary()
		return &expr.Unary{
			Op:   op,
			Expr: right,
		}
	}

	return p.primary()
}

func (p *Parser) primary() expr.Expr {
	if p.match(scanner.True) {
		return &expr.Literal{
			Value: true,
		}
	}
	if p.match(scanner.False) {
		return &expr.Literal{
			Value: false,
		}
	}
	if p.match(scanner.Nil) {
		return &expr.Literal{
			Value: nil,
		}
	}
	if p.match(scanner.Number, scanner.String) {
		return &expr.Literal{
			Value: p.previous().Literal,
		}
	}

	if p.match(scanner.LeftParen) {
		e := p.expression()
		p.consume(scanner.RightParen, "Expecting ')' after expression")
		return &expr.Grouping{
			Expr: e,
		}
	}

	p.error(p.peek(), "Expecting expression.")
	return nil
}

func (p *Parser) consume(t scanner.TokenType, message string) scanner.Token {
	if p.check(t) {
		return p.advance()
	}

	p.error(p.peek(), message)
	return scanner.Token{}
}

// Consumes a token if it matches one of the given types
func (p *Parser) match(types ...scanner.TokenType) bool {
	matched := slices.ContainsFunc(types, func(t scanner.TokenType) bool {
		return p.check(t)
	})
	if matched {
		p.advance()
		return true
	}
	return false
}

// Checks if a token matches the given type
func (p *Parser) check(t scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	peek := p.peek()
	return peek.Type == t
}

// Returns the current token without consuming it
func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

// Returns the previous token without consuming it
func (p *Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}

// Checks for end of file, this time checking tokens
// The scanner was concerned about characters
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.Eof
}

// Consumes the current token and returns it
func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}
