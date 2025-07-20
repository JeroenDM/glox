package chunk

import (
	"testing"

	"github.com/huandu/go-assert"
)

func TestValuesEqualSameType(t *testing.T) {
	assert.AssertEqual(t, ValuesEqual(NewNil(), NewNil()), true)
	assert.AssertEqual(t, ValuesEqual(NewBool(true), NewBool(true)), true)
	assert.AssertEqual(t, ValuesEqual(NewBool(true), NewBool(false)), false)
	assert.AssertEqual(t, ValuesEqual(NewNumber(3.0), NewNumber(3.0)), true)
	assert.AssertEqual(t, ValuesEqual(NewNumber(3.0), NewNumber(4.0)), false)
}

func TestValuesEqualDifferentType(t *testing.T) {
	assert.AssertEqual(t, ValuesEqual(NewNil(), NewBool(true)), false)
	assert.AssertEqual(t, ValuesEqual(NewNil(), NewNumber(3.0)), false)
	assert.AssertEqual(t, ValuesEqual(NewBool(true), NewNumber(4.0)), false)
}
