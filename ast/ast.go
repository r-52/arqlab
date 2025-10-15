package ast

import (
	"fmt"
	"strings"
)

// Position represents a precise offset within the source text.
type Position struct {
	Offset int // zero-based byte offset within the source
	Line   int // one-based source line number
	Column int // zero-based column count using UTF-16 code units per ECMAScript
}

// Location models the start and end positions of a node.
type Location struct {
	Start Position
	End   Position
}

// NodeKind enumerates the kinds of AST nodes.
type NodeKind string

// Node is the base interface implemented by all AST nodes.
type Node interface {
	Kind() NodeKind
	Loc() Location
	String() string
	node()
}

// Statement represents an executable statement node.
type Statement interface {
	Node
	statement()
}

// Expression represents an evaluatable expression node.
type Expression interface {
	Node
	expression()
}

// Pattern marks nodes that can appear where a binding pattern is required.
type Pattern interface {
	Node
	pattern()
}

// Declaration marks nodes that introduce bindings (e.g. functions, classes, variables).
type Declaration interface {
	Statement
	declaration()
}

// BaseNode provides reusable behaviour for AST nodes.
type BaseNode struct {
	kind NodeKind
	loc  Location
}

// NewBaseNode constructs a BaseNode from its kind and location.
func NewBaseNode(kind NodeKind, loc Location) BaseNode {
	return BaseNode{kind: kind, loc: loc}
}

// Kind returns the node kind discriminator.
func (n BaseNode) Kind() NodeKind { return n.kind }

// Loc returns the source location covered by the node.
func (n BaseNode) Loc() Location { return n.loc }

// Position returns the start position of the node (alias for Loc().Start).
func (n BaseNode) Position() Position { return n.loc.Start }

// End returns the end position of the node (alias for Loc().End).
func (n BaseNode) End() Position { return n.loc.End }

// SetLoc updates the location metadata.
func (n *BaseNode) SetLoc(loc Location) { n.loc = loc }

// SetKind updates the node kind discriminator (useful when reusing structs in builders).
func (n *BaseNode) SetKind(kind NodeKind) { n.kind = kind }

// Location utilities -------------------------------------------------------

// IsValid returns true when both endpoints have non-negative offsets.
func (l Location) IsValid() bool {
	return l.Start.Offset >= 0 && l.End.Offset >= l.Start.Offset
}

// Span reports the width of the location in bytes.
func (l Location) Span() int {
	if l.End.Offset < l.Start.Offset {
		return 0
	}
	return l.End.Offset - l.Start.Offset
}

// String renders the location in "start-end" format (line:column pairs).
func (l Location) String() string {
	return fmt.Sprintf("%s-%s", l.Start, l.End)
}

// String renders a position as "line:column" (1-based columns for readability).
func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column+1)
}

// Visitor enables node traversal. Implementations can return an error to abort traversal.
type Visitor interface {
	Visit(Node) error
}

// Walk performs a depth-first traversal using the provided walker function.
// Concrete node types should provide specific child visitation logic via helper methods.
func Walk(v Visitor, n Node) error {
	if n == nil || v == nil {
		return nil
	}
	return v.Visit(n)
}

// DebugString prints a compact textual representation of a node tree using a supplied writer function.
func DebugString(n Node, children func(Node) []Node) string {
	if n == nil {
		return "<nil>"
	}
	var b strings.Builder
	var walk func(Node, int)
	walk = func(cur Node, depth int) {
		if cur == nil {
			return
		}
		b.WriteString(strings.Repeat("  ", depth))
		b.WriteString(string(cur.Kind()))
		b.WriteString(" ")
		b.WriteString(cur.Loc().String())
		b.WriteByte('\n')
		for _, child := range children(cur) {
			walk(child, depth+1)
		}
	}
	walk(n, 0)
	return b.String()
}
