package main

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
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
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

func makeChunk() Chunk {
	const initCapacity = 100
	return Chunk{
		make([]uint8, 0, initCapacity),
		make([]Value, 0, initCapacity),
		make([]int, 0, initCapacity),
	}
}

func (chunk *Chunk) simpleInstruction(label string, offset int) int {
	fmt.Println(label)
	return offset + 1
}

func (chunk *Chunk) constantInstruction(name string, offset int) int {
	constant := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d '", name, constant)
	printValue(chunk.Constants[constant])
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
		offset = chunk.constantInstruction("OP_CONSTANT", offset)
	case OP_ADD:
		offset = chunk.simpleInstruction("ADD", offset)
	case OP_SUBTRACT:
		offset = chunk.simpleInstruction("OP_SUBTRACT", offset)
	case OP_MULTIPLY:
		offset = chunk.simpleInstruction("OP_MULTIPLY", offset)
	case OP_DIVIDE:
		offset = chunk.simpleInstruction("OP_DIVIDE", offset)
	case OP_NEGATE:
		offset = chunk.simpleInstruction("OP_NEGATE", offset)
	case OP_RETURN:
		offset = chunk.simpleInstruction("OP_RETURN", offset)
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

func (chunk *Chunk) addConstant(x Value) uint8 {
	chunk.Constants = append(chunk.Constants, x)
	return uint8(len(chunk.Constants)) - 1
}

func parseByteCode(r io.Reader) (Chunk, error) {
	chunk := makeChunk()
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
			chunk.Constants = append(chunk.Constants, Value(num))
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
