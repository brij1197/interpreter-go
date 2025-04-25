package main

import "fmt"

type RuntimeError struct {
	token   Token
	message string
}

func NewRuntimeError(token Token, message string) *RuntimeError {
	return &RuntimeError{
		token:   token,
		message: message,
	}
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", e.message, e.token.Line)
}

type ParseError struct {
	token   Token
	message string
}

func NewParseError(token Token, message string) *ParseError {
	return &ParseError{
		token:   token,
		message: message,
	}
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("[line %d] Error at '%s': %s", e.token.Line, e.token.Lexeme, e.message)
}
