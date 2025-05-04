package main

type LoxClass struct {
	name string
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{name: name}
}

func (c *LoxClass) String() string {
	return c.name
}

func (c *LoxClass) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	instance := NewLoxInstance(c)
	return instance
}

func (c *LoxClass) Arity() int {
	return 0
}
