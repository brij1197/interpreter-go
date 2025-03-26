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
	LESS        rune = '<'
	GREATER     rune = '>'
	SLASH       rune = '/'
	SPACE       rune = ' '
	TAB         rune = '\t'
	NEWLINE     rune = '\n'
)

func scanTokens(fileContents string) bool {
	line := 1
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
		case LESS:
			if i+1 < len(runes) && runes[i+1] == '=' {
				fmt.Println("LESS_EQUAL <= null")
				i++
			} else {
				fmt.Println("LESS < null")
			}
		case GREATER:
			if i+1 < len(runes) && runes[i+1] == '=' {
				fmt.Println("GREATER_EQUAL >= null")
				i++
			} else {
				fmt.Println("GREATER > null")
			}
		case SLASH:
			if i+1 < len(runes) && runes[i+1] == '/' {
				for {
					if i >= len(runes) {
						break
					}
					if runes[i] == '\n' {
						line++
						i++
						break
					}
					i++
				}
				// if i < len(runes) {
				// 	i++
				// }
			} else {
				fmt.Println("SLASH / null")
			}
		case NEWLINE:
			line++
		case SPACE, TAB:
			continue
		default:
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line, current)
			hasError = true
		}
	}
	fmt.Println("EOF  null")
	return hasError
}
