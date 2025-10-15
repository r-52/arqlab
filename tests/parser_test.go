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
