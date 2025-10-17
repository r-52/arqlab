package vm

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"es6-interpreter/ast"
)

// Interpreter evaluates ECMAScript AST nodes to produce runtime values.
type Interpreter struct {
	global *Environment
}

// NewInterpreter constructs a fresh interpreter instance with an empty global scope.
func NewInterpreter() *Interpreter {
	global := NewEnvironment(nil)
	return &Interpreter{global: global}
}

// Execute runs the supplied program and returns the completion value produced by
// the final statement. Scripts that do not yield a value return undefined.
func Execute(program *ast.Program) (Value, error) {
	intr := NewInterpreter()
	comp, err := intr.evalProgram(program)
	if err != nil {
		return Value{}, err
	}
	return comp.value, nil
}

type completionType int

const (
	completionNormal completionType = iota
	completionReturn
	completionBreak
	completionContinue
)

type completion struct {
	kind  completionType
	value Value
	label string
}

func normalCompletion(v Value) completion {
	return completion{kind: completionNormal, value: v}
}

func (c completion) withValue(v Value) completion {
	c.value = v
	return c
}

func (i *Interpreter) evalProgram(program *ast.Program) (completion, error) {
	var last Value = Undefined
	for _, stmt := range program.Body {
		comp, err := i.evalStatement(i.global, stmt)
		if err != nil {
			return completion{}, err
		}
		switch comp.kind {
		case completionNormal:
			last = comp.value
		case completionReturn:
			return comp, nil
		case completionBreak, completionContinue:
			return completion{}, fmt.Errorf("runtime error: unexpected %s outside of loop", i.describeCompletion(comp))
		default:
			return completion{}, fmt.Errorf("runtime error: unsupported completion type %d", comp.kind)
		}
	}
	return normalCompletion(last), nil
}

func (i *Interpreter) evalStatement(env *Environment, stmt ast.Statement) (completion, error) {
	switch s := stmt.(type) {
	case *ast.BlockStatement:
		blockEnv := NewEnvironment(env)
		return i.evalStatementList(blockEnv, s.Body)
	case *ast.ExpressionStatement:
		val, err := i.evalExpression(env, s.Expression)
		if err != nil {
			return completion{}, err
		}
		return normalCompletion(val), nil
	case *ast.EmptyStatement:
		return normalCompletion(Undefined), nil
	case *ast.VariableDeclaration:
		if err := i.evalVariableDeclaration(env, s); err != nil {
			return completion{}, err
		}
		return normalCompletion(Undefined), nil
	case *ast.IfStatement:
		return i.evalIfStatement(env, s)
	case *ast.WhileStatement:
		return i.evalWhileStatement(env, s)
	case *ast.ForStatement:
		return i.evalForStatement(env, s)
	case *ast.BreakStatement:
		label := ""
		if s.Label != nil {
			label = s.Label.Name
		}
		return completion{kind: completionBreak, label: label}, nil
	case *ast.ContinueStatement:
		label := ""
		if s.Label != nil {
			label = s.Label.Name
		}
		return completion{kind: completionContinue, label: label}, nil
	case *ast.ReturnStatement:
		val := Undefined
		if s.Argument != nil {
			result, err := i.evalExpression(env, s.Argument)
			if err != nil {
				return completion{}, err
			}
			val = result
		}
		return completion{kind: completionReturn, value: val}, nil
	case *ast.LabeledStatement:
		comp, err := i.evalStatement(env, s.Body)
		if err != nil {
			return completion{}, err
		}
		if comp.kind == completionBreak && comp.label == s.Label.Name {
			return normalCompletion(comp.value), nil
		}
		return comp, nil
	default:
		return completion{}, fmt.Errorf("runtime error: statement %T not supported", s)
	}
}

func (i *Interpreter) evalStatementList(env *Environment, stmts []ast.Statement) (completion, error) {
	var last Value = Undefined
	for _, stmt := range stmts {
		comp, err := i.evalStatement(env, stmt)
		if err != nil {
			return completion{}, err
		}
		switch comp.kind {
		case completionNormal:
			last = comp.value
		case completionBreak, completionContinue, completionReturn:
			return comp, nil
		default:
			return completion{}, fmt.Errorf("runtime error: unsupported completion type %d", comp.kind)
		}
	}
	return normalCompletion(last), nil
}

func (i *Interpreter) evalIfStatement(env *Environment, stmt *ast.IfStatement) (completion, error) {
	testVal, err := i.evalExpression(env, stmt.Test)
	if err != nil {
		return completion{}, err
	}
	if ToBoolean(testVal) {
		return i.evalStatement(env, stmt.Consequent)
	}
	if stmt.Alternate != nil {
		return i.evalStatement(env, stmt.Alternate)
	}
	return normalCompletion(Undefined), nil
}

func (i *Interpreter) evalWhileStatement(env *Environment, stmt *ast.WhileStatement) (completion, error) {
	var last Value = Undefined
	for {
		testVal, err := i.evalExpression(env, stmt.Test)
		if err != nil {
			return completion{}, err
		}
		if !ToBoolean(testVal) {
			return normalCompletion(last), nil
		}

		bodyComp, err := i.evalStatement(env, stmt.Body)
		if err != nil {
			return completion{}, err
		}

		switch bodyComp.kind {
		case completionNormal:
			last = bodyComp.value
		case completionReturn:
			return bodyComp, nil
		case completionBreak:
			if bodyComp.label == "" {
				return normalCompletion(bodyComp.value), nil
			}
			return bodyComp, nil
		case completionContinue:
			if bodyComp.label != "" {
				return bodyComp, nil
			}
			continue
		default:
			return completion{}, fmt.Errorf("runtime error: unsupported completion in while body: %d", bodyComp.kind)
		}
	}
}

func (i *Interpreter) evalForStatement(env *Environment, stmt *ast.ForStatement) (completion, error) {
	loopEnv := NewEnvironment(env)
	if stmt.Init != nil {
		switch init := stmt.Init.(type) {
		case ast.Expression:
			if _, err := i.evalExpression(loopEnv, init); err != nil {
				return completion{}, err
			}
		case *ast.VariableDeclaration:
			if err := i.evalVariableDeclaration(loopEnv, init); err != nil {
				return completion{}, err
			}
		default:
			return completion{}, fmt.Errorf("runtime error: unsupported for-loop initializer %T", init)
		}
	}

	var last Value = Undefined
	for {
		if stmt.Test != nil {
			testVal, err := i.evalExpression(loopEnv, stmt.Test)
			if err != nil {
				return completion{}, err
			}
			if !ToBoolean(testVal) {
				return normalCompletion(last), nil
			}
		}

		bodyComp, err := i.evalStatement(loopEnv, stmt.Body)
		if err != nil {
			return completion{}, err
		}

		skipUpdate := false
		switch bodyComp.kind {
		case completionNormal:
			last = bodyComp.value
		case completionReturn:
			return bodyComp, nil
		case completionBreak:
			if bodyComp.label == "" {
				return normalCompletion(bodyComp.value), nil
			}
			return bodyComp, nil
		case completionContinue:
			if bodyComp.label != "" {
				return bodyComp, nil
			}
			skipUpdate = false
		default:
			return completion{}, fmt.Errorf("runtime error: unsupported completion in for body: %d", bodyComp.kind)
		}

		if stmt.Update != nil && !skipUpdate {
			if _, err := i.evalExpression(loopEnv, stmt.Update); err != nil {
				return completion{}, err
			}
		}
	}
}

func (i *Interpreter) evalVariableDeclaration(env *Environment, decl *ast.VariableDeclaration) error {
	var kind BindingKind
	switch decl.DeclareKind {
	case ast.VarKind:
		kind = BindingVar
	case ast.LetKind:
		kind = BindingLet
	case ast.ConstKind:
		kind = BindingConst
	default:
		return fmt.Errorf("runtime error: unsupported variable kind %s", decl.DeclareKind)
	}

	for _, d := range decl.Declarations {
		ident, ok := d.ID.(*ast.Identifier)
		if !ok {
			return fmt.Errorf("runtime error: destructuring bindings are not implemented yet (%T)", d.ID)
		}

		target := env
		if kind == BindingVar {
			target = env.VarParent()
		}

		if err := target.Declare(ident.Name, kind); err != nil {
			return err
		}

		if d.Init != nil {
			initVal, err := i.evalExpression(env, d.Init)
			if err != nil {
				return err
			}
			switch kind {
			case BindingVar:
				if err := target.Set(ident.Name, initVal); err != nil {
					return err
				}
			case BindingLet, BindingConst:
				if err := target.Initialize(ident.Name, initVal); err != nil {
					return err
				}
			}
		} else if kind == BindingConst {
			return fmt.Errorf("TypeError: const declaration %q requires an initializer", ident.Name)
		}
	}

	return nil
}

func (i *Interpreter) evalExpression(env *Environment, expr ast.Expression) (Value, error) {
	switch e := expr.(type) {
	case *ast.NumberLiteral:
		return i.evalNumberLiteral(e)
	case *ast.StringLiteral:
		return NewString(e.Value), nil
	case *ast.BooleanLiteral:
		return NewBoolean(e.Value), nil
	case *ast.NullLiteral:
		return Null, nil
	case *ast.Identifier:
		val, err := env.Get(e.Name)
		if err != nil {
			return Value{}, err
		}
		return val, nil
	case *ast.BinaryExpression:
		left, err := i.evalExpression(env, e.Left)
		if err != nil {
			return Value{}, err
		}
		right, err := i.evalExpression(env, e.Right)
		if err != nil {
			return Value{}, err
		}
		return i.applyBinary(e.Operator, left, right)
	case *ast.AssignmentExpression:
		return i.evalAssignmentExpression(env, e)
	case *ast.LogicalExpression:
		return i.evalLogicalExpression(env, e)
	case *ast.UnaryExpression:
		return i.evalUnaryExpression(env, e)
	case *ast.UpdateExpression:
		return i.evalUpdateExpression(env, e)
	case *ast.ConditionalExpression:
		test, err := i.evalExpression(env, e.Test)
		if err != nil {
			return Value{}, err
		}
		if ToBoolean(test) {
			return i.evalExpression(env, e.Consequent)
		}
		return i.evalExpression(env, e.Alternate)
	case *ast.SequenceExpression:
		var last Value = Undefined
		for _, inner := range e.Expressions {
			val, err := i.evalExpression(env, inner)
			if err != nil {
				return Value{}, err
			}
			last = val
		}
		return last, nil
	default:
		return Value{}, fmt.Errorf("runtime error: expression %T not supported", e)
	}
}

func (i *Interpreter) evalNumberLiteral(lit *ast.NumberLiteral) (Value, error) {
	num, err := parseNumericLiteral(lit.Value)
	if err != nil {
		return Value{}, fmt.Errorf("runtime error: invalid numeric literal %q: %v", lit.Value, err)
	}
	return NewNumber(num), nil
}

func (i *Interpreter) evalAssignmentExpression(env *Environment, expr *ast.AssignmentExpression) (Value, error) {
	target, ok := expr.Left.(*ast.Identifier)
	if !ok {
		return Value{}, fmt.Errorf("runtime error: assignment target %T not supported", expr.Left)
	}

	right, err := i.evalExpression(env, expr.Right)
	if err != nil {
		return Value{}, err
	}

	switch expr.Operator {
	case "=":
		if err := env.Set(target.Name, right); err != nil {
			return Value{}, err
		}
		return right, nil
	case "+=", "-=", "*=", "/=", "%=":
		current, err := env.Get(target.Name)
		if err != nil {
			return Value{}, err
		}
		op := expr.Operator[:len(expr.Operator)-1]
		result, err := i.applyBinary(op, current, right)
		if err != nil {
			return Value{}, err
		}
		if err := env.Set(target.Name, result); err != nil {
			return Value{}, err
		}
		return result, nil
	default:
		return Value{}, fmt.Errorf("runtime error: assignment operator %q not implemented", expr.Operator)
	}
}

func (i *Interpreter) evalLogicalExpression(env *Environment, expr *ast.LogicalExpression) (Value, error) {
	left, err := i.evalExpression(env, expr.Left)
	if err != nil {
		return Value{}, err
	}

	switch expr.Operator {
	case "&&":
		if !ToBoolean(left) {
			return left, nil
		}
		return i.evalExpression(env, expr.Right)
	case "||":
		if ToBoolean(left) {
			return left, nil
		}
		return i.evalExpression(env, expr.Right)
	case "??":
		if left.Kind() != UndefinedKind && left.Kind() != NullKind {
			return left, nil
		}
		return i.evalExpression(env, expr.Right)
	default:
		return Value{}, fmt.Errorf("runtime error: logical operator %q not supported", expr.Operator)
	}
}

func (i *Interpreter) evalUnaryExpression(env *Environment, expr *ast.UnaryExpression) (Value, error) {
	arg, err := i.evalExpression(env, expr.Argument)
	if err != nil {
		return Value{}, err
	}

	switch expr.Operator {
	case "!":
		return NewBoolean(!ToBoolean(arg)), nil
	case "+":
		n := ToNumber(arg)
		return n, nil
	case "-":
		n := ToNumber(arg)
		return NewNumber(-n.Number()), nil
	case "typeof":
		return NewString(i.typeOfValue(arg)), nil
	case "void":
		return Undefined, nil
	default:
		return Value{}, fmt.Errorf("runtime error: unary operator %q not implemented", expr.Operator)
	}
}

func (i *Interpreter) evalUpdateExpression(env *Environment, expr *ast.UpdateExpression) (Value, error) {
	target, ok := expr.Argument.(*ast.Identifier)
	if !ok {
		return Value{}, fmt.Errorf("runtime error: update target %T not supported", expr.Argument)
	}

	current, err := env.Get(target.Name)
	if err != nil {
		return Value{}, err
	}

	n := ToNumber(current)
	value := n.Number()

	var next float64
	switch expr.Operator {
	case "++":
		next = value + 1
	case "--":
		next = value - 1
	default:
		return Value{}, fmt.Errorf("runtime error: update operator %q not supported", expr.Operator)
	}

	updated := NewNumber(next)
	if err := env.Set(target.Name, updated); err != nil {
		return Value{}, err
	}

	if expr.Prefix {
		return updated, nil
	}
	return current, nil
}

func (i *Interpreter) applyBinary(op string, left, right Value) (Value, error) {
	switch op {
	case "+":
		if left.Kind() == StringKind || right.Kind() == StringKind {
			ls := ToString(left)
			rs := ToString(right)
			return NewString(ls.StringValue() + rs.StringValue()), nil
		}
		ln := ToNumber(left)
		rn := ToNumber(right)
		return NewNumber(ln.Number() + rn.Number()), nil
	case "-":
		ln := ToNumber(left)
		rn := ToNumber(right)
		return NewNumber(ln.Number() - rn.Number()), nil
	case "*":
		ln := ToNumber(left)
		rn := ToNumber(right)
		return NewNumber(ln.Number() * rn.Number()), nil
	case "/":
		ln := ToNumber(left)
		rn := ToNumber(right)
		return NewNumber(ln.Number() / rn.Number()), nil
	case "%":
		ln := ToNumber(left)
		rn := ToNumber(right)
		return NewNumber(math.Mod(ln.Number(), rn.Number())), nil
	case "===":
		return NewBoolean(StrictEquals(left, right)), nil
	case "!==":
		return NewBoolean(!StrictEquals(left, right)), nil
	case "==":
		return NewBoolean(StrictEquals(left, right)), nil
	case "!=":
		return NewBoolean(!StrictEquals(left, right)), nil
	case "<":
		ln := ToNumber(left)
		rn := ToNumber(right)
		if math.IsNaN(ln.Number()) || math.IsNaN(rn.Number()) {
			return NewBoolean(false), nil
		}
		return NewBoolean(ln.Number() < rn.Number()), nil
	case "<=":
		ln := ToNumber(left)
		rn := ToNumber(right)
		if math.IsNaN(ln.Number()) || math.IsNaN(rn.Number()) {
			return NewBoolean(false), nil
		}
		return NewBoolean(ln.Number() <= rn.Number()), nil
	case ">":
		ln := ToNumber(left)
		rn := ToNumber(right)
		if math.IsNaN(ln.Number()) || math.IsNaN(rn.Number()) {
			return NewBoolean(false), nil
		}
		return NewBoolean(ln.Number() > rn.Number()), nil
	case ">=":
		ln := ToNumber(left)
		rn := ToNumber(right)
		if math.IsNaN(ln.Number()) || math.IsNaN(rn.Number()) {
			return NewBoolean(false), nil
		}
		return NewBoolean(ln.Number() >= rn.Number()), nil
	default:
		return Value{}, fmt.Errorf("runtime error: binary operator %q not implemented", op)
	}
}

func (i *Interpreter) typeOfValue(v Value) string {
	switch v.Kind() {
	case UndefinedKind:
		return "undefined"
	case NullKind:
		return "object"
	case BooleanKind:
		return "boolean"
	case NumberKind:
		return "number"
	case StringKind:
		return "string"
	default:
		return "object"
	}
}

func (i *Interpreter) describeCompletion(c completion) string {
	switch c.kind {
	case completionBreak:
		return "break"
	case completionContinue:
		return "continue"
	case completionReturn:
		return "return"
	default:
		return "normal"
	}
}

func parseNumericLiteral(raw string) (float64, error) {
	s := strings.ReplaceAll(raw, "_", "")
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		v, err := strconv.ParseUint(s[2:], 16, 64)
		if err != nil {
			return 0, err
		}
		return float64(v), nil
	}
	if strings.HasPrefix(s, "0o") || strings.HasPrefix(s, "0O") {
		v, err := strconv.ParseUint(s[2:], 8, 64)
		if err != nil {
			return 0, err
		}
		return float64(v), nil
	}
	if strings.HasPrefix(s, "0b") || strings.HasPrefix(s, "0B") {
		v, err := strconv.ParseUint(s[2:], 2, 64)
		if err != nil {
			return 0, err
		}
		return float64(v), nil
	}
	if strings.HasPrefix(s, "0") && len(s) > 1 && s[1] >= '0' && s[1] <= '7' && !strings.ContainsAny(s, "89") {
		v, err := strconv.ParseUint(s[1:], 8, 64)
		if err == nil {
			return float64(v), nil
		}
	}
	if strings.HasSuffix(s, "n") {
		return 0, fmt.Errorf("bigint literals are not supported")
	}
	return strconv.ParseFloat(s, 64)
}
