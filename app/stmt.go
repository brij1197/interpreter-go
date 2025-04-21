package main

type StmtVisitor interface {
	VisitPrintStmt(stmt *Print) interface{}
	VisitExpressionStmt(stmt *Expression) interface{}
	VisitVarStmt(stmt *Var) interface{}
	VisitBlockStmt(stmt *Block) interface{}
	VisitIfStmt(stmt *If) interface{}
	VisitWhileStmt(stmt *While) interface{}
	VisitFunctionStmt(stmt *Function) interface{}
	VisitReturnStmt(stmt *ReturnStmt) interface{}
}

type Stmt interface {
	Accept(visitor StmtVisitor) interface{}
}

type Print struct {
	Expression Expr
}

type Expression struct {
	Expression Expr
}

type Var struct {
	Name        Token
	initializer Expr
}

type Block struct {
	Statements []Stmt
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type While struct {
	Condition Expr
	Body      Stmt
}

func (stmt *Print) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitPrintStmt(stmt)
}

func (s *Expression) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExpressionStmt(s)
}

func (v *Var) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVarStmt(v)
}

func (b *Block) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitBlockStmt(b)
}

func (i *If) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitIfStmt(i)
}

func (w *While) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitWhileStmt(w)
}

func (s *Function) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitFunctionStmt(s)
}

func (r *ReturnStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitReturnStmt(r)
}
