package ast

const (
	ArrayPatternKind      NodeKind = "ArrayPattern"
	ObjectPatternKind     NodeKind = "ObjectPattern"
	AssignmentPatternKind NodeKind = "AssignmentPattern"
	RestElementKind       NodeKind = "RestElement"
	ObjectPatternPropKind NodeKind = "ObjectPatternProperty"
)

// PatternList captures a slice of patterns while preserving holes for elided elements.
type PatternList []Pattern

// ArrayPattern models ECMAScript array destructuring patterns and parameter lists.
type ArrayPattern struct {
	BaseNode
	Elements PatternList
	Rest     *RestElement
}

func NewArrayPattern(elements PatternList, rest *RestElement, loc Location) *ArrayPattern {
	return &ArrayPattern{BaseNode: NewBaseNode(ArrayPatternKind, loc), Elements: elements, Rest: rest}
}

func (a *ArrayPattern) node()    {}
func (a *ArrayPattern) pattern() {}
func (a *ArrayPattern) String() string {
	return "ArrayPattern"
}

// ObjectPatternProperty associates a property key with a binding target.
type ObjectPatternProperty struct {
	BaseNode
	Key       Expression
	Value     Pattern
	Computed  bool
	Shorthand bool
}

func NewObjectPatternProperty(key Expression, value Pattern, computed, shorthand bool, loc Location) *ObjectPatternProperty {
	return &ObjectPatternProperty{BaseNode: NewBaseNode(ObjectPatternPropKind, loc), Key: key, Value: value, Computed: computed, Shorthand: shorthand}
}

func (p *ObjectPatternProperty) node() {}
func (p *ObjectPatternProperty) String() string {
	return "ObjectPatternProperty"
}

// ObjectPattern represents object destructuring patterns.
type ObjectPattern struct {
	BaseNode
	Properties []*ObjectPatternProperty
	Rest       *RestElement
}

func NewObjectPattern(props []*ObjectPatternProperty, rest *RestElement, loc Location) *ObjectPattern {
	return &ObjectPattern{BaseNode: NewBaseNode(ObjectPatternKind, loc), Properties: props, Rest: rest}
}

func (o *ObjectPattern) node()    {}
func (o *ObjectPattern) pattern() {}
func (o *ObjectPattern) String() string {
	return "ObjectPattern"
}

// AssignmentPattern models default values in destructuring (e.g., [a = 1]).
type AssignmentPattern struct {
	BaseNode
	Left  Pattern
	Right Expression
}

func NewAssignmentPattern(left Pattern, right Expression, loc Location) *AssignmentPattern {
	return &AssignmentPattern{BaseNode: NewBaseNode(AssignmentPatternKind, loc), Left: left, Right: right}
}

func (a *AssignmentPattern) node()    {}
func (a *AssignmentPattern) pattern() {}
func (a *AssignmentPattern) String() string {
	return "AssignmentPattern"
}

// RestElement represents ...target in patterns.
type RestElement struct {
	BaseNode
	Argument Pattern
}

func NewRestElement(argument Pattern, loc Location) *RestElement {
	return &RestElement{BaseNode: NewBaseNode(RestElementKind, loc), Argument: argument}
}

func (r *RestElement) node()    {}
func (r *RestElement) pattern() {}
func (r *RestElement) String() string {
	return "RestElement"
}

var (
	_ Pattern = (*ArrayPattern)(nil)
	_ Pattern = (*ObjectPattern)(nil)
	_ Pattern = (*AssignmentPattern)(nil)
	_ Pattern = (*RestElement)(nil)
)
