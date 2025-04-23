package main

import "fmt"

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

func NewLoxFunction(declaration *Function, closure *Environment) *LoxFunction {
	return &LoxFunction{
		declaration: declaration,
		closure:     closure,
	}
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	environment := NewEnvironment(f.closure)

	for i := 0; i < len(f.declaration.Params); i++ {
		environment.Define(f.declaration.Params[i].Lexeme, arguments[i])
	}

	result := interpreter.executeBlock(f.declaration.Body, environment)
	if returnValue, ok := result.(ReturnValue); ok {
		return returnValue.Value
	}

	return nil
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}
