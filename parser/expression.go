package parser

import (
	"errors"
	"strconv"

	"es6-interpreter/ast"
	"es6-interpreter/lexer"
)

func (p *Parser) registerPrefixFns() {
	p.registerPrefix(lexer.Identifier, p.parseIdentifier)
	p.registerPrefix(lexer.Number, p.parseNumberLiteral)
	p.registerPrefix(lexer.String, p.parseStringLiteral)
	p.registerPrefix(lexer.TrueLiteral, p.parseBooleanLiteral)
	p.registerPrefix(lexer.FalseLiteral, p.parseBooleanLiteral)
	p.registerPrefix(lexer.NullLiteral, p.parseNullLiteral)
	p.registerPrefix(lexer.LParen, p.parseGroupedExpression)
	p.registerPrefix(lexer.LogicalNot, p.parsePrefixExpression)
	p.registerPrefix(lexer.BitwiseNot, p.parsePrefixExpression)
	p.registerPrefix(lexer.Minus, p.parsePrefixExpression)
	p.registerPrefix(lexer.Plus, p.parsePrefixExpression)
	p.registerPrefix(lexer.Increment, p.parsePrefixExpression)
	p.registerPrefix(lexer.Decrement, p.parsePrefixExpression)
	p.registerPrefix(lexer.KeywordThis, p.parseThisExpression)
	p.registerPrefix(lexer.KeywordSuper, p.parseSuperExpression)
}

func (p *Parser) registerInfixFns() {
	p.registerInfix(lexer.Plus, p.parseInfixExpression)
	p.registerInfix(lexer.Minus, p.parseInfixExpression)
	p.registerInfix(lexer.Multiply, p.parseInfixExpression)
	p.registerInfix(lexer.Divide, p.parseInfixExpression)
	p.registerInfix(lexer.Assign, p.parseAssignmentExpression)
	p.registerInfix(lexer.PlusAssign, p.parseAssignmentExpression)
	p.registerInfix(lexer.MinusAssign, p.parseAssignmentExpression)
	p.registerInfix(lexer.MultiplyAssign, p.parseAssignmentExpression)
	p.registerInfix(lexer.DivideAssign, p.parseAssignmentExpression)
	p.registerInfix(lexer.ModuloAssign, p.parseAssignmentExpression)
	p.registerInfix(lexer.ShiftLeftAssign, p.parseAssignmentExpression)
	p.registerInfix(lexer.ShiftRightAssign, p.parseAssignmentExpression)
	p.registerInfix(lexer.UnsignedShiftAssign, p.parseAssignmentExpression)
	p.registerInfix(lexer.BitwiseAndAssign, p.parseAssignmentExpression)
	p.registerInfix(lexer.BitwiseOrAssign, p.parseAssignmentExpression)
	p.registerInfix(lexer.BitwiseXorAssign, p.parseAssignmentExpression)
	p.registerInfix(lexer.LParen, p.parseCallExpression)
	p.registerInfix(lexer.Dot, p.parseMemberExpression)
	p.registerInfix(lexer.LBracket, p.parseComputedMemberExpression)
	p.registerInfix(lexer.Increment, p.parsePostfixExpression)
	p.registerInfix(lexer.Decrement, p.parsePostfixExpression)
	p.registerInfix(lexer.LogicalAnd, p.parseLogicalExpression)
	p.registerInfix(lexer.LogicalOr, p.parseLogicalExpression)
	p.registerInfix(lexer.Equal, p.parseInfixExpression)
	p.registerInfix(lexer.NotEqual, p.parseInfixExpression)
	p.registerInfix(lexer.StrictEqual, p.parseInfixExpression)
	p.registerInfix(lexer.StrictNotEqual, p.parseInfixExpression)
	p.registerInfix(lexer.LessThan, p.parseInfixExpression)
	p.registerInfix(lexer.LessEqual, p.parseInfixExpression)
	p.registerInfix(lexer.GreaterThan, p.parseInfixExpression)
	p.registerInfix(lexer.GreaterEqual, p.parseInfixExpression)
	p.registerInfix(lexer.KeywordIn, p.parseInfixExpression)
	p.registerInfix(lexer.KeywordInstanceof, p.parseInfixExpression)
	p.registerInfix(lexer.ShiftLeft, p.parseInfixExpression)
	p.registerInfix(lexer.ShiftRight, p.parseInfixExpression)
	p.registerInfix(lexer.UnsignedShiftRight, p.parseInfixExpression)
	p.registerInfix(lexer.BitwiseAnd, p.parseInfixExpression)
	p.registerInfix(lexer.BitwiseOr, p.parseInfixExpression)
	p.registerInfix(lexer.BitwiseXor, p.parseInfixExpression)
}

func (p *Parser) registerPrefix(tt lexer.TokenType, fn prefixParseFn) {
	p.prefixFns[tt] = fn
}

func (p *Parser) registerInfix(tt lexer.TokenType, fn infixParseFn) {
	p.infixFns[tt] = fn
}

func (p *Parser) parseExpression(pre precedence) ast.Expression {
	prefix := p.prefixFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(lexer.Semicolon) && pre < p.peekPrecedence() {
		infix := p.infixFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	tok := p.curToken
	return ast.NewIdentifier(tok.Literal, p.tokenLocation(tok))
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	tok := p.curToken
	return ast.NewNumberLiteral(tok.Literal, p.tokenLocation(tok))
}

func (p *Parser) parseStringLiteral() ast.Expression {
	tok := p.curToken
	val, err := strconv.Unquote(tok.Literal)
	if err != nil {
		p.errors = append(p.errors, err)
		val = tok.Literal
	}
	return ast.NewStringLiteral(val, p.tokenLocation(tok))
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	tok := p.curToken
	value := tok.Type == lexer.TrueLiteral
	return ast.NewBooleanLiteral(value, p.tokenLocation(tok))
}

func (p *Parser) parseNullLiteral() ast.Expression {
	tok := p.curToken
	return ast.NewNullLiteral(p.tokenLocation(tok))
}

func (p *Parser) parseThisExpression() ast.Expression {
	tok := p.curToken
	return ast.NewThisExpression(p.tokenLocation(tok))
}

func (p *Parser) parseSuperExpression() ast.Expression {
	tok := p.curToken
	return ast.NewSuper(p.tokenLocation(tok))
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	start := p.curToken.Start
	p.nextToken()
	exp := p.parseExpression(lowest)
	if exp == nil {
		return nil
	}
	if !p.expectPeek(lexer.RParen) {
		return nil
	}
	loc := ast.Location{Start: convertPosition(start), End: convertPosition(p.curToken.End)}
	p.setNodeLocation(exp, loc)
	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	tok := p.curToken
	operator := tok.Literal

	p.nextToken()
	right := p.parseExpression(prefixPrec)
	if right == nil {
		return nil
	}

	loc := ast.Location{Start: convertPosition(tok.Start), End: right.Loc().End}
	switch tok.Type {
	case lexer.Increment, lexer.Decrement:
		if !isAssignable(right) {
			p.errors = append(p.errors, errors.New("invalid update target"))
			return nil
		}
		return ast.NewUpdateExpression(operator, right, true, loc)
	default:
		return ast.NewUnaryExpression(operator, right, true, loc)
	}
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	operator := p.curToken.Literal
	if !isAssignable(left) {
		p.errors = append(p.errors, errors.New("invalid update target"))
		return nil
	}
	loc := ast.Location{Start: left.Loc().Start, End: convertPosition(p.curToken.End)}
	return ast.NewUpdateExpression(operator, left, false, loc)
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	operator := p.curToken.Literal
	precedence := p.curPrecedence()

	p.nextToken()
	right := p.parseExpression(precedence)
	if right == nil {
		return nil
	}

	loc := ast.Location{Start: left.Loc().Start, End: right.Loc().End}
	return ast.NewBinaryExpression(operator, left, right, loc)
}

func (p *Parser) parseLogicalExpression(left ast.Expression) ast.Expression {
	operator := p.curToken.Literal
	precedence := p.curPrecedence()

	p.nextToken()
	right := p.parseExpression(precedence)
	if right == nil {
		return nil
	}

	loc := ast.Location{Start: left.Loc().Start, End: right.Loc().End}
	return ast.NewLogicalExpression(operator, left, right, loc)
}

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	if !isAssignable(left) {
		p.errors = append(p.errors, errors.New("invalid assignment target"))
		return nil
	}

	operator := p.curToken.Literal
	precedence := p.curPrecedence()

	p.nextToken()
	right := p.parseExpression(precedence - 1)
	if right == nil {
		return nil
	}

	loc := ast.Location{Start: left.Loc().Start, End: right.Loc().End}
	return ast.NewAssignmentExpression(operator, left, right, loc)
}

func (p *Parser) parseCallExpression(callee ast.Expression) ast.Expression {
	start := callee.Loc().Start
	p.nextToken()
	var args []ast.Expression
	if !p.curTokenIs(lexer.RParen) {
		for {
			arg := p.parseExpression(lowest)
			if arg == nil {
				return nil
			}
			args = append(args, arg)

			if !p.peekTokenIs(lexer.Comma) {
				break
			}
			p.nextToken() // move to comma
			p.nextToken() // move to next argument
		}
		if !p.expectPeek(lexer.RParen) {
			p.errors = append(p.errors, errors.New("unterminated call expression"))
			return nil
		}
	}
	end := convertPosition(p.curToken.End)
	loc := ast.Location{Start: start, End: end}
	return ast.NewCallExpression(callee, args, loc)
}

func (p *Parser) parseMemberExpression(object ast.Expression) ast.Expression {
	start := object.Loc().Start
	if !p.expectPeek(lexer.Identifier) {
		return nil
	}
	property := ast.NewIdentifier(p.curToken.Literal, p.tokenLocation(p.curToken))
	loc := ast.Location{Start: start, End: property.Loc().End}
	return ast.NewMemberExpression(object, property, false, loc)
}

func (p *Parser) parseComputedMemberExpression(object ast.Expression) ast.Expression {
	start := object.Loc().Start
	p.nextToken()
	property := p.parseExpression(lowest)
	if property == nil {
		return nil
	}
	if !p.expectPeek(lexer.RBracket) {
		p.errors = append(p.errors, errors.New("unterminated computed member expression"))
		return nil
	}
	loc := ast.Location{Start: start, End: convertPosition(p.curToken.End)}
	return ast.NewMemberExpression(object, property, true, loc)
}

func (p *Parser) noPrefixParseFnError(tt lexer.TokenType) {
	msg := "no prefix parse function for " + string(tt)
	p.errors = append(p.errors, errors.New(msg))
}

func (p *Parser) setNodeLocation(node ast.Node, loc ast.Location) {
	if loc.IsValid() {
		switch n := node.(type) {
		case interface{ SetLoc(ast.Location) }:
			n.SetLoc(loc)
		}
	}
}

func isAssignable(expr ast.Expression) bool {
	switch expr.(type) {
	case *ast.Identifier, *ast.MemberExpression:
		return true
	default:
		return false
	}
}
