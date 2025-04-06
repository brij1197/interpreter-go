package main

type Stmt interface {
	Accept(visitor StmtVisitor) interface{}
}

type StmtVisitor interface {
	VisitPrintStmt(stmt *Print) interface{}
	VisitExpressionStmt(stmt *Expression) interface{}
	VisitVarStmt(stmt *Var) interface{}
}

type Print struct {
	Expression Expr
}

func (stmt *Print) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitPrintStmt(stmt)
}

type Expression struct {
	Expression Expr
}

func (s *Expression) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExpressionStmt(s)
}

type Var struct {
	Name        Token
	initializer Expr
}

func (v *Var) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVarStmt(v)
}
