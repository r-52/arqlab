package vm

import "testing"

func TestEnvironmentLetLifecycle(t *testing.T) {
	env := NewEnvironment(nil)
	if err := env.Declare("x", BindingLet); err != nil {
		t.Fatalf("declare let: %v", err)
	}
	if _, err := env.Get("x"); err == nil {
		t.Fatalf("expected temporal dead zone error")
	}
	if err := env.Initialize("x", NewNumber(42)); err != nil {
		t.Fatalf("initialize let: %v", err)
	}
	v, err := env.Get("x")
	if err != nil {
		t.Fatalf("get after initialize: %v", err)
	}
	if !StrictEquals(v, NewNumber(42)) {
		t.Fatalf("expected 42, got %s", v.Inspect())
	}
}

func TestEnvironmentConstAssignment(t *testing.T) {
	env := NewEnvironment(nil)
	if err := env.Declare("answer", BindingConst); err != nil {
		t.Fatalf("declare const: %v", err)
	}
	if err := env.Initialize("answer", NewNumber(42)); err != nil {
		t.Fatalf("initialize const: %v", err)
	}
	if err := env.Set("answer", NewNumber(99)); err == nil {
		t.Fatalf("expected error assigning to const")
	}
	v, err := env.Get("answer")
	if err != nil {
		t.Fatalf("get const: %v", err)
	}
	if !StrictEquals(v, NewNumber(42)) {
		t.Fatalf("const mutated unexpectedly: %s", v.Inspect())
	}
}

func TestEnvironmentVarRedeclaration(t *testing.T) {
	env := NewEnvironment(nil)
	if err := env.Declare("foo", BindingVar); err != nil {
		t.Fatalf("declare var: %v", err)
	}
	if err := env.Declare("foo", BindingVar); err != nil {
		t.Fatalf("redeclaration of var should succeed: %v", err)
	}
	if err := env.Set("foo", NewString("bar")); err != nil {
		t.Fatalf("set var: %v", err)
	}
	v, err := env.Get("foo")
	if err != nil {
		t.Fatalf("get var: %v", err)
	}
	if !StrictEquals(v, NewString("bar")) {
		t.Fatalf("expected bar, got %s", v.Inspect())
	}
}

func TestEnvironmentScopeLookupAndAssignment(t *testing.T) {
	global := NewEnvironment(nil)
	if err := global.Declare("count", BindingVar); err != nil {
		t.Fatalf("declare global var: %v", err)
	}

	block := NewEnvironment(global)
	if err := block.Declare("local", BindingLet); err != nil {
		t.Fatalf("declare local let: %v", err)
	}
	if err := block.Initialize("local", NewNumber(10)); err != nil {
		t.Fatalf("initialize local: %v", err)
	}

	if err := block.Set("count", NewNumber(7)); err != nil {
		t.Fatalf("set outer var: %v", err)
	}

	v, err := global.Get("count")
	if err != nil {
		t.Fatalf("get global: %v", err)
	}
	if !StrictEquals(v, NewNumber(7)) {
		t.Fatalf("expected global count to be 7, got %s", v.Inspect())
	}

	local, err := block.Get("local")
	if err != nil {
		t.Fatalf("get local: %v", err)
	}
	if !StrictEquals(local, NewNumber(10)) {
		t.Fatalf("expected local 10, got %s", local.Inspect())
	}
}

func TestEnvironmentSetBeforeInitialization(t *testing.T) {
	env := NewEnvironment(nil)
	if err := env.Declare("value", BindingLet); err != nil {
		t.Fatalf("declare let: %v", err)
	}
	if err := env.Set("value", NewNumber(1)); err == nil {
		t.Fatalf("expected error when assigning before initialization")
	}
}
