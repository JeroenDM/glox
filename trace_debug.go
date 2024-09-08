//go:build debug
// +build debug

package main

import "fmt"

func traceInstruction(vm *VM, offset int) {
	fmt.Printf("          ")
	for _, value := range vm.stack[:vm.stackTop] {
		fmt.Printf("[ ")
		printValue(value)
		fmt.Printf(" ]")
	}
	fmt.Printf("\n")
	vm.chunk.DisassembleInstruction(offset)
}
