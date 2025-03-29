package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	filename := os.Args[2]

	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "tokenize":
		runTokenize(string(bytes))
	case "parse":
		runParse(string(bytes))
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(64)
	}
}

func runParse(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := NewParser(tokens)
	expression, err := parser.parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(65)
	}
	printer := &AstPrinter{}
	fmt.Println(printer.Print(expression))
}

func runTokenize(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	for _, token := range tokens {
		var literalStr string
		if token.Literal == nil {
			literalStr = "null"
		} else if token.Type == NUMBER {
			if num, ok := token.Literal.(float64); ok {
				literalStr = fmt.Sprintf("%g", num)
			} else if num, ok := token.Literal.(int); ok {
				literalStr = fmt.Sprintf("%.1f", float64(num))
			} else {
				literalStr = fmt.Sprintf("%v", token.Literal)
			}
		} else {
			literalStr = fmt.Sprintf("%v", token.Literal)
		}
		fmt.Printf("%s %s %s\n", token.Type, token.Lexeme, literalStr)
	}
}
