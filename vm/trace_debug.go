//go:build debug
// +build debug

package vm

import (
	"fmt"

	"github.com/jeroendm/glox/chunk"
)

func traceInstruction(vm *VM, offset int) {
	fmt.Printf("          ")
	for _, value := range vm.stack[:vm.stackTop] {
		fmt.Printf("[ ")
		chunk.PrintValue(value)
		fmt.Printf(" ]")
	}
	fmt.Printf("\n")
	vm.chunk.DisassembleInstruction(offset)
}
