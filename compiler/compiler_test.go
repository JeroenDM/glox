package compiler

import (
	"testing"

	"github.com/jeroendm/glox/chunk"
)

func TestSmallExpression(t *testing.T) {
	chunk := chunk.MakeChunk()
	source := "-1;"
	hasError := Compile([]byte(source), &chunk)
	if hasError {
		t.Fatal("failed to compile")
	}
}
