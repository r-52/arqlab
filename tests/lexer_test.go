package tests

import (
	"testing"

	"es6-interpreter/lexer"
)

type tokenExpectation struct {
	typ     lexer.TokenType
	literal string
}

func collectTokens(t *testing.T, l *lexer.Lexer) []lexer.Token {
	t.Helper()
	var tokens []lexer.Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == lexer.EOF || tok.Type == lexer.Illegal {
			break
		}
	}
	return tokens
}

func assertTokens(t *testing.T, tokens []lexer.Token, want []tokenExpectation) {
	t.Helper()
	if len(tokens) != len(want) {
		t.Fatalf("token length mismatch: got %d, want %d\n%v", len(tokens), len(want), tokens)
	}
	for i, tok := range tokens {
		if tok.Type != want[i].typ {
			t.Fatalf("token %d type mismatch: got %s, want %s", i, tok.Type, want[i].typ)
		}
		if want[i].literal != "" && tok.Literal != want[i].literal {
			t.Fatalf("token %d literal mismatch: got %q, want %q", i, tok.Literal, want[i].literal)
		}
	}
}

func TestLexerIdentifiersAndKeywords(t *testing.T) {
	source := "let answer = 42; const truth = true;"
	l := lexer.New(source)
	got := collectTokens(t, l)
	want := []tokenExpectation{
		{lexer.KeywordLet, "let"},
		{lexer.Identifier, "answer"},
		{lexer.Assign, "="},
		{lexer.Number, "42"},
		{lexer.Semicolon, ";"},
		{lexer.KeywordConst, "const"},
		{lexer.Identifier, "truth"},
		{lexer.Assign, "="},
		{lexer.TrueLiteral, "true"},
		{lexer.Semicolon, ";"},
		{lexer.EOF, ""},
	}
	assertTokens(t, got, want)
}

func TestLexerNumberVariants(t *testing.T) {
	source := "0 123 12.34 6.02e23 0xFF 0o755 0b1010"
	l := lexer.New(source)
	got := collectTokens(t, l)
	want := []tokenExpectation{
		{lexer.Number, "0"},
		{lexer.Number, "123"},
		{lexer.Number, "12.34"},
		{lexer.Number, "6.02e23"},
		{lexer.Number, "0xFF"},
		{lexer.Number, "0o755"},
		{lexer.Number, "0b1010"},
		{lexer.EOF, ""},
	}
	assertTokens(t, got, want)
}

func TestLexerStringLiterals(t *testing.T) {
	source := "'single \\'quoted\\'' \"double \\\"quoted\\\"\""
	l := lexer.New(source)
	got := collectTokens(t, l)
	want := []tokenExpectation{
		{lexer.String, "'single \\'quoted\\''"},
		{lexer.String, "\"double \\\"quoted\\\"\""},
		{lexer.EOF, ""},
	}
	assertTokens(t, got, want)
}

func TestLexerTemplateLiteral(t *testing.T) {
	source := "`hello ${name}!`"
	l := lexer.New(source)
	got := collectTokens(t, l)
	want := []tokenExpectation{
		{lexer.TemplateHead, "hello "},
		{lexer.TemplateExprStart, "${"},
		{lexer.Identifier, "name"},
		{lexer.TemplateExprEnd, "}"},
		{lexer.TemplateTail, "!"},
		{lexer.EOF, ""},
	}
	assertTokens(t, got, want)
}

func TestLexerRegularExpression(t *testing.T) {
	source := "var r = /a[b-d]+/gi; r.test('abc');"
	l := lexer.New(source)
	got := collectTokens(t, l)
	want := []tokenExpectation{
		{lexer.KeywordVar, "var"},
		{lexer.Identifier, "r"},
		{lexer.Assign, "="},
		{lexer.Regex, "/a[b-d]+/gi"},
		{lexer.Semicolon, ";"},
		{lexer.Identifier, "r"},
		{lexer.Dot, "."},
		{lexer.Identifier, "test"},
		{lexer.LParen, "("},
		{lexer.String, "'abc'"},
		{lexer.RParen, ")"},
		{lexer.Semicolon, ";"},
		{lexer.EOF, ""},
	}
	assertTokens(t, got, want)
}

func TestLineTerminatorHandling(t *testing.T) {
	source := "return\n/x/"
	l := lexer.New(source)
	got := collectTokens(t, l)
	want := []tokenExpectation{
		{lexer.KeywordReturn, "return"},
		{lexer.Regex, "/x/"},
		{lexer.EOF, ""},
	}
	assertTokens(t, got, want)
}

func TestUnterminatedStringProducesIllegal(t *testing.T) {
	source := "\"unterminated"
	l := lexer.New(source)
	tokens := collectTokens(t, l)
	last := tokens[len(tokens)-1]
	if last.Type != lexer.Illegal {
		t.Fatalf("expected last token to be ILLEGAL, got %s", last.Type)
	}
}

func TestUnterminatedCommentProducesIllegal(t *testing.T) {
	source := "/* comment"
	l := lexer.New(source)
	tokens := collectTokens(t, l)
	last := tokens[len(tokens)-1]
	if last.Type != lexer.Illegal {
		t.Fatalf("expected ILLEGAL token, got %s", last.Type)
	}
}

func TestTemplateWithNestedBraces(t *testing.T) {
	source := "`value ${1 + {nested: true}.nested}`"
	l := lexer.New(source)
	got := collectTokens(t, l)
	want := []tokenExpectation{
		{lexer.TemplateHead, "value "},
		{lexer.TemplateExprStart, "${"},
		{lexer.Number, "1"},
		{lexer.Plus, "+"},
		{lexer.LBrace, "{"},
		{lexer.Identifier, "nested"},
		{lexer.Colon, ":"},
		{lexer.TrueLiteral, "true"},
		{lexer.RBrace, "}"},
		{lexer.Dot, "."},
		{lexer.Identifier, "nested"},
		{lexer.TemplateExprEnd, "}"},
		{lexer.TemplateTail, ""},
		{lexer.EOF, ""},
	}
	assertTokens(t, got, want)
}

func TestIllegalCharacter(t *testing.T) {
	source := "let x = #;"
	l := lexer.New(source)
	tokens := collectTokens(t, l)
	last := tokens[len(tokens)-1]
	if last.Type != lexer.Illegal {
		t.Fatalf("expected ILLEGAL token for #, got %s", last.Type)
	}
}
