package main

import (
	"fmt"
	"os"
	"time"
)

type Interpreter struct {
	environment *Environment
	globals     *Environment
	locals      map[Expr]int
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment(nil)
	i := &Interpreter{
		environment: globals,
		globals:     globals,
		locals:      make(map[Expr]int),
	}

	i.globals.Define("clock", &NativeFunction{
		name:  "clock",
		arity: 0,
		function: func(arguments []interface{}) interface{} {
			return float64(time.Now().Unix())
		},
	})
	return i
}

func (i *Interpreter) Evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitVarStmt(stmt *Var) interface{} {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.Evaluate(stmt.Initializer)
	}
	fmt.Fprintf(os.Stderr, "Defining %s in environment %p (parent %p)\n", stmt.Name.Lexeme, i.environment, i.environment.enclosing)
	i.environment.Define(stmt.Name.Lexeme, value)
	return nil
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
		}
		if lNum, lOk := left.(float64); lOk {
			if rNum, rOk := right.(float64); rOk {
				return lNum + rNum
			}
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

func (i *Interpreter) Interpret(statements []Stmt) error {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(*RuntimeError); ok {
				fmt.Fprintf(os.Stderr, "%s\n [line %d]\n", err.message, err.token.Line)
				os.Exit(70)
			}
		}
	}()

	for _, statement := range statements {
		i.Execute(statement)
	}

	return nil
}

func (i *Interpreter) Execute(stmt Stmt) interface{} {
	return stmt.Accept(i)
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

func (i *Interpreter) stringify(obj interface{}) string {
	if obj == nil {
		return "nil"
	}

	if num, ok := obj.(float64); ok {
		if num == float64(int64(num)) {
			return fmt.Sprintf("%d", int64(num))
		}
		return fmt.Sprintf("%g", num)
	}

	return fmt.Sprintf("%v", obj)
}

func (i *Interpreter) VisitExpressionStmt(stmt *Expression) interface{} {
	i.Evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *Print) interface{} {
	value := i.Evaluate(stmt.Expression)
	fmt.Println(i.stringify(value))
	return value
}

func (i *Interpreter) VisitVariableExpr(expr *Variable) interface{} {
	return i.lookupVariable(expr.Name, expr)
}

func (i *Interpreter) VisitAssignExpr(expr *Assign) interface{} {
	value := i.Evaluate(expr.Value)

	if distance, ok := i.locals[expr]; ok {
		i.environment.AssignAt(distance, expr.Name, value)
	} else {
		err := i.globals.Assign(expr.Name, value)
		if err != nil {
			panic(err)
		}
	}

	return value
}

func (i *Interpreter) VisitBlockStmt(stmt *Block) interface{} {
	return i.executeBlock(stmt.Statements, NewEnvironment(i.environment))
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) interface{} {
	fmt.Fprintf(os.Stderr, "Entering block. Env: %p, Parent: %p, Keys: %v\n", i.environment, i.environment.enclosing, i.environment.values)
	previous := i.environment
	i.environment = environment

	defer func() {
		i.environment = previous
	}()

	for _, statement := range statements {
		i.Execute(statement)
	}
	fmt.Fprintf(os.Stderr, "Leaving block. Env: %p\n", i.environment)

	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *If) interface{} {
	if i.isTruthy(i.Evaluate(stmt.Condition)) {
		return i.Execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.Execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitLogicalExpr(expr *Logical) interface{} {
	left := i.Evaluate(expr.Left)
	if expr.Operator.Type == OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}
	return i.Evaluate(expr.Right)
}

func (i *Interpreter) VisitWhileStmt(stmt *While) interface{} {
	for i.isTruthy(i.Evaluate(stmt.Condition)) {
		result := i.Execute(stmt.Body)
		if _, ok := result.(ReturnValue); ok {
			return result
		}
	}
	return nil
}

func (i *Interpreter) VisitCallExpr(expr *Call) interface{} {
	callee := i.Evaluate(expr.Callee)

	arguments := make([]interface{}, 0, len(expr.Arguments))
	for _, argument := range expr.Arguments {
		arguments = append(arguments, i.Evaluate(argument))
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		panic(&RuntimeError{
			token:   expr.Paren,
			message: "Can only call functions and classes.",
		})
	}

	if len(arguments) != function.Arity() {
		panic(&RuntimeError{
			token:   expr.Paren,
			message: fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments)),
		})
	}
	return function.Call(i, arguments)
}

func (i *Interpreter) VisitFunctionStmt(stmt *Function) interface{} {
	function := &LoxFunction{
		declaration:   stmt,
		closure:       i.environment,
		isInitializer: false,
	}
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ReturnStmt) interface{} {
	var value interface{}
	if stmt.Value != nil {
		value = i.Evaluate(stmt.Value)
	}
	panic(&ReturnValue{Value: value})
}

func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) lookupVariable(name Token, expr Expr) interface{} {
	fmt.Fprintf(os.Stderr, "Looking up %s in env %p (parent %p), keys: %v\n", name.Lexeme, i.environment, i.environment.enclosing, i.environment.values)

	if distance, ok := i.locals[expr]; ok {
		return i.environment.GetAt(distance, name.Lexeme)
	}
	val, err := i.globals.Get(name.Lexeme)
	if err != nil {
		panic(&RuntimeError{token: name, message: err.Error()})
	}
	return val
}

func (i *Interpreter) VisitResolverStmt(stmt *Resolver) interface{} {
	return nil
}

func (i *Interpreter) VisitFunctionExpr(expr *FunctionExpr) interface{} {
	function := &Function{
		Name:   expr.Name,
		Params: expr.Params,
		Body:   expr.Body,
	}
	return NewLoxFunction(function, i.environment, false)
}

func (i *Interpreter) VisitClassStmt(stmt *Class) interface{} {
	fmt.Fprintf(os.Stderr, "Defining class %s in env %p (parent %p), keys: %v\n", stmt.Name.Lexeme, i.environment, i.environment.enclosing, i.environment.values)

	var superclass *LoxClass = nil
	if stmt.Superclass != nil {
		value := i.Evaluate(stmt.Superclass)
		var ok bool
		superclass, ok = value.(*LoxClass)
		if !ok {
			panic(&RuntimeError{stmt.Name, "Superclass must be a class."})
		}
	}

	i.environment.Define(stmt.Name.Lexeme, nil)

	if stmt.Superclass != nil {
		i.environment = NewEnvironment(i.environment)
		i.environment.Define("super", superclass)
	}

	methods := make(map[string]*LoxFunction)
	for _, method := range stmt.Methods {
		function := method.(*Function)
		isInitializer := function.Name.Lexeme == "init"
		methods[function.Name.Lexeme] = &LoxFunction{
			declaration:   function,
			closure:       i.environment,
			isInitializer: isInitializer,
		}
	}

	class := &LoxClass{
		name:       stmt.Name.Lexeme,
		superclass: superclass,
		methods:    methods,
	}

	if stmt.Superclass != nil {
		i.environment = i.environment.enclosing
	}

	i.environment.Assign(stmt.Name, class)
	return nil
}

func (i *Interpreter) VisitGetExpr(expr *Get) interface{} {
	object := i.Evaluate(expr.Object)

	if instance, ok := object.(*LoxInstance); ok {
		return instance.Get(expr.Name)
	}
	panic(RuntimeError{expr.Name, "Only instances have properties."})
}

func (i *Interpreter) VisitSetExpr(expr *Set) interface{} {
	object := i.Evaluate(expr.Object)

	if instance, ok := object.(*LoxInstance); ok {
		value := i.Evaluate(expr.Value)
		instance.Set(expr.Name, value)
		return value
	}
	panic(&RuntimeError{
		token:   expr.Name,
		message: "Only instances have fields.",
	})
}

func (i *Interpreter) VisitThisExpr(expr *This) interface{} {
	return i.lookupVariable(expr.Keyword, expr)
}

func (i *Interpreter) VisitSuperExpr(expr *Super) interface{} {
	distance := i.locals[expr]
	superclass := i.environment.GetAt(distance, "super").(*LoxClass)
	object := i.environment.GetAt(distance-1, "this").(*LoxInstance)
	method := superclass.FindMethod(expr.Method.Lexeme)
	if method == nil {
		panic(&RuntimeError{
			token:   expr.Method,
			message: fmt.Sprintf("Undefined property '%s'.", expr.Method.Lexeme),
		})
	}
	return method.Bind(object)
}
