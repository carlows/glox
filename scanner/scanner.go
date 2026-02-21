package scanner

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	source       string
	tokens       []Token
	start        int
	current      int
	line         int
	errorHandler func(line int, message string)
}

var keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"fun":    Fun,
	"for":    For,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

func NewScanner(source string, errorHandler func(line int, message string)) *Scanner {
	return &Scanner{
		source:       source,
		tokens:       make([]Token, 0),
		start:        0,
		current:      0,
		line:         1,
		errorHandler: errorHandler,
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
		if s.match('=') {
			s.addToken(BangEqual)
		} else {
			s.addToken(Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(EqualEqual)
		} else {
			s.addToken(Equal)
		}
	case '>':
		if s.match('=') {
			s.addToken(GreaterEqual)
		} else {
			s.addToken(Greater)
		}
	case '<':
		if s.match('=') {
			s.addToken(LessEqual)
		} else {
			s.addToken(Less)
		}
	case '/':
		// Keeps consuming until we reach the new line character
		// because we reached a comment, such as this current one :)
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(Slash)
		}
	case '"':
		s.string()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.number()
	case ' ', '\t', '\r':
		// Ignore whitespace
	case '\n':
		s.line++
	default:
		if s.isAlpha(char) {
			s.identifier()
		} else {
			s.errorHandler(s.line, fmt.Sprintf("Unexpected character: %s", string(char)))
		}
	}
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// Important: we don't want to consume the dot unless it's actually a decimal
	// Because you might otherwise not be able to support 123.abs()
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	decimal := s.source[s.start:s.current]

	conversion, err := strconv.ParseFloat(decimal, 64)
	if err != nil {
		s.errorHandler(s.line, fmt.Sprintf("Invalid number: %s", decimal))
	}

	s.addTokenWithLiteral(Number, conversion)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	if keyword, ok := keywords[s.source[s.start:s.current]]; ok {
		s.addToken(keyword)
		return
	}

	s.addToken(Identifier)
}

func (s *Scanner) isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		char == '_'
}

func (s *Scanner) isAlphaNumeric(char byte) bool {
	return s.isAlpha(char) || s.isDigit(char)
}

func (s *Scanner) isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.errorHandler(s.line, "Unterminated string")
	}

	s.advance()
	s.addTokenWithLiteral(String, s.source[s.start+1:s.current-1])
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

// "Consumes" the current character and returns it
func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

// Does not "consume" the current character and returns it
func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\x00'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return s.source[s.current+1]
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
