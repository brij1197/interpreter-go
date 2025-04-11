package main

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

type Variable struct {
	Name Token
}

type Literal struct {
	Value interface{}
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

type Grouping struct {
	Expression Expr
}

type Unary struct {
	Operator Token
	Right    Expr
}

type Assign struct {
	Name  Token
	Value Expr
}

type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitVariableExpr(expr *Variable) interface{}
	VisitAssignExpr(expr *Assign) interface{}
	VisitLogicalExpr(expr *Logical) interface{}
}

func (b *Binary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

func (l *Literal) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

func (g *Grouping) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(g)
}

func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

func (v *Variable) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariableExpr(v)
}

func (v *Assign) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitAssignExpr(v)
}

func (l *Logical) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLogicalExpr(l)
}
