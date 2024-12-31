package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jeroendm/glox/chunk"
	"github.com/jeroendm/glox/compiler"
	"github.com/jeroendm/glox/vm"
)

func main() {
	args := os.Args[1:]
	if len(args) > 2 {
		const help = `Usage: glox [-b] [script]
  -b      Run byte code file.
  script  Filename for lox or bytecode file.`
		fmt.Println(help)
		os.Exit(64)
	} else if len(args) == 1 {
		runFile(args[0])
		fmt.Println()
	} else if len(args) == 2 {
		runByteCode(args[1])
	} else {
		runPrompt()
	}
}

func runFile(filename string) {
	fmt.Println("Running file: ", filename)
	content, e := os.ReadFile(filename)
	if e != nil {
		fmt.Printf("ERROR: %s", e)
		fmt.Printf("Failed to open file: '%s'\n", filename)
	}
	run(content)
}

func runByteCode(filename string) {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error opening file: %s", err)
		return
	}
	defer file.Close()

	chunk, err := chunk.ParseByteCode(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	chunk.Disassemble(filename)
	fmt.Printf("\n --- running ---\n")
	vm := vm.MakeVM()
	vm.InterpretChunk(&chunk)
}

func runPrompt() {
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, e := input.ReadString('\n')
		if e != nil {
			panic(e)
		}
		if text == "\n" {
			break
		}
		run([]uint8(text))
	}
}

func run(source []uint8) {
	c := chunk.MakeChunk()

	hadError := compiler.Compile(source, &c)
	if hadError {
		panic("Failed to compile.")
	}

	vm1 := vm.MakeVM()
	err := vm1.Interpret(&c)
	if err != nil {
		panic("Runtime error.")
	}
}
