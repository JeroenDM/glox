package chunk

import "fmt"

type Number float64

type ValueKind byte

const (
	VAL_BOOL ValueKind = iota
	VAL_NIL
	VAL_NUMBER
)

type Value struct {
	kind ValueKind
	data Number
}

func NewBool(value bool) Value {
	return Value{
		kind: VAL_BOOL,
		data: bool2Number(value),
	}
}

func NewNumber(value Number) Value {
	return Value{
		kind: VAL_NUMBER,
		data: value,
	}
}

func NewNil() Value {
	return Value{kind: VAL_NIL}
}

func (v Value) AsBool() bool {
	if !v.IsBool() {
		panic("Value is not a boolean.")
	}
	return number2Bool(v.data)
}

func (v Value) AsNumber() Number {
	if !v.IsNumber() {
		panic("Value is not a number.")
	}
	return v.data
}

func (v Value) IsBool() bool {
	return v.kind == VAL_BOOL
}

func (v Value) IsNumber() bool {
	return v.kind == VAL_NUMBER
}

func (v Value) IsNil() bool {
	return v.kind == VAL_NIL
}

func bool2Number(b bool) Number {
	// The compiler currently only optimizes this form.
	// https://dev.to/chigbeef_77/bool-int-but-stupid-in-go-3jb3
	var i Number
	if b {
		i = 1.0
	} else {
		i = 0.0
	}
	return i
}

func number2Bool(f Number) bool {
	return !(f == 0.0)
}

func ValuesEqual(a, b Value) bool {
	if a.kind != b.kind {
		return false
	}
	switch a.kind {
	case VAL_BOOL:
		return a.AsBool() == b.AsBool()
	case VAL_NIL:
		return true
	case VAL_NUMBER:
		return a.AsNumber() == b.AsNumber()
	default:
		panic("Should be unreachable (valuesEqual).")
	}
}

func PrintValue(x Value) {
	switch x.kind {
	case VAL_BOOL:
		fmt.Printf("%t", x.AsBool())
	case VAL_NIL:
		fmt.Printf("nil")
	case VAL_NUMBER:
		fmt.Printf("%g", x.AsNumber())
	default:
		panic("Unknown value type.")
	}
}
