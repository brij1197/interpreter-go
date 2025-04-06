package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (a *AstPrinter) Print(expr Expr) string {
	return expr.Accept(a).(string)
}

func (a *AstPrinter) VisitVariableExpr(expr *Variable) interface{} {
	return expr.Name.Lexeme
}

func (a *AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	switch v := expr.Value.(type) {
	case string:
		return v
	case float64:
		str := fmt.Sprintf("%v", v)
		if !strings.Contains(str, ".") {
			return fmt.Sprintf("%.1f", v)
		}
		return str
	default:
		return fmt.Sprintf("%v", expr.Value)
	}
}

func (a *AstPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
	return a.parenthesize("group", expr.Expression)
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
