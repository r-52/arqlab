package vm

import "fmt"

// BindingKind represents how an identifier was declared.
type BindingKind int

const (
	BindingVar BindingKind = iota
	BindingLet
	BindingConst
)

type binding struct {
	value       Value
	mutable     bool
	initialized bool
	kind        BindingKind
}

// Environment models a lexical environment (scope) with an optional outer scope.
type Environment struct {
	outer     *Environment
	record    map[string]*binding
	varParent *Environment
	isVarEnv  bool
}

// NewEnvironment creates a new environment with the provided outer environment.
func NewEnvironment(outer *Environment) *Environment {
	env := &Environment{
		outer:  outer,
		record: make(map[string]*binding),
	}
	if outer == nil {
		env.varParent = env
		env.isVarEnv = true
	} else {
		env.varParent = outer.varParent
	}
	return env
}

// NewVariableEnvironment constructs an environment that introduces a new var scope
// (e.g. function execution context).
func NewVariableEnvironment(outer *Environment) *Environment {
	env := &Environment{
		outer:    outer,
		record:   make(map[string]*binding),
		isVarEnv: true,
	}
	env.varParent = env
	return env
}

// Outer returns the parent environment.
func (e *Environment) Outer() *Environment { return e.outer }

// VarParent returns the var-binding container associated with this environment.
func (e *Environment) VarParent() *Environment {
	if e.varParent != nil {
		return e.varParent
	}
	return e
}

// HasOwn reports whether the current environment contains a binding for name.
func (e *Environment) HasOwn(name string) bool {
	_, ok := e.record[name]
	return ok
}

func (e *Environment) targetFor(kind BindingKind) *Environment {
	if kind == BindingVar {
		if e.varParent != nil {
			return e.varParent
		}
		return e
	}
	return e
}

// Declare creates a binding following the semantics of the provided kind.
// Redeclaring a var in the same scope is a no-op.
func (e *Environment) Declare(name string, kind BindingKind) error {
	target := e.targetFor(kind)
	if existing, ok := target.record[name]; ok {
		if kind == BindingVar && existing.kind == BindingVar {
			return nil
		}
		return fmt.Errorf("SyntaxError: identifier %q has already been declared", name)
	}

	b := &binding{kind: kind}
	switch kind {
	case BindingVar:
		b.mutable = true
		b.initialized = true
		b.value = Undefined
	case BindingLet:
		b.mutable = true
	case BindingConst:
		b.mutable = false
	default:
		return fmt.Errorf("internal error: unknown binding kind %d", kind)
	}

	target.record[name] = b
	return nil
}

// Initialize assigns the first value to a previously declared binding in the
// current environment. It is primarily used for let/const declarations.
func (e *Environment) Initialize(name string, value Value) error {
	b, ok := e.record[name]
	if !ok {
		return fmt.Errorf("ReferenceError: %s is not defined", name)
	}
	if b.initialized {
		return fmt.Errorf("TypeError: identifier %q has already been initialized", name)
	}

	b.value = value
	b.initialized = true
	return nil
}

// Get returns the value bound to name, searching outward through parent
// environments.
func (e *Environment) Get(name string) (Value, error) {
	if b, ok := e.record[name]; ok {
		if !b.initialized {
			return Value{}, fmt.Errorf("ReferenceError: Cannot access '%s' before initialization", name)
		}
		return b.value, nil
	}
	if e.outer != nil {
		return e.outer.Get(name)
	}
	return Value{}, fmt.Errorf("ReferenceError: %s is not defined", name)
}

// Set updates the value bound to name, searching outward through parent
// environments. Attempting to update an immutable binding yields an error.
func (e *Environment) Set(name string, value Value) error {
	if b, ok := e.record[name]; ok {
		if !b.initialized {
			return fmt.Errorf("ReferenceError: Cannot access '%s' before initialization", name)
		}
		if !b.mutable {
			return fmt.Errorf("TypeError: Assignment to constant variable %q", name)
		}
		b.value = value
		return nil
	}
	if e.outer != nil {
		return e.outer.Set(name, value)
	}
	return fmt.Errorf("ReferenceError: %s is not defined", name)
}

// Resolve finds the binding entry for name, searching through outer environments.
func (e *Environment) Resolve(name string) (*binding, bool) {
	if b, ok := e.record[name]; ok {
		return b, true
	}
	if e.outer != nil {
		return e.outer.Resolve(name)
	}
	return nil, false
}
