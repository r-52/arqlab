package parser

import (
	"errors"
	"fmt"

	"es6-interpreter/ast"
	"es6-interpreter/lexer"
)

func (p *Parser) parseBindingElement(allowDefault bool) ast.Pattern {
	start := p.curToken.Start

	primary := p.parseBindingPrimary()
	if primary == nil {
		return nil
	}

	if allowDefault && p.peekTokenIs(lexer.Assign) {
		p.nextToken() // move to '='
		p.nextToken() // advance to initializer expression
		right := p.parseExpression(sequencePrec)
		if right == nil {
			return nil
		}
		loc := p.locFrom(start, p.curToken.End)
		return ast.NewAssignmentPattern(primary, right, loc)
	}

	return primary
}

func (p *Parser) parseBindingPrimary() ast.Pattern {
	switch p.curToken.Type {
	case lexer.Identifier:
		return ast.NewIdentifier(p.curToken.Literal, p.tokenLocation(p.curToken))
	case lexer.LBracket:
		return p.parseArrayPattern()
	case lexer.LBrace:
		return p.parseObjectPattern()
	default:
		msg := fmt.Sprintf("unsupported binding pattern starting with %s", p.curToken.Type)
		p.errors = append(p.errors, errors.New(msg))
		return nil
	}
}

func (p *Parser) parseArrayPattern() ast.Pattern {
	start := p.curToken.Start
	var elements ast.PatternList
	var rest *ast.RestElement

	if p.peekTokenIs(lexer.RBracket) {
		p.nextToken() // move to closing bracket
	} else {
		p.nextToken() // move to first element
		for !p.curTokenIs(lexer.RBracket) && !p.curTokenIs(lexer.EOF) {
			if p.curTokenIs(lexer.Comma) {
				elements = append(elements, nil)
				p.nextToken()
				continue
			}

			if p.curTokenIs(lexer.Ellipsis) {
				if rest != nil {
					p.errors = append(p.errors, errors.New("duplicate rest element in array pattern"))
					return nil
				}
				restStart := p.curToken.Start
				p.nextToken()
				arg := p.parseBindingElement(false)
				if arg == nil {
					return nil
				}
				rest = ast.NewRestElement(arg, p.locFrom(restStart, p.curToken.End))
				if !p.peekTokenIs(lexer.RBracket) {
					p.errors = append(p.errors, errors.New("rest element must be last in array pattern"))
					return nil
				}
				p.nextToken() // move to closing bracket
				break
			}

			elem := p.parseBindingElement(true)
			if elem == nil {
				return nil
			}
			elements = append(elements, elem)

			if p.peekTokenIs(lexer.Comma) {
				p.nextToken() // move to comma
				if p.peekTokenIs(lexer.RBracket) {
					elements = append(elements, nil)
					p.nextToken() // move to closing bracket
					break
				}
				p.nextToken() // move to next element
			} else {
				p.nextToken()
			}
		}
	}

	if !p.curTokenIs(lexer.RBracket) {
		p.errors = append(p.errors, errors.New("unterminated array pattern"))
		return nil
	}

	loc := p.locFrom(start, p.curToken.End)
	return ast.NewArrayPattern(elements, rest, loc)
}

func (p *Parser) parseObjectPattern() ast.Pattern {
	start := p.curToken.Start
	var props []*ast.ObjectPatternProperty
	var rest *ast.RestElement

	if p.peekTokenIs(lexer.RBrace) {
		p.nextToken()
	} else {
		p.nextToken()
		for !p.curTokenIs(lexer.RBrace) && !p.curTokenIs(lexer.EOF) {
			if p.curTokenIs(lexer.Ellipsis) {
				if rest != nil {
					p.errors = append(p.errors, errors.New("duplicate rest element in object pattern"))
					return nil
				}
				restStart := p.curToken.Start
				p.nextToken()
				arg := p.parseBindingElement(false)
				if arg == nil {
					return nil
				}
				rest = ast.NewRestElement(arg, p.locFrom(restStart, p.curToken.End))
				if !p.peekTokenIs(lexer.RBrace) {
					p.errors = append(p.errors, errors.New("rest element must be last in object pattern"))
					return nil
				}
				p.nextToken()
				break
			}

			prop := p.parseObjectPatternProperty()
			if prop == nil {
				return nil
			}
			props = append(props, prop)

			if p.peekTokenIs(lexer.Comma) {
				p.nextToken()
				if p.peekTokenIs(lexer.RBrace) {
					p.nextToken()
					break
				}
				p.nextToken()
			} else {
				p.nextToken()
			}
		}
	}

	if !p.curTokenIs(lexer.RBrace) {
		p.errors = append(p.errors, errors.New("unterminated object pattern"))
		return nil
	}

	loc := p.locFrom(start, p.curToken.End)
	return ast.NewObjectPattern(props, rest, loc)
}

func (p *Parser) parseObjectPatternProperty() *ast.ObjectPatternProperty {
	start := p.curToken.Start

	switch p.curToken.Type {
	case lexer.Identifier:
		keyTok := p.curToken
		key := ast.NewIdentifier(keyTok.Literal, p.tokenLocation(keyTok))
		basePattern := ast.NewIdentifier(keyTok.Literal, p.tokenLocation(keyTok))
		value := ast.Pattern(basePattern)
		shorthand := true

		if p.peekTokenIs(lexer.Colon) {
			shorthand = false
			p.nextToken() // move to ':'
			p.nextToken() // move to value start
			value = p.parseBindingElement(true)
			if value == nil {
				return nil
			}
		} else if p.peekTokenIs(lexer.Assign) {
			p.nextToken() // move to '='
			p.nextToken() // move to initializer expression
			right := p.parseExpression(sequencePrec)
			if right == nil {
				return nil
			}
			value = ast.NewAssignmentPattern(basePattern, right, p.locFrom(start, p.curToken.End))
		}

		loc := p.locFrom(start, p.curToken.End)
		return ast.NewObjectPatternProperty(key, value, false, shorthand, loc)
	default:
		msg := fmt.Sprintf("unsupported object pattern property starting with %s", p.curToken.Type)
		p.errors = append(p.errors, errors.New(msg))
		return nil
	}
}
