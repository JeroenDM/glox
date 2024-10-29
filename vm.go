package main

import (
	"fmt"
)

type InterpretError uint8

type BinaryOp func(Value, Value) Value

const (
	INTERPRET_COMPILE_ERROR InterpretError = iota
	INTERPRET_RUNTIME_ERROR
)

func (e InterpretError) Error() string {
	switch e {
	case INTERPRET_COMPILE_ERROR:
		return "compile error"
	case INTERPRET_RUNTIME_ERROR:
		return "runtime error"
	}
	// Can this happen?
	return "unknown interpret error"
}

const STACK_MAX = 256

type VM struct {
	chunk    *Chunk
	ip       int
	stack    []Value
	stackTop uint8
}

func MakeVM() VM {
	return VM{nil, 0, make([]Value, STACK_MAX), 0}
}

func (vm *VM) InterpretChunk(chunk *Chunk) error {
	vm.chunk = chunk
	vm.ip = 0
	return vm.run()
}

func (vm *VM) Interpret(source []uint8) error {
	c := makeChunk()

	err := compile(source, &c)
	if err != nil {
		return INTERPRET_COMPILE_ERROR
	}
	return nil

	// vm.chunk = &c
	// vm.ip = 0
	// return vm.run()
}

func (vm *VM) run() error {
	for {
		traceInstruction(vm, vm.ip)
		switch OpCode(vm.readByte()) {
		case OP_CONSTANT:
			constant := vm.readConstant()
			vm.push(constant)
		case OP_NEGATE:
			vm.push(-vm.pop())
		case OP_ADD:
			vm.binary(func(a Value, b Value) Value { return a + b })
		case OP_SUBTRACT:
			vm.binary(func(a Value, b Value) Value { return a - b })
		case OP_MULTIPLY:
			vm.binary(func(a Value, b Value) Value { return a * b })
		case OP_DIVIDE:
			vm.binary(func(a Value, b Value) Value { return a / b })
		case OP_RETURN:
			printValue(vm.pop())
			fmt.Printf("\n")
			return nil
		default:

		}
	}
}

func (vm *VM) readByte() uint8 {
	i := vm.ip
	vm.ip++
	return vm.chunk.Code[i]
}

func (vm *VM) readConstant() Value {
	value := vm.chunk.Constants[vm.readByte()]
	return value
}

func (vm *VM) resetStack() {
	vm.stackTop = 0
}

func (vm *VM) push(value Value) {
	vm.stack[vm.stackTop] = value
	vm.stackTop++
}

func (vm *VM) pop() Value {
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

func (vm *VM) binary(op BinaryOp) {
	// Order of pops is important!
	b := vm.pop()
	a := vm.pop()
	vm.push(op(a, b))
}
