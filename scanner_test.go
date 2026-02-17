package main

import (
	"reflect"
	"testing"
)

func TestScanTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "Empty source",
			input: "",
			expected: []Token{
				{Type: Eof, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			name:  "Single character tokens",
			input: "(){},.-+;*",
			expected: []Token{
				{Type: LeftParen, Lexeme: "(", Literal: nil, Line: 1},
				{Type: RightParen, Lexeme: ")", Literal: nil, Line: 1},
				{Type: LeftBrace, Lexeme: "{", Literal: nil, Line: 1},
				{Type: RightBrace, Lexeme: "}", Literal: nil, Line: 1},
				{Type: Comma, Lexeme: ",", Literal: nil, Line: 1},
				{Type: Dot, Lexeme: ".", Literal: nil, Line: 1},
				{Type: Minus, Lexeme: "-", Literal: nil, Line: 1},
				{Type: Plus, Lexeme: "+", Literal: nil, Line: 1},
				{Type: Semicolon, Lexeme: ";", Literal: nil, Line: 1},
				{Type: Star, Lexeme: "*", Literal: nil, Line: 1},
				{Type: Eof, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			name:  "Numbers - Integer",
			input: "123",
			expected: []Token{
				{Type: Number, Lexeme: "123", Literal: 123.0, Line: 1},
				{Type: Eof, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			name:  "Numbers - Decimal",
			input: "123.456",
			expected: []Token{
				{Type: Number, Lexeme: "123.456", Literal: 123.456, Line: 1},
				{Type: Eof, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			name:  "Numbers - Multiple with whitespace",
			input: "123 456",
			expected: []Token{
				{Type: Number, Lexeme: "123", Literal: 123.0, Line: 1},
				{Type: Number, Lexeme: "456", Literal: 456.0, Line: 1},
				{Type: Eof, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			name:  "Numbers - Followed by operator",
			input: "123+456",
			expected: []Token{
				{Type: Number, Lexeme: "123", Literal: 123.0, Line: 1},
				{Type: Plus, Lexeme: "+", Literal: nil, Line: 1},
				{Type: Number, Lexeme: "456", Literal: 456.0, Line: 1},
				{Type: Eof, Lexeme: "", Literal: nil, Line: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner(tt.input)
			tokens := scanner.ScanTokens()

			if len(tokens) != len(tt.expected) {
				t.Fatalf("Expected %d tokens, got %d", len(tt.expected), len(tokens))
			}

			for i, token := range tokens {
				expected := tt.expected[i]

				if token.Type != expected.Type {
					t.Errorf("Token %d: expected type %v, got %v", i, expected.Type, token.Type)
				}
				if token.Lexeme != expected.Lexeme {
					t.Errorf("Token %d: expected lexeme %q, got %q", i, expected.Lexeme, token.Lexeme)
				}
				if expected.Literal != nil {
					// Compare literals if expected is not nil
					if !reflect.DeepEqual(token.Literal, expected.Literal) {
						t.Errorf("Token %d: expected literal %v, got %v", i, expected.Literal, token.Literal)
					}
				}
				if token.Line != expected.Line {
					t.Errorf("Token %d: expected line %d, got %d", i, expected.Line, token.Line)
				}
			}
		})
	}
}
