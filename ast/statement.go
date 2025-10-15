package ast

const (
	ProgramKind             NodeKind = "Program"
	BlockStatementKind      NodeKind = "BlockStatement"
	ExpressionStatementKind NodeKind = "ExpressionStatement"
	EmptyStatementKind      NodeKind = "EmptyStatement"
	DebuggerStatementKind   NodeKind = "DebuggerStatement"
	ReturnStatementKind     NodeKind = "ReturnStatement"
	BreakStatementKind      NodeKind = "BreakStatement"
	ContinueStatementKind   NodeKind = "ContinueStatement"
	ThrowStatementKind      NodeKind = "ThrowStatement"
	IfStatementKind         NodeKind = "IfStatement"
	SwitchStatementKind     NodeKind = "SwitchStatement"
	SwitchCaseKind          NodeKind = "SwitchCase"
	WhileStatementKind      NodeKind = "WhileStatement"
	DoWhileStatementKind    NodeKind = "DoWhileStatement"
	ForStatementKind        NodeKind = "ForStatement"
	ForInStatementKind      NodeKind = "ForInStatement"
	ForOfStatementKind      NodeKind = "ForOfStatement"
	WithStatementKind       NodeKind = "WithStatement"
	LabeledStatementKind    NodeKind = "LabeledStatement"
	TryStatementKind        NodeKind = "TryStatement"
	CatchClauseKind         NodeKind = "CatchClause"
	VariableDeclarationKind NodeKind = "VariableDeclaration"
	VariableDeclaratorKind  NodeKind = "VariableDeclarator"
	FunctionDeclarationKind NodeKind = "FunctionDeclaration"
)

// SourceType differentiates between classic scripts and ECMAScript modules.
type SourceType string

const (
	SourceTypeScript SourceType = "script"
	SourceTypeModule SourceType = "module"
)

// Program represents the root of the AST.
type Program struct {
	BaseNode
	Body       []Statement
	SourceType SourceType
}

func NewProgram(body []Statement, sourceType SourceType, loc Location) *Program {
	return &Program{BaseNode: NewBaseNode(ProgramKind, loc), Body: body, SourceType: sourceType}
}

func (p *Program) node()          {}
func (p *Program) String() string { return "Program" }

// BlockStatement represents a list of statements enclosed in braces.
type BlockStatement struct {
	BaseNode
	Body []Statement
}

func NewBlockStatement(body []Statement, loc Location) *BlockStatement {
	return &BlockStatement{BaseNode: NewBaseNode(BlockStatementKind, loc), Body: body}
}

func (b *BlockStatement) node()      {}
func (b *BlockStatement) statement() {}
func (b *BlockStatement) String() string {
	return "BlockStatement"
}

// ExpressionStatement wraps an expression in statement position.
type ExpressionStatement struct {
	BaseNode
	Expression Expression
}

func NewExpressionStatement(expr Expression, loc Location) *ExpressionStatement {
	return &ExpressionStatement{BaseNode: NewBaseNode(ExpressionStatementKind, loc), Expression: expr}
}

func (e *ExpressionStatement) node()      {}
func (e *ExpressionStatement) statement() {}
func (e *ExpressionStatement) String() string {
	return "ExpressionStatement"
}

// EmptyStatement corresponds to a solitary semicolon.
type EmptyStatement struct {
	BaseNode
}

func NewEmptyStatement(loc Location) *EmptyStatement {
	return &EmptyStatement{BaseNode: NewBaseNode(EmptyStatementKind, loc)}
}

func (e *EmptyStatement) node()      {}
func (e *EmptyStatement) statement() {}
func (e *EmptyStatement) String() string {
	return "EmptyStatement"
}

// DebuggerStatement represents a debugger; usage.
type DebuggerStatement struct {
	BaseNode
}

func NewDebuggerStatement(loc Location) *DebuggerStatement {
	return &DebuggerStatement{BaseNode: NewBaseNode(DebuggerStatementKind, loc)}
}

func (d *DebuggerStatement) node()      {}
func (d *DebuggerStatement) statement() {}
func (d *DebuggerStatement) String() string {
	return "DebuggerStatement"
}

// ReturnStatement models return keyword usage.
type ReturnStatement struct {
	BaseNode
	Argument Expression // may be nil
}

func NewReturnStatement(argument Expression, loc Location) *ReturnStatement {
	return &ReturnStatement{BaseNode: NewBaseNode(ReturnStatementKind, loc), Argument: argument}
}

func (r *ReturnStatement) node()      {}
func (r *ReturnStatement) statement() {}
func (r *ReturnStatement) String() string {
	return "ReturnStatement"
}

// BreakStatement models break statements with optional labels.
type BreakStatement struct {
	BaseNode
	Label *Identifier
}

func NewBreakStatement(label *Identifier, loc Location) *BreakStatement {
	return &BreakStatement{BaseNode: NewBaseNode(BreakStatementKind, loc), Label: label}
}

func (b *BreakStatement) node()      {}
func (b *BreakStatement) statement() {}
func (b *BreakStatement) String() string {
	return "BreakStatement"
}

// ContinueStatement models continue statements with optional labels.
type ContinueStatement struct {
	BaseNode
	Label *Identifier
}

func NewContinueStatement(label *Identifier, loc Location) *ContinueStatement {
	return &ContinueStatement{BaseNode: NewBaseNode(ContinueStatementKind, loc), Label: label}
}

func (c *ContinueStatement) node()      {}
func (c *ContinueStatement) statement() {}
func (c *ContinueStatement) String() string {
	return "ContinueStatement"
}

// ThrowStatement models throwing of exceptions.
type ThrowStatement struct {
	BaseNode
	Argument Expression
}

func NewThrowStatement(argument Expression, loc Location) *ThrowStatement {
	return &ThrowStatement{BaseNode: NewBaseNode(ThrowStatementKind, loc), Argument: argument}
}

func (t *ThrowStatement) node()      {}
func (t *ThrowStatement) statement() {}
func (t *ThrowStatement) String() string {
	return "ThrowStatement"
}

// IfStatement represents conditional branching.
type IfStatement struct {
	BaseNode
	Test       Expression
	Consequent Statement
	Alternate  Statement // may be nil
}

func NewIfStatement(test Expression, consequent, alternate Statement, loc Location) *IfStatement {
	return &IfStatement{BaseNode: NewBaseNode(IfStatementKind, loc), Test: test, Consequent: consequent, Alternate: alternate}
}

func (i *IfStatement) node()      {}
func (i *IfStatement) statement() {}
func (i *IfStatement) String() string {
	return "IfStatement"
}

// WhileStatement models while loops.
type WhileStatement struct {
	BaseNode
	Test Expression
	Body Statement
}

func NewWhileStatement(test Expression, body Statement, loc Location) *WhileStatement {
	return &WhileStatement{BaseNode: NewBaseNode(WhileStatementKind, loc), Test: test, Body: body}
}

func (w *WhileStatement) node()      {}
func (w *WhileStatement) statement() {}
func (w *WhileStatement) String() string {
	return "WhileStatement"
}

// DoWhileStatement models do { } while (test) loops.
type DoWhileStatement struct {
	BaseNode
	Body Statement
	Test Expression
}

func NewDoWhileStatement(body Statement, test Expression, loc Location) *DoWhileStatement {
	return &DoWhileStatement{BaseNode: NewBaseNode(DoWhileStatementKind, loc), Body: body, Test: test}
}

func (d *DoWhileStatement) node()      {}
func (d *DoWhileStatement) statement() {}
func (d *DoWhileStatement) String() string {
	return "DoWhileStatement"
}

// ForStatement models classic for(initializer; test; update) loops.
type ForStatement struct {
	BaseNode
	Init   Node // nil, *VariableDeclaration, or Expression
	Test   Expression
	Update Expression
	Body   Statement
}

func NewForStatement(init Node, test, update Expression, body Statement, loc Location) *ForStatement {
	return &ForStatement{BaseNode: NewBaseNode(ForStatementKind, loc), Init: init, Test: test, Update: update, Body: body}
}

func (f *ForStatement) node()      {}
func (f *ForStatement) statement() {}
func (f *ForStatement) String() string {
	return "ForStatement"
}

// ForInStatement models for (lhs in rhs) loops.
type ForInStatement struct {
	BaseNode
	Left  Node // *VariableDeclaration or Pattern
	Right Expression
	Body  Statement
}

func NewForInStatement(left Node, right Expression, body Statement, loc Location) *ForInStatement {
	return &ForInStatement{BaseNode: NewBaseNode(ForInStatementKind, loc), Left: left, Right: right, Body: body}
}

func (f *ForInStatement) node()      {}
func (f *ForInStatement) statement() {}
func (f *ForInStatement) String() string {
	return "ForInStatement"
}

// ForOfStatement models for (lhs of rhs) loops.
type ForOfStatement struct {
	BaseNode
	Left  Node // *VariableDeclaration or Pattern
	Right Expression
	Body  Statement
	Await bool
}

func NewForOfStatement(left Node, right Expression, body Statement, await bool, loc Location) *ForOfStatement {
	return &ForOfStatement{BaseNode: NewBaseNode(ForOfStatementKind, loc), Left: left, Right: right, Body: body, Await: await}
}

func (f *ForOfStatement) node()      {}
func (f *ForOfStatement) statement() {}
func (f *ForOfStatement) String() string {
	return "ForOfStatement"
}

// SwitchCase represents a case clause in a switch statement.
type SwitchCase struct {
	BaseNode
	Test       Expression // nil for default case
	Consequent []Statement
}

func NewSwitchCase(test Expression, consequent []Statement, loc Location) *SwitchCase {
	return &SwitchCase{BaseNode: NewBaseNode(SwitchCaseKind, loc), Test: test, Consequent: consequent}
}

func (s *SwitchCase) node() {}
func (s *SwitchCase) String() string {
	return "SwitchCase"
}

// SwitchStatement represents switch constructs.
type SwitchStatement struct {
	BaseNode
	Discriminant Expression
	Cases        []*SwitchCase
}

func NewSwitchStatement(discriminant Expression, cases []*SwitchCase, loc Location) *SwitchStatement {
	return &SwitchStatement{BaseNode: NewBaseNode(SwitchStatementKind, loc), Discriminant: discriminant, Cases: cases}
}

func (s *SwitchStatement) node()      {}
func (s *SwitchStatement) statement() {}
func (s *SwitchStatement) String() string {
	return "SwitchStatement"
}

// WithStatement represents with (object) body constructs.
type WithStatement struct {
	BaseNode
	Object Expression
	Body   Statement
}

func NewWithStatement(object Expression, body Statement, loc Location) *WithStatement {
	return &WithStatement{BaseNode: NewBaseNode(WithStatementKind, loc), Object: object, Body: body}
}

func (w *WithStatement) node()      {}
func (w *WithStatement) statement() {}
func (w *WithStatement) String() string {
	return "WithStatement"
}

// LabeledStatement models label: statement.
type LabeledStatement struct {
	BaseNode
	Label *Identifier
	Body  Statement
}

func NewLabeledStatement(label *Identifier, body Statement, loc Location) *LabeledStatement {
	return &LabeledStatement{BaseNode: NewBaseNode(LabeledStatementKind, loc), Label: label, Body: body}
}

func (l *LabeledStatement) node()      {}
func (l *LabeledStatement) statement() {}
func (l *LabeledStatement) String() string {
	return "LabeledStatement"
}

// TryStatement represents try/catch/finally constructs.
type TryStatement struct {
	BaseNode
	Block     *BlockStatement
	Handler   *CatchClause
	Finalizer *BlockStatement
}

func NewTryStatement(block *BlockStatement, handler *CatchClause, finalizer *BlockStatement, loc Location) *TryStatement {
	return &TryStatement{BaseNode: NewBaseNode(TryStatementKind, loc), Block: block, Handler: handler, Finalizer: finalizer}
}

func (t *TryStatement) node()      {}
func (t *TryStatement) statement() {}
func (t *TryStatement) String() string {
	return "TryStatement"
}

// CatchClause represents the catch (binding) { } handler.
type CatchClause struct {
	BaseNode
	Param Pattern
	Body  *BlockStatement
}

func NewCatchClause(param Pattern, body *BlockStatement, loc Location) *CatchClause {
	return &CatchClause{BaseNode: NewBaseNode(CatchClauseKind, loc), Param: param, Body: body}
}

func (c *CatchClause) node() {}
func (c *CatchClause) String() string {
	return "CatchClause"
}

// VariableKind distinguishes var / let / const declarations.
type VariableKind string

const (
	VarKind   VariableKind = "var"
	LetKind   VariableKind = "let"
	ConstKind VariableKind = "const"
)

// VariableDeclarator couples an identifier/pattern with an initializer.
type VariableDeclarator struct {
	BaseNode
	ID   Pattern
	Init Expression
}

func NewVariableDeclarator(id Pattern, init Expression, loc Location) *VariableDeclarator {
	return &VariableDeclarator{BaseNode: NewBaseNode(VariableDeclaratorKind, loc), ID: id, Init: init}
}

func (v *VariableDeclarator) node() {}
func (v *VariableDeclarator) String() string {
	return "VariableDeclarator"
}

// VariableDeclaration models var/let/const declarations.
type VariableDeclaration struct {
	BaseNode
	Declarations []*VariableDeclarator
	DeclareKind  VariableKind
}

func NewVariableDeclaration(kind VariableKind, decls []*VariableDeclarator, loc Location) *VariableDeclaration {
	return &VariableDeclaration{BaseNode: NewBaseNode(VariableDeclarationKind, loc), Declarations: decls, DeclareKind: kind}
}

func (v *VariableDeclaration) node()        {}
func (v *VariableDeclaration) statement()   {}
func (v *VariableDeclaration) declaration() {}
func (v *VariableDeclaration) String() string {
	return "VariableDeclaration"
}

// FunctionDeclaration represents function keyword declarations.
type FunctionDeclaration struct {
	BaseNode
	ID        *Identifier
	Params    []Pattern
	Body      *BlockStatement
	Generator bool
}

func NewFunctionDeclaration(id *Identifier, params []Pattern, body *BlockStatement, generator bool, loc Location) *FunctionDeclaration {
	return &FunctionDeclaration{BaseNode: NewBaseNode(FunctionDeclarationKind, loc), ID: id, Params: params, Body: body, Generator: generator}
}

func (f *FunctionDeclaration) node()        {}
func (f *FunctionDeclaration) statement()   {}
func (f *FunctionDeclaration) declaration() {}
func (f *FunctionDeclaration) String() string {
	return "FunctionDeclaration"
}

var (
	_ Statement = (*BlockStatement)(nil)
	_ Statement = (*ExpressionStatement)(nil)
	_ Statement = (*EmptyStatement)(nil)
	_ Statement = (*DebuggerStatement)(nil)
	_ Statement = (*ReturnStatement)(nil)
	_ Statement = (*BreakStatement)(nil)
	_ Statement = (*ContinueStatement)(nil)
	_ Statement = (*ThrowStatement)(nil)
	_ Statement = (*IfStatement)(nil)
	_ Statement = (*WhileStatement)(nil)
	_ Statement = (*DoWhileStatement)(nil)
	_ Statement = (*ForStatement)(nil)
	_ Statement = (*ForInStatement)(nil)
	_ Statement = (*ForOfStatement)(nil)
	_ Statement = (*SwitchStatement)(nil)
	_ Statement = (*WithStatement)(nil)
	_ Statement = (*LabeledStatement)(nil)
	_ Statement = (*TryStatement)(nil)
	_ Statement = (*VariableDeclaration)(nil)
	_ Statement = (*FunctionDeclaration)(nil)

	_ Declaration = (*VariableDeclaration)(nil)
	_ Declaration = (*FunctionDeclaration)(nil)
)
