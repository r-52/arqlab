package vm

import (
	"testing"

	"es6-interpreter/parser"
)

func executeSnippet(t *testing.T, src string) Value {
	t.Helper()
	p := parser.New(src)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	result, err := Execute(program)
	if err != nil {
		t.Fatalf("execute error: %v", err)
	}
	return result
}

func executeSnippetExpectError(t *testing.T, src string) error {
	t.Helper()
	p := parser.New(src)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	_, err = Execute(program)
	if err == nil {
		t.Fatalf("expected execution error but got nil")
	}
	return err
}

func TestInterpreterArithmetic(t *testing.T) {
	result := executeSnippet(t, "let x = 1 + 2 * 3; x;")
	if result.Kind() != NumberKind {
		t.Fatalf("expected number result, got %v", result.Kind())
	}
	if got := result.Number(); got != 7 {
		t.Fatalf("expected 7, got %v", got)
	}
}

func TestInterpreterBlockScoping(t *testing.T) {
	result := executeSnippet(t, `
let x = 1;
{
  let x = 2;
  x + 1;
}
x;
`)
	if result.Kind() != NumberKind || result.Number() != 1 {
		t.Fatalf("expected outer x to remain 1, got %s", result.Inspect())
	}
}

func TestInterpreterIfElse(t *testing.T) {
	result := executeSnippet(t, `
let value = 0;
if (1 < 2) {
  value = 5;
} else {
  value = 10;
}
value;
`)
	if result.Kind() != NumberKind || result.Number() != 5 {
		t.Fatalf("expected value 5, got %s", result.Inspect())
	}
}

func TestInterpreterWhileLoop(t *testing.T) {
	result := executeSnippet(t, `
let sum = 0;
let i = 0;
while (i < 3) {
  sum += i;
  i = i + 1;
}
sum;
`)
	if result.Kind() != NumberKind || result.Number() != 3 {
		t.Fatalf("expected sum 3, got %s", result.Inspect())
	}
}

func TestInterpreterLogicalShortCircuit(t *testing.T) {
	result := executeSnippet(t, `
let x = 0;
(0 && (x = 1));
(1 || (x = 2));
x;
`)
	if result.Kind() != NumberKind || result.Number() != 0 {
		t.Fatalf("expected x to remain 0, got %s", result.Inspect())
	}
}

func TestInterpreterConstReassignmentError(t *testing.T) {
	executeSnippetExpectError(t, `
const answer = 42;
answer = 7;
`)
}

func TestInterpreterConstRequiresInitializer(t *testing.T) {
	executeSnippetExpectError(t, `const missing;`)
}

func TestInterpreterUndefinedIdentifier(t *testing.T) {
	executeSnippetExpectError(t, `unknown;`)
}
