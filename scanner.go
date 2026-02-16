package main

import "fmt"

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  make([]Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	eofToken := Token{Type: Eof, Lexeme: "", Literal: nil, Line: s.line}
	s.tokens = append(s.tokens, eofToken)

	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
	case '(':
		s.addToken(LeftParen)
	case ')':
		s.addToken(RightParen)
	case '{':
		s.addToken(LeftBrace)
	case '}':
		s.addToken(RightBrace)
	case ',':
		s.addToken(Comma)
	case '.':
		s.addToken(Dot)
	case '-':
		s.addToken(Minus)
	case '+':
		s.addToken(Plus)
	case ';':
		s.addToken(Semicolon)
	case '*':
		s.addToken(Star)
	case '!':
		if s.match('=') { s.addToken(BangEqual) } else { s.addToken(Bang) }
	case '=':
		if s.match('=') { s.addToken(EqualEqual) } else { s.addToken(Equal) }
	case '>':
		if s.match('=') { s.addToken(GreaterEqual) } else { s.addToken(Greater) }
	case '<':
		if s.match('=') { s.addToken(LessEqual) } else { s.addToken(Less) }
	default:
		Error(s.line, fmt.Sprintf("Unexpected character: %s", string(char)))
	}
}

func (s *Scanner) match(char byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != char {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal any) {
	token := Token{
		Type:    tokenType,
		Lexeme:  string(s.source[s.start:s.current]),
		Literal: literal,
		Line:    s.line,
	}

	s.tokens = append(s.tokens, token)
}
