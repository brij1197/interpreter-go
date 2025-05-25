package main

type LoxClass struct {
	name    string
	methods map[string]*LoxFunction
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
	return instance
}

func (c *LoxClass) Arity() int {
	return 0
}

func (i *Interpreter) VisitThisExpr(expr *This) interface{} {
	return i.lookupVariable(expr.Keyword, expr)
}
