package main

import (
	"fmt"
	"os"
)

type InterpretError uint8

type BinaryOp func(Number, Number) Number

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

	vm.chunk = &c
	vm.ip = 0
	return vm.run()
}

func (vm *VM) run() error {
	for {
		traceInstruction(vm, vm.ip)
		var err error
		switch OpCode(vm.readByte()) {
		case OP_CONSTANT:
			constant := vm.readConstant()
			vm.push(constant)
		case OP_NEGATE:
			if !(vm.peek(0).IsNumber()) {
				vm.runtimeError("Operand must be a number.")
				err = INTERPRET_RUNTIME_ERROR
			}
			vm.push(NewNumber(-vm.pop().AsNumber()))
		case OP_NIL:
			vm.push(NewNil())
		case OP_TRUE:
			vm.push(NewBool(true))
		case OP_FALSE:
			vm.push(NewBool(false))
		case OP_ADD:
			err = vm.binary(NewNumber, func(a Number, b Number) Number { return a + b })
		case OP_SUBTRACT:
			err = vm.binary(NewNumber, func(a Number, b Number) Number { return a - b })
		case OP_MULTIPLY:
			err = vm.binary(NewNumber, func(a Number, b Number) Number { return a * b })
		case OP_DIVIDE:
			err = vm.binary(NewNumber, func(a Number, b Number) Number { return a / b })
		case OP_RETURN:
			printValue(vm.pop())
			fmt.Printf("\n")
			return nil
		default:
			panic("Unknown opcode.")
		}
		// Break from for loop if we have an error.
		if err != nil {
			return err
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

func (vm *VM) runtimeError(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintf(os.Stderr, "\n")

	line := vm.chunk.Lines[vm.ip-1]
	fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
	vm.resetStack()
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

func (vm *VM) peek(distance uint8) Value {
	return vm.stack[vm.stackTop-distance-1]
}

// TODO binary op should probably be generic to support multiple value types.
func (vm *VM) binary(toValue func(Number) Value, op BinaryOp) error {
	if !vm.peek(0).IsNumber() || !vm.peek(1).IsNumber() {
		vm.runtimeError("Operands must be numbers.")
		return INTERPRET_RUNTIME_ERROR
	}
	// Order of pops is important!
	b := vm.pop().AsNumber()
	a := vm.pop().AsNumber()
	vm.push(toValue(op(a, b)))

	return nil
}
