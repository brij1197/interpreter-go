package main

import (
	"fmt"
	"os"
)

const (
	LEFT_PAREN  rune = '('
	RIGHT_PAREN rune = ')'
	LEFT_BRACE  rune = '{'
	RIGHT_BRACE rune = '}'
	COMMA       rune = ','
	DOT         rune = '.'
	MINUS       rune = '-'
	PLUS        rune = '+'
	SEMICOLON   rune = ';'
	STAR        rune = '*'
	EQUAL       rune = '='
	BANG        rune = '!'
)

func scanTokens(fileContents string) bool {
	hasError := false
	runes := []rune(fileContents)
	for i := 0; i < len(runes); i++ {
		current := runes[i]
		switch current {
		case LEFT_PAREN:
			fmt.Println("LEFT_PAREN ( null")
		case RIGHT_PAREN:
			fmt.Println("RIGHT_PAREN ) null")
		case LEFT_BRACE:
			fmt.Println("LEFT_BRACE { null")
		case RIGHT_BRACE:
			fmt.Println("RIGHT_BRACE } null")
		case COMMA:
			fmt.Println("COMMA , null")
		case DOT:
			fmt.Println("DOT . null")
		case MINUS:
			fmt.Println("MINUS - null")
		case PLUS:
			fmt.Println("PLUS + null")
		case SEMICOLON:
			fmt.Println("SEMICOLON ; null")
		case STAR:
			fmt.Println("STAR * null")
		case BANG:
			if i+1 < len(runes) && runes[i+1] == '=' {
				fmt.Println("BANG_EQUAL != null")
				i++
			} else {
				fmt.Println("BANG ! null")
			}
		case EQUAL:
			if i+1 < len(runes) && runes[i+1] == '=' {
				fmt.Println("EQUAL_EQUAL == null")
				i++
			} else {
				fmt.Println("EQUAL = null")
			}
		default:
			fmt.Fprintf(os.Stderr, "[line 1] Error: Unexpected character: %c\n", current)
			hasError = true
		}
	}
	fmt.Println("EOF  null")
	return hasError
}
