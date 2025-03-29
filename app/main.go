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
	tokens, err := scanner.ScanTokens()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(65)
	}
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
	tokens, err := scanner.ScanTokens()

	for _, token := range tokens {
		var literalStr string
		if token.Literal == nil {
			literalStr = "null"
		} else if token.Type == NUMBER {
			switch v := token.Literal.(type) {
			case float64:
				if v == float64(int(v)) {
					literalStr = fmt.Sprintf("%.1f", v)
				} else {
					literalStr = fmt.Sprintf("%g", v)
				}
			case int:
				literalStr = fmt.Sprintf("%.1f", float64(v))
			default:
				literalStr = fmt.Sprintf("%v", token.Literal)
			}
		} else {
			literalStr = fmt.Sprintf("%v", token.Literal)
		}
		fmt.Printf("%s %s %s\n", token.Type, token.Lexeme, literalStr)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}
}
