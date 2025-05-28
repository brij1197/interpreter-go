package main

import (
	"fmt"
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
	if value, ok := instance.fields[name.Lexeme]; ok {
		return value
	}
	method := instance.class.FindMethod(name.Lexeme)
	if method != nil {
		return method.Bind(instance)
	}
	panic(&RuntimeError{name, "Undefined property '" + name.Lexeme + "'."})

}

func (i *LoxInstance) Set(name Token, value interface{}) {
	i.fields[name.Lexeme] = value
}
