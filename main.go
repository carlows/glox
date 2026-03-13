package main

import (
	"bufio"
	"fmt"
	"os"

	"glox/interpreter"
	"glox/parser"
	"glox/scanner"
)

var hadError bool = false

func main() {
	if len(os.Args) > 3 {
		fmt.Println("Usage: glox <script>")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runInteractive()
	}
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	run(string(content))

	if hadError {
		os.Exit(65)
	}
}

func runInteractive() {
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !sc.Scan() {
			break
		}

		line := sc.Text()
		run(line)
		hadError = false
	}
}

func run(line string) {
	s := scanner.NewScanner(line, Error)
	tokens := s.ScanTokens()

	p := parser.NewParser(tokens, ErrorAtToken)
	expr := p.Parse()

	if hadError {
		return
	}

	i := interpreter.NewInterpreter(ErrorAtToken)
	i.Interpret(expr)
}

func Error(line int, message string) {
	fmt.Printf("Error: %s [line %d]\n", message, line)
	hadError = true
}

func ErrorAtToken(token scanner.Token, message string) {
	fmt.Printf("Error: %s [line %d] lexeme: %s\n", message, token.Line, token.Lexeme)
	hadError = true
}
