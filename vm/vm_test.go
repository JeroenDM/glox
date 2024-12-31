package vm

import (
	"strings"
	"testing"

	"github.com/jeroendm/glox/chunk"
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

	chunk, err := chunk.ParseByteCode(strings.NewReader(asm))
	if err != nil {
		t.Fatal(err)
	}
	vm := MakeVM()
	if err := vm.InterpretChunk(&chunk); err != nil {
		t.Fatal(err)
	}
}
