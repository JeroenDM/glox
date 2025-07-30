package chunk

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/huandu/go-assert"
)

func TestValuesEqualSameType(t *testing.T) {
	assert.AssertEqual(t, ValuesEqual(NewNil(), NewNil()), true)
	assert.AssertEqual(t, ValuesEqual(NewBool(true), NewBool(true)), true)
	assert.AssertEqual(t, ValuesEqual(NewBool(true), NewBool(false)), false)
	assert.AssertEqual(t, ValuesEqual(NewNumber(3.0), NewNumber(3.0)), true)
	assert.AssertEqual(t, ValuesEqual(NewNumber(3.0), NewNumber(4.0)), false)
	assert.AssertEqual(t, ValuesEqual(NewObjString([]byte("hello")), NewObjString([]byte("hello"))), true)
	assert.AssertEqual(t, ValuesEqual(NewObjString([]byte("a")), NewObjString([]byte("aa"))), false)
	assert.AssertEqual(t, ValuesEqual(NewObjString([]byte("hello")), NewObjString([]byte("dummy"))), false)
}

func TestValuesEqualDifferentType(t *testing.T) {
	s1 := NewObjString([]byte("hello"))
	assert.AssertEqual(t, ValuesEqual(NewNil(), NewBool(true)), false)
	assert.AssertEqual(t, ValuesEqual(NewNil(), NewNumber(3.0)), false)
	assert.AssertEqual(t, ValuesEqual(NewNil(), s1), false)

	assert.AssertEqual(t, ValuesEqual(NewBool(true), NewNumber(4.0)), false)
	assert.AssertEqual(t, ValuesEqual(NewBool(true), s1), false)

	assert.AssertEqual(t, ValuesEqual(NewNumber(4.0), s1), false)
}

func TestObjValue(t *testing.T) {
	obj1 := Obj{kind: OBJ_STRING}
	value1 := NewObj(&obj1)
	assert.AssertEqual(t, value1.ObjKind(), OBJ_STRING)
	assert.Assert(t, value1.IsString())
}

func TestDoublePtr(t *testing.T) {
	num1 := Number(3.0)
	value1 := NewNumber(num1)
	num1 = Number(4.0)
	// The pointer inside 'value1' does not point to the original number.
	// This is as indented in the book, where 'primitive' types are
	// stored by value in the union, only for objects to we store a pointer
	// to the actual object, not to a copy.
	assert.AssertEqual(t, fmt.Sprintf("%.1f", value1.AsNumber()), "3.0")
}

func TestCopyString(t *testing.T) {
	s1 := []byte("hello")
	obj_str := CopyString(s1)
	value := NewObj((*Obj)(unsafe.Pointer(&obj_str)))
	assert.Assert(t, value.IsString())
	assert.AssertEqual(t, value.AsGoString(), "hello")
	s1[0] = 'T'
	assert.AssertEqual(t, value.AsGoString(), "hello")
}

func TestTakeString(t *testing.T) {
	s1 := []byte("hello")
	obj_str := TakeString(s1)
	value := NewObj((*Obj)(unsafe.Pointer(&obj_str)))
	assert.Assert(t, value.IsString())
	assert.AssertEqual(t, value.AsGoString(), "hello")
	s1[0] = 'T'
	assert.AssertEqual(t, value.AsGoString(), "Tello")
}
