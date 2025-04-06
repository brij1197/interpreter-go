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
	return nil, fmt.Errorf("undefined variable '%s'", name.Lexeme)
}
