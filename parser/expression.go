package parser

import (
	"errors"
	"strconv"
	"strings"

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
	p.registerPrefix(lexer.LBracket, p.parseArrayLiteral)
	p.registerPrefix(lexer.LBrace, p.parseObjectLiteral)
	p.registerPrefix(lexer.Regex, p.parseRegExpLiteral)
	p.registerPrefix(lexer.KeywordThis, p.parseThisExpression)
	p.registerPrefix(lexer.KeywordSuper, p.parseSuperExpression)
	p.registerPrefix(lexer.KeywordTypeof, p.parsePrefixExpression)
	p.registerPrefix(lexer.KeywordVoid, p.parsePrefixExpression)
	p.registerPrefix(lexer.KeywordDelete, p.parsePrefixExpression)
	p.registerPrefix(lexer.KeywordNew, p.parseNewExpression)
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
	p.registerInfix(lexer.Question, p.parseConditionalExpression)
	p.registerInfix(lexer.Comma, p.parseSequenceExpression)
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

func (p *Parser) parseNewExpression() ast.Expression {
	newTok := p.curToken
	start := newTok.Start

	// Advance to the token following `new` for further inspection.
	p.nextToken()

	// Handle meta property new.target explicitly.
	if p.curTokenIs(lexer.Dot) {
		if !p.expectPeek(lexer.Identifier) {
			return nil
		}
		identTok := p.curToken
		if identTok.Literal != "target" {
			p.errors = append(p.errors, errors.New("expected target after new"))
			return nil
		}
		meta := ast.NewIdentifier("new", p.locFrom(newTok.Start, newTok.End))
		property := ast.NewIdentifier(identTok.Literal, p.tokenLocation(identTok))
		loc := p.locFrom(newTok.Start, identTok.End)
		return ast.NewMetaProperty(meta, property, loc)
	}

	expr := p.parseExpression(postfixPrec)
	if expr == nil {
		return nil
	}

	return p.wrapNewExpression(expr, start)
}

func (p *Parser) parseConditionalExpression(test ast.Expression) ast.Expression {
	start := test.Loc().Start

	// parse consequent expression after '?'
	p.nextToken()
	consequent := p.parseExpression(lowest)
	if consequent == nil {
		return nil
	}

	if !p.expectPeek(lexer.Colon) {
		return nil
	}

	p.nextToken()
	alternate := p.parseExpression(conditionalPrec - 1)
	if alternate == nil {
		return nil
	}

	loc := ast.Location{Start: start, End: alternate.Loc().End}
	return ast.NewConditionalExpression(test, consequent, alternate, loc)
}

func (p *Parser) parseCallExpression(callee ast.Expression) ast.Expression {
	start := callee.Loc().Start
	p.nextToken()
	var args []ast.Expression
	if !p.curTokenIs(lexer.RParen) {
		for {
			arg := p.parseExpression(sequencePrec)
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

func (p *Parser) parseSequenceExpression(left ast.Expression) ast.Expression {
	start := left.Loc().Start
	expressions := []ast.Expression{left}

	for {
		p.nextToken()
		expr := p.parseExpression(sequencePrec - 1)
		if expr == nil {
			return nil
		}
		if nested, ok := expr.(*ast.SequenceExpression); ok {
			expressions = append(expressions, nested.Expressions...)
		} else {
			expressions = append(expressions, expr)
		}

		if !p.peekTokenIs(lexer.Comma) {
			break
		}

		p.nextToken()
	}

	last := expressions[len(expressions)-1]
	loc := ast.Location{Start: start, End: last.Loc().End}
	return ast.NewSequenceExpression(expressions, loc)
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	start := p.curToken.Start
	var elements []ast.Expression

	if p.peekTokenIs(lexer.RBracket) {
		p.nextToken()
	} else {
		p.nextToken()
		for !p.curTokenIs(lexer.RBracket) && !p.curTokenIs(lexer.EOF) {
			if p.curTokenIs(lexer.Comma) {
				elements = append(elements, nil)
				p.nextToken()
				continue
			}

			var element ast.Expression
			if p.curTokenIs(lexer.Ellipsis) {
				spreadStart := p.curToken.Start
				p.nextToken()
				arg := p.parseExpression(sequencePrec)
				if arg == nil {
					return nil
				}
				element = ast.NewSpreadElement(arg, p.locFrom(spreadStart, p.curToken.End))
			} else {
				element = p.parseExpression(sequencePrec)
				if element == nil {
					return nil
				}
			}

			elements = append(elements, element)

			if p.peekTokenIs(lexer.Comma) {
				p.nextToken()
				if p.peekTokenIs(lexer.RBracket) {
					elements = append(elements, nil)
					p.nextToken()
					break
				}
				p.nextToken()
				continue
			}

			p.nextToken()
		}
	}

	if !p.curTokenIs(lexer.RBracket) {
		p.errors = append(p.errors, errors.New("unterminated array literal"))
		return nil
	}

	loc := p.locFrom(start, p.curToken.End)
	return ast.NewArrayLiteral(elements, loc)
}

func (p *Parser) parseObjectLiteral() ast.Expression {
	start := p.curToken.Start
	var properties []ast.Property

	if p.peekTokenIs(lexer.RBrace) {
		p.nextToken()
	} else {
		p.nextToken()
		for !p.curTokenIs(lexer.RBrace) && !p.curTokenIs(lexer.EOF) {
			if p.curTokenIs(lexer.Ellipsis) {
				spreadStart := p.curToken.Start
				p.nextToken()
				arg := p.parseExpression(sequencePrec)
				if arg == nil {
					return nil
				}
				properties = append(properties, ast.NewSpreadElement(arg, p.locFrom(spreadStart, p.curToken.End)))
			} else {
				prop := p.parseObjectProperty()
				if prop == nil {
					return nil
				}
				properties = append(properties, prop)
			}

			if p.peekTokenIs(lexer.Comma) {
				p.nextToken()
				if p.peekTokenIs(lexer.RBrace) {
					p.nextToken()
					break
				}
				p.nextToken()
				continue
			}

			p.nextToken()
		}
	}

	if !p.curTokenIs(lexer.RBrace) {
		p.errors = append(p.errors, errors.New("unterminated object literal"))
		return nil
	}

	loc := p.locFrom(start, p.curToken.End)
	return ast.NewObjectLiteral(properties, loc)
}

func (p *Parser) parseObjectProperty() ast.Property {
	start := p.curToken.Start

	if p.curTokenIs(lexer.Ellipsis) {
		spreadStart := p.curToken.Start
		p.nextToken()
		arg := p.parseExpression(lowest)
		if arg == nil {
			return nil
		}
		return ast.NewSpreadElement(arg, p.locFrom(spreadStart, p.curToken.End))
	}

	computed := false
	var key ast.Expression

	switch p.curToken.Type {
	case lexer.Identifier:
		key = ast.NewIdentifier(p.curToken.Literal, p.tokenLocation(p.curToken))
	case lexer.String:
		val, err := strconv.Unquote(p.curToken.Literal)
		if err != nil {
			p.errors = append(p.errors, err)
			val = p.curToken.Literal
		}
		key = ast.NewStringLiteral(val, p.tokenLocation(p.curToken))
	case lexer.Number:
		key = ast.NewNumberLiteral(p.curToken.Literal, p.tokenLocation(p.curToken))
	case lexer.LBracket:
		computed = true
		p.nextToken()
		expr := p.parseExpression(lowest)
		if expr == nil {
			return nil
		}
		key = expr
		if !p.expectPeek(lexer.RBracket) {
			return nil
		}
	default:
		msg := "unexpected token " + string(p.curToken.Type) + " in object literal property"
		p.errors = append(p.errors, errors.New(msg))
		return nil
	}

	// shorthand property for identifiers only
	if !computed {
		if ident, ok := key.(*ast.Identifier); ok {
			if p.peekTokenIs(lexer.Comma) || p.peekTokenIs(lexer.RBrace) {
				loc := p.locFrom(start, p.curToken.End)
				return ast.NewObjectProperty(key, ident, ast.PropertyInit, false, true, false, loc)
			}
		}
	}

	if !p.expectPeek(lexer.Colon) {
		return nil
	}

	p.nextToken()
	value := p.parseExpression(sequencePrec)
	if value == nil {
		return nil
	}

	loc := p.locFrom(start, p.curToken.End)
	return ast.NewObjectProperty(key, value, ast.PropertyInit, computed, false, false, loc)
}

func (p *Parser) wrapNewExpression(expr ast.Expression, start lexer.Position) ast.Expression {
	newStart := convertPosition(start)
	switch e := expr.(type) {
	case *ast.CallExpression:
		switch callee := e.Callee.(type) {
		case *ast.CallExpression:
			wrapped := p.wrapNewExpression(callee, start)
			if wrapped != e.Callee {
				e.Callee = wrapped
				p.extendNodeStart(e, newStart)
				return e
			}
		case *ast.MemberExpression:
			if containsCallExpression(callee.Object) {
				original := callee.Object
				p.wrapNewExpression(callee, start)
				if original != callee.Object {
					p.extendNodeStart(callee, newStart)
					p.extendNodeStart(e, newStart)
					return e
				}
			}
		}
		loc := ast.Location{Start: newStart, End: e.Loc().End}
		return ast.NewNewExpression(e.Callee, e.Arguments, loc)
	case *ast.MemberExpression:
		wrapped := p.wrapNewExpression(e.Object, start)
		if wrapped != e.Object {
			e.Object = wrapped
			p.extendNodeStart(e, newStart)
			return e
		}
		loc := ast.Location{Start: newStart, End: e.Loc().End}
		return ast.NewNewExpression(e, nil, loc)
	case *ast.NewExpression:
		p.extendNodeStart(e, newStart)
		return e
	default:
		loc := ast.Location{Start: newStart, End: expr.Loc().End}
		return ast.NewNewExpression(expr, nil, loc)
	}
}

func containsCallExpression(expr ast.Expression) bool {
	switch e := expr.(type) {
	case *ast.CallExpression:
		return true
	case *ast.MemberExpression:
		return containsCallExpression(e.Object)
	default:
		return false
	}
}

func (p *Parser) extendNodeStart(node ast.Node, start ast.Position) {
	if start.Offset < 0 {
		return
	}
	loc := node.Loc()
	if loc.Start.Offset > start.Offset {
		loc.Start = start
		p.setNodeLocation(node, loc)
	}
}

func (p *Parser) parseRegExpLiteral() ast.Expression {
	tok := p.curToken
	lit := tok.Literal
	pattern := ""
	flags := ""

	if len(lit) >= 2 && strings.HasPrefix(lit, "/") {
		lastSlash := strings.LastIndex(lit, "/")
		if lastSlash > 0 {
			pattern = lit[1:lastSlash]
			flags = lit[lastSlash+1:]
		}
	}

	return ast.NewRegExpLiteral(pattern, flags, p.tokenLocation(tok))
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
