package main

type Callable interface {
	Call(interpreter *Interpreter, arguments []interface{}) interface{}
	Arity() int
}

type LoxCallable interface {
	Call(interpreter *Interpreter, arguments []interface{}) interface{}
	String() string
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
