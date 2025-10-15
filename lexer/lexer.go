package lexer

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type templateContext struct {
	braceDepth int
}

// Lexer transforms ECMAScript source text into a stream of tokens.
type Lexer struct {
	src                  string
	ch                   rune
	chPos                Position
	nextPos              Position
	buffer               []Token
	contexts             []templateContext
	continueTemplate     bool
	canStartRegex        bool
	lineTerminatorBefore bool
	lastTokenType        TokenType
	err                  error
}

// New creates a new lexer with the provided ECMAScript source code.
func New(src string) *Lexer {
	l := &Lexer{
		src:           src,
		nextPos:       Position{Line: 1, Column: 0, Offset: 0},
		canStartRegex: true,
		lastTokenType: Illegal,
	}
	l.advance()
	return l
}

// NextToken returns the next token from the input stream.
func (l *Lexer) NextToken() Token {
	for {
		if len(l.buffer) > 0 {
			tok := l.buffer[0]
			l.buffer = l.buffer[1:]
			l.updateAfterToken(tok)
			return tok
		}

		if l.err != nil {
			tok := Token{Type: Illegal, Literal: l.err.Error(), Start: l.chPos, End: l.chPos}
			l.err = nil
			l.updateAfterToken(tok)
			return tok
		}

		if l.continueTemplate {
			if err := l.lexTemplateChunk(false); err != nil {
				l.err = err
				continue
			}
			continue
		}

		l.skipWhitespaceAndComments()

		if l.err != nil {
			continue
		}

		start := l.chPos

		switch l.ch {
		case 0:
			tok := Token{Type: EOF, Start: start, End: start}
			l.updateAfterToken(tok)
			return tok
		case '`':
			if err := l.lexTemplateChunk(true); err != nil {
				l.err = err
				continue
			}
			continue
		case '/':
			if tok, ok := l.scanSlash(start); ok {
				l.updateAfterToken(tok)
				return tok
			}
		case '\'', '"':
			if tok, ok := l.scanString(start, l.ch); ok {
				l.updateAfterToken(tok)
				return tok
			}
		case '.':
			tok := l.scanDot(start)
			l.updateAfterToken(tok)
			return tok
		case '+', '-', '*', '%', '&', '|', '^', '!', '=', '<', '>', '?', ':':
			tok := l.scanOperator(start)
			l.updateAfterToken(tok)
			return tok
		case '{', '}', '(', ')', '[', ']', ',', ';':
			tok := l.scanPunctuation(start)
			l.updateAfterToken(tok)
			return tok
		default:
			if l.isIdentifierStart(l.ch) {
				tok := l.scanIdentifier(start)
				l.updateAfterToken(tok)
				return tok
			}
			if unicode.IsDigit(l.ch) {
				tok := l.scanNumber(start)
				l.updateAfterToken(tok)
				return tok
			}

			literal := string(l.ch)
			l.advance()
			tok := Token{Type: Illegal, Literal: fmt.Sprintf("unexpected character %q", literal), Start: start, End: l.chPos}
			l.updateAfterToken(tok)
			return tok
		}
	}
}

func (l *Lexer) scanSlash(start Position) (Token, bool) {
	next := l.peekRune()
	if next == '/' {
		l.advance()
		l.advance()
		l.consumeLineComment()
		return Token{}, false
	}
	if next == '*' {
		l.advance()
		l.advance()
		if err := l.consumeBlockComment(); err != nil {
			l.err = err
		}
		return Token{}, false
	}

	if l.canStartRegex {
		tok, err := l.scanRegularExpression(start)
		if err != nil {
			l.err = err
			return Token{}, false
		}
		return tok, true
	}

	l.advance()
	tokType := Divide
	literal := "/"
	if next == '=' {
		l.advance()
		tokType = DivideAssign
		literal = "/="
	}
	return Token{Type: tokType, Literal: literal, Start: start, End: l.chPos}, true
}

func (l *Lexer) scanString(start Position, quote rune) (Token, bool) {
	l.advance()
	for {
		switch l.ch {
		case 0, '\n':
			l.err = fmt.Errorf("unterminated string literal")
			return Token{}, false
		case '\\':
			l.advance()
			if l.ch == 0 {
				l.err = fmt.Errorf("unterminated escape sequence")
				return Token{}, false
			}
			l.advance()
		case quote:
			l.advance()
			literal := l.slice(start, l.chPos)
			tok := Token{Type: String, Literal: literal, Start: start, End: l.chPos}
			return tok, true
		default:
			l.advance()
		}
	}
}

func (l *Lexer) scanDot(start Position) Token {
	if l.peekRune() == '.' && l.peekRuneN(1) == '.' {
		l.advance()
		l.advance()
		l.advance()
		return Token{Type: Ellipsis, Literal: "...", Start: start, End: l.chPos}
	}

	if unicode.IsDigit(l.peekRune()) {
		return l.scanNumber(start)
	}

	l.advance()
	return Token{Type: Dot, Literal: ".", Start: start, End: l.chPos}
}

func (l *Lexer) scanOperator(start Position) Token {
	switch l.ch {
	case '+':
		l.advance()
		if l.ch == '+' {
			l.advance()
			return Token{Type: Increment, Literal: "++", Start: start, End: l.chPos}
		}
		if l.ch == '=' {
			l.advance()
			return Token{Type: PlusAssign, Literal: "+=", Start: start, End: l.chPos}
		}
		return Token{Type: Plus, Literal: "+", Start: start, End: l.chPos}
	case '-':
		l.advance()
		if l.ch == '-' {
			l.advance()
			return Token{Type: Decrement, Literal: "--", Start: start, End: l.chPos}
		}
		if l.ch == '=' {
			l.advance()
			return Token{Type: MinusAssign, Literal: "-=", Start: start, End: l.chPos}
		}
		return Token{Type: Minus, Literal: "-", Start: start, End: l.chPos}
	case '*':
		l.advance()
		if l.ch == '=' {
			l.advance()
			return Token{Type: MultiplyAssign, Literal: "*=", Start: start, End: l.chPos}
		}
		return Token{Type: Multiply, Literal: "*", Start: start, End: l.chPos}
	case '%':
		l.advance()
		if l.ch == '=' {
			l.advance()
			return Token{Type: ModuloAssign, Literal: "%=", Start: start, End: l.chPos}
		}
		return Token{Type: Modulo, Literal: "%", Start: start, End: l.chPos}
	case '&':
		l.advance()
		if l.ch == '&' {
			l.advance()
			return Token{Type: LogicalAnd, Literal: "&&", Start: start, End: l.chPos}
		}
		if l.ch == '=' {
			l.advance()
			return Token{Type: BitwiseAndAssign, Literal: "&=", Start: start, End: l.chPos}
		}
		return Token{Type: BitwiseAnd, Literal: "&", Start: start, End: l.chPos}
	case '|':
		l.advance()
		if l.ch == '|' {
			l.advance()
			return Token{Type: LogicalOr, Literal: "||", Start: start, End: l.chPos}
		}
		if l.ch == '=' {
			l.advance()
			return Token{Type: BitwiseOrAssign, Literal: "|=", Start: start, End: l.chPos}
		}
		return Token{Type: BitwiseOr, Literal: "|", Start: start, End: l.chPos}
	case '^':
		l.advance()
		if l.ch == '=' {
			l.advance()
			return Token{Type: BitwiseXorAssign, Literal: "^=", Start: start, End: l.chPos}
		}
		return Token{Type: BitwiseXor, Literal: "^", Start: start, End: l.chPos}
	case '!':
		l.advance()
		if l.ch == '=' {
			l.advance()
			if l.ch == '=' {
				l.advance()
				return Token{Type: StrictNotEqual, Literal: "!==", Start: start, End: l.chPos}
			}
			return Token{Type: NotEqual, Literal: "!=", Start: start, End: l.chPos}
		}
		return Token{Type: LogicalNot, Literal: "!", Start: start, End: l.chPos}
	case '=':
		l.advance()
		if l.ch == '=' {
			l.advance()
			if l.ch == '=' {
				l.advance()
				return Token{Type: StrictEqual, Literal: "===", Start: start, End: l.chPos}
			}
			return Token{Type: Equal, Literal: "==", Start: start, End: l.chPos}
		}
		if l.ch == '>' {
			l.advance()
			return Token{Type: Arrow, Literal: "=>", Start: start, End: l.chPos}
		}
		return Token{Type: Assign, Literal: "=", Start: start, End: l.chPos}
	case '<':
		l.advance()
		if l.ch == '<' {
			l.advance()
			if l.ch == '=' {
				l.advance()
				return Token{Type: ShiftLeftAssign, Literal: "<<=", Start: start, End: l.chPos}
			}
			return Token{Type: ShiftLeft, Literal: "<<", Start: start, End: l.chPos}
		}
		if l.ch == '=' {
			l.advance()
			return Token{Type: LessEqual, Literal: "<=", Start: start, End: l.chPos}
		}
		return Token{Type: LessThan, Literal: "<", Start: start, End: l.chPos}
	case '>':
		l.advance()
		if l.ch == '>' {
			l.advance()
			if l.ch == '>' {
				l.advance()
				if l.ch == '=' {
					l.advance()
					return Token{Type: UnsignedShiftAssign, Literal: ">>>=", Start: start, End: l.chPos}
				}
				return Token{Type: UnsignedShiftRight, Literal: ">>>", Start: start, End: l.chPos}
			}
			if l.ch == '=' {
				l.advance()
				return Token{Type: ShiftRightAssign, Literal: ">>=", Start: start, End: l.chPos}
			}
			return Token{Type: ShiftRight, Literal: ">>", Start: start, End: l.chPos}
		}
		if l.ch == '=' {
			l.advance()
			return Token{Type: GreaterEqual, Literal: ">=", Start: start, End: l.chPos}
		}
		return Token{Type: GreaterThan, Literal: ">", Start: start, End: l.chPos}
	case '?':
		l.advance()
		return Token{Type: Question, Literal: "?", Start: start, End: l.chPos}
	case ':':
		l.advance()
		return Token{Type: Colon, Literal: ":", Start: start, End: l.chPos}
	}
	l.advance()
	return Token{Type: Illegal, Literal: l.slice(start, l.chPos), Start: start, End: l.chPos}
}

func (l *Lexer) scanPunctuation(start Position) Token {
	switch l.ch {
	case '{':
		l.advance()
		return Token{Type: LBrace, Literal: "{", Start: start, End: l.chPos}
	case '}':
		if len(l.contexts) > 0 && l.contexts[len(l.contexts)-1].braceDepth == 0 {
			l.advance()
			tok := Token{Type: TemplateExprEnd, Literal: "}", Start: start, End: l.chPos}
			l.continueTemplate = true
			return tok
		}
		if len(l.contexts) > 0 {
			l.contexts[len(l.contexts)-1].braceDepth--
		}
		l.advance()
		return Token{Type: RBrace, Literal: "}", Start: start, End: l.chPos}
	case '(':
		l.advance()
		return Token{Type: LParen, Literal: "(", Start: start, End: l.chPos}
	case ')':
		l.advance()
		return Token{Type: RParen, Literal: ")", Start: start, End: l.chPos}
	case '[':
		l.advance()
		return Token{Type: LBracket, Literal: "[", Start: start, End: l.chPos}
	case ']':
		l.advance()
		return Token{Type: RBracket, Literal: "]", Start: start, End: l.chPos}
	case ',':
		l.advance()
		return Token{Type: Comma, Literal: ",", Start: start, End: l.chPos}
	case ';':
		l.advance()
		return Token{Type: Semicolon, Literal: ";", Start: start, End: l.chPos}
	}
	l.advance()
	return Token{Type: Illegal, Literal: l.slice(start, l.chPos), Start: start, End: l.chPos}
}

func (l *Lexer) scanIdentifier(start Position) Token {
	var b strings.Builder
	for l.isIdentifierPart(l.ch) {
		b.WriteRune(l.ch)
		l.advance()
	}
	literal := b.String()
	typ := LookupIdentifier(literal)
	return Token{Type: typ, Literal: literal, Start: start, End: l.chPos}
}

func (l *Lexer) scanNumber(start Position) Token {
	literal, typ, err := l.readNumberLiteral()
	if err != nil {
		l.err = err
		return Token{Type: Illegal, Literal: literal, Start: start, End: l.chPos}
	}
	return Token{Type: typ, Literal: literal, Start: start, End: l.chPos}
}

func (l *Lexer) readNumberLiteral() (string, TokenType, error) {
	start := l.chPos
	if l.ch == '0' {
		next := l.peekRune()
		switch next {
		case 'x', 'X':
			l.advance()
			l.advance()
			if !l.consumeDigits(func(r rune) bool { return unicode.Is(unicode.Hex_Digit, r) }) {
				return l.slice(start, l.chPos), Illegal, fmt.Errorf("invalid hexadecimal literal")
			}
			literal := l.slice(start, l.chPos)
			return literal, Number, nil
		case 'o', 'O':
			l.advance()
			l.advance()
			if !l.consumeDigits(isOctalDigit) {
				return l.slice(start, l.chPos), Illegal, fmt.Errorf("invalid octal literal")
			}
			return l.slice(start, l.chPos), Number, nil
		case 'b', 'B':
			l.advance()
			l.advance()
			if !l.consumeDigits(isBinaryDigit) {
				return l.slice(start, l.chPos), Illegal, fmt.Errorf("invalid binary literal")
			}
			return l.slice(start, l.chPos), Number, nil
		}
	}

	l.consumeDigits(unicode.IsDigit)

	if l.ch == '.' {
		l.advance()
		if !l.consumeDigits(unicode.IsDigit) {
			return l.slice(start, l.chPos), Illegal, fmt.Errorf("invalid floating-point literal")
		}
	}

	if l.ch == 'e' || l.ch == 'E' {
		l.advance()
		if l.ch == '+' || l.ch == '-' {
			l.advance()
		}
		if !l.consumeDigits(unicode.IsDigit) {
			return l.slice(start, l.chPos), Illegal, fmt.Errorf("invalid exponent in numeric literal")
		}
	}

	return l.slice(start, l.chPos), Number, nil
}

func (l *Lexer) consumeDigits(match func(rune) bool) bool {
	count := 0
	for match(l.ch) {
		count++
		l.advance()
	}
	return count > 0
}

func (l *Lexer) lexTemplateChunk(startWithBacktick bool) error {
	chunkStart := l.chPos
	if startWithBacktick {
		l.advance()
		chunkStart = l.chPos
	}

	for {
		switch l.ch {
		case 0:
			return fmt.Errorf("unterminated template literal")
		case '`':
			lit := l.slice(chunkStart, l.chPos)
			tail := Token{Type: TemplateTail, Literal: lit, Start: chunkStart, End: l.chPos}
			l.advance()
			tail.End = l.chPos
			l.buffer = append(l.buffer, tail)
			if len(l.contexts) > 0 {
				l.contexts = l.contexts[:len(l.contexts)-1]
			}
			l.continueTemplate = false
			return nil
		case '$':
			if l.peekRune() == '{' {
				lit := l.slice(chunkStart, l.chPos)
				tokType := TemplateHead
				if !startWithBacktick {
					tokType = TemplateMiddle
				} else if len(l.contexts) > 0 {
					tokType = TemplateMiddle
				}
				head := Token{Type: tokType, Literal: lit, Start: chunkStart, End: l.chPos}
				l.advance()
				braceStart := l.chPos
				l.advance()
				exprStart := Token{Type: TemplateExprStart, Literal: "${", Start: braceStart, End: l.chPos}
				l.buffer = append(l.buffer, head, exprStart)
				l.contexts = append(l.contexts, templateContext{braceDepth: 0})
				l.continueTemplate = false
				return nil
			}
			l.advance()
		case '\\':
			l.advance()
			if l.ch == 0 {
				return fmt.Errorf("unterminated escape in template literal")
			}
			l.advance()
		default:
			if l.ch == '{' && len(l.contexts) > 0 {
				l.contexts[len(l.contexts)-1].braceDepth++
			}
			l.advance()
		}
	}
}

func (l *Lexer) scanRegularExpression(start Position) (Token, error) {
	literalStart := start
	l.advance()
	inClass := false
	for {
		switch l.ch {
		case 0, '\n':
			return Token{}, errors.New("unterminated regular expression literal")
		case '/':
			if !inClass {
				l.advance()
				goto flags
			}
			l.advance()
		case '[':
			inClass = true
			l.advance()
		case ']':
			inClass = false
			l.advance()
		case '\\':
			l.advance()
			if l.ch == 0 {
				return Token{}, errors.New("unterminated regular expression escape")
			}
			l.advance()
		default:
			l.advance()
		}
	}

flags:
	for l.isIdentifierPart(l.ch) {
		l.advance()
	}
	literal := l.slice(literalStart, l.chPos)
	return Token{Type: Regex, Literal: literal, Start: start, End: l.chPos}, nil
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		progressed := false
		switch l.ch {
		case ' ', '\t', '\f', '\v', '\u00a0':
			l.advance()
			progressed = true
		case '\n':
			l.lineTerminatorBefore = true
			l.advance()
			progressed = true
		default:
			if l.ch == '/' {
				next := l.peekRune()
				if next == '/' {
					l.advance()
					l.advance()
					l.consumeLineComment()
					progressed = true
					continue
				}
				if next == '*' {
					l.advance()
					l.advance()
					if err := l.consumeBlockComment(); err != nil {
						l.err = err
						return
					}
					progressed = true
					continue
				}
			}
		}
		if !progressed {
			return
		}
	}
}

func (l *Lexer) consumeLineComment() {
	for l.ch != 0 && l.ch != '\n' {
		l.advance()
	}
}

func (l *Lexer) consumeBlockComment() error {
	for {
		if l.ch == 0 {
			return errors.New("unterminated block comment")
		}
		if l.ch == '*' && l.peekRune() == '/' {
			l.advance()
			l.advance()
			return nil
		}
		if l.ch == '\n' {
			l.lineTerminatorBefore = true
		}
		l.advance()
	}
}

func (l *Lexer) updateAfterToken(tok Token) {
	switch tok.Type {
	case LBrace:
		if len(l.contexts) > 0 {
			l.contexts[len(l.contexts)-1].braceDepth++
		}
		l.canStartRegex = true
	case RBrace:
		if len(l.contexts) > 0 && l.contexts[len(l.contexts)-1].braceDepth > 0 {
			l.contexts[len(l.contexts)-1].braceDepth--
		}
		l.canStartRegex = false
	case Identifier, Number, String, TrueLiteral, FalseLiteral, NullLiteral, TemplateTail, RParen, RBracket:
		l.canStartRegex = false
	case Increment, Decrement:
		l.canStartRegex = true
	case TemplateHead, TemplateMiddle, TemplateExprStart:
		l.canStartRegex = true
		l.continueTemplate = false
	case TemplateExprEnd:
		l.canStartRegex = false
		l.continueTemplate = true
	default:
		l.canStartRegex = true
	}

	l.lastTokenType = tok.Type
	l.lineTerminatorBefore = false
}

func (l *Lexer) isIdentifierStart(r rune) bool {
	return r == '_' || r == '$' || unicode.IsLetter(r)
}

func (l *Lexer) isIdentifierPart(r rune) bool {
	return l.isIdentifierStart(r) || unicode.IsDigit(r)
}

func (l *Lexer) peekRune() rune {
	if l.nextPos.Offset >= len(l.src) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.src[l.nextPos.Offset:])
	if r == '\r' {
		return '\n'
	}
	return r
}

func (l *Lexer) peekRuneN(n int) rune {
	offset := l.nextPos.Offset
	var r rune
	for i := 0; i < n; i++ {
		if offset >= len(l.src) {
			return 0
		}
		r, size := utf8.DecodeRuneInString(l.src[offset:])
		if r == '\r' {
			r = '\n'
			size = 1
			if offset+1 < len(l.src) && l.src[offset+1] == '\n' {
				size++
			}
		}
		offset += size
	}
	if offset >= len(l.src) {
		return 0
	}
	r, _ = utf8.DecodeRuneInString(l.src[offset:])
	if r == '\r' {
		return '\n'
	}
	return r
}

func (l *Lexer) advance() {
	pos := l.nextPos
	if pos.Offset >= len(l.src) {
		l.ch = 0
		l.chPos = pos
		return
	}

	r, size := utf8.DecodeRuneInString(l.src[pos.Offset:])
	offset := pos.Offset + size
	if r == '\r' {
		r = '\n'
		if offset < len(l.src) && l.src[offset] == '\n' {
			offset++
		}
	}

	l.ch = r
	l.chPos = Position{Offset: pos.Offset, Line: pos.Line, Column: pos.Column}
	if r == '\n' {
		l.nextPos = Position{Offset: offset, Line: pos.Line + 1, Column: 0}
	} else {
		l.nextPos = Position{Offset: offset, Line: pos.Line, Column: pos.Column + 1}
	}
}

func (l *Lexer) slice(start, end Position) string {
	return l.src[start.Offset:end.Offset]
}

func isOctalDigit(r rune) bool {
	return r >= '0' && r <= '7'
}

func isBinaryDigit(r rune) bool {
	return r == '0' || r == '1'
}
