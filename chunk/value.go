package chunk

import (
	"bytes"
	"fmt"
	"unsafe"
)

type Number float64

type ValueKind byte

const (
	VAL_BOOL ValueKind = iota
	VAL_NIL
	VAL_NUMBER
	VAL_OBJ
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

func NewObj(value *Obj) Value {
	return Value{
		kind: VAL_OBJ,
		data: unsafe.Pointer(value),
	}
}

func NewObjString(s []byte) Value {
	obj_str := CopyString(s)
	return NewObj((*Obj)(unsafe.Pointer(&obj_str)))
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

func (v Value) AsObj() Obj {
	if !v.IsObj() {
		panic("Value is not an object.")
	}
	return *(*Obj)(v.data)
}

func (v Value) AsString() *ObjString {
	if !v.IsString() {
		panic("Value is not a string.")
	}
	return (*ObjString)(v.data)
}

func (v Value) AsGoString() string {
	if !v.IsString() {
		panic("Value is not a string.")
	}
	obj_str_ptr := (*ObjString)(v.data)
	return string(obj_str_ptr.Bytes)
}

func (v Value) ObjKind() ObjKind {
	return v.AsObj().kind
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

func (v Value) IsObj() bool {
	return v.kind == VAL_OBJ
}

func (v Value) IsString() bool {
	return v.IsObj() && (v.ObjKind() == OBJ_STRING)
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
	case VAL_OBJ:
		s1 := a.AsString()
		s2 := b.AsString()
		return s1.Length == s2.Length && bytes.Equal(s1.Bytes, s2.Bytes)
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
	case VAL_OBJ:
		// TODO check obj kind.
		fmt.Printf("%s", x.AsGoString())
	default:
		panic("Unknown value type.")
	}
}
