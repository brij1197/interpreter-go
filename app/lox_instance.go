package main

import "fmt"

type LoxInstance struct {
	class *LoxClass
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{class: class}
}

func (i *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", i.class.name)
}
