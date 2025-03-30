package main

import (
	"fmt"
)

type AstPrinter struct{}

func (a *AstPrinter) Print(expr Expr) string {
	return expr.Accept(a).(string)
}

func (a *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	if num, ok := expr.Value.(float64); ok {
		return fmt.Sprintf("%.1f", num)
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var result string
	result += "(" + name
	for _, expr := range exprs {
		result += " "
		result += expr.Accept(a).(string)
	}
	result += ")"
	return result
}
