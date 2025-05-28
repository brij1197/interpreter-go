package main

type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]*LoxFunction
}

func NewLoxClass(name string, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{name: name, methods: methods}
}

func (c *LoxClass) String() string {
	return c.name
}

func (c *LoxClass) FindMethod(name string) *LoxFunction {
	if method, ok := c.methods[name]; ok {
		return method
	}
	return nil
}

func (c *LoxClass) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	instance := NewLoxInstance(c)
	if initializer, ok := c.methods["init"]; ok {
		initializer.Bind(instance).Call(interpreter, arguments)
	}
	return instance
}

func (c *LoxClass) Arity() int {
	if initializer := c.FindMethod("init"); initializer != nil {
		return initializer.Arity()
	}
	return 0
}
