package main

import (
	"fmt"
	"reflect"
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
	globals         map[string]bool
	inInitializer   map[string]bool
}

type FunctionType int

const (
	NONE FunctionType = iota
	FUNCTION
)

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          make([]map[string]bool, 0),
		currentFunction: NONE,
		globals:         make(map[string]bool),
		inInitializer:   make(map[string]bool),
	}
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	if len(r.scopes) > 0 {
		r.scopes = r.scopes[:len(r.scopes)-1]
	}
}

func (r *Resolver) declare(name *Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	if _, exists := scope[name.Lexeme]; exists {
		panic(
			&ParseError{
				token:   *name,
				message: "Variable already declared in this scope.",
			},
		)
	}
	scope[name.Lexeme] = false
	r.inInitializer[name.Lexeme] = true
}

func (r *Resolver) define(name *Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.Lexeme] = true
	delete(r.inInitializer, name.Lexeme)

}

func (r *Resolver) resolveLocal(expr Expr, name *Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			if len(r.scopes) > 0 && r.inInitializer[name.Lexeme] {
				panic(&ParseError{
					token:   *name,
					message: "Can't read local variable in its own initializer.",
				})
			}
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) VisitBinaryExpr(expr *Binary) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *Grouping) interface{} {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *Literal) interface{} {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *Logical) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *Unary) interface{} {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitVariableStmt(stmt *Var) interface{} {
	r.declare(&stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(&stmt.Name)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *Variable) interface{} {
	if len(r.scopes) > 0 {
		if initialized, ok := r.scopes[len(r.scopes)-1][expr.Name.Lexeme]; ok && !initialized {
			panic(NewParseError(expr.Name, "Can't read local variable in its own initializer."))
		}
	}
	r.resolveLocal(expr, &expr.Name)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr *Assign) interface{} {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, &expr.Name)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *Call) interface{} {
	r.resolveExpr(expr.Callee)
	for _, argument := range expr.Arguments {
		r.resolveExpr(argument)
	}
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *Expression) interface{} {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *Function) interface{} {
	r.declare(&stmt.Name)
	r.define(&stmt.Name)

	r.resolveFunction(stmt)
	return nil
}

func (r *Resolver) resolveFunction(function *Function) {
	enclosingFunction := r.currentFunction
	r.currentFunction = FUNCTION

	r.beginScope()

	for _, param := range function.Params {
		r.declare(&param)
		r.define(&param)
	}

	r.resolveStatements(function.Body)

	r.endScope()
	r.currentFunction = enclosingFunction
}

func (r *Resolver) VisitBlockStmt(stmt *Block) interface{} {
	r.beginScope()
	r.Resolve(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitVarStmt(stmt *Var) interface{} {
	r.declare(&stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(&stmt.Name)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *If) interface{} {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *Print) interface{} {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *ReturnStmt) interface{} {
	if r.currentFunction == NONE {
		panic(NewParseError(stmt.Keyword, "Can't return from top-level code."))
	}
	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *While) interface{} {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil
}

func (r *Resolver) Resolve(statements interface{}) {
	switch v := statements.(type) {
	case []Stmt:
		for _, statement := range v {
			r.resolveStmt(statement)
		}
	case Stmt:
		r.resolveStmt(v)
	case Expr:
		r.resolveExpr(v)
	default:
		fmt.Println("Unknown type in resolver:", reflect.TypeOf(statements))
	}
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) VisitResolverStmt(stmt *Resolver) interface{} {
	return nil
}

func (r *Resolver) resolveStatements(statements []Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *Resolver) VisitFunctionExpr(expr *FunctionExpr) interface{} {
	enclosingFunction := r.currentFunction
	r.currentFunction = FUNCTION

	r.beginScope()
	if expr.Name.Lexeme != "" {
		r.declare(&expr.Name)
		r.define(&expr.Name)
	}

	for _, param := range expr.Params {
		r.declare(&param)
		r.define(&param)
	}

	r.resolveStatements(expr.Body)
	r.endScope()

	r.currentFunction = enclosingFunction
	return nil
}

func (r *Resolver) VisitClassStmt(stmt *Class) interface{} {
	r.declare(&stmt.Name)
	r.define(&stmt.Name)
	return nil
}

func (r *Resolver) VisitGetExpr(expr *Get) interface{} {
	r.resolveExpr(expr.Object)
	return nil
}

func (r *Resolver) VisitSetExpr(expr *Set) interface{} {
	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)
	return nil
}
