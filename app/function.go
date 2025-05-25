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
	return &LoxFunction{
		declaration:   declaration,
		closure:       closure,
		isInitializer: isInitializer,
	}
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) (ret interface{}) {
	newEnv := NewEnvironment(f.closure)
	for i, param := range f.declaration.Params {
		newEnv.Define(param.Lexeme, arguments[i])
	}

	prev := interpreter.environment
	interpreter.environment = newEnv
	defer func() { interpreter.environment = prev }()

	defer func() {
		if r := recover(); r != nil {
			if retVal, ok := r.(*ReturnValue); ok {
				ret = retVal.Value
			} else {
				panic(r)
			}
		}
	}()

	// Don't use executeBlock here!
	for _, stmt := range f.declaration.Body {
		ret = interpreter.Execute(stmt)
	}

	if f.isInitializer {
		thisVal, _ := newEnv.Get("this")
		return thisVal
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
	return NewLoxFunction(f.declaration, env, f.isInitializer)
}
