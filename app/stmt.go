package main

type StmtVisitor interface {
	VisitPrintStmt(stmt *Print) interface{}
	VisitExpressionStmt(stmt *Expression) interface{}
	VisitVarStmt(stmt *Var) interface{}
	VisitBlockStmt(stmt *Block) interface{}
	VisitIfStmt(stmt *If) interface{}
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
