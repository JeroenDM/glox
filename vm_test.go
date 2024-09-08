package main

import (
	"strings"
	"testing"
)

func TestByteCodeSmall(t *testing.T) {
	const asm = `.data
1.2
3.4
5.6

.text
constant 0
constant 1
add
constant 2
divide
negate
return`

	chunk, err := parseByteCode(strings.NewReader(asm))
	if err != nil {
		t.Fatal(err)
	}
	vm := MakeVM()
	if err := vm.InterpretChunk(&chunk); err != nil {
		t.Fatal(err)
	}
}
