package main

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

type Literal struct {
	Value interface{}
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
}

func (b *Binary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

func (l *Literal) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}
