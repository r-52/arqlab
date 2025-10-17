package vm

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ValueKind identifies the concrete type of an ECMAScript value.
type ValueKind int

const (
	UndefinedKind ValueKind = iota
	NullKind
	BooleanKind
	NumberKind
	StringKind
)

// Value holds one ECMAScript value. Non-primitive forms (objects, functions,
// arrays) will be modeled in future iterations.
type Value struct {
	kind ValueKind
	num  float64
	str  string
	bool bool
}

// Common singleton values reused across the VM.
var (
	Undefined = Value{kind: UndefinedKind}
	Null      = Value{kind: NullKind}
	True      = Value{kind: BooleanKind, bool: true}
	False     = Value{kind: BooleanKind, bool: false}
)

// NewBoolean returns a boolean value.
func NewBoolean(b bool) Value {
	if b {
		return True
	}
	return False
}

// NewNumber returns a numeric value.
func NewNumber(f float64) Value {
	return Value{kind: NumberKind, num: f}
}

// NewString returns a string value.
func NewString(s string) Value {
	return Value{kind: StringKind, str: s}
}

// Kind exposes the underlying ValueKind.
func (v Value) Kind() ValueKind { return v.kind }

// Bool retrieves the boolean payload, panicking if the kind mismatches.
func (v Value) Bool() bool {
	if v.kind != BooleanKind {
		panic(fmt.Sprintf("vm: Bool() on non-boolean value %s", v.Inspect()))
	}
	return v.bool
}

// Number retrieves the numeric payload, panicking if the kind mismatches.
func (v Value) Number() float64 {
	if v.kind != NumberKind {
		panic(fmt.Sprintf("vm: Number() on non-number value %s", v.Inspect()))
	}
	return v.num
}

// StringValue retrieves the string payload, panicking if the kind mismatches.
func (v Value) StringValue() string {
	if v.kind != StringKind {
		panic(fmt.Sprintf("vm: StringValue() on non-string value %s", v.Inspect()))
	}
	return v.str
}

// String implements fmt.Stringer and returns a descriptive representation.
func (v Value) String() string { return v.Inspect() }

// Inspect returns a descriptive string for debugging.
func (v Value) Inspect() string {
	switch v.kind {
	case UndefinedKind:
		return "undefined"
	case NullKind:
		return "null"
	case BooleanKind:
		if v.bool {
			return "true"
		}
		return "false"
	case NumberKind:
		if math.IsNaN(v.num) {
			return "NaN"
		}
		if math.IsInf(v.num, 1) {
			return "Infinity"
		}
		if math.IsInf(v.num, -1) {
			return "-Infinity"
		}
		return strconv.FormatFloat(v.num, 'g', -1, 64)
	case StringKind:
		return strconv.Quote(v.str)
	default:
		return "<unknown>"
	}
}

// StrictEquals implements the === operator for the supported types.
func StrictEquals(a, b Value) bool {
	if a.kind != b.kind {
		return false
	}
	switch a.kind {
	case UndefinedKind, NullKind:
		return true
	case BooleanKind:
		return a.bool == b.bool
	case NumberKind:
		if math.IsNaN(a.num) || math.IsNaN(b.num) {
			return false
		}
		return a.num == b.num
	case StringKind:
		return a.str == b.str
	default:
		return false
	}
}

// ToBoolean performs JavaScript truthiness conversion for the supported types.
func ToBoolean(v Value) bool {
	switch v.kind {
	case UndefinedKind, NullKind:
		return false
	case BooleanKind:
		return v.bool
	case NumberKind:
		if v.num == 0 || math.IsNaN(v.num) {
			return false
		}
		return true
	case StringKind:
		return len(v.str) > 0
	default:
		return false
	}
}

// ToNumber converts a value to a number following simplified ECMAScript rules.
func ToNumber(v Value) Value {
	switch v.kind {
	case UndefinedKind:
		return NewNumber(math.NaN())
	case NullKind:
		return NewNumber(0)
	case BooleanKind:
		if v.bool {
			return NewNumber(1)
		}
		return NewNumber(0)
	case NumberKind:
		return v
	case StringKind:
		s := strings.TrimSpace(v.str)
		if s == "" {
			return NewNumber(0)
		}
		if strings.EqualFold(s, "NaN") {
			return NewNumber(math.NaN())
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return NewNumber(math.NaN())
		}
		return NewNumber(f)
	default:
		return NewNumber(math.NaN())
	}
}

// ToString converts a value to a string value.
func ToString(v Value) Value {
	switch v.kind {
	case UndefinedKind:
		return NewString("undefined")
	case NullKind:
		return NewString("null")
	case BooleanKind:
		if v.bool {
			return NewString("true")
		}
		return NewString("false")
	case NumberKind:
		if math.IsNaN(v.num) {
			return NewString("NaN")
		}
		if math.IsInf(v.num, 1) {
			return NewString("Infinity")
		}
		if math.IsInf(v.num, -1) {
			return NewString("-Infinity")
		}
		return NewString(strconv.FormatFloat(v.num, 'g', -1, 64))
	case StringKind:
		return v
	default:
		return NewString("<unknown>")
	}
}

// ToPrimitiveNumber prepares a Value for numeric operations by returning the
// float64 representation along with a success flag.
func ToPrimitiveNumber(v Value) (float64, bool) {
	n := ToNumber(v)
	if n.kind != NumberKind {
		return math.NaN(), false
	}
	return n.num, true
}
