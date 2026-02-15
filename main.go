package main

import (
	"bufio"
	"fmt"
	"os"
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
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		run(line)
		hadError = false
	}
}

func run(line string) {
	// scanner := NewScanner(line)
	// tokens := scanner.ScanTokens()
	//
	// for _, token := range tokens {
	// 	fmt.Printf("%s %s\n", token.Type, token.Literal)
	// }
}

func error(line int, message string) {
	fmt.Printf("Error: %s [line %d]\n", message, line)
	hadError = true
}
