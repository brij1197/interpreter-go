package main

import "fmt"

type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    make(map[string]interface{}),
		enclosing: enclosing,
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Get(name Token) (interface{}, error) {
	if val, ok := e.values[name.Lexeme]; ok {
		return val, nil
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, &RuntimeError{
		token:   name,
		message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
	}
}

func (e *Environment) Assign(name Token, value interface{}) error {
	if _, exists := e.values[name.Lexeme]; exists {
		e.values[name.Lexeme] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}
	return NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
}
