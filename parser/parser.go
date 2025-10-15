package parser

import (
	"errors"

	"es6-interpreter/ast"
	"es6-interpreter/lexer"
)

type prefixParseFn func() ast.Expression

type infixParseFn func(ast.Expression) ast.Expression

// Parser consumes tokens produced by the lexer and constructs an AST.
type Parser struct {
	lex *lexer.Lexer

	curToken  lexer.Token
	peekToken lexer.Token

	errors []error

	prefixFns map[lexer.TokenType]prefixParseFn
	infixFns  map[lexer.TokenType]infixParseFn
}

// New returns a parser initialised from ECMAScript source text.
func New(src string) *Parser {
	return NewFromLexer(lexer.New(src))
}

// NewFromLexer returns a parser that pulls tokens directly from the supplied lexer.
func NewFromLexer(l *lexer.Lexer) *Parser {
	p := &Parser{
		lex:       l,
		prefixFns: make(map[lexer.TokenType]prefixParseFn),
		infixFns:  make(map[lexer.TokenType]infixParseFn),
	}

	// prime tokens
	p.nextToken()
	p.nextToken()

	p.registerPrefixFns()
	p.registerInfixFns()

	return p
}

// Errors returns the list of all parsing errors encountered.
func (p *Parser) Errors() []error {
	return p.errors
}

// ParseProgram parses the entire input into a Program node.
func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := ast.NewProgram(nil, ast.SourceTypeScript, ast.Location{})

	for !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Body = append(program.Body, stmt)
		}
		p.nextToken()
	}

	if len(program.Body) > 0 {
		first := program.Body[0].Loc()
		last := program.Body[len(program.Body)-1].Loc()
		program.SetLoc(ast.Location{Start: first.Start, End: last.End})
	}

	if len(p.errors) > 0 {
		return nil, errors.Join(p.errors...)
	}

	return program, nil
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) curTokenIs(tt lexer.TokenType) bool {
	return p.curToken.Type == tt
}

func (p *Parser) peekTokenIs(tt lexer.TokenType) bool {
	return p.peekToken.Type == tt
}

func (p *Parser) expectPeek(tt lexer.TokenType) bool {
	if p.peekTokenIs(tt) {
		p.nextToken()
		return true
	}
	p.peekError(tt)
	return false
}

func (p *Parser) peekError(tt lexer.TokenType) {
	msg := "expected next token to be " + string(tt) + ", got " + string(p.peekToken.Type)
	p.errors = append(p.errors, errors.New(msg))
}

func (p *Parser) curLoc() ast.Location {
	return ast.Location{
		Start: convertPosition(p.curToken.Start),
		End:   convertPosition(p.curToken.End),
	}
}

func (p *Parser) locFrom(start lexer.Position, end lexer.Position) ast.Location {
	return ast.Location{
		Start: convertPosition(start),
		End:   convertPosition(end),
	}
}

func (p *Parser) tokenLocation(tok lexer.Token) ast.Location {
	return p.locFrom(tok.Start, tok.End)
}

func convertPosition(pos lexer.Position) ast.Position {
	return ast.Position{Offset: pos.Offset, Line: pos.Line, Column: pos.Column}
}
