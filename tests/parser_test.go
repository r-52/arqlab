package tests

import (
	"testing"

	"es6-interpreter/ast"
	"es6-interpreter/parser"
)

func parseProgram(t *testing.T, src string) *ast.Program {
	t.Helper()
	p := parser.New(src)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	return program
}

func TestParseExpressionStatement(t *testing.T) {
	prog := parseProgram(t, "1 + 2 * 3;")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	exprStmt, ok := prog.Body[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", prog.Body[0])
	}

	binary, ok := exprStmt.Expression.(*ast.BinaryExpression)
	if !ok {
		t.Fatalf("expected BinaryExpression, got %T", exprStmt.Expression)
	}

	if binary.Operator != "+" {
		t.Fatalf("unexpected operator: got %s", binary.Operator)
	}

	leftNum, ok := binary.Left.(*ast.NumberLiteral)
	if !ok || leftNum.Value != "1" {
		t.Fatalf("unexpected left operand: %#v", binary.Left)
	}

	rightBinary, ok := binary.Right.(*ast.BinaryExpression)
	if !ok || rightBinary.Operator != "*" {
		t.Fatalf("unexpected right operand: %#v", binary.Right)
	}

	rightLeft, ok := rightBinary.Left.(*ast.NumberLiteral)
	if !ok || rightLeft.Value != "2" {
		t.Fatalf("unexpected nested left operand: %#v", rightBinary.Left)
	}

	rightRight, ok := rightBinary.Right.(*ast.NumberLiteral)
	if !ok || rightRight.Value != "3" {
		t.Fatalf("unexpected nested right operand: %#v", rightBinary.Right)
	}
}

func TestParseNewExpressionSimple(t *testing.T) {
	prog := parseProgram(t, "new Foo(1, 2);")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	exprStmt, ok := prog.Body[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", prog.Body[0])
	}

	newExpr, ok := exprStmt.Expression.(*ast.NewExpression)
	if !ok {
		t.Fatalf("expected NewExpression, got %T", exprStmt.Expression)
	}

	callee, ok := newExpr.Callee.(*ast.Identifier)
	if !ok || callee.Name != "Foo" {
		t.Fatalf("unexpected callee: %#v", newExpr.Callee)
	}

	if len(newExpr.Arguments) != 2 {
		t.Fatalf("expected 2 arguments, got %d", len(newExpr.Arguments))
	}

	for i, want := range []string{"1", "2"} {
		num, ok := newExpr.Arguments[i].(*ast.NumberLiteral)
		if !ok || num.Value != want {
			t.Fatalf("argument %d mismatch: %#v", i, newExpr.Arguments[i])
		}
	}
}

func TestParseNewExpressionChainedCall(t *testing.T) {
	prog := parseProgram(t, "new Foo()();")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	exprStmt, ok := prog.Body[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", prog.Body[0])
	}

	call, ok := exprStmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expected CallExpression, got %T", exprStmt.Expression)
	}

	innerNew, ok := call.Callee.(*ast.NewExpression)
	if !ok {
		t.Fatalf("expected inner callee to be NewExpression, got %T", call.Callee)
	}

	if len(innerNew.Arguments) != 0 {
		t.Fatalf("expected new Foo() to have no arguments, got %d", len(innerNew.Arguments))
	}

	if _, ok := innerNew.Callee.(*ast.Identifier); !ok {
		t.Fatalf("expected identifier callee in new expression, got %T", innerNew.Callee)
	}
}

func TestParseNewExpressionMemberAccess(t *testing.T) {
	prog := parseProgram(t, "new Foo().bar();")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	exprStmt, ok := prog.Body[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", prog.Body[0])
	}

	call, ok := exprStmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("unexpected expression: %#v", exprStmt.Expression)
	}

	member, ok := call.Callee.(*ast.MemberExpression)
	if !ok {
		t.Fatalf("expected member callee, got %T", call.Callee)
	}

	innerNew, ok := member.Object.(*ast.NewExpression)
	if !ok {
		t.Fatalf("expected member object to be NewExpression, got %T", member.Object)
	}

	if _, ok := member.Property.(*ast.Identifier); !ok {
		t.Fatalf("expected identifier property, got %T", member.Property)
	}

	if len(innerNew.Arguments) != 0 {
		t.Fatalf("expected no arguments in new Foo(), got %d", len(innerNew.Arguments))
	}
}

func TestParseVariableDeclaration(t *testing.T) {
	prog := parseProgram(t, "const answer = 42;")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	decl, ok := prog.Body[0].(*ast.VariableDeclaration)
	if !ok {
		t.Fatalf("expected VariableDeclaration, got %T", prog.Body[0])
	}

	if decl.DeclareKind != ast.ConstKind {
		t.Fatalf("expected const kind, got %s", decl.DeclareKind)
	}

	if len(decl.Declarations) != 1 {
		t.Fatalf("expected 1 declarator, got %d", len(decl.Declarations))
	}

	declarator := decl.Declarations[0]
	ident, ok := declarator.ID.(*ast.Identifier)
	if !ok || ident.Name != "answer" {
		t.Fatalf("unexpected identifier: %#v", declarator.ID)
	}

	num, ok := declarator.Init.(*ast.NumberLiteral)
	if !ok || num.Value != "42" {
		t.Fatalf("unexpected initializer: %#v", declarator.Init)
	}
}

func TestParseLetWithoutInitializer(t *testing.T) {
	prog := parseProgram(t, "let value;")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	decl, ok := prog.Body[0].(*ast.VariableDeclaration)
	if !ok {
		t.Fatalf("expected VariableDeclaration, got %T", prog.Body[0])
	}

	if decl.DeclareKind != ast.LetKind {
		t.Fatalf("expected let kind, got %s", decl.DeclareKind)
	}

	if len(decl.Declarations) != 1 {
		t.Fatalf("expected 1 declarator, got %d", len(decl.Declarations))
	}

	declarator := decl.Declarations[0]
	ident, ok := declarator.ID.(*ast.Identifier)
	if !ok || ident.Name != "value" {
		t.Fatalf("unexpected identifier: %#v", declarator.ID)
	}

	if declarator.Init != nil {
		t.Fatalf("expected nil initializer, got %#v", declarator.Init)
	}
}

func TestParseArrayDestructuring(t *testing.T) {
	prog := parseProgram(t, "const [a, b = 2, ...rest] = source;")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	decl, ok := prog.Body[0].(*ast.VariableDeclaration)
	if !ok {
		t.Fatalf("expected VariableDeclaration, got %T", prog.Body[0])
	}

	if len(decl.Declarations) != 1 {
		t.Fatalf("expected 1 declarator, got %d", len(decl.Declarations))
	}

	declarator := decl.Declarations[0]
	arrayPat, ok := declarator.ID.(*ast.ArrayPattern)
	if !ok {
		t.Fatalf("expected ArrayPattern, got %T", declarator.ID)
	}

	if len(arrayPat.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(arrayPat.Elements))
	}

	first, ok := arrayPat.Elements[0].(*ast.Identifier)
	if !ok || first.Name != "a" {
		t.Fatalf("unexpected first element: %#v", arrayPat.Elements[0])
	}

	assign, ok := arrayPat.Elements[1].(*ast.AssignmentPattern)
	if !ok {
		t.Fatalf("second element should be AssignmentPattern, got %T", arrayPat.Elements[1])
	}

	leftIdent, ok := assign.Left.(*ast.Identifier)
	if !ok || leftIdent.Name != "b" {
		t.Fatalf("unexpected assignment left: %#v", assign.Left)
	}

	rightNum, ok := assign.Right.(*ast.NumberLiteral)
	if !ok || rightNum.Value != "2" {
		t.Fatalf("unexpected assignment right: %#v", assign.Right)
	}

	if arrayPat.Rest == nil {
		t.Fatalf("expected rest element in array pattern")
	}

	restIdent, ok := arrayPat.Rest.Argument.(*ast.Identifier)
	if !ok || restIdent.Name != "rest" {
		t.Fatalf("unexpected rest argument: %#v", arrayPat.Rest.Argument)
	}

	initIdent, ok := declarator.Init.(*ast.Identifier)
	if !ok || initIdent.Name != "source" {
		t.Fatalf("unexpected initializer: %#v", declarator.Init)
	}
}

func TestParseObjectDestructuring(t *testing.T) {
	prog := parseProgram(t, "let {a, b: alias = 3, ...rest} = obj;")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	decl, ok := prog.Body[0].(*ast.VariableDeclaration)
	if !ok {
		t.Fatalf("expected VariableDeclaration, got %T", prog.Body[0])
	}

	if len(decl.Declarations) != 1 {
		t.Fatalf("expected 1 declarator, got %d", len(decl.Declarations))
	}

	declarator := decl.Declarations[0]
	objPat, ok := declarator.ID.(*ast.ObjectPattern)
	if !ok {
		t.Fatalf("expected ObjectPattern, got %T", declarator.ID)
	}

	if len(objPat.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(objPat.Properties))
	}

	propA := objPat.Properties[0]
	if !propA.Shorthand {
		t.Fatalf("first property should be shorthand")
	}
	if _, ok := propA.Value.(*ast.Identifier); !ok {
		t.Fatalf("first property value should be Identifier, got %T", propA.Value)
	}

	propB := objPat.Properties[1]
	if propB.Shorthand {
		t.Fatalf("second property should not be shorthand")
	}
	assign, ok := propB.Value.(*ast.AssignmentPattern)
	if !ok {
		t.Fatalf("second property value should be AssignmentPattern, got %T", propB.Value)
	}
	left, ok := assign.Left.(*ast.Identifier)
	if !ok || left.Name != "alias" {
		t.Fatalf("unexpected assignment left: %#v", assign.Left)
	}
	right, ok := assign.Right.(*ast.NumberLiteral)
	if !ok || right.Value != "3" {
		t.Fatalf("unexpected assignment right: %#v", assign.Right)
	}

	if objPat.Rest == nil {
		t.Fatalf("expected rest element in object pattern")
	}
	restIdent, ok := objPat.Rest.Argument.(*ast.Identifier)
	if !ok || restIdent.Name != "rest" {
		t.Fatalf("unexpected rest argument: %#v", objPat.Rest.Argument)
	}

	initIdent, ok := declarator.Init.(*ast.Identifier)
	if !ok || initIdent.Name != "obj" {
		t.Fatalf("unexpected initializer: %#v", declarator.Init)
	}
}

func TestParseBlockStatement(t *testing.T) {
	prog := parseProgram(t, "{ let x = 1; x; }")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	block, ok := prog.Body[0].(*ast.BlockStatement)
	if !ok {
		t.Fatalf("expected BlockStatement, got %T", prog.Body[0])
	}

	if len(block.Body) != 2 {
		t.Fatalf("expected 2 statements in block, got %d", len(block.Body))
	}

	if _, ok := block.Body[0].(*ast.VariableDeclaration); !ok {
		t.Fatalf("expected first block statement to be VariableDeclaration, got %T", block.Body[0])
	}

	exprStmt, ok := block.Body[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected second block statement to be ExpressionStatement, got %T", block.Body[1])
	}

	ident, ok := exprStmt.Expression.(*ast.Identifier)
	if !ok || ident.Name != "x" {
		t.Fatalf("unexpected expression in block: %#v", exprStmt.Expression)
	}
}

func TestParseReturnStatement(t *testing.T) {
	prog := parseProgram(t, "return answer;")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	ret, ok := prog.Body[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("expected ReturnStatement, got %T", prog.Body[0])
	}

	arg, ok := ret.Argument.(*ast.Identifier)
	if !ok || arg.Name != "answer" {
		t.Fatalf("unexpected return argument: %#v", ret.Argument)
	}
}

func TestParseIfElseStatement(t *testing.T) {
	prog := parseProgram(t, "if (flag) foo(); else bar();")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	ifStmt, ok := prog.Body[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("expected IfStatement, got %T", prog.Body[0])
	}

	testIdent, ok := ifStmt.Test.(*ast.Identifier)
	if !ok || testIdent.Name != "flag" {
		t.Fatalf("unexpected test expression: %#v", ifStmt.Test)
	}

	consExpr, ok := ifStmt.Consequent.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected consequent ExpressionStatement, got %T", ifStmt.Consequent)
	}

	consCall, ok := consExpr.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expected consequent CallExpression, got %T", consExpr.Expression)
	}

	consCallee, ok := consCall.Callee.(*ast.Identifier)
	if !ok || consCallee.Name != "foo" {
		t.Fatalf("unexpected consequent callee: %#v", consCall.Callee)
	}

	altExpr, ok := ifStmt.Alternate.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected alternate ExpressionStatement, got %T", ifStmt.Alternate)
	}

	altCall, ok := altExpr.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expected alternate CallExpression, got %T", altExpr.Expression)
	}

	altCallee, ok := altCall.Callee.(*ast.Identifier)
	if !ok || altCallee.Name != "bar" {
		t.Fatalf("unexpected alternate callee: %#v", altCall.Callee)
	}
}

func TestParseWhileStatement(t *testing.T) {
	prog := parseProgram(t, "while (active) { active--; }")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	whileStmt, ok := prog.Body[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("expected WhileStatement, got %T", prog.Body[0])
	}

	testIdent, ok := whileStmt.Test.(*ast.Identifier)
	if !ok || testIdent.Name != "active" {
		t.Fatalf("unexpected while test: %#v", whileStmt.Test)
	}

	body, ok := whileStmt.Body.(*ast.BlockStatement)
	if !ok {
		t.Fatalf("expected while body BlockStatement, got %T", whileStmt.Body)
	}

	if len(body.Body) != 1 {
		t.Fatalf("expected 1 statement in while body, got %d", len(body.Body))
	}

	updateExpr, ok := body.Body[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected expression statement in while body, got %T", body.Body[0])
	}

	update, ok := updateExpr.Expression.(*ast.UpdateExpression)
	if !ok || update.Operator != "--" {
		t.Fatalf("unexpected update expression: %#v", updateExpr.Expression)
	}
}

func TestParseDoWhileStatement(t *testing.T) {
	prog := parseProgram(t, "do foo(); while (condition);")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	doStmt, ok := prog.Body[0].(*ast.DoWhileStatement)
	if !ok {
		t.Fatalf("expected DoWhileStatement, got %T", prog.Body[0])
	}

	body, ok := doStmt.Body.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected do body ExpressionStatement, got %T", doStmt.Body)
	}

	call, ok := body.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expected call expression in body, got %T", body.Expression)
	}

	callee, ok := call.Callee.(*ast.Identifier)
	if !ok || callee.Name != "foo" {
		t.Fatalf("unexpected body callee: %#v", call.Callee)
	}

	test, ok := doStmt.Test.(*ast.Identifier)
	if !ok || test.Name != "condition" {
		t.Fatalf("unexpected do-while test: %#v", doStmt.Test)
	}
}

func TestParseForStatement(t *testing.T) {
	prog := parseProgram(t, "for (let i = 0; i < limit; i++) total += i;")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	forStmt, ok := prog.Body[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("expected ForStatement, got %T", prog.Body[0])
	}

	init, ok := forStmt.Init.(*ast.VariableDeclaration)
	if !ok {
		t.Fatalf("expected initializer VariableDeclaration, got %T", forStmt.Init)
	}

	if len(init.Declarations) != 1 {
		t.Fatalf("expected 1 initializer declarator, got %d", len(init.Declarations))
	}

	decl := init.Declarations[0]
	id, ok := decl.ID.(*ast.Identifier)
	if !ok || id.Name != "i" {
		t.Fatalf("unexpected initializer identifier: %#v", decl.ID)
	}

	if decl.Init == nil {
		t.Fatalf("expected initializer expression")
	}

	test, ok := forStmt.Test.(*ast.BinaryExpression)
	if !ok || test.Operator != "<" {
		t.Fatalf("unexpected test expression: %#v", forStmt.Test)
	}

	update, ok := forStmt.Update.(*ast.UpdateExpression)
	if !ok || update.Operator != "++" {
		t.Fatalf("unexpected update expression: %#v", forStmt.Update)
	}

	body, ok := forStmt.Body.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected for body ExpressionStatement, got %T", forStmt.Body)
	}

	assign, ok := body.Expression.(*ast.AssignmentExpression)
	if !ok || assign.Operator != "+=" {
		t.Fatalf("unexpected body assignment: %#v", body.Expression)
	}
}

func TestParseBreakStatement(t *testing.T) {
	prog := parseProgram(t, "break done;")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	brk, ok := prog.Body[0].(*ast.BreakStatement)
	if !ok {
		t.Fatalf("expected BreakStatement, got %T", prog.Body[0])
	}

	if brk.Label == nil || brk.Label.Name != "done" {
		t.Fatalf("unexpected break label: %#v", brk.Label)
	}
}

func TestParseContinueStatement(t *testing.T) {
	prog := parseProgram(t, "continue;")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	cont, ok := prog.Body[0].(*ast.ContinueStatement)
	if !ok {
		t.Fatalf("expected ContinueStatement, got %T", prog.Body[0])
	}

	if cont.Label != nil {
		t.Fatalf("expected nil continue label, got %#v", cont.Label)
	}
}

func TestParseThrowStatement(t *testing.T) {
	prog := parseProgram(t, "throw error;")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	throwStmt, ok := prog.Body[0].(*ast.ThrowStatement)
	if !ok {
		t.Fatalf("expected ThrowStatement, got %T", prog.Body[0])
	}

	arg, ok := throwStmt.Argument.(*ast.Identifier)
	if !ok || arg.Name != "error" {
		t.Fatalf("unexpected throw argument: %#v", throwStmt.Argument)
	}
}

func TestParseDebuggerStatement(t *testing.T) {
	prog := parseProgram(t, "debugger;")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	if _, ok := prog.Body[0].(*ast.DebuggerStatement); !ok {
		t.Fatalf("expected DebuggerStatement, got %T", prog.Body[0])
	}
}

func TestParseSwitchStatement(t *testing.T) {
	prog := parseProgram(t, "switch (value) { case 1: foo(); break; default: bar(); }")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	switchStmt, ok := prog.Body[0].(*ast.SwitchStatement)
	if !ok {
		t.Fatalf("expected SwitchStatement, got %T", prog.Body[0])
	}

	disc, ok := switchStmt.Discriminant.(*ast.Identifier)
	if !ok || disc.Name != "value" {
		t.Fatalf("unexpected discriminant: %#v", switchStmt.Discriminant)
	}

	if len(switchStmt.Cases) != 2 {
		t.Fatalf("expected 2 cases, got %d", len(switchStmt.Cases))
	}

	first := switchStmt.Cases[0]
	caseTest, ok := first.Test.(*ast.NumberLiteral)
	if !ok || caseTest.Value != "1" {
		t.Fatalf("unexpected first case test: %#v", first.Test)
	}

	if len(first.Consequent) != 2 {
		t.Fatalf("expected 2 statements in first case, got %d", len(first.Consequent))
	}

	if _, ok := first.Consequent[0].(*ast.ExpressionStatement); !ok {
		t.Fatalf("expected call expression statement in first case, got %T", first.Consequent[0])
	}

	if _, ok := first.Consequent[1].(*ast.BreakStatement); !ok {
		t.Fatalf("expected break statement in first case, got %T", first.Consequent[1])
	}

	second := switchStmt.Cases[1]
	if second.Test != nil {
		t.Fatalf("expected default case to have nil test, got %#v", second.Test)
	}

	if len(second.Consequent) != 1 {
		t.Fatalf("expected 1 statement in default case, got %d", len(second.Consequent))
	}

	if _, ok := second.Consequent[0].(*ast.ExpressionStatement); !ok {
		t.Fatalf("expected expression statement in default case, got %T", second.Consequent[0])
	}
}

func TestParseWithStatement(t *testing.T) {
	prog := parseProgram(t, "with (ctx) { ctx.run(); }")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	withStmt, ok := prog.Body[0].(*ast.WithStatement)
	if !ok {
		t.Fatalf("expected WithStatement, got %T", prog.Body[0])
	}

	obj, ok := withStmt.Object.(*ast.Identifier)
	if !ok || obj.Name != "ctx" {
		t.Fatalf("unexpected with object: %#v", withStmt.Object)
	}

	body, ok := withStmt.Body.(*ast.BlockStatement)
	if !ok {
		t.Fatalf("expected with body BlockStatement, got %T", withStmt.Body)
	}

	if len(body.Body) != 1 {
		t.Fatalf("expected 1 statement in with body, got %d", len(body.Body))
	}
}

func TestParseLabeledStatement(t *testing.T) {
	prog := parseProgram(t, "loop: while (true) break loop;")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	labeled, ok := prog.Body[0].(*ast.LabeledStatement)
	if !ok {
		t.Fatalf("expected LabeledStatement, got %T", prog.Body[0])
	}

	if labeled.Label == nil || labeled.Label.Name != "loop" {
		t.Fatalf("unexpected label: %#v", labeled.Label)
	}

	body, ok := labeled.Body.(*ast.WhileStatement)
	if !ok {
		t.Fatalf("expected labeled body WhileStatement, got %T", labeled.Body)
	}

	brk, ok := body.Body.(*ast.BreakStatement)
	if !ok {
		t.Fatalf("expected break in loop body, got %T", body.Body)
	}

	if brk.Label == nil || brk.Label.Name != "loop" {
		t.Fatalf("unexpected break label: %#v", brk.Label)
	}
}

func TestParseTryCatchFinally(t *testing.T) {
	prog := parseProgram(t, "try { risky(); } catch (err) { handle(err); } finally { cleanup(); }")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	tryStmt, ok := prog.Body[0].(*ast.TryStatement)
	if !ok {
		t.Fatalf("expected TryStatement, got %T", prog.Body[0])
	}

	if tryStmt.Block == nil || len(tryStmt.Block.Body) != 1 {
		t.Fatalf("unexpected try block: %#v", tryStmt.Block)
	}

	if tryStmt.Handler == nil {
		t.Fatalf("expected catch handler")
	}

	param, ok := tryStmt.Handler.Param.(*ast.Identifier)
	if !ok || param.Name != "err" {
		t.Fatalf("unexpected catch parameter: %#v", tryStmt.Handler.Param)
	}

	if tryStmt.Handler.Body == nil || len(tryStmt.Handler.Body.Body) != 1 {
		t.Fatalf("unexpected catch body: %#v", tryStmt.Handler.Body)
	}

	if tryStmt.Finalizer == nil || len(tryStmt.Finalizer.Body) != 1 {
		t.Fatalf("unexpected finally block: %#v", tryStmt.Finalizer)
	}
}

func TestParseFunctionDeclaration(t *testing.T) {
	prog := parseProgram(t, "function greet(name, title = \"Dr\") { return name; }")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	fn, ok := prog.Body[0].(*ast.FunctionDeclaration)
	if !ok {
		t.Fatalf("expected FunctionDeclaration, got %T", prog.Body[0])
	}

	if fn.Generator {
		t.Fatalf("expected non-generator function")
	}

	if fn.ID == nil || fn.ID.Name != "greet" {
		t.Fatalf("unexpected function name: %#v", fn.ID)
	}

	if len(fn.Params) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(fn.Params))
	}

	if _, ok := fn.Params[0].(*ast.Identifier); !ok {
		t.Fatalf("expected first parameter Identifier, got %T", fn.Params[0])
	}

	assign, ok := fn.Params[1].(*ast.AssignmentPattern)
	if !ok {
		t.Fatalf("expected second parameter AssignmentPattern, got %T", fn.Params[1])
	}

	right, ok := assign.Right.(*ast.StringLiteral)
	if !ok || right.Value != "Dr" {
		t.Fatalf("unexpected default value: %#v", assign.Right)
	}

	if fn.Body == nil || len(fn.Body.Body) != 1 {
		t.Fatalf("unexpected function body: %#v", fn.Body)
	}
}

func TestParseGeneratorFunctionDeclaration(t *testing.T) {
	prog := parseProgram(t, "function* iterate(...items) { return items; }")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	fn, ok := prog.Body[0].(*ast.FunctionDeclaration)
	if !ok {
		t.Fatalf("expected FunctionDeclaration, got %T", prog.Body[0])
	}

	if !fn.Generator {
		t.Fatalf("expected generator function")
	}

	if len(fn.Params) != 1 {
		t.Fatalf("expected rest parameter only, got %d", len(fn.Params))
	}

	rest, ok := fn.Params[0].(*ast.RestElement)
	if !ok {
		t.Fatalf("expected RestElement, got %T", fn.Params[0])
	}

	if ident, ok := rest.Argument.(*ast.Identifier); !ok || ident.Name != "items" {
		t.Fatalf("unexpected rest argument: %#v", rest.Argument)
	}
}

func TestParseArrayLiteralExpression(t *testing.T) {
	prog := parseProgram(t, "const arr = [1,,2,...more];")

	decl, ok := prog.Body[0].(*ast.VariableDeclaration)
	if !ok {
		t.Fatalf("expected VariableDeclaration, got %T", prog.Body[0])
	}

	arrInit, ok := decl.Declarations[0].Init.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expected ArrayLiteral, got %T", decl.Declarations[0].Init)
	}

	if len(arrInit.Elements) != 4 {
		t.Fatalf("expected 4 elements, got %d", len(arrInit.Elements))
	}

	if _, ok := arrInit.Elements[0].(*ast.NumberLiteral); !ok {
		t.Fatalf("expected first element NumberLiteral, got %T", arrInit.Elements[0])
	}

	if arrInit.Elements[1] != nil {
		t.Fatalf("expected hole in second position, got %#v", arrInit.Elements[1])
	}

	if num, ok := arrInit.Elements[2].(*ast.NumberLiteral); !ok || num.Value != "2" {
		t.Fatalf("unexpected third element: %#v", arrInit.Elements[2])
	}

	spread, ok := arrInit.Elements[3].(*ast.SpreadElement)
	if !ok {
		t.Fatalf("expected spread element, got %T", arrInit.Elements[3])
	}

	if ident, ok := spread.Argument.(*ast.Identifier); !ok || ident.Name != "more" {
		t.Fatalf("unexpected spread argument: %#v", spread.Argument)
	}
}

func TestParseObjectLiteralExpression(t *testing.T) {
	prog := parseProgram(t, "const obj = { foo, bar: 2, [\"baz\"]: value, ...rest };")

	decl, ok := prog.Body[0].(*ast.VariableDeclaration)
	if !ok {
		t.Fatalf("expected VariableDeclaration, got %T", prog.Body[0])
	}

	objInit, ok := decl.Declarations[0].Init.(*ast.ObjectLiteral)
	if !ok {
		t.Fatalf("expected ObjectLiteral, got %T", decl.Declarations[0].Init)
	}

	if len(objInit.Properties) != 4 {
		t.Fatalf("expected 4 properties, got %d", len(objInit.Properties))
	}

	prop0, ok := objInit.Properties[0].(*ast.ObjectProperty)
	if !ok || !prop0.Shorthand {
		t.Fatalf("expected shorthand property, got %#v", objInit.Properties[0])
	}

	prop1, ok := objInit.Properties[1].(*ast.ObjectProperty)
	if !ok {
		t.Fatalf("expected object property, got %T", objInit.Properties[1])
	}

	if num, ok := prop1.Value.(*ast.NumberLiteral); !ok || num.Value != "2" {
		t.Fatalf("unexpected value for bar: %#v", prop1.Value)
	}

	prop2, ok := objInit.Properties[2].(*ast.ObjectProperty)
	if !ok || !prop2.Computed {
		t.Fatalf("expected computed property, got %#v", objInit.Properties[2])
	}

	if lit, ok := prop2.Key.(*ast.StringLiteral); !ok || lit.Value != "baz" {
		t.Fatalf("unexpected computed key: %#v", prop2.Key)
	}

	spread, ok := objInit.Properties[3].(*ast.SpreadElement)
	if !ok {
		t.Fatalf("expected spread element, got %T", objInit.Properties[3])
	}

	if ident, ok := spread.Argument.(*ast.Identifier); !ok || ident.Name != "rest" {
		t.Fatalf("unexpected spread argument: %#v", spread.Argument)
	}
}

func TestParseConditionalExpression(t *testing.T) {
	prog := parseProgram(t, "const result = cond ? left() : right;")

	decl, ok := prog.Body[0].(*ast.VariableDeclaration)
	if !ok {
		t.Fatalf("expected VariableDeclaration, got %T", prog.Body[0])
	}

	condExpr, ok := decl.Declarations[0].Init.(*ast.ConditionalExpression)
	if !ok {
		t.Fatalf("expected ConditionalExpression, got %T", decl.Declarations[0].Init)
	}

	if ident, ok := condExpr.Test.(*ast.Identifier); !ok || ident.Name != "cond" {
		t.Fatalf("unexpected test: %#v", condExpr.Test)
	}

	if _, ok := condExpr.Consequent.(*ast.CallExpression); !ok {
		t.Fatalf("expected call expression consequent, got %T", condExpr.Consequent)
	}

	if alt, ok := condExpr.Alternate.(*ast.Identifier); !ok || alt.Name != "right" {
		t.Fatalf("unexpected alternate: %#v", condExpr.Alternate)
	}
}

func TestParseSequenceExpression(t *testing.T) {
	prog := parseProgram(t, "a(), b = 2, c + d;")

	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}

	exprStmt, ok := prog.Body[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", prog.Body[0])
	}

	seq, ok := exprStmt.Expression.(*ast.SequenceExpression)
	if !ok {
		t.Fatalf("expected SequenceExpression, got %T", exprStmt.Expression)
	}

	if len(seq.Expressions) != 3 {
		t.Fatalf("expected 3 expressions in sequence, got %d", len(seq.Expressions))
	}

	if _, ok := seq.Expressions[0].(*ast.CallExpression); !ok {
		t.Fatalf("expected call expression first, got %T", seq.Expressions[0])
	}

	if _, ok := seq.Expressions[1].(*ast.AssignmentExpression); !ok {
		t.Fatalf("expected assignment second, got %T", seq.Expressions[1])
	}

	if _, ok := seq.Expressions[2].(*ast.BinaryExpression); !ok {
		t.Fatalf("expected binary expression third, got %T", seq.Expressions[2])
	}
}
