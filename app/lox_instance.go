package main

import "fmt"

type LoxInstance struct {
	class  *LoxClass
	fields map[string]interface{}
}

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
	panic(&RuntimeError{
		token:   name,
		message: fmt.Sprintf("Undefined property '%s'.", name.Lexeme),
	})
}

func (i *LoxInstance) Set(name Token, value interface{}) {
	i.fields[name.Lexeme] = value
}
