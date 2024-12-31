package chunk

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type OpCode uint8

const (
	OP_CONSTANT OpCode = iota
	OP_NIL             // push nil literal on stack
	OP_TRUE            // push true literal on stack
	OP_FALSE           // push false literal on stack
	OP_EQUAL
	OP_GREATER
	OP_LESS
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_NOT
	OP_NEGATE
	OP_RETURN
)

type Chunk struct {
	Code      []uint8
	Constants []Value
	Lines     []int
}

func (chunk *Chunk) Write(code uint8, line int) {
	chunk.Code = append(chunk.Code, uint8(code))
	chunk.Lines = append(chunk.Lines, line)
}

func MakeChunk() Chunk {
	const initCapacity = 100
	return Chunk{
		make([]uint8, 0, initCapacity),
		make([]Value, 0, initCapacity),
		make([]int, 0, initCapacity),
	}
}

func (chunk *Chunk) printSimpleInstruction(label string, offset int) int {
	fmt.Println(label)
	return offset + 1
}

func (chunk *Chunk) printConstantInstruction(name string, offset int) int {
	constant := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d '", name, constant)
	PrintValue(chunk.Constants[constant])
	fmt.Printf("'\n")
	return offset + 2
}

func (chunk *Chunk) DisassembleInstruction(offset int) int {
	// fmt.Printf("constants: %v\n", chunk.Constants)
	fmt.Printf("%04d ", offset)

	if offset > 0 && chunk.Lines[offset] == chunk.Lines[offset-1] {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", chunk.Lines[offset])
	}

	c := OpCode(chunk.Code[offset])
	switch c {
	case OP_CONSTANT:
		offset = chunk.printConstantInstruction("OP_CONSTANT", offset)
	case OP_NIL:
		offset = chunk.printSimpleInstruction("OP_NIL", offset)
	case OP_TRUE:
		offset = chunk.printSimpleInstruction("OP_TRUE", offset)
	case OP_FALSE:
		offset = chunk.printSimpleInstruction("OP_FALSE", offset)
	case OP_EQUAL:
		offset = chunk.printSimpleInstruction("OP_EQUAL", offset)
	case OP_GREATER:
		offset = chunk.printSimpleInstruction("OP_GREATER", offset)
	case OP_LESS:
		offset = chunk.printSimpleInstruction("OP_LESS", offset)
	case OP_ADD:
		offset = chunk.printSimpleInstruction("ADD", offset)
	case OP_SUBTRACT:
		offset = chunk.printSimpleInstruction("OP_SUBTRACT", offset)
	case OP_MULTIPLY:
		offset = chunk.printSimpleInstruction("OP_MULTIPLY", offset)
	case OP_DIVIDE:
		offset = chunk.printSimpleInstruction("OP_DIVIDE", offset)
	case OP_NOT:
		offset = chunk.printSimpleInstruction("OP_NOT", offset)
	case OP_NEGATE:
		offset = chunk.printSimpleInstruction("OP_NEGATE", offset)
	case OP_RETURN:
		offset = chunk.printSimpleInstruction("OP_RETURN", offset)
	default:
		fmt.Printf("Unknown opcode %d\n", c)
		offset += 1
	}
	return offset
}

func (chunk *Chunk) Disassemble(name string) {
	fmt.Printf("== %s ==\n", name)
	for offset := 0; offset < len(chunk.Code); {
		offset = chunk.DisassembleInstruction(offset)
	}

}

func (chunk *Chunk) AddConstant(x Value) uint8 {
	chunk.Constants = append(chunk.Constants, x)
	return uint8(len(chunk.Constants)) - 1
}

func ParseByteCode(r io.Reader) (Chunk, error) {
	chunk := MakeChunk()
	scanner := bufio.NewScanner(r)
	section := ""
	lineNumber := -1
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, ".") {
			section = line[1:]
			continue
		}

		switch section {
		case "data":
			num, err := strconv.ParseFloat(line, 64)
			if err != nil {
				return chunk, err
			}
			chunk.Constants = append(chunk.Constants, NewNumber(Number(num)))
		case "text":
			parts := strings.Split(line, " ")
			switch parts[0] {
			case "constant":
				if len(parts) != 2 {
					return chunk, fmt.Errorf("wrong number of arguments for constant instruction, expected %d, got %d", 1, len(parts)-1)
				}
				chunk.Write(uint8(OP_CONSTANT), lineNumber)
				c, err := strconv.ParseUint(parts[1], 10, 8)
				if err != nil {
					return chunk, err
				}
				chunk.Write(uint8(c), lineNumber)
			case "add":
				chunk.Write(uint8(OP_ADD), lineNumber)
			case "subtract":
				chunk.Write(uint8(OP_SUBTRACT), lineNumber)
			case "multiply":
				chunk.Write(uint8(OP_MULTIPLY), lineNumber)
			case "divide":
				chunk.Write(uint8(OP_DIVIDE), lineNumber)
			case "negate":
				chunk.Write(uint8(OP_NEGATE), lineNumber)
			case "return":
				chunk.Write(uint8(OP_RETURN), lineNumber)
			default:
				return chunk, fmt.Errorf("unknown instruction %s", line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return chunk, fmt.Errorf("error reading file: %s", err)
	}

	return chunk, nil
}
