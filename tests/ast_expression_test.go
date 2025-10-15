package tests

import (
	"testing"

	"es6-interpreter/ast"
)

func sampleLoc() ast.Location {
	return ast.Location{
		Start: ast.Position{Offset: 0, Line: 1, Column: 0},
		End:   ast.Position{Offset: 1, Line: 1, Column: 1},
	}
}

func TestIdentifierNodeContracts(t *testing.T) {
	loc := sampleLoc()
	id := ast.NewIdentifier("answer", loc)

	if got := id.Kind(); got != ast.IdentifierKind {
		t.Fatalf("identifier kind mismatch: got %q", got)
	}
	if id.Name != "answer" {
		t.Fatalf("identifier name mismatch: got %q", id.Name)
	}
	if id.Loc() != loc {
		t.Fatalf("identifier location mismatch: got %+v", id.Loc())
	}

	if _, ok := any(id).(ast.Expression); !ok {
		t.Fatalf("identifier should satisfy ast.Expression")
	}
	if _, ok := any(id).(ast.Pattern); !ok {
		t.Fatalf("identifier should satisfy ast.Pattern")
	}
}

func TestMemberAndCallExpressions(t *testing.T) {
	loc := sampleLoc()
	obj := ast.NewIdentifier("obj", loc)
	prop := ast.NewIdentifier("prop", loc)
	member := ast.NewMemberExpression(obj, prop, false, loc)

	if member.Kind() != ast.MemberExpressionKind {
		t.Fatalf("member kind mismatch: got %q", member.Kind())
	}
	if member.Object != obj || member.Property != prop {
		t.Fatalf("member child nodes not preserved")
	}
	if member.Computed {
		t.Fatalf("member should not be computed")
	}

	call := ast.NewCallExpression(member, []ast.Expression{prop}, loc)
	if call.Kind() != ast.CallExpressionKind {
		t.Fatalf("call kind mismatch: got %q", call.Kind())
	}
	if len(call.Arguments) != 1 || call.Arguments[0] != prop {
		t.Fatalf("call arguments not preserved")
	}
}

func TestAssignmentAndBinaryExpressions(t *testing.T) {
	loc := sampleLoc()
	left := ast.NewIdentifier("left", loc)
	right := ast.NewIdentifier("right", loc)

	assign := ast.NewAssignmentExpression("=", left, right, loc)
	if assign.Kind() != ast.AssignmentExpressionKind {
		t.Fatalf("assignment kind mismatch: got %q", assign.Kind())
	}
	if assign.Operator != "=" {
		t.Fatalf("assignment operator mismatch: got %q", assign.Operator)
	}
	if assign.Left != left || assign.Right != right {
		t.Fatalf("assignment children mismatch")
	}

	binary := ast.NewBinaryExpression("+", left, right, loc)
	if binary.Kind() != ast.BinaryExpressionKind {
		t.Fatalf("binary kind mismatch: got %q", binary.Kind())
	}
	if binary.Operator != "+" {
		t.Fatalf("binary operator mismatch: got %q", binary.Operator)
	}
}

func TestSpreadAndTemplateHelpers(t *testing.T) {
	loc := sampleLoc()
	spread := ast.NewSpreadElement(ast.NewIdentifier("arg", loc), loc)

	if _, ok := any(spread).(ast.Expression); !ok {
		t.Fatalf("spread should satisfy ast.Expression")
	}
	if _, ok := any(spread).(ast.Property); !ok {
		t.Fatalf("spread should satisfy ast.Property")
	}

	quasi := ast.NewTemplateElement("raw", "cooked", true, loc)
	tmpl := ast.NewTemplateLiteral([]*ast.TemplateElement{quasi}, nil, loc)
	tagged := ast.NewTaggedTemplateExpression(ast.NewIdentifier("tag", loc), tmpl, loc)

	if tagged.Kind() != ast.TaggedTemplateExpressionKind {
		t.Fatalf("tagged template kind mismatch: got %q", tagged.Kind())
	}
	if tagged.Quasi != tmpl {
		t.Fatalf("tagged template literal not preserved")
	}
}
