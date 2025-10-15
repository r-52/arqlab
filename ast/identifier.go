package ast

import "fmt"

const (
	IdentifierKind     NodeKind = "Identifier"
	ThisExpressionKind NodeKind = "ThisExpression"
	SuperKind          NodeKind = "Super"
	MetaPropertyKind   NodeKind = "MetaProperty"
)

// Identifier represents an ECMAScript IdentifierName token used in expressions and bindings.
type Identifier struct {
	BaseNode
	Name string
}

func NewIdentifier(name string, loc Location) *Identifier {
	return &Identifier{BaseNode: NewBaseNode(IdentifierKind, loc), Name: name}
}

func (i *Identifier) node()       {}
func (i *Identifier) expression() {}
func (i *Identifier) pattern()    {}

func (i *Identifier) String() string {
	return fmt.Sprintf("Identifier(%s)", i.Name)
}

// ThisExpression models usage of the `this` keyword.
type ThisExpression struct {
	BaseNode
}

func NewThisExpression(loc Location) *ThisExpression {
	return &ThisExpression{BaseNode: NewBaseNode(ThisExpressionKind, loc)}
}

func (t *ThisExpression) node()       {}
func (t *ThisExpression) expression() {}
func (t *ThisExpression) String() string {
	return "ThisExpression"
}

// Super represents the `super` keyword in member access and calls.
type Super struct {
	BaseNode
}

func NewSuper(loc Location) *Super {
	return &Super{BaseNode: NewBaseNode(SuperKind, loc)}
}

func (s *Super) node()       {}
func (s *Super) expression() {}
func (s *Super) String() string {
	return "Super"
}

// MetaProperty represents constructs such as `new.target`.
type MetaProperty struct {
	BaseNode
	Meta     *Identifier
	Property *Identifier
}

func NewMetaProperty(meta, property *Identifier, loc Location) *MetaProperty {
	return &MetaProperty{BaseNode: NewBaseNode(MetaPropertyKind, loc), Meta: meta, Property: property}
}

func (m *MetaProperty) node()       {}
func (m *MetaProperty) expression() {}
func (m *MetaProperty) String() string {
	return "MetaProperty"
}
