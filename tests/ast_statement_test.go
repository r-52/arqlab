package tests

import (
	"testing"

	"es6-interpreter/ast"
)

func TestProgramAndBlock(t *testing.T) {
	loc := sampleLoc()
	body := []ast.Statement{
		ast.NewExpressionStatement(ast.NewIdentifier("x", loc), loc),
	}
	block := ast.NewBlockStatement(body, loc)
	prog := ast.NewProgram([]ast.Statement{block}, ast.SourceTypeScript, loc)

	if prog.Kind() != ast.ProgramKind {
		t.Fatalf("program kind mismatch: got %q", prog.Kind())
	}
	if prog.SourceType != ast.SourceTypeScript {
		t.Fatalf("program source type mismatch: got %q", prog.SourceType)
	}
	if len(prog.Body) != 1 || prog.Body[0] != block {
		t.Fatalf("program body not preserved")
	}
}

func TestVariableDeclarationShape(t *testing.T) {
	loc := sampleLoc()
	id := ast.NewIdentifier("value", loc)
	init := ast.NewNumberLiteral("1", loc)
	declarator := ast.NewVariableDeclarator(id, init, loc)
	decl := ast.NewVariableDeclaration(ast.ConstKind, []*ast.VariableDeclarator{declarator}, loc)

	if decl.Kind() != ast.VariableDeclarationKind {
		t.Fatalf("variable declaration kind mismatch: got %q", decl.Kind())
	}
	if decl.DeclareKind != ast.ConstKind {
		t.Fatalf("variable declaration kind value mismatch: got %q", decl.DeclareKind)
	}
	if len(decl.Declarations) != 1 || decl.Declarations[0] != declarator {
		t.Fatalf("variable declaration members mismatch")
	}
}

func TestControlFlowNodes(t *testing.T) {
	loc := sampleLoc()
	testExpr := ast.NewIdentifier("cond", loc)
	body := ast.NewBlockStatement(nil, loc)

	ifStmt := ast.NewIfStatement(testExpr, body, nil, loc)
	if ifStmt.Kind() != ast.IfStatementKind {
		t.Fatalf("if statement kind mismatch: got %q", ifStmt.Kind())
	}

	forStmt := ast.NewForStatement(nil, testExpr, nil, body, loc)
	if forStmt.Kind() != ast.ForStatementKind {
		t.Fatalf("for statement kind mismatch: got %q", forStmt.Kind())
	}

	whileStmt := ast.NewWhileStatement(testExpr, body, loc)
	if whileStmt.Kind() != ast.WhileStatementKind {
		t.Fatalf("while statement kind mismatch: got %q", whileStmt.Kind())
	}

	returnStmt := ast.NewReturnStatement(testExpr, loc)
	if returnStmt.Argument != testExpr {
		t.Fatalf("return statement argument mismatch")
	}
}

func TestFunctionDeclarationNode(t *testing.T) {
	loc := sampleLoc()
	id := ast.NewIdentifier("fn", loc)
	params := []ast.Pattern{ast.NewIdentifier("x", loc)}
	body := ast.NewBlockStatement(nil, loc)
	fn := ast.NewFunctionDeclaration(id, params, body, false, loc)

	if fn.Kind() != ast.FunctionDeclarationKind {
		t.Fatalf("function declaration kind mismatch: got %q", fn.Kind())
	}
	if fn.ID != id {
		t.Fatalf("function declaration id mismatch")
	}
	if len(fn.Params) != 1 {
		t.Fatalf("function declaration params mismatch: got %d", len(fn.Params))
	}
}
