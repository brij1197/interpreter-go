package main

type ReturnValue struct {
	Value interface{}
}

type ReturnStmt struct {
	Keyword Token
	Value   Expr
}

func (r *ReturnValue) Error() string {
	return "return"
}
