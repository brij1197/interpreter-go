package main

import "fmt"

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
	environment := NewEnvironment(f.closure)

	for i, param := range f.declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	result := interpreter.executeBlock(f.declaration.Body, environment)
	if ret, ok := result.(ReturnValue); ok {
		return ret.Value
	}
	return nil
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}
