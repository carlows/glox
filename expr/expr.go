package expr

import "glox/scanner"

type Expr interface {
	Accept(Visitor) any
}

type Binary struct {
	Left Expr
	Op   scanner.TokenType
	Right Expr
}

func (b *Binary) Accept(v Visitor) any {
	return v.VisitBinary(b)
}

type Grouping struct {
	Expr Expr
}

func (g *Grouping) Accept(v Visitor) any {
	return v.VisitGrouping(g)
}

type Literal struct {
	Value any
}

func (l *Literal) Accept(v Visitor) any {
	return v.VisitLiteral(l)
}

type Unary struct {
	Op scanner.TokenType
	Expr Expr
}

func (u *Unary) Accept(v Visitor) any {
	return v.VisitUnary(u)
}
