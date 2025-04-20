package main

import "fmt"

type Callable interface {
	Call(interpreter *Interpreter, arguments []interface{}) interface{}
	Arity() int
}

type LoxCallable interface {
	Call(interpreter *Interpreter, arguments []interface{}) interface{}
	Arity() int
}

type NativeFunction struct {
	name     string
	function func([]interface{}) interface{}
	arity    int
}

func (n *NativeFunction) Call(_ *Interpreter, arguments []interface{}) interface{} {
	return n.function(arguments)
}

func (n *NativeFunction) Arity() int {
	return n.arity
}

func (n *NativeFunction) String() string {
	return fmt.Sprintf("<native fn %s>", n.name)
}
