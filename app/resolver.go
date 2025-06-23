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
	INITIALIZER
	METHOD
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
	//fmt.Println("BEGIN SCOPE", len(r.scopes))
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	//fmt.Println("END SCOPE", len(r.scopes)-1)
	if len(r.scopes) > 0 {
		r.scopes = r.scopes[:len(r.scopes)-1]
	}
}

func (r *Resolver) declare(name *Token) {
	if len(r.scopes) == 0 {
		return
	}
	//fmt.Printf("DECLARE %s in scope %d\n", name.Lexeme, len(r.scopes)-1)
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
	//fmt.Printf("DEFINE %s in scope %d\n", name.Lexeme, len(r.scopes)-1)
	scope := r.scopes[len(r.scopes)-1]
	scope[name.Lexeme] = true
	delete(r.inInitializer, name.Lexeme)
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			// Correct debug:
			fmt.Fprintf(os.Stderr, "Resolved %s at distance %d\n", name.Lexeme, len(r.scopes)-1-i)
			return
		}
	}
	fmt.Fprintf(os.Stderr, "Did NOT resolve %s\n", name.Lexeme)
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

func (r *Resolver) VisitVarStmt(stmt *Var) interface{} {
	r.declare(&stmt.Name)

	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}

	r.define(&stmt.Name)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *Variable) interface{} {
	if len(r.scopes) > 0 {
		if val, ok := r.scopes[len(r.scopes)-1][expr.Name.Lexeme]; ok && !val {
			panic("Can't read local variable in its own initializer.")
		}
	}
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr *Assign) interface{} {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
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
	r.resolveFunction(stmt, FUNCTION)
	return nil
}

func (r *Resolver) resolveFunction(function *Function, funcType FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = funcType

	r.beginScope()
	for _, param := range function.Params {
		r.declare(&param)
		r.define(&param)
	}

	if function.Body != nil {
		r.Resolve(function.Body)
	}

	r.endScope()

	r.currentFunction = enclosingFunction
}

func (r *Resolver) VisitBlockStmt(stmt *Block) interface{} {
	r.beginScope()
	r.Resolve(stmt.Statements)
	r.endScope()
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
		panic(fmt.Errorf("Can't return from top-level code."))
	}
	if r.currentFunction == INITIALIZER && stmt.Value != nil {
		panic(fmt.Errorf("[line %d] Error at 'return': Can't return a value from an initializer.", stmt.Keyword.Line))
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

	r.declare(&stmt.Name)
	r.define(&stmt.Name)

	if stmt.Superclass != nil {
		if superVar, ok := stmt.Superclass.(*Variable); ok {
			if stmt.Name.Lexeme == superVar.Name.Lexeme {
				panic(fmt.Errorf("A class can't inherit from itself."))
			}
		}
		r.resolveExpr(stmt.Superclass)
	}

	r.beginScope()
	r.scopes[len(r.scopes)-1]["this"] = true

	for _, method := range stmt.Methods {
		function, _ := method.(*Function)
		declaration := METHOD
		if function.Name.Lexeme == "init" {
			declaration = INITIALIZER
		}
		r.resolveFunction(function, declaration)
	}

	r.endScope()
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
		panic("Cannot use 'this' outside of a class method.")
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil
}
