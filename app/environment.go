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

func (e *Environment) Get(name string) (interface{}, error) {
	if value, ok := e.values[name]; ok {
		return value, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return nil, fmt.Errorf("Undefined variable '%s'.", name)
}

func (e *Environment) Assign(name Token, value interface{}) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}

	if e.enclosing == nil {
		e.values[name.Lexeme] = value
		return nil
	}

	return &RuntimeError{
		token:   name,
		message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
	}
}

func (e *Environment) GetAt(distance int, name string) interface{} {
	return e.ancestor(distance).values[name]
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}

func (e *Environment) AssignAt(distance int, name Token, value interface{}) {
	e.ancestor(distance).values[name.Lexeme] = value
}
