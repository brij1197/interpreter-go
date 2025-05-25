package main

import (
	"fmt"
	"os"
)

type LoxFunction struct {
	declaration *Function
	closure     *Environment
}

type Function struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

func NewLoxFunction(declaration *Function, closure *Environment) *LoxFunction {
	return &LoxFunction{
		declaration: declaration,
		closure:     closure,
	}
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	env := f.closure

	for i, param := range f.declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	var result interface{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				if ret, ok := r.(*ReturnValue); ok {
					result = ret.Value
				} else {
					panic(r)
				}
			}
		}()
		interpreter.executeBlock(f.declaration.Body, env)
	}()

	return result
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	environment := NewEnvironment(f.closure)
	environment.Define("this", instance)

	fmt.Fprintf(os.Stderr, "DEBUG: Binding this to %v\n", instance)
	fmt.Fprintf(os.Stderr, "DEBUG: f.closure = %p, newEnv = %p\n", f.closure, environment)

	return &LoxFunction{
		declaration: f.declaration,
		closure:     environment,
	}
}
