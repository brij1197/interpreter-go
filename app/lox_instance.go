package main

import (
	"fmt"
	"os"
)

type LoxInstance struct {
	class  *LoxClass
	fields map[string]interface{}
}

var _ fmt.Stringer = (*LoxInstance)(nil)

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		class:  class,
		fields: make(map[string]interface{}),
	}
}

func (i *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", i.class.name)
}

func (i *LoxInstance) Get(name Token) interface{} {
	if value, ok := i.fields[name.Lexeme]; ok {
		return value
	}

	method := i.class.FindMethod(name.Lexeme)
	if method != nil {
		bound := method.Bind(i)
		fmt.Fprintf(os.Stderr, "DEBUG: Returning bound method with closure = %p\n", bound.closure)
		return bound
	}

	panic(&RuntimeError{
		token:   name,
		message: fmt.Sprintf("Undefined property '%s'.", name.Lexeme),
	})
}

func (i *LoxInstance) Set(name Token, value interface{}) {
	i.fields[name.Lexeme] = value
}
