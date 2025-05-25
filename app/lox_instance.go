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

func (instance *LoxInstance) Get(name Token) interface{} {
	// Check fields first
	if value, ok := instance.fields[name.Lexeme]; ok {
		fmt.Fprintf(os.Stderr, "DEBUG: Get field %s = %v\n", name.Lexeme, value)
		return value
	}

	// Then check methods
	if method := instance.class.FindMethod(name.Lexeme); method != nil {
		// If it's a method, bind 'this' to the instance
		return method.Bind(instance)
	}

	panic(&RuntimeError{
		token:   name,
		message: fmt.Sprintf("Undefined property '%s'.", name.Lexeme),
	})
}

func (i *LoxInstance) Set(name Token, value interface{}) {
	i.fields[name.Lexeme] = value
}
