package main

import (
	"fmt"
)

type Parser struct {
	tokens           []Token
	current          int
	currentClassName *Token
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(CLASS) {
		return p.classDeclaration()
	}
	if p.match(FUN) {
		return p.function("function")
	}
	if p.match(VAR) {
		return p.varDeclaration()
	}
	stmt, err := p.statement()
	if err != nil {
		p.synchronize()
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "expect variable name")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(EQUAL) {
		init, err := p.expression()
		if err != nil {
			return nil, err
		}
		initializer = init
	}

	_, err = p.consume(SEMICOLON, "expect ';' after variable declaration")
	if err != nil {
		return nil, err
	}

	return &Var{
		Name:        *name,
		Initializer: initializer,
	}, nil

}

func (p *Parser) consume(tokenType TokenType, message string) (*Token, error) {
	if p.check(tokenType) {
		token := p.advance()
		return &token, nil
	}
	return nil, fmt.Errorf(message)
}

func (p *Parser) parse() ([]Stmt, error) {
	var statements []Stmt
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

func (p *Parser) parseExpression() (Expr, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(LEFT_BRACE) {
		return p.block()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) ifStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *Parser) printStatement() (Stmt, error) {

	if p.match(SEMICOLON) {
		return nil, fmt.Errorf("expect expression after 'print'")
	}

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "expect ';' after value")
	if err != nil {
		return nil, err
	}

	return &Print{Expression: expr}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "expect ';' after expression")
	if err != nil {
		return nil, err
	}

	return &Expression{Expression: expr}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &Logical{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &Logical{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(EQUAL_EQUAL, BANG_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(PLUS, MINUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(STAR, SLASH) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Unary{
			Operator: operator,
			Right:    right,
		}, nil
	}
	return p.call()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(TRUE) {
		return &Literal{Value: true}, nil
	}
	if p.match(FALSE) {
		return &Literal{Value: false}, nil
	}
	if p.match(NIL) {
		return &Literal{Value: nil}, nil
	}
	if p.match(NUMBER, STRING) {
		return &Literal{Value: p.previous().Literal}, nil
	}
	if p.match(IDENTIFIER) {
		token := p.previous()
		if p.currentClassName != nil && token.Lexeme == p.currentClassName.Lexeme {
			token = *p.currentClassName
		}
		return &Variable{Name: token}, nil
	}
	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &Grouping{Expression: expr}, nil
	}
	if p.match(THIS) {
		return &This{Keyword: p.previous()}, nil
	}
	if p.match(SUPER) {
		keyword := p.previous()
		_, err := p.consume(DOT, "Expect '.' after 'super'.")
		if err != nil {
			return nil, err
		}
		method, err := p.consume(IDENTIFIER, "Expect superclass method name.")
		if err != nil {
			return nil, err
		}
		return &Super{Keyword: keyword, Method: *method}, nil
	}

	return nil, fmt.Errorf("Expect expression.")
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.match(EQUAL) {
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if get, ok := expr.(*Get); ok {
			return &Set{
				Object: get.Object,
				Name:   get.Name,
				Value:  value,
			}, nil
		} else if v, ok := expr.(*Variable); ok {
			return &Assign{
				Name:  v.Name,
				Value: value,
			}, nil
		}
		return nil, fmt.Errorf("invalid assignment target")
	}
	return expr, nil
}

func (p *Parser) block() (Stmt, error) {
	var statements []Stmt

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		decl, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, decl)
	}
	_, err := p.consume(RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}
	return &Block{Statements: statements}, nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &While{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) forStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}
	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = &Block{
			Statements: []Stmt{
				body,
				&Expression{Expression: increment},
			},
		}
	}
	if condition == nil {
		condition = &Literal{Value: true}
	}
	body = &While{
		Condition: condition,
		Body:      body,
	}
	if initializer != nil {
		body = &Block{
			Statements: []Stmt{
				initializer,
				body,
			},
		}
	}
	return body, nil

}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(DOT) {
			name, err := p.consume(IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = &Get{
				Object: expr,
				Name:   *name,
			}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	var arguments []Expr

	if !p.check(RIGHT_PAREN) {
		for {

			if len(arguments) >= 255 {
				return nil, fmt.Errorf("Can't have more than 255 arguments.")
			}

			expr, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, expr)
			if !p.match(COMMA) {
				break
			}
		}
	}
	paren, err := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &Call{
		Callee:    callee,
		Paren:     *paren,
		Arguments: arguments,
	}, nil

}

func (p *Parser) function(kind string) (*Function, error) {
	var name Token
	var err error

	if p.check(IDENTIFIER) {
		nameToken, err := p.consume(IDENTIFIER, "Expect "+kind+" name.")
		if err != nil {
			return nil, err
		}
		name = *nameToken
	}

	_, err = p.consume(LEFT_PAREN, "Expect '(' after "+kind+" name.")
	if err != nil {
		return nil, err
	}

	parameters := make([]Token, 0)
	if !p.check(RIGHT_PAREN) {
		for {
			param, err := p.consume(IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, *param)

			if !p.match(COMMA) {
				break
			}
		}
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_BRACE, "Expect '{' before "+kind+" body.")
	if err != nil {
		return nil, err
	}

	blockStmt, err := p.block()
	if err != nil {
		return nil, err
	}

	return &Function{
		Name:   name,
		Params: parameters,
		Body:   blockStmt.(*Block).Statements,
	}, nil
}

func (p *Parser) returnStatement() (Stmt, error) {
	keyword := p.previous()
	var value Expr

	if !p.check(SEMICOLON) {
		if p.match(FUN) {
			funcStmt, err := p.function("function")
			if err != nil {
				return nil, err
			}
			value = &FunctionExpr{
				Name:   funcStmt.Name,
				Params: funcStmt.Params,
				Body:   funcStmt.Body,
			}
		} else {
			var err error
			value, err = p.expression()
			if err != nil {
				return nil, err
			}
		}
	}

	_, err := p.consume(SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}

	return &ReturnStmt{Keyword: keyword, Value: value}, nil
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) classDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect class name.")
	if err != nil {
		return nil, err
	}

	previousClass := p.currentClassName
	p.currentClassName = name

	var superclass Expr = nil
	if p.match(LESS) {
		p.consume(IDENTIFIER, "Expect superclass name.")
		superclass = &Variable{
			Name: p.previous(),
		}
	}

	_, err = p.consume(LEFT_BRACE, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}

	methods := make([]Stmt, 0)
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		method, err := p.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}

	_, err = p.consume(RIGHT_BRACE, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}

	p.currentClassName = previousClass

	return &Class{
		Name:       *name,
		Superclass: superclass,
		Methods:    methods,
	}, nil
}
