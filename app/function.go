package main

import (
	"fmt"
)

type LoxFunction struct {
	declaration   *Function
	closure       *Environment
	isInitializer bool
}

type Function struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

func NewLoxFunction(declaration *Function, closure *Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{declaration: declaration, closure: closure, isInitializer: isInitializer}
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) (result interface{}) {
	environment := NewEnvironment(f.closure)
	for i := 0; i < len(f.declaration.Params); i++ {
		environment.Define(f.declaration.Params[i].Lexeme, arguments[i])
	}

	defer func() {
		if r := recover(); r != nil {
			if ret, ok := r.(*ReturnValue); ok {
				result = ret.Value
			} else {
				panic(r)
			}
		}
	}()

	interpreter.executeBlock(f.declaration.Body, environment)

	if f.isInitializer {
		result = f.closure.GetAt(0, "this")
	}
	return
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	env := NewEnvironment(f.closure)
	env.Define("this", instance)
	return &LoxFunction{
		declaration:   f.declaration,
		closure:       env,
		isInitializer: f.isInitializer,
	}
}
