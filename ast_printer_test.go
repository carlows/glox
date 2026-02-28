package main

import (
	"glox/expr"
	"glox/scanner"
	"testing"
)

func TestAstPrinter(t *testing.T) {
	expression := expr.Binary{
		Left: &expr.Unary{
			Op: scanner.Token{
				Type:    scanner.Minus,
				Lexeme:  "-",
				Literal: nil,
				Line:    0,
			},
			Expr: &expr.Literal{
				Value: 123,
			},
		},
		Op: scanner.Token{
			Type:    scanner.Star,
			Lexeme:  "*",
			Literal: nil,
			Line:    0,
		},
		Right: &expr.Grouping{
			Expr: &expr.Literal{
				Value: 45.67,
			},
		},
	}

	printer := AstPrinter{}
	result := printer.Print(&expression)

	expected := "(Star (Minus 123) (Grouping 45.67))"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
