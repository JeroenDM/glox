package chunk

import (
	"fmt"
	"unsafe"
)

type Number float64

type ValueKind byte

const (
	VAL_BOOL ValueKind = iota
	VAL_NIL
	VAL_NUMBER
)

type Value struct {
	kind ValueKind
	data unsafe.Pointer
}

func NewBool(value bool) Value {
	return Value{
		kind: VAL_BOOL,
		data: unsafe.Pointer(&value),
	}
}

func NewNumber(value Number) Value {
	return Value{
		kind: VAL_NUMBER,
		data: unsafe.Pointer(&value),
	}
}

func NewNil() Value {
	return Value{kind: VAL_NIL}
}

func (v Value) AsBool() bool {
	if !v.IsBool() {
		panic("Value is not a boolean.")
	}
	return *(*bool)(v.data)
}

func (v Value) AsNumber() Number {
	if !v.IsNumber() {
		panic("Value is not a number.")
	}
	return *(*Number)(v.data)
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
