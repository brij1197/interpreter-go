package main

type StmtVisitor interface {
	VisitPrintStmt(stmt *Print) interface{}
	VisitExpressionStmt(stmt *Expression) interface{}
	VisitVarStmt(stmt *Var) interface{}
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

func (stmt *Print) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitPrintStmt(stmt)
}

func (s *Expression) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExpressionStmt(s)
}

func (v *Var) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVarStmt(v)
}
