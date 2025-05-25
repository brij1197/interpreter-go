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
	environment := NewEnvironment(f.closure)

	for i := 0; i < len(f.declaration.Params); i++ {
		environment.Define(f.declaration.Params[i].Lexeme, arguments[i])
	}

	if val, ok := environment.values["this"]; ok {
		fmt.Fprintf(os.Stderr, "DEBUG: this = %v\n", val)
	}

	previousEnv := interpreter.environment
	interpreter.environment = environment

	defer func() {
		interpreter.environment = previousEnv
	}()

	result := interpreter.executeBlock(f.declaration.Body, environment)

	if ret, ok := result.(ReturnValue); ok {
		return ret.Value
	}

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
	return &LoxFunction{
		declaration: f.declaration,
		closure:     environment,
	}
}
