package main

type Resolver struct {
	interpreter     *Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
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
		panic("Variable with this name already declared in this scope.")
	}
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name *Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[len(r.scopes)-1][name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name *Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) VisitBinaryExpr(expr *Binary) interface{} {
	r.Resolve(expr.Left)
	r.Resolve(expr.Right)
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *Grouping) interface{} {
	r.Resolve(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *Literal) interface{} {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *Logical) interface{} {
	r.Resolve(expr.Left)
	r.Resolve(expr.Right)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *Unary) interface{} {
	r.Resolve(expr.Right)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *Variable) interface{} {
	if len(r.scopes) > 0 {
		scope := r.scopes[len(r.scopes)-1]
		if val, ok := scope[expr.Name.Lexeme]; ok && !val {
			panic("Can't read local variable in its own initializer.")
		}
	}
	r.resolveLocal(expr, &expr.Name)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr *Assign) interface{} {
	r.Resolve(expr.Value)
	r.resolveLocal(expr, &expr.Name)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *Call) interface{} {
	r.Resolve(expr.Callee)
	for _, argument := range expr.Arguments {
		r.Resolve(argument)
	}
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *Expression) interface{} {
	r.Resolve(stmt.Expression)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *Function) interface{} {
	enclosingFunction := r.currentFunction
	r.currentFunction = FUNCTION

	r.beginScope()

	for _, param := range stmt.Params {
		r.declare(&param)
		r.define(&param)
	}
	for _, statement := range stmt.Body {
		r.resolveStmt(statement)
	}
	r.endScope()

	r.currentFunction = enclosingFunction

	return nil
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
		r.Resolve(stmt.Initializer)
	}
	r.define(&stmt.Name)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *If) interface{} {
	r.Resolve(stmt.Condition)
	r.Resolve(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.Resolve(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *Print) interface{} {
	r.Resolve(stmt.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *ReturnStmt) interface{} {
	if r.currentFunction == NONE {
		panic("Can't return from top-level code.")
	}
	if stmt.Value != nil {
		r.Resolve(stmt.Value)
	}
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *While) interface{} {
	r.Resolve(stmt.Condition)
	r.Resolve(stmt.Body)
	return nil
}

func (r *Resolver) Resolve(statements interface{}) {
	switch v := statements.(type) {
	case []Stmt:
		for _, statement := range v {
			r.resolveStmt(statement)
		}
	case Stmt:
		v.Accept(r)
	case Expr:
		v.Accept(r)
	}
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.Accept(r)
}
