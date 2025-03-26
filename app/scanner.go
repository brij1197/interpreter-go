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
	QUOTE       rune = '"'
)

func scanTokens(fileContents string) bool {
	line := 1
	hasError := false
	runes := []rune(fileContents)
	var errors []string
	var tokens []string

	for i := 0; i < len(runes); i++ {
		current := runes[i]
		switch current {
		case LEFT_PAREN:
			tokens = append(tokens, "LEFT_PAREN ( null")
		case RIGHT_PAREN:
			tokens = append(tokens, "RIGHT_PAREN ) null")
		case LEFT_BRACE:
			tokens = append(tokens, "LEFT_BRACE { null")
		case RIGHT_BRACE:
			tokens = append(tokens, "RIGHT_BRACE } null")
		case COMMA:
			tokens = append(tokens, "COMMA , null")
		case DOT:
			tokens = append(tokens, "DOT . null")
		case MINUS:
			tokens = append(tokens, "MINUS - null")
		case PLUS:
			tokens = append(tokens, "PLUS + null")
		case SEMICOLON:
			tokens = append(tokens, "SEMICOLON ; null")
		case STAR:
			tokens = append(tokens, "STAR * null")
		case BANG:
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, "BANG_EQUAL != null")
				i++
			} else {
				tokens = append(tokens, "BANG ! null")
			}
		case EQUAL:
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, "EQUAL_EQUAL == null")
				i++
			} else {
				tokens = append(tokens, "EQUAL = null")
			}
		case LESS:
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, "LESS_EQUAL <= null")
				i++
			} else {
				tokens = append(tokens, "LESS < null")
			}
		case GREATER:
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, "GREATER_EQUAL >= null")
				i++
			} else {
				tokens = append(tokens, "GREATER > null")
			}
		case SLASH:
			if i+1 < len(runes) && runes[i+1] == '/' {
				i++
				for i < len(runes) && runes[i] != NEWLINE {
					i++
				}
				i--
			} else {
				tokens = append(tokens, "SLASH / null")
			}
		case QUOTE:
			start := i + 1
			i++
			for i < len(runes) && runes[i] != QUOTE && runes[i] != NEWLINE {
				i++
			}
			if i >= len(runes) || runes[i] != QUOTE {
				errors = append(errors, fmt.Sprintf("[line %d] Error: Unterminated string", line))
				hasError = true
				break
			}
			value := string(runes[start:i])
			tokens = append(tokens, fmt.Sprintf("STRING \"%s\" %s", value, value))
		case NEWLINE:
			line++
		case SPACE, TAB:
			continue
		default:
			errors = append(errors, fmt.Sprintf("[line %d] Error: Unexpected character: %c", line, current))
			hasError = true
		}
	}

	for _, err := range errors {
		fmt.Fprintln(os.Stderr, err)
	}

	for _, token := range tokens {
		fmt.Println(token)
	}

	fmt.Println("EOF  null")
	return hasError
}
