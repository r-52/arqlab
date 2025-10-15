package ast

import "fmt"

const (
	NumberLiteralKind   NodeKind = "NumberLiteral"
	StringLiteralKind   NodeKind = "StringLiteral"
	BooleanLiteralKind  NodeKind = "BooleanLiteral"
	NullLiteralKind     NodeKind = "NullLiteral"
	RegExpLiteralKind   NodeKind = "RegExpLiteral"
	TemplateLiteralKind NodeKind = "TemplateLiteral"
	TemplateElementKind NodeKind = "TemplateElement"
	ArrayLiteralKind    NodeKind = "ArrayLiteral"
	ObjectLiteralKind   NodeKind = "ObjectLiteral"
	ObjectPropertyKind  NodeKind = "ObjectProperty"
	SpreadElementKind   NodeKind = "SpreadElement"
)

// Literal marks expression nodes that evaluate to literal values.
type Literal interface {
	Expression
	literal()
}

// NumberLiteral represents numeric literal values of any radix.
type NumberLiteral struct {
	BaseNode
	Value string
}

func NewNumberLiteral(value string, loc Location) *NumberLiteral {
	return &NumberLiteral{BaseNode: NewBaseNode(NumberLiteralKind, loc), Value: value}
}

func (n *NumberLiteral) node()          {}
func (n *NumberLiteral) expression()    {}
func (n *NumberLiteral) literal()       {}
func (n *NumberLiteral) String() string { return fmt.Sprintf("NumberLiteral(%s)", n.Value) }

// StringLiteral represents quoted string literals.
type StringLiteral struct {
	BaseNode
	Value string
}

func NewStringLiteral(value string, loc Location) *StringLiteral {
	return &StringLiteral{BaseNode: NewBaseNode(StringLiteralKind, loc), Value: value}
}

func (s *StringLiteral) node()       {}
func (s *StringLiteral) expression() {}
func (s *StringLiteral) literal()    {}
func (s *StringLiteral) String() string {
	return fmt.Sprintf("StringLiteral(%q)", s.Value)
}

// BooleanLiteral represents the keywords true and false.
type BooleanLiteral struct {
	BaseNode
	Value bool
}

func NewBooleanLiteral(value bool, loc Location) *BooleanLiteral {
	return &BooleanLiteral{BaseNode: NewBaseNode(BooleanLiteralKind, loc), Value: value}
}

func (b *BooleanLiteral) node()       {}
func (b *BooleanLiteral) expression() {}
func (b *BooleanLiteral) literal()    {}
func (b *BooleanLiteral) String() string {
	return fmt.Sprintf("BooleanLiteral(%t)", b.Value)
}

// NullLiteral represents the keyword null.
type NullLiteral struct {
	BaseNode
}

func NewNullLiteral(loc Location) *NullLiteral {
	return &NullLiteral{BaseNode: NewBaseNode(NullLiteralKind, loc)}
}

func (n *NullLiteral) node()       {}
func (n *NullLiteral) expression() {}
func (n *NullLiteral) literal()    {}
func (n *NullLiteral) String() string {
	return "NullLiteral"
}

// RegExpLiteral represents /pattern/flags expressions.
type RegExpLiteral struct {
	BaseNode
	Pattern string
	Flags   string
}

func NewRegExpLiteral(pattern, flags string, loc Location) *RegExpLiteral {
	return &RegExpLiteral{BaseNode: NewBaseNode(RegExpLiteralKind, loc), Pattern: pattern, Flags: flags}
}

func (r *RegExpLiteral) node()       {}
func (r *RegExpLiteral) expression() {}
func (r *RegExpLiteral) literal()    {}
func (r *RegExpLiteral) String() string {
	if r.Flags == "" {
		return fmt.Sprintf("RegExpLiteral(%s)", r.Pattern)
	}
	return fmt.Sprintf("RegExpLiteral(%s/%s)", r.Pattern, r.Flags)
}

// TemplateLiteral represents a template literal with quasis and expressions.
type TemplateLiteral struct {
	BaseNode
	Quasis      []*TemplateElement
	Expressions []Expression
}

func NewTemplateLiteral(quasis []*TemplateElement, exprs []Expression, loc Location) *TemplateLiteral {
	return &TemplateLiteral{
		BaseNode:    NewBaseNode(TemplateLiteralKind, loc),
		Quasis:      quasis,
		Expressions: exprs,
	}
}

func (t *TemplateLiteral) node()       {}
func (t *TemplateLiteral) expression() {}
func (t *TemplateLiteral) literal()    {}
func (t *TemplateLiteral) String() string {
	return "TemplateLiteral"
}

// TemplateElement represents the static portion of a template literal.
type TemplateElement struct {
	BaseNode
	Raw    string
	Cooked string
	Tail   bool
}

func NewTemplateElement(raw, cooked string, tail bool, loc Location) *TemplateElement {
	return &TemplateElement{BaseNode: NewBaseNode(TemplateElementKind, loc), Raw: raw, Cooked: cooked, Tail: tail}
}

func (t *TemplateElement) node() {}
func (t *TemplateElement) String() string {
	if t.Cooked != "" {
		return fmt.Sprintf("TemplateElement(raw=%q,cooked=%q,tail=%t)", t.Raw, t.Cooked, t.Tail)
	}
	return fmt.Sprintf("TemplateElement(raw=%q,tail=%t)", t.Raw, t.Tail)
}

// ArrayLiteral represents array literals [a, b, ...c].
type ArrayLiteral struct {
	BaseNode
	Elements []Expression // nil entries indicate elided elements (holes).
}

func NewArrayLiteral(elements []Expression, loc Location) *ArrayLiteral {
	return &ArrayLiteral{BaseNode: NewBaseNode(ArrayLiteralKind, loc), Elements: elements}
}

func (a *ArrayLiteral) node()       {}
func (a *ArrayLiteral) expression() {}
func (a *ArrayLiteral) literal()    {}
func (a *ArrayLiteral) String() string {
	return "ArrayLiteral"
}

// PropertyKind identifies the semantics of an object property.
type PropertyKind string

const (
	PropertyInit   PropertyKind = "init"
	PropertyGet    PropertyKind = "get"
	PropertySet    PropertyKind = "set"
	PropertyMethod PropertyKind = "method"
)

// Property represents an entry in an object literal.
type Property interface {
	Node
	property()
}

// ObjectProperty represents key/value pairs inside object literals.
type ObjectProperty struct {
	BaseNode
	Key       Expression
	Value     Expression
	PropKind  PropertyKind
	Computed  bool
	Shorthand bool
	Method    bool
}

func NewObjectProperty(key, value Expression, kind PropertyKind, computed, shorthand, method bool, loc Location) *ObjectProperty {
	return &ObjectProperty{
		BaseNode:  NewBaseNode(ObjectPropertyKind, loc),
		Key:       key,
		Value:     value,
		PropKind:  kind,
		Computed:  computed,
		Shorthand: shorthand,
		Method:    method,
	}
}

func (p *ObjectProperty) node()     {}
func (p *ObjectProperty) property() {}
func (p *ObjectProperty) String() string {
	return fmt.Sprintf("ObjectProperty(kind=%s)", p.PropKind)
}

// SpreadElement represents ...expr within arrays or objects.
type SpreadElement struct {
	BaseNode
	Argument Expression
}

func NewSpreadElement(arg Expression, loc Location) *SpreadElement {
	return &SpreadElement{BaseNode: NewBaseNode(SpreadElementKind, loc), Argument: arg}
}

func (s *SpreadElement) node()       {}
func (s *SpreadElement) expression() {}
func (s *SpreadElement) property()   {}
func (s *SpreadElement) String() string {
	return "SpreadElement"
}

// ObjectLiteral aggregates multiple properties.
type ObjectLiteral struct {
	BaseNode
	Properties []Property
}

func NewObjectLiteral(props []Property, loc Location) *ObjectLiteral {
	return &ObjectLiteral{BaseNode: NewBaseNode(ObjectLiteralKind, loc), Properties: props}
}

func (o *ObjectLiteral) node()       {}
func (o *ObjectLiteral) expression() {}
func (o *ObjectLiteral) literal()    {}
func (o *ObjectLiteral) String() string {
	return "ObjectLiteral"
}

var (
	_ Literal    = (*NumberLiteral)(nil)
	_ Literal    = (*StringLiteral)(nil)
	_ Literal    = (*BooleanLiteral)(nil)
	_ Literal    = (*NullLiteral)(nil)
	_ Literal    = (*RegExpLiteral)(nil)
	_ Literal    = (*TemplateLiteral)(nil)
	_ Literal    = (*ArrayLiteral)(nil)
	_ Literal    = (*ObjectLiteral)(nil)
	_ Property   = (*ObjectProperty)(nil)
	_ Property   = (*SpreadElement)(nil)
	_ Expression = (*SpreadElement)(nil)
)
