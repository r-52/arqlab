package parser

import "es6-interpreter/lexer"

type precedence int

const (
	lowest precedence = iota
	sequencePrec
	assignmentPrec
	conditionalPrec
	logicalOrPrec
	logicalAndPrec
	bitwiseOrPrec
	bitwiseXorPrec
	bitwiseAndPrec
	equalityPrec
	relationalPrec
	shiftPrec
	additivePrec
	multiplicativePrec
	prefixPrec
	postfixPrec
	callPrec
)

var precedences = map[lexer.TokenType]precedence{
	lexer.Comma:               sequencePrec,
	lexer.Assign:              assignmentPrec,
	lexer.PlusAssign:          assignmentPrec,
	lexer.MinusAssign:         assignmentPrec,
	lexer.MultiplyAssign:      assignmentPrec,
	lexer.DivideAssign:        assignmentPrec,
	lexer.ModuloAssign:        assignmentPrec,
	lexer.ShiftLeftAssign:     assignmentPrec,
	lexer.ShiftRightAssign:    assignmentPrec,
	lexer.UnsignedShiftAssign: assignmentPrec,
	lexer.BitwiseAndAssign:    assignmentPrec,
	lexer.BitwiseOrAssign:     assignmentPrec,
	lexer.BitwiseXorAssign:    assignmentPrec,
	lexer.Question:            conditionalPrec,
	lexer.LogicalOr:           logicalOrPrec,
	lexer.LogicalAnd:          logicalAndPrec,
	lexer.BitwiseOr:           bitwiseOrPrec,
	lexer.BitwiseXor:          bitwiseXorPrec,
	lexer.BitwiseAnd:          bitwiseAndPrec,
	lexer.Equal:               equalityPrec,
	lexer.NotEqual:            equalityPrec,
	lexer.StrictEqual:         equalityPrec,
	lexer.StrictNotEqual:      equalityPrec,
	lexer.LessThan:            relationalPrec,
	lexer.LessEqual:           relationalPrec,
	lexer.GreaterThan:         relationalPrec,
	lexer.GreaterEqual:        relationalPrec,
	lexer.KeywordIn:           relationalPrec,
	lexer.KeywordInstanceof:   relationalPrec,
	lexer.ShiftLeft:           shiftPrec,
	lexer.ShiftRight:          shiftPrec,
	lexer.UnsignedShiftRight:  shiftPrec,
	lexer.Plus:                additivePrec,
	lexer.Minus:               additivePrec,
	lexer.Multiply:            multiplicativePrec,
	lexer.Divide:              multiplicativePrec,
	lexer.Modulo:              multiplicativePrec,
	lexer.Increment:           postfixPrec,
	lexer.Decrement:           postfixPrec,
	lexer.LParen:              callPrec,
	lexer.LBracket:            callPrec,
	lexer.Dot:                 callPrec,
}

func (p *Parser) peekPrecedence() precedence {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return lowest
}

func (p *Parser) curPrecedence() precedence {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return lowest
}
