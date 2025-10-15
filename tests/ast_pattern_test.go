package tests

import (
	"testing"

	"es6-interpreter/ast"
)

func TestArrayPatternConstruction(t *testing.T) {
	loc := sampleLoc()
	elem := ast.NewIdentifier("value", loc)
	rest := ast.NewRestElement(ast.NewIdentifier("rest", loc), loc)
	pattern := ast.NewArrayPattern(ast.PatternList{elem, nil}, rest, loc)

	if pattern.Kind() != ast.ArrayPatternKind {
		t.Fatalf("array pattern kind mismatch: got %q", pattern.Kind())
	}
	if len(pattern.Elements) != 2 {
		t.Fatalf("array pattern elements length mismatch: got %d", len(pattern.Elements))
	}
	if pattern.Rest == nil || pattern.Rest.Argument.String() != "Identifier(rest)" {
		t.Fatalf("array pattern rest element missing or incorrect")
	}
}

func TestObjectPatternProperty(t *testing.T) {
	loc := sampleLoc()
	key := ast.NewIdentifier("key", loc)
	value := ast.NewIdentifier("alias", loc)
	prop := ast.NewObjectPatternProperty(key, value, false, true, loc)

	if prop.Kind() != ast.ObjectPatternPropKind {
		t.Fatalf("object pattern property kind mismatch: got %q", prop.Kind())
	}
	if !prop.Shorthand {
		t.Fatalf("object pattern property should be shorthand")
	}

	pattern := ast.NewObjectPattern([]*ast.ObjectPatternProperty{prop}, nil, loc)
	if pattern.Kind() != ast.ObjectPatternKind {
		t.Fatalf("object pattern kind mismatch: got %q", pattern.Kind())
	}
	if len(pattern.Properties) != 1 || pattern.Properties[0] != prop {
		t.Fatalf("object pattern properties not preserved")
	}
}

func TestAssignmentPatternDefaults(t *testing.T) {
	loc := sampleLoc()
	left := ast.NewIdentifier("value", loc)
	right := ast.NewNumberLiteral("42", loc)
	assign := ast.NewAssignmentPattern(left, right, loc)

	if assign.Kind() != ast.AssignmentPatternKind {
		t.Fatalf("assignment pattern kind mismatch: got %q", assign.Kind())
	}
	if assign.Left != left || assign.Right != right {
		t.Fatalf("assignment pattern children mismatch")
	}
}

func TestRestElementPattern(t *testing.T) {
	loc := sampleLoc()
	arg := ast.NewIdentifier("rest", loc)
	rest := ast.NewRestElement(arg, loc)

	if rest.Kind() != ast.RestElementKind {
		t.Fatalf("rest element kind mismatch: got %q", rest.Kind())
	}
	if rest.Argument != arg {
		t.Fatalf("rest element argument mismatch")
	}
}
