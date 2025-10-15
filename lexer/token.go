package lexer

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// TokenType represents the classification of a lexical token.
type TokenType string

// Token encapsulates a lexical token including its literal value and source span.
type Token struct {
	Type    TokenType
	Literal string
	Start   Position
	End     Position
}

// Position tracks a byte offset and human readable coordinates within the source.
type Position struct {
	Offset int // zero-based byte offset into the source text
	Line   int // one-based line number
	Column int // zero-based UTF-16 column per ECMAScript convention
}

// Sentinel and generic token types.
const (
	Illegal TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"
	Comment TokenType = "COMMENT"

	Identifier TokenType = "IDENT"
	Number     TokenType = "NUMBER"
	String     TokenType = "STRING"
	Regex      TokenType = "REGEXP"
)

// Template component token types used while scanning template literals.
const (
	TemplateHead      TokenType = "TEMPLATE_HEAD"
	TemplateMiddle    TokenType = "TEMPLATE_MIDDLE"
	TemplateTail      TokenType = "TEMPLATE_TAIL"
	TemplateExprStart TokenType = "TEMPLATE_EXPR_START"
	TemplateExprEnd   TokenType = "TEMPLATE_EXPR_END"
)

// Literal tokens that evaluate to intrinsic values.
const (
	NullLiteral  TokenType = "NULL"
	TrueLiteral  TokenType = "TRUE"
	FalseLiteral TokenType = "FALSE"
)

// Punctuation tokens.
const (
	LParen    TokenType = "LPAREN"
	RParen    TokenType = "RPAREN"
	LBrace    TokenType = "LBRACE"
	RBrace    TokenType = "RBRACE"
	LBracket  TokenType = "LBRACKET"
	RBracket  TokenType = "RBRACKET"
	Semicolon TokenType = "SEMICOLON"
	Comma     TokenType = "COMMA"
	Colon     TokenType = "COLON"
	Dot       TokenType = "DOT"
	Question  TokenType = "QUESTION"
	Backtick  TokenType = "BACKTICK"
)

// Operator tokens covering arithmetic, comparison, logical, and assignment operators.
const (
	Assign TokenType = "ASSIGN"

	Plus       TokenType = "PLUS"
	Minus      TokenType = "MINUS"
	Multiply   TokenType = "MULTIPLY"
	Divide     TokenType = "DIVIDE"
	Modulo     TokenType = "MODULO"
	Increment  TokenType = "INCREMENT"
	Decrement  TokenType = "DECREMENT"
	BitwiseNot TokenType = "BITWISE_NOT"
	LogicalNot TokenType = "LOGICAL_NOT"

	ShiftLeft          TokenType = "SHIFT_LEFT"
	ShiftRight         TokenType = "SHIFT_RIGHT"
	UnsignedShiftRight TokenType = "UNSIGNED_SHIFT_RIGHT"

	BitwiseAnd TokenType = "BITWISE_AND"
	BitwiseOr  TokenType = "BITWISE_OR"
	BitwiseXor TokenType = "BITWISE_XOR"

	LogicalAnd TokenType = "LOGICAL_AND"
	LogicalOr  TokenType = "LOGICAL_OR"

	Equal          TokenType = "EQUAL"
	StrictEqual    TokenType = "STRICT_EQUAL"
	NotEqual       TokenType = "NOT_EQUAL"
	StrictNotEqual TokenType = "STRICT_NOT_EQUAL"
	LessThan       TokenType = "LESS_THAN"
	LessEqual      TokenType = "LESS_EQUAL"
	GreaterThan    TokenType = "GREATER_THAN"
	GreaterEqual   TokenType = "GREATER_EQUAL"

	PlusAssign          TokenType = "PLUS_ASSIGN"
	MinusAssign         TokenType = "MINUS_ASSIGN"
	MultiplyAssign      TokenType = "MULTIPLY_ASSIGN"
	DivideAssign        TokenType = "DIVIDE_ASSIGN"
	ModuloAssign        TokenType = "MODULO_ASSIGN"
	ShiftLeftAssign     TokenType = "SHIFT_LEFT_ASSIGN"
	ShiftRightAssign    TokenType = "SHIFT_RIGHT_ASSIGN"
	UnsignedShiftAssign TokenType = "UNSIGNED_SHIFT_ASSIGN"
	BitwiseAndAssign    TokenType = "BITWISE_AND_ASSIGN"
	BitwiseOrAssign     TokenType = "BITWISE_OR_ASSIGN"
	BitwiseXorAssign    TokenType = "BITWISE_XOR_ASSIGN"

	Arrow    TokenType = "ARROW"
	Ellipsis TokenType = "ELLIPSIS"
)

// Keyword tokens. Async/await intentionally omitted.
const (
	KeywordBreak      TokenType = "BREAK"
	KeywordCase       TokenType = "CASE"
	KeywordCatch      TokenType = "CATCH"
	KeywordClass      TokenType = "CLASS"
	KeywordConst      TokenType = "CONST"
	KeywordContinue   TokenType = "CONTINUE"
	KeywordDebugger   TokenType = "DEBUGGER"
	KeywordDefault    TokenType = "DEFAULT"
	KeywordDelete     TokenType = "DELETE"
	KeywordDo         TokenType = "DO"
	KeywordElse       TokenType = "ELSE"
	KeywordEnum       TokenType = "ENUM"
	KeywordExport     TokenType = "EXPORT"
	KeywordExtends    TokenType = "EXTENDS"
	KeywordFinally    TokenType = "FINALLY"
	KeywordFor        TokenType = "FOR"
	KeywordFunction   TokenType = "FUNCTION"
	KeywordIf         TokenType = "IF"
	KeywordImport     TokenType = "IMPORT"
	KeywordIn         TokenType = "IN"
	KeywordInstanceof TokenType = "INSTANCEOF"
	KeywordLet        TokenType = "LET"
	KeywordNew        TokenType = "NEW"
	KeywordReturn     TokenType = "RETURN"
	KeywordSuper      TokenType = "SUPER"
	KeywordSwitch     TokenType = "SWITCH"
	KeywordThis       TokenType = "THIS"
	KeywordThrow      TokenType = "THROW"
	KeywordTry        TokenType = "TRY"
	KeywordTypeof     TokenType = "TYPEOF"
	KeywordVar        TokenType = "VAR"
	KeywordVoid       TokenType = "VOID"
	KeywordWhile      TokenType = "WHILE"
	KeywordWith       TokenType = "WITH"
	KeywordYield      TokenType = "YIELD"
	KeywordPackage    TokenType = "PACKAGE"
	KeywordPrivate    TokenType = "PRIVATE"
	KeywordProtected  TokenType = "PROTECTED"
	KeywordPublic     TokenType = "PUBLIC"
	KeywordInterface  TokenType = "INTERFACE"
	KeywordImplements TokenType = "IMPLEMENTS"
)

var keywords = map[string]TokenType{
	"break":      KeywordBreak,
	"case":       KeywordCase,
	"catch":      KeywordCatch,
	"class":      KeywordClass,
	"const":      KeywordConst,
	"continue":   KeywordContinue,
	"debugger":   KeywordDebugger,
	"default":    KeywordDefault,
	"delete":     KeywordDelete,
	"do":         KeywordDo,
	"else":       KeywordElse,
	"enum":       KeywordEnum,
	"export":     KeywordExport,
	"extends":    KeywordExtends,
	"finally":    KeywordFinally,
	"for":        KeywordFor,
	"function":   KeywordFunction,
	"if":         KeywordIf,
	"import":     KeywordImport,
	"in":         KeywordIn,
	"instanceof": KeywordInstanceof,
	"let":        KeywordLet,
	"new":        KeywordNew,
	"return":     KeywordReturn,
	"super":      KeywordSuper,
	"switch":     KeywordSwitch,
	"this":       KeywordThis,
	"throw":      KeywordThrow,
	"try":        KeywordTry,
	"typeof":     KeywordTypeof,
	"var":        KeywordVar,
	"void":       KeywordVoid,
	"while":      KeywordWhile,
	"with":       KeywordWith,
	"yield":      KeywordYield,
	"package":    KeywordPackage,
	"private":    KeywordPrivate,
	"protected":  KeywordProtected,
	"public":     KeywordPublic,
	"interface":  KeywordInterface,
	"implements": KeywordImplements,
	"null":       NullLiteral,
	"true":       TrueLiteral,
	"false":      FalseLiteral,
}

// LookupIdentifier returns the token type for a given identifier or keyword.
func LookupIdentifier(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Identifier
}

// Keywords returns the sorted list of reserved words recognised by the lexer.
func Keywords() []string {
	keys := make([]string, 0, len(keywords))
	for k := range keywords {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// IsKeyword reports whether the provided string is a reserved keyword.
func IsKeyword(word string) bool {
	_, ok := keywords[strings.ToLower(word)]
	return ok
}

// String implements fmt.Stringer for tokens, aiding debugging and logging.
func (t Token) String() string {
	if t.Literal != "" {
		return fmt.Sprintf("%s(%s)", t.Type, strconv.Quote(t.Literal))
	}
	return string(t.Type)
}

// String renders the position in a human friendly format.
func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column+1)
}
