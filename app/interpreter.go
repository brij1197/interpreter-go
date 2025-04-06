package main

import (
	"fmt"
	"os"
	"strings"
)

type Interpreter struct{}

type RuntimeError struct {
	token   Token
	message string
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", e.message, e.token.Line)
}

func (i *Interpreter) Evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) interface{} {
	return i.Evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) interface{} {
	right := i.Evaluate(expr.Right)
	switch expr.Operator.Type {
	case MINUS:
		if _, ok := right.(float64); !ok {
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operand must be a number.",
			})
		}
		return -(right.(float64))
	case BANG:
		return !i.isTruthy(right)
	}
	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) interface{} {
	left := i.Evaluate(expr.Left)
	right := i.Evaluate(expr.Right)

	switch expr.Operator.Type {
	case STAR:
		if _, ok := left.(float64); !ok {
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be numbers.",
			})
		}
		if _, ok := right.(float64); !ok {
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be numbers.",
			})
		}
		return left.(float64) * right.(float64)
	case SLASH:
		if _, ok := left.(float64); !ok {
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be numbers.",
			})
		}
		if _, ok := right.(float64); !ok {
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be numbers.",
			})
		}
		return left.(float64) / right.(float64)
	case EQUAL_EQUAL:
		return i.isEqual(left, right)
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	case GREATER:
		if _, ok := left.(float64); !ok {
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be numbers.",
			})
		}
		if _, ok := right.(float64); !ok {
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be numbers.",
			})
		}
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		if _, ok := left.(float64); !ok {
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be numbers.",
			})
		}
		if _, ok := right.(float64); !ok {
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be numbers.",
			})
		}
		return left.(float64) >= right.(float64)
	case LESS:
		if _, ok := left.(float64); !ok {
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be numbers.",
			})
		}
		if _, ok := right.(float64); !ok {
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be numbers.",
			})
		}
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		if _, ok := left.(float64); !ok {
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be numbers.",
			})
		}
		if _, ok := right.(float64); !ok {
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be numbers.",
			})
		}
		return left.(float64) <= right.(float64)
	case PLUS:
		if lStr, lOk := left.(string); lOk {
			if rStr, rOk := right.(string); rOk {
				return lStr + rStr
			}
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be two numbers or two strings.",
			})
		}
		if lNum, lOk := left.(float64); lOk {
			if rNum, rOk := right.(float64); rOk {
				return lNum + rNum
			}
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be two numbers or two strings.",
			})
		}
		panic(&RuntimeError{
			token:   expr.Operator,
			message: "Operands must be two numbers or two strings.",
		})
	case MINUS:
		if _, ok := left.(float64); !ok {
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be numbers.",
			})
		}
		if _, ok := right.(float64); !ok {
			panic(&RuntimeError{
				token:   expr.Operator,
				message: "Operands must be numbers.",
			})
		}
		return left.(float64) - right.(float64)
	}
	return nil
}

func (i *Interpreter) Interpret(statements []Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(*RuntimeError); ok {
				fmt.Fprintln(os.Stderr, runtimeErr.Error())
				os.Exit(70)
			} else {
				panic(r)
			}
		}
	}()
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) isEqual(left, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}
	if aStr, aOk := left.(string); aOk {
		if bStr, bOk := right.(string); bOk {
			return aStr == bStr
		}
		return false
	}
	if aNum, aOk := left.(float64); aOk {
		if bNum, bOk := right.(float64); bOk {
			return aNum == bNum
		}
		return false
	}
	if aBool, aOk := left.(bool); aOk {
		if bBool, bOk := right.(bool); bOk {
			return aBool == bBool
		}
		return false
	}
	return false
}

func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if b, ok := object.(bool); ok {
		return b
	}
	return true
}

func (i *Interpreter) stringify(value interface{}) string {
	if value == nil {
		return "nil"
	}

	switch v := value.(type) {
	case float64:
		text := fmt.Sprintf("%g", v)
		text = strings.TrimSuffix(text, ".0")
		return text
	case string:
		return v
	case bool:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (i *Interpreter) VisitExpressionStmt(stmt *Expression) interface{} {
	i.Evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *Print) interface{} {
	value := i.Evaluate(stmt.Expression)
	fmt.Println(i.stringify(value))
	return nil
}
