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

func (a *AstPrinter) VisitAssignExpr(expr *Assign) interface{} {
	return a.parenthesize("=", &Variable{Name: expr.Name}, expr.Value)
}

func (a *AstPrinter) VisitLogicalExpr(expr *Logical) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitCallExpr(expr *Call) interface{} {
	calleeStr := expr.Callee.Accept(a).(string)

	var allExprs []Expr
	allExprs = append(allExprs, expr.Arguments...)

	return a.parenthesize(calleeStr, allExprs...)
}

func (a *AstPrinter) VisitFunctionExpr(expr *FunctionExpr) interface{} {
	var builder strings.Builder

	if expr.Name.Lexeme != "" {
		builder.WriteString(expr.Name.Lexeme)
	} else {
		builder.WriteString("anonymous")
	}

	builder.WriteString("(")

	for i, param := range expr.Params {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(param.Lexeme)
	}

	builder.WriteString(") ")
	builder.WriteString("{ ... }")

	return builder.String()
}

func (a *AstPrinter) VisitGetExpr(expr *Get) interface{} {
	objectStr := expr.Object.Accept(a).(string)
	return fmt.Sprintf("%s.%s", objectStr, expr.Name.Lexeme)
}

func (a *AstPrinter) VisitSetExpr(expr *Set) interface{} {
	objectStr := expr.Object.Accept(a).(string)
	valueStr := expr.Value.Accept(a).(string)
	return fmt.Sprintf("(%s.%s = %s)", objectStr, expr.Name.Lexeme, valueStr)
}
