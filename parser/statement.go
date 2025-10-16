package parser

import (
	"errors"

	"es6-interpreter/ast"
	"es6-interpreter/lexer"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.KeywordVar, lexer.KeywordLet, lexer.KeywordConst:
		return p.parseVariableStatement()
	case lexer.Semicolon:
		return p.parseEmptyStatement()
	case lexer.LBrace:
		return p.parseBlockStatement()
	case lexer.KeywordReturn:
		return p.parseReturnStatement()
	case lexer.KeywordIf:
		return p.parseIfStatement()
	case lexer.KeywordWhile:
		return p.parseWhileStatement()
	case lexer.KeywordDo:
		return p.parseDoWhileStatement()
	case lexer.KeywordFor:
		return p.parseForStatement()
	case lexer.KeywordBreak:
		return p.parseBreakStatement()
	case lexer.KeywordContinue:
		return p.parseContinueStatement()
	case lexer.KeywordThrow:
		return p.parseThrowStatement()
	case lexer.KeywordDebugger:
		return p.parseDebuggerStatement()
	case lexer.KeywordSwitch:
		return p.parseSwitchStatement()
	case lexer.KeywordWith:
		return p.parseWithStatement()
	case lexer.Identifier:
		if p.peekTokenIs(lexer.Colon) {
			return p.parseLabeledStatement()
		}
		return p.parseExpressionStatement()
	case lexer.KeywordTry:
		return p.parseTryStatement()
	case lexer.KeywordFunction:
		return p.parseFunctionDeclaration()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseEmptyStatement() ast.Statement {
	loc := p.tokenLocation(p.curToken)
	return ast.NewEmptyStatement(loc)
}

func (p *Parser) parseBlockStatement() ast.Statement {
	start := p.curToken.Start

	// Move inside the block body.
	p.nextToken()

	var body []ast.Statement
	for !p.curTokenIs(lexer.RBrace) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
		p.nextToken()
	}

	if !p.curTokenIs(lexer.RBrace) {
		p.errors = append(p.errors, errors.New("unterminated block statement"))
		return nil
	}

	loc := p.locFrom(start, p.curToken.End)
	return ast.NewBlockStatement(body, loc)
}

func (p *Parser) parseReturnStatement() ast.Statement {
	start := p.curToken.Start

	// No argument if the next token is a semicolon, closing brace, or EOF.
	switch {
	case p.peekTokenIs(lexer.Semicolon):
		p.nextToken()
		loc := p.locFrom(start, p.curToken.End)
		return ast.NewReturnStatement(nil, loc)
	case p.peekTokenIs(lexer.RBrace) || p.peekTokenIs(lexer.EOF):
		loc := p.locFrom(start, p.curToken.End)
		return ast.NewReturnStatement(nil, loc)
	}

	// Parse return argument expression.
	p.nextToken()
	argument := p.parseExpression(lowest)
	if argument == nil {
		return nil
	}

	end := argument.Loc().End
	if p.peekTokenIs(lexer.Semicolon) {
		p.nextToken()
		end = convertPosition(p.curToken.End)
	}

	loc := ast.Location{Start: convertPosition(start), End: end}
	return ast.NewReturnStatement(argument, loc)
}

func (p *Parser) parseIfStatement() ast.Statement {
	start := p.curToken.Start

	if !p.expectPeek(lexer.LParen) {
		return nil
	}

	p.nextToken()
	test := p.parseExpression(lowest)
	if test == nil {
		return nil
	}

	if !p.expectPeek(lexer.RParen) {
		return nil
	}

	// Move to the first token of the consequent statement.
	p.nextToken()
	consequent := p.parseStatement()
	if consequent == nil {
		return nil
	}

	altLoc := consequent.Loc().End
	var alternate ast.Statement
	if p.peekTokenIs(lexer.KeywordElse) {
		p.nextToken() // move to 'else'
		p.nextToken() // move to alternate statement start
		alternate = p.parseStatement()
		if alternate == nil {
			return nil
		}
		altLoc = alternate.Loc().End
	}

	loc := ast.Location{Start: convertPosition(start), End: altLoc}
	return ast.NewIfStatement(test, consequent, alternate, loc)
}

func (p *Parser) parseWhileStatement() ast.Statement {
	start := p.curToken.Start

	if !p.expectPeek(lexer.LParen) {
		return nil
	}

	p.nextToken()
	test := p.parseExpression(lowest)
	if test == nil {
		return nil
	}

	if !p.expectPeek(lexer.RParen) {
		return nil
	}

	p.nextToken()
	body := p.parseStatement()
	if body == nil {
		return nil
	}

	loc := ast.Location{Start: convertPosition(start), End: body.Loc().End}
	return ast.NewWhileStatement(test, body, loc)
}

func (p *Parser) parseDoWhileStatement() ast.Statement {
	start := p.curToken.Start

	p.nextToken()
	body := p.parseStatement()
	if body == nil {
		return nil
	}

	if !p.expectPeek(lexer.KeywordWhile) {
		return nil
	}

	if !p.expectPeek(lexer.LParen) {
		return nil
	}

	p.nextToken()
	test := p.parseExpression(lowest)
	if test == nil {
		return nil
	}

	if !p.expectPeek(lexer.RParen) {
		return nil
	}

	end := convertPosition(p.curToken.End)
	if p.peekTokenIs(lexer.Semicolon) {
		p.nextToken()
		end = convertPosition(p.curToken.End)
	}

	loc := ast.Location{Start: convertPosition(start), End: end}
	return ast.NewDoWhileStatement(body, test, loc)
}

func (p *Parser) parseBreakStatement() ast.Statement {
	start := p.curToken.Start
	end := p.curToken.End

	var label *ast.Identifier
	if p.peekTokenIs(lexer.Identifier) && p.peekToken.Start.Line == p.curToken.End.Line {
		p.nextToken()
		tok := p.curToken
		label = ast.NewIdentifier(tok.Literal, p.tokenLocation(tok))
		end = tok.End
	}

	if p.peekTokenIs(lexer.Semicolon) {
		p.nextToken()
		end = p.curToken.End
	}

	loc := p.locFrom(start, end)
	return ast.NewBreakStatement(label, loc)
}

func (p *Parser) parseContinueStatement() ast.Statement {
	start := p.curToken.Start
	end := p.curToken.End

	var label *ast.Identifier
	if p.peekTokenIs(lexer.Identifier) && p.peekToken.Start.Line == p.curToken.End.Line {
		p.nextToken()
		tok := p.curToken
		label = ast.NewIdentifier(tok.Literal, p.tokenLocation(tok))
		end = tok.End
	}

	if p.peekTokenIs(lexer.Semicolon) {
		p.nextToken()
		end = p.curToken.End
	}

	loc := p.locFrom(start, end)
	return ast.NewContinueStatement(label, loc)
}

func (p *Parser) parseThrowStatement() ast.Statement {
	start := p.curToken.Start

	if p.peekToken.Start.Line != p.curToken.End.Line {
		p.errors = append(p.errors, errors.New("illegal newline after throw"))
		return nil
	}

	p.nextToken()
	argument := p.parseExpression(lowest)
	if argument == nil {
		return nil
	}

	end := argument.Loc().End
	if p.peekTokenIs(lexer.Semicolon) {
		p.nextToken()
		end = convertPosition(p.curToken.End)
	}

	loc := ast.Location{Start: convertPosition(start), End: end}
	return ast.NewThrowStatement(argument, loc)
}

func (p *Parser) parseDebuggerStatement() ast.Statement {
	start := p.curToken.Start
	end := p.curToken.End

	if p.peekTokenIs(lexer.Semicolon) {
		p.nextToken()
		end = p.curToken.End
	}

	loc := p.locFrom(start, end)
	return ast.NewDebuggerStatement(loc)
}

func (p *Parser) parseSwitchStatement() ast.Statement {
	start := p.curToken.Start

	if !p.expectPeek(lexer.LParen) {
		return nil
	}

	p.nextToken()
	discriminant := p.parseExpression(lowest)
	if discriminant == nil {
		return nil
	}

	if !p.expectPeek(lexer.RParen) {
		return nil
	}

	if !p.expectPeek(lexer.LBrace) {
		return nil
	}

	// move inside switch body
	p.nextToken()

	var cases []*ast.SwitchCase
	seenDefault := false

	for !p.curTokenIs(lexer.RBrace) && !p.curTokenIs(lexer.EOF) {
		caseStart := p.curToken.Start
		var test ast.Expression

		switch p.curToken.Type {
		case lexer.KeywordCase:
			p.nextToken()
			test = p.parseExpression(lowest)
			if test == nil {
				return nil
			}
			if !p.expectPeek(lexer.Colon) {
				return nil
			}
		case lexer.KeywordDefault:
			if seenDefault {
				p.errors = append(p.errors, errors.New("multiple default clauses in switch"))
				return nil
			}
			seenDefault = true
			if !p.expectPeek(lexer.Colon) {
				return nil
			}
		default:
			p.errors = append(p.errors, errors.New("expected case or default clause"))
			return nil
		}

		end := p.curToken.End // colon end

		// move to first statement in clause (if any)
		p.nextToken()

		var consequent []ast.Statement
		for !p.curTokenIs(lexer.KeywordCase) && !p.curTokenIs(lexer.KeywordDefault) && !p.curTokenIs(lexer.RBrace) && !p.curTokenIs(lexer.EOF) {
			stmt := p.parseStatement()
			if stmt == nil {
				return nil
			}
			consequent = append(consequent, stmt)
			loc := stmt.Loc().End
			end = lexer.Position{Offset: loc.Offset, Line: loc.Line, Column: loc.Column}
			p.nextToken()
		}

		loc := p.locFrom(caseStart, end)
		cases = append(cases, ast.NewSwitchCase(test, consequent, loc))
	}

	if !p.curTokenIs(lexer.RBrace) {
		p.errors = append(p.errors, errors.New("unterminated switch statement"))
		return nil
	}

	loc := p.locFrom(start, p.curToken.End)
	return ast.NewSwitchStatement(discriminant, cases, loc)
}

func (p *Parser) parseWithStatement() ast.Statement {
	start := p.curToken.Start

	if !p.expectPeek(lexer.LParen) {
		return nil
	}

	p.nextToken()
	object := p.parseExpression(lowest)
	if object == nil {
		return nil
	}

	if !p.expectPeek(lexer.RParen) {
		return nil
	}

	p.nextToken()
	body := p.parseStatement()
	if body == nil {
		return nil
	}

	loc := ast.Location{Start: convertPosition(start), End: body.Loc().End}
	return ast.NewWithStatement(object, body, loc)
}

func (p *Parser) parseLabeledStatement() ast.Statement {
	start := p.curToken.Start
	labelTok := p.curToken
	label := ast.NewIdentifier(labelTok.Literal, p.tokenLocation(labelTok))

	if !p.expectPeek(lexer.Colon) {
		return nil
	}

	p.nextToken()
	body := p.parseStatement()
	if body == nil {
		return nil
	}

	loc := ast.Location{Start: convertPosition(start), End: body.Loc().End}
	return ast.NewLabeledStatement(label, body, loc)
}

func (p *Parser) parseTryStatement() ast.Statement {
	start := p.curToken.Start

	if !p.expectPeek(lexer.LBrace) {
		return nil
	}

	blockStmt := p.parseBlockStatement()
	if blockStmt == nil {
		return nil
	}

	tryBlock, ok := blockStmt.(*ast.BlockStatement)
	if !ok {
		p.errors = append(p.errors, errors.New("try block did not produce BlockStatement"))
		return nil
	}

	end := p.curToken.End

	var handler *ast.CatchClause
	var finalizer *ast.BlockStatement

	if p.peekTokenIs(lexer.KeywordCatch) {
		p.nextToken()
		handler = p.parseCatchClause()
		if handler == nil {
			return nil
		}
		end = p.curToken.End
	}

	if p.peekTokenIs(lexer.KeywordFinally) {
		p.nextToken()
		if !p.expectPeek(lexer.LBrace) {
			return nil
		}
		finalizerStmt := p.parseBlockStatement()
		if finalizerStmt == nil {
			return nil
		}
		var ok bool
		finalizer, ok = finalizerStmt.(*ast.BlockStatement)
		if !ok {
			p.errors = append(p.errors, errors.New("finally block did not produce BlockStatement"))
			return nil
		}
		end = p.curToken.End
	}

	if handler == nil && finalizer == nil {
		p.errors = append(p.errors, errors.New("try statement requires catch or finally"))
		return nil
	}

	loc := p.locFrom(start, end)
	return ast.NewTryStatement(tryBlock, handler, finalizer, loc)
}

func (p *Parser) parseCatchClause() *ast.CatchClause {
	start := p.curToken.Start

	if !p.expectPeek(lexer.LParen) {
		return nil
	}

	p.nextToken()
	param := p.parseBindingElement(false)
	if param == nil {
		return nil
	}

	if !p.expectPeek(lexer.RParen) {
		return nil
	}

	if !p.expectPeek(lexer.LBrace) {
		return nil
	}

	bodyStmt := p.parseBlockStatement()
	if bodyStmt == nil {
		return nil
	}

	body, ok := bodyStmt.(*ast.BlockStatement)
	if !ok {
		p.errors = append(p.errors, errors.New("catch body did not produce BlockStatement"))
		return nil
	}

	loc := p.locFrom(start, p.curToken.End)
	return ast.NewCatchClause(param, body, loc)
}

func (p *Parser) parseFunctionDeclaration() ast.Statement {
	start := p.curToken.Start

	isGenerator := false
	if p.peekTokenIs(lexer.Multiply) {
		p.nextToken()
		isGenerator = true
	}

	if !p.expectPeek(lexer.Identifier) {
		return nil
	}

	nameTok := p.curToken
	id := ast.NewIdentifier(nameTok.Literal, p.tokenLocation(nameTok))

	if !p.expectPeek(lexer.LParen) {
		return nil
	}

	params, ok := p.parseFunctionParams()
	if !ok {
		return nil
	}

	if !p.expectPeek(lexer.LBrace) {
		return nil
	}

	bodyStmt := p.parseBlockStatement()
	if bodyStmt == nil {
		return nil
	}

	body, ok2 := bodyStmt.(*ast.BlockStatement)
	if !ok2 {
		p.errors = append(p.errors, errors.New("function body did not produce BlockStatement"))
		return nil
	}

	loc := p.locFrom(start, p.curToken.End)
	return ast.NewFunctionDeclaration(id, params, body, isGenerator, loc)
}

func (p *Parser) parseFunctionParams() ([]ast.Pattern, bool) {
	var params []ast.Pattern

	if p.peekTokenIs(lexer.RParen) {
		p.nextToken()
		return params, true
	}

	// move to first parameter token
	p.nextToken()

	restSeen := false
	for !p.curTokenIs(lexer.RParen) && !p.curTokenIs(lexer.EOF) {
		if restSeen {
			p.errors = append(p.errors, errors.New("parameters not allowed after rest element"))
			return nil, false
		}

		if p.curTokenIs(lexer.Ellipsis) {
			restStart := p.curToken.Start
			p.nextToken()
			arg := p.parseBindingElement(false)
			if arg == nil {
				return nil, false
			}
			rest := ast.NewRestElement(arg, p.locFrom(restStart, p.curToken.End))
			params = append(params, rest)
			restSeen = true
			if !p.expectPeek(lexer.RParen) {
				return nil, false
			}
			break
		}

		param := p.parseBindingElement(true)
		if param == nil {
			return nil, false
		}
		params = append(params, param)

		if p.peekTokenIs(lexer.Comma) {
			p.nextToken()
			if p.peekTokenIs(lexer.RParen) {
				p.errors = append(p.errors, errors.New("trailing comma without parameter"))
				return nil, false
			}
			p.nextToken()
			continue
		}

		if p.peekTokenIs(lexer.RParen) {
			p.nextToken()
			break
		}

		p.errors = append(p.errors, errors.New("unexpected token in parameter list"))
		return nil, false
	}

	return params, true
}

func (p *Parser) parseForStatement() ast.Statement {
	start := p.curToken.Start

	if !p.expectPeek(lexer.LParen) {
		return nil
	}

	p.nextToken()

	var init ast.Node
	if !p.curTokenIs(lexer.Semicolon) {
		switch p.curToken.Type {
		case lexer.KeywordVar, lexer.KeywordLet, lexer.KeywordConst:
			decl := p.parseVariableStatement()
			if decl == nil {
				return nil
			}
			init = decl
		default:
			expr := p.parseExpression(lowest)
			if expr == nil {
				return nil
			}
			init = expr
			if !p.peekTokenIs(lexer.Semicolon) {
				p.errors = append(p.errors, errors.New("expected semicolon after for-loop initializer"))
				return nil
			}
		}
	}

	if !p.curTokenIs(lexer.Semicolon) {
		if !p.expectPeek(lexer.Semicolon) {
			return nil
		}
	}

	p.nextToken()

	var test ast.Expression
	if !p.curTokenIs(lexer.Semicolon) {
		test = p.parseExpression(lowest)
		if test == nil {
			return nil
		}
	}

	if !p.curTokenIs(lexer.Semicolon) {
		if !p.expectPeek(lexer.Semicolon) {
			return nil
		}
	}

	p.nextToken()

	var update ast.Expression
	if p.curTokenIs(lexer.RParen) {
		// no update expression
	} else {
		update = p.parseExpression(lowest)
		if update == nil {
			return nil
		}
		if !p.expectPeek(lexer.RParen) {
			return nil
		}
	}

	if !p.curTokenIs(lexer.RParen) {
		p.errors = append(p.errors, errors.New("unterminated for-loop clause"))
		return nil
	}

	p.nextToken()
	body := p.parseStatement()
	if body == nil {
		return nil
	}

	loc := ast.Location{Start: convertPosition(start), End: body.Loc().End}
	return ast.NewForStatement(init, test, update, body, loc)
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	expr := p.parseExpression(lowest)
	if expr == nil {
		return nil
	}

	stmt := ast.NewExpressionStatement(expr, expr.Loc())

	if p.peekTokenIs(lexer.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseVariableStatement() ast.Statement {
	kind := ast.VarKind
	switch p.curToken.Type {
	case lexer.KeywordConst:
		kind = ast.ConstKind
	case lexer.KeywordLet:
		kind = ast.LetKind
	}

	start := p.curToken.Start

	// advance to first binding token
	p.nextToken()

	var declarators []*ast.VariableDeclarator
	for {
		if p.curToken.Type == lexer.Semicolon {
			p.errors = append(p.errors, errors.New("missing binding in variable declaration"))
			return nil
		}

		decl := p.parseVariableDeclarator()
		if decl == nil {
			return nil
		}
		declarators = append(declarators, decl)

		if !p.peekTokenIs(lexer.Comma) {
			break
		}

		p.nextToken() // consume comma
		p.nextToken() // advance to next binding token
	}

	end := p.curToken.End

	if p.peekTokenIs(lexer.Semicolon) {
		p.nextToken()
		end = p.curToken.End
	}

	return ast.NewVariableDeclaration(kind, declarators, p.locFrom(start, end))
}

func (p *Parser) parseVariableDeclarator() *ast.VariableDeclarator {
	start := p.curToken.Start

	pattern := p.parseBindingElement(false)
	if pattern == nil {
		return nil
	}

	var init ast.Expression
	if p.peekTokenIs(lexer.Assign) {
		p.nextToken() // move to '='
		p.nextToken() // advance to initializer expression
		init = p.parseExpression(lowest)
		if init == nil {
			return nil
		}
	}

	end := p.curToken.End

	return ast.NewVariableDeclarator(pattern, init, p.locFrom(start, end))
}
