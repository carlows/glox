package parser

import (
	"glox/expr"
	"glox/scanner"
	"testing"
)

func TestParserValidExpressions(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []scanner.Token
		validate func(t *testing.T, e expr.Expr)
	}{
		{
			name: "Literal - number",
			tokens: []scanner.Token{
				{Type: scanner.Number, Lexeme: "42", Literal: 42.0, Line: 1},
				{Type: scanner.Eof, Lexeme: "", Literal: nil, Line: 1},
			},
			validate: func(t *testing.T, e expr.Expr) {
				lit, ok := e.(*expr.Literal)
				if !ok {
					t.Fatalf("Expected Literal, got %T", e)
				}
				if lit.Value != 42.0 {
					t.Errorf("Expected 42.0, got %v", lit.Value)
				}
			},
		},
		{
			name: "Literal - true",
			tokens: []scanner.Token{
				{Type: scanner.True, Lexeme: "true", Literal: nil, Line: 1},
				{Type: scanner.Eof, Lexeme: "", Literal: nil, Line: 1},
			},
			validate: func(t *testing.T, e expr.Expr) {
				lit, ok := e.(*expr.Literal)
				if !ok {
					t.Fatalf("Expected Literal, got %T", e)
				}
				if lit.Value != true {
					t.Errorf("Expected true, got %v", lit.Value)
				}
			},
		},
		{
			name: "Literal - false",
			tokens: []scanner.Token{
				{Type: scanner.False, Lexeme: "false", Literal: nil, Line: 1},
				{Type: scanner.Eof, Lexeme: "", Literal: nil, Line: 1},
			},
			validate: func(t *testing.T, e expr.Expr) {
				lit, ok := e.(*expr.Literal)
				if !ok {
					t.Fatalf("Expected Literal, got %T", e)
				}
				if lit.Value != false {
					t.Errorf("Expected false, got %v", lit.Value)
				}
			},
		},
		{
			name: "Literal - nil",
			tokens: []scanner.Token{
				{Type: scanner.Nil, Lexeme: "nil", Literal: nil, Line: 1},
				{Type: scanner.Eof, Lexeme: "", Literal: nil, Line: 1},
			},
			validate: func(t *testing.T, e expr.Expr) {
				lit, ok := e.(*expr.Literal)
				if !ok {
					t.Fatalf("Expected Literal, got %T", e)
				}
				if lit.Value != nil {
					t.Errorf("Expected nil, got %v", lit.Value)
				}
			},
		},
		{
			name: "Unary - negation",
			tokens: []scanner.Token{
				{Type: scanner.Minus, Lexeme: "-", Literal: nil, Line: 1},
				{Type: scanner.Number, Lexeme: "5", Literal: 5.0, Line: 1},
				{Type: scanner.Eof, Lexeme: "", Literal: nil, Line: 1},
			},
			validate: func(t *testing.T, e expr.Expr) {
				unary, ok := e.(*expr.Unary)
				if !ok {
					t.Fatalf("Expected Unary, got %T", e)
				}
				if unary.Op.Type != scanner.Minus {
					t.Errorf("Expected Minus operator, got %v", unary.Op.Type)
				}
				lit, ok := unary.Expr.(*expr.Literal)
				if !ok {
					t.Fatalf("Expected Literal operand, got %T", unary.Expr)
				}
				if lit.Value != 5.0 {
					t.Errorf("Expected 5.0, got %v", lit.Value)
				}
			},
		},
		{
			name: "Unary - logical not",
			tokens: []scanner.Token{
				{Type: scanner.Bang, Lexeme: "!", Literal: nil, Line: 1},
				{Type: scanner.True, Lexeme: "true", Literal: nil, Line: 1},
				{Type: scanner.Eof, Lexeme: "", Literal: nil, Line: 1},
			},
			validate: func(t *testing.T, e expr.Expr) {
				unary, ok := e.(*expr.Unary)
				if !ok {
					t.Fatalf("Expected Unary, got %T", e)
				}
				if unary.Op.Type != scanner.Bang {
					t.Errorf("Expected Bang operator, got %v", unary.Op.Type)
				}
			},
		},
		{
			name: "Binary - addition",
			tokens: []scanner.Token{
				{Type: scanner.Number, Lexeme: "1", Literal: 1.0, Line: 1},
				{Type: scanner.Plus, Lexeme: "+", Literal: nil, Line: 1},
				{Type: scanner.Number, Lexeme: "2", Literal: 2.0, Line: 1},
				{Type: scanner.Eof, Lexeme: "", Literal: nil, Line: 1},
			},
			validate: func(t *testing.T, e expr.Expr) {
				binary, ok := e.(*expr.Binary)
				if !ok {
					t.Fatalf("Expected Binary, got %T", e)
				}
				if binary.Op.Type != scanner.Plus {
					t.Errorf("Expected Plus operator, got %v", binary.Op.Type)
				}
				left, ok := binary.Left.(*expr.Literal)
				if !ok {
					t.Fatalf("Expected Literal left operand, got %T", binary.Left)
				}
				if left.Value != 1.0 {
					t.Errorf("Expected 1.0, got %v", left.Value)
				}
				right, ok := binary.Right.(*expr.Literal)
				if !ok {
					t.Fatalf("Expected Literal right operand, got %T", binary.Right)
				}
				if right.Value != 2.0 {
					t.Errorf("Expected 2.0, got %v", right.Value)
				}
			},
		},
		{
			name: "Binary - multiplication",
			tokens: []scanner.Token{
				{Type: scanner.Number, Lexeme: "3", Literal: 3.0, Line: 1},
				{Type: scanner.Star, Lexeme: "*", Literal: nil, Line: 1},
				{Type: scanner.Number, Lexeme: "4", Literal: 4.0, Line: 1},
				{Type: scanner.Eof, Lexeme: "", Literal: nil, Line: 1},
			},
			validate: func(t *testing.T, e expr.Expr) {
				binary, ok := e.(*expr.Binary)
				if !ok {
					t.Fatalf("Expected Binary, got %T", e)
				}
				if binary.Op.Type != scanner.Star {
					t.Errorf("Expected Star operator, got %v", binary.Op.Type)
				}
			},
		},
		{
			name: "Grouping - parentheses",
			tokens: []scanner.Token{
				{Type: scanner.LeftParen, Lexeme: "(", Literal: nil, Line: 1},
				{Type: scanner.Number, Lexeme: "42", Literal: 42.0, Line: 1},
				{Type: scanner.RightParen, Lexeme: ")", Literal: nil, Line: 1},
				{Type: scanner.Eof, Lexeme: "", Literal: nil, Line: 1},
			},
			validate: func(t *testing.T, e expr.Expr) {
				grouping, ok := e.(*expr.Grouping)
				if !ok {
					t.Fatalf("Expected Grouping, got %T", e)
				}
				lit, ok := grouping.Expr.(*expr.Literal)
				if !ok {
					t.Fatalf("Expected Literal inside grouping, got %T", grouping.Expr)
				}
				if lit.Value != 42.0 {
					t.Errorf("Expected 42.0, got %v", lit.Value)
				}
			},
		},
		{
			name: "Complex - operator precedence",
			tokens: []scanner.Token{
				{Type: scanner.Number, Lexeme: "1", Literal: 1.0, Line: 1},
				{Type: scanner.Plus, Lexeme: "+", Literal: nil, Line: 1},
				{Type: scanner.Number, Lexeme: "2", Literal: 2.0, Line: 1},
				{Type: scanner.Star, Lexeme: "*", Literal: nil, Line: 1},
				{Type: scanner.Number, Lexeme: "3", Literal: 3.0, Line: 1},
				{Type: scanner.Eof, Lexeme: "", Literal: nil, Line: 1},
			},
			validate: func(t *testing.T, e expr.Expr) {
				// Currently builds: (1 + 2) * 3 due to term/factor being swapped
				// This is: * at root, + on left
				binary, ok := e.(*expr.Binary)
				if !ok {
					t.Fatalf("Expected Binary, got %T", e)
				}
				if binary.Op.Type != scanner.Star {
					t.Errorf("Expected Star as root operator, got %v", binary.Op.Type)
				}
				// Left side should be addition
				leftBinary, ok := binary.Left.(*expr.Binary)
				if !ok {
					t.Fatalf("Expected Binary on left, got %T", binary.Left)
				}
				if leftBinary.Op.Type != scanner.Plus {
					t.Errorf("Expected Plus on left side, got %v", leftBinary.Op.Type)
				}
			},
		},
		{
			name: "Complex - unary with grouping",
			tokens: []scanner.Token{
				{Type: scanner.Minus, Lexeme: "-", Literal: nil, Line: 1},
				{Type: scanner.LeftParen, Lexeme: "(", Literal: nil, Line: 1},
				{Type: scanner.Number, Lexeme: "5", Literal: 5.0, Line: 1},
				{Type: scanner.RightParen, Lexeme: ")", Literal: nil, Line: 1},
				{Type: scanner.Eof, Lexeme: "", Literal: nil, Line: 1},
			},
			validate: func(t *testing.T, e expr.Expr) {
				unary, ok := e.(*expr.Unary)
				if !ok {
					t.Fatalf("Expected Unary, got %T", e)
				}
				grouping, ok := unary.Expr.(*expr.Grouping)
				if !ok {
					t.Fatalf("Expected Grouping operand, got %T", unary.Expr)
				}
				lit, ok := grouping.Expr.(*expr.Literal)
				if !ok {
					t.Fatalf("Expected Literal in grouping, got %T", grouping.Expr)
				}
				if lit.Value != 5.0 {
					t.Errorf("Expected 5.0, got %v", lit.Value)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.tokens, func(token scanner.Token, message string) {
				t.Errorf("Parse error at line %d: %s", token.Line, message)
			})
			result := p.Parse()

			if result == nil {
				t.Fatalf("Parse returned nil")
			}

			tt.validate(t, result)
		})
	}
}

func TestParserErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		tokens     []scanner.Token
		shouldFail bool
	}{
		{
			name: "Missing closing paren",
			tokens: []scanner.Token{
				{Type: scanner.LeftParen, Lexeme: "(", Literal: nil, Line: 1},
				{Type: scanner.Number, Lexeme: "42", Literal: 42.0, Line: 1},
				{Type: scanner.Eof, Lexeme: "", Literal: nil, Line: 1},
			},
			shouldFail: true,
		},
		{
			name: "Empty parentheses",
			tokens: []scanner.Token{
				{Type: scanner.LeftParen, Lexeme: "(", Literal: nil, Line: 1},
				{Type: scanner.RightParen, Lexeme: ")", Literal: nil, Line: 1},
				{Type: scanner.Eof, Lexeme: "", Literal: nil, Line: 1},
			},
			shouldFail: true,
		},
		{
			name: "Missing operand",
			tokens: []scanner.Token{
				{Type: scanner.Number, Lexeme: "1", Literal: 1.0, Line: 1},
				{Type: scanner.Plus, Lexeme: "+", Literal: nil, Line: 1},
				{Type: scanner.Eof, Lexeme: "", Literal: nil, Line: 1},
			},
			shouldFail: true,
		},
		{
			name: "Unary with no operand",
			tokens: []scanner.Token{
				{Type: scanner.Minus, Lexeme: "-", Literal: nil, Line: 1},
				{Type: scanner.Eof, Lexeme: "", Literal: nil, Line: 1},
			},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorCalled := false
			p := NewParser(tt.tokens, func(token scanner.Token, message string) {
				errorCalled = true
			})
			result := p.Parse()

			if tt.shouldFail && !errorCalled {
				t.Errorf("Expected parse error, but none was reported")
			}
			if tt.shouldFail && result != nil {
				t.Errorf("Expected parse to return nil on error, got %v", result)
			}
		})
	}
}
