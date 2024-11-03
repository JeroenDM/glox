package main

import "testing"

func TestSmallExpression(t *testing.T) {
	chunk := makeChunk()
	source := "-1"
	hasError := compile([]byte(source), &chunk)
	if hasError {
		t.Fatal("failed to compile")
	}
}
