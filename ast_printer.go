package main

import (
	"fmt"
	"glox/expr"
	"strings"
)

type AstPrinter struct{}

func (p *AstPrinter) Print(expr expr.Expr) string {
	return expr.Accept(p).(string)
}

func (p *AstPrinter) VisitBinary(expr *expr.Binary) any {
	return p.parenthesize(expr.Op.String(), expr.Left, expr.Right)
}

func (p *AstPrinter) VisitGrouping(expr *expr.Grouping) any {
	return p.parenthesize("Grouping", expr.Expr)
}

func (p *AstPrinter) VisitLiteral(expr *expr.Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprint(expr.Value)
}

func (p *AstPrinter) VisitUnary(expr *expr.Unary) any {
	return p.parenthesize(expr.Op.String(), expr.Expr)
}

func (p *AstPrinter) parenthesize(name string, exprs ...expr.Expr) string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(name)

	for _, expr := range exprs {
		sb.WriteString(" ")
		sb.WriteString(expr.Accept(p).(string))
	}

	sb.WriteString(")")

	return sb.String()
}
