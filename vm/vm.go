package vm

import (
	"fmt"
	"os"

	"github.com/jeroendm/glox/chunk"
)

type InterpretError uint8

type BinaryOp func(chunk.Number, chunk.Number) chunk.Number

const (
	INTERPRET_COMPILE_ERROR InterpretError = iota
	INTERPRET_RUNTIME_ERROR
)

var (
	LESS = func(a chunk.Number, b chunk.Number) bool { return a > b }
	GREATER = func(a chunk.Number, b chunk.Number) bool { return a > b }
)

var (
	PLUS = func(a chunk.Number, b chunk.Number) chunk.Number { return a + b }
	SUBTRACT = func(a chunk.Number, b chunk.Number) chunk.Number { return a - b }
	MULTIPLY = func(a chunk.Number, b chunk.Number) chunk.Number { return a * b }
	DIVIDE = func(a chunk.Number, b chunk.Number) chunk.Number { return a / b }
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

func isFalsey(value chunk.Value) bool {
	return value.IsNil() || value.IsBool() && !value.AsBool()
}

// One less than 256, stackTop points to the next empty element,
// therefore we need 256 as a token to signify that the stack is full,
// otherwise the uint8 would overflow to zero and it would look like the stack is empty.
const STACK_MAX = 255

type VM struct {
	chunk    *chunk.Chunk
	ip       int
	stack    []chunk.Value
	stackTop uint8
}

func MakeVM() VM {
	return VM{nil, 0, make([]chunk.Value, STACK_MAX), 0}
}

func (vm *VM) InterpretChunk(chunk *chunk.Chunk) error {
	vm.chunk = chunk
	vm.ip = 0
	return vm.run()
}

func (vm *VM) Interpret(c *chunk.Chunk) error {
	// c := chunk.MakeChunk()

	// hadError := compile(source, &c)
	// if hadError {
	// 	return INTERPRET_COMPILE_ERROR
	// }

	// vm.chunk = &c
	vm.chunk = c
	vm.ip = 0
	return vm.run()
}

func (vm *VM) run() error {
	for {
		traceInstruction(vm, vm.ip)
		var err error
		switch chunk.OpCode(vm.readByte()) {
		case chunk.OP_CONSTANT:
			constant := vm.readConstant()
			vm.push(constant)
		case chunk.OP_NEGATE:
			if !(vm.peek(0).IsNumber()) {
				vm.runtimeError("Operand must be a number.")
				err = INTERPRET_RUNTIME_ERROR
			} else {
				vm.push(chunk.NewNumber(-vm.pop().AsNumber()))
			}
		case chunk.OP_NIL:
			vm.push(chunk.NewNil())
		case chunk.OP_TRUE:
			vm.push(chunk.NewBool(true))
		case chunk.OP_FALSE:
			vm.push(chunk.NewBool(false))
		case chunk.OP_EQUAL:
			b := vm.pop()
			a := vm.pop()
			vm.push(chunk.NewBool(chunk.ValuesEqual(a, b)))
		case chunk.OP_GREATER:
			err = vm.binaryBool(chunk.NewBool, GREATER)
		case chunk.OP_LESS:
			err = vm.binaryBool(chunk.NewBool, LESS)
		case chunk.OP_ADD:
			err = vm.binary(chunk.NewNumber, PLUS)
		case chunk.OP_SUBTRACT:
			err = vm.binary(chunk.NewNumber, SUBTRACT)
		case chunk.OP_MULTIPLY:
			err = vm.binary(chunk.NewNumber, MULTIPLY)
		case chunk.OP_DIVIDE:
			err = vm.binary(chunk.NewNumber, DIVIDE)
		case chunk.OP_NOT:
			vm.push(chunk.NewBool(isFalsey(vm.pop())))
		case chunk.OP_RETURN:
			chunk.PrintValue(vm.pop())
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

func (vm *VM) readConstant() chunk.Value {
	value := vm.chunk.Constants[vm.readByte()]
	return value
}

func (vm *VM) resetStack() {
	vm.stackTop = 0
}

func (vm *VM) runtimeError(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintf(os.Stderr, "\n")

	// Minus one because the interpreter advances past and instruction
	// before executing it.
	line := vm.chunk.Lines[vm.ip-1]
	fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
	vm.resetStack()
}

func (vm *VM) push(value chunk.Value) {
	if vm.stackTop == STACK_MAX {
		msg := fmt.Sprintf("Stack overflow! Max stack size (%d) reached.", STACK_MAX)
		panic(msg)
	}
	vm.stack[vm.stackTop] = value
	vm.stackTop++
}

func (vm *VM) pop() chunk.Value {
	if vm.stackTop == 0 {
		panic("Cannot pop value from an empty stack.")
	}
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

func (vm *VM) peek(distance uint8) chunk.Value {
	return vm.stack[vm.stackTop-distance-1]
}

// TODO binary op should probably be generic to support multiple value types.
func (vm *VM) binary(toValue func(chunk.Number) chunk.Value, op BinaryOp) error {
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

func (vm *VM) binaryBool(toValue func(bool) chunk.Value, op func(chunk.Number, chunk.Number) bool) error {
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
