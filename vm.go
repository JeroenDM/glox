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

// One less than 256, stackTop points to the next empty element,
// therefore we need 256 as a token to signify that the stack is full,
// otherwise the uint8 would overflow to zero and it would look like the stack is empty.
const STACK_MAX = 255

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

	hadError := compile(source, &c)
	if hadError {
		return INTERPRET_COMPILE_ERROR
	}

	c.Write(uint8(OP_CONSTANT), 0)
	c.Write(c.addConstant(Value(1234.5)), 0)
	c.Write(uint8(OP_RETURN), 1)

	vm.chunk = &c
	vm.ip = 0
	return vm.run()
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
	if vm.stackTop == STACK_MAX {
		msg := fmt.Sprintf("Stack overflow! Max stack size (%d) reached.", STACK_MAX)
		panic(msg)
	}
	vm.stack[vm.stackTop] = value
	vm.stackTop++
}

func (vm *VM) pop() Value {
	if vm.stackTop == 0 {
		panic("Cannot pop value from an empty stack.")
	}
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

func (vm *VM) binary(op BinaryOp) {
	// Order of pops is important!
	b := vm.pop()
	a := vm.pop()
	vm.push(op(a, b))
}
