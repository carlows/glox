package main

import (
	"glox/expr"
	"glox/scanner"
	"testing"
)

func TestAstPrinter(t *testing.T) {
	expression := expr.Binary{
		Left: &expr.Unary{
			Op: scanner.Minus,
			Expr: &expr.Literal{
				Value: 123,
			},
		},
		Op: scanner.Star,
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
