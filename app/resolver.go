package main

import (
	"fmt"
	"os"
	"reflect"
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
	globals         map[string]bool
	inInitializer   map[string]bool
	currentClass    ClassType
}

type FunctionType int

type ClassType int

const (
	NO_CLASS ClassType = iota
	IN_CLASS
	IN_SUBCLASS
)

const (
	NONE FunctionType = iota
	FUNCTION
	CLASS_TYPE
)

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          make([]map[string]bool, 0),
		currentFunction: NONE,
		currentClass:    NO_CLASS,
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
			if r.inInitializer[name.Lexeme] {
				panic(&ParseError{
					token:   *name,
					message: "Can't read local variable in its own initializer.",
				})
			}
			depth := len(r.scopes) - 1 - i
			fmt.Fprintf(os.Stderr, "DEBUG: resolveLocal %s at depth %d\n", name.Lexeme, depth)
			r.interpreter.resolve(expr, depth)
			return
		}

	}
	// fallback: global scope
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
	r.resolveFunction(stmt)
	return nil
}

func (r *Resolver) resolveFunction(function *Function) {
	enclosingFunction := r.currentFunction
	r.currentFunction = FUNCTION

	r.beginScope()

	if r.currentClass != NO_CLASS {
		r.scopes[len(r.scopes)-1]["this"] = true
	}

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

	// Pre-declare and define functions to capture outer scope correctly
	for _, s := range stmt.Statements {
		if fn, ok := s.(*Function); ok {
			r.declare(&fn.Name)
			r.define(&fn.Name)
			fmt.Fprintf(os.Stderr, "DEBUG: Pre-defining function %s\n", fn.Name.Lexeme)
		}
	}

	for _, s := range stmt.Statements {
		r.resolveStmt(s)
	}

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
	enclosingClass := r.currentClass
	r.currentClass = IN_CLASS

	// Declare and define the class name in the current scope
	r.declare(&stmt.Name)
	r.define(&stmt.Name)

	// Do not push another scope here!
	// Each method gets its own scope (handled in resolveFunction)

	for _, method := range stmt.Methods {
		if function, ok := method.(*Function); ok {
			r.resolveFunction(function)
		}
	}

	r.currentClass = enclosingClass
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

func (r *Resolver) VisitThisExpr(expr *This) interface{} {
	if r.currentClass == NO_CLASS {
		panic(NewParseError(expr.Keyword, "Can't use 'this' outside of a class method."))
	}
	r.resolveLocal(expr, &expr.Keyword)
	return nil
}
