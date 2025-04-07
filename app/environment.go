package main

import "fmt"

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]interface{}),
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Get(name Token) (interface{}, error) {
	if val, ok := e.values[name.Lexeme]; ok {
		return val, nil
	}
	return nil, NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
}

func (e *Environment) Assign(name Token, value interface{}) error {
	if _, exists := e.values[name.Lexeme]; exists {
		e.values[name.Lexeme] = value
		return nil
	}
	return NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
}
