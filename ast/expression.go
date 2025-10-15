package ast

import "fmt"

const (
	MemberExpressionKind         NodeKind = "MemberExpression"
	CallExpressionKind           NodeKind = "CallExpression"
	NewExpressionKind            NodeKind = "NewExpression"
	TaggedTemplateExpressionKind NodeKind = "TaggedTemplateExpression"
	BinaryExpressionKind         NodeKind = "BinaryExpression"
	LogicalExpressionKind        NodeKind = "LogicalExpression"
	AssignmentExpressionKind     NodeKind = "AssignmentExpression"
	UnaryExpressionKind          NodeKind = "UnaryExpression"
	UpdateExpressionKind         NodeKind = "UpdateExpression"
	ConditionalExpressionKind    NodeKind = "ConditionalExpression"
	SequenceExpressionKind       NodeKind = "SequenceExpression"
)

// MemberExpression represents property access such as obj.prop or obj[expr].
type MemberExpression struct {
	BaseNode
	Object   Expression
	Property Expression
	Computed bool
}

func NewMemberExpression(object, property Expression, computed bool, loc Location) *MemberExpression {
	return &MemberExpression{
		BaseNode: NewBaseNode(MemberExpressionKind, loc),
		Object:   object,
		Property: property,
		Computed: computed,
	}
}

func (m *MemberExpression) node()       {}
func (m *MemberExpression) expression() {}
func (m *MemberExpression) String() string {
	if m.Computed {
		return "MemberExpression[computed]"
	}
	return "MemberExpression"
}

// CallExpression models function invocation.
type CallExpression struct {
	BaseNode
	Callee    Expression
	Arguments []Expression
}

func NewCallExpression(callee Expression, args []Expression, loc Location) *CallExpression {
	return &CallExpression{BaseNode: NewBaseNode(CallExpressionKind, loc), Callee: callee, Arguments: args}
}

func (c *CallExpression) node()       {}
func (c *CallExpression) expression() {}
func (c *CallExpression) String() string {
	return "CallExpression"
}

// NewExpression represents the `new` operator with arguments.
type NewExpression struct {
	BaseNode
	Callee    Expression
	Arguments []Expression
}

func NewNewExpression(callee Expression, args []Expression, loc Location) *NewExpression {
	return &NewExpression{BaseNode: NewBaseNode(NewExpressionKind, loc), Callee: callee, Arguments: args}
}

func (n *NewExpression) node()       {}
func (n *NewExpression) expression() {}
func (n *NewExpression) String() string {
	return "NewExpression"
}

// TaggedTemplateExpression models tag`template` constructs.
type TaggedTemplateExpression struct {
	BaseNode
	Tag   Expression
	Quasi *TemplateLiteral
}

func NewTaggedTemplateExpression(tag Expression, quasi *TemplateLiteral, loc Location) *TaggedTemplateExpression {
	return &TaggedTemplateExpression{BaseNode: NewBaseNode(TaggedTemplateExpressionKind, loc), Tag: tag, Quasi: quasi}
}

func (t *TaggedTemplateExpression) node()       {}
func (t *TaggedTemplateExpression) expression() {}
func (t *TaggedTemplateExpression) String() string {
	return "TaggedTemplateExpression"
}

// BinaryExpression covers binary operators such as +, -, *, etc.
type BinaryExpression struct {
	BaseNode
	Operator string
	Left     Expression
	Right    Expression
}

func NewBinaryExpression(operator string, left, right Expression, loc Location) *BinaryExpression {
	return &BinaryExpression{BaseNode: NewBaseNode(BinaryExpressionKind, loc), Operator: operator, Left: left, Right: right}
}

func (b *BinaryExpression) node()       {}
func (b *BinaryExpression) expression() {}
func (b *BinaryExpression) String() string {
	return fmt.Sprintf("BinaryExpression(%s)", b.Operator)
}

// LogicalExpression covers logical operators &&, ||.
type LogicalExpression struct {
	BaseNode
	Operator string
	Left     Expression
	Right    Expression
}

func NewLogicalExpression(operator string, left, right Expression, loc Location) *LogicalExpression {
	return &LogicalExpression{BaseNode: NewBaseNode(LogicalExpressionKind, loc), Operator: operator, Left: left, Right: right}
}

func (l *LogicalExpression) node()       {}
func (l *LogicalExpression) expression() {}
func (l *LogicalExpression) String() string {
	return fmt.Sprintf("LogicalExpression(%s)", l.Operator)
}

// AssignmentExpression models compound assignment forms.
type AssignmentExpression struct {
	BaseNode
	Operator string
	Left     Expression
	Right    Expression
}

func NewAssignmentExpression(operator string, left, right Expression, loc Location) *AssignmentExpression {
	return &AssignmentExpression{BaseNode: NewBaseNode(AssignmentExpressionKind, loc), Operator: operator, Left: left, Right: right}
}

func (a *AssignmentExpression) node()       {}
func (a *AssignmentExpression) expression() {}
func (a *AssignmentExpression) String() string {
	return fmt.Sprintf("AssignmentExpression(%s)", a.Operator)
}

// UnaryExpression models unary operators like typeof, void, delete, -.
type UnaryExpression struct {
	BaseNode
	Operator string
	Argument Expression
	Prefix   bool
}

func NewUnaryExpression(operator string, argument Expression, prefix bool, loc Location) *UnaryExpression {
	return &UnaryExpression{BaseNode: NewBaseNode(UnaryExpressionKind, loc), Operator: operator, Argument: argument, Prefix: prefix}
}

func (u *UnaryExpression) node()       {}
func (u *UnaryExpression) expression() {}
func (u *UnaryExpression) String() string {
	return fmt.Sprintf("UnaryExpression(%s)", u.Operator)
}

// UpdateExpression models ++/-- in prefix or postfix form.
type UpdateExpression struct {
	BaseNode
	Operator string
	Argument Expression
	Prefix   bool
}

func NewUpdateExpression(operator string, argument Expression, prefix bool, loc Location) *UpdateExpression {
	return &UpdateExpression{BaseNode: NewBaseNode(UpdateExpressionKind, loc), Operator: operator, Argument: argument, Prefix: prefix}
}

func (u *UpdateExpression) node()       {}
func (u *UpdateExpression) expression() {}
func (u *UpdateExpression) String() string {
	return fmt.Sprintf("UpdateExpression(%s)", u.Operator)
}

// ConditionalExpression models ternary expressions test ? consequent : alternate.
type ConditionalExpression struct {
	BaseNode
	Test       Expression
	Consequent Expression
	Alternate  Expression
}

func NewConditionalExpression(test, consequent, alternate Expression, loc Location) *ConditionalExpression {
	return &ConditionalExpression{BaseNode: NewBaseNode(ConditionalExpressionKind, loc), Test: test, Consequent: consequent, Alternate: alternate}
}

func (c *ConditionalExpression) node()       {}
func (c *ConditionalExpression) expression() {}
func (c *ConditionalExpression) String() string {
	return "ConditionalExpression"
}

// SequenceExpression represents comma-separated expressions evaluated left-to-right.
type SequenceExpression struct {
	BaseNode
	Expressions []Expression
}

func NewSequenceExpression(exprs []Expression, loc Location) *SequenceExpression {
	return &SequenceExpression{BaseNode: NewBaseNode(SequenceExpressionKind, loc), Expressions: exprs}
}

func (s *SequenceExpression) node()       {}
func (s *SequenceExpression) expression() {}
func (s *SequenceExpression) String() string {
	return "SequenceExpression"
}
