package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const MEM_SIZE int = 30000

type bfOperation func()

type Brainfckr struct {
	reader     io.Reader
	writer     io.Writer
	stack      Stack
	code       []byte
	codePtr    int
	mem        []byte
	memPtr     int
	executed   []byte
	operations map[byte]bfOperation
}

func NewBrainfckr(reader io.Reader, writer io.Writer) *Brainfckr {

	bf := new(Brainfckr)

	bf.opsMapSetup()
	bf.mem = make([]byte, MEM_SIZE)
	bf.reader = reader
	bf.writer = writer

	return bf
}

func (bf *Brainfckr) opsMapSetup() {
	bf.operations = make(map[byte]bfOperation)
	bf.operations['+'] = func() {
		bf.mem[bf.memPtr]++
		bf.executed = append(bf.executed, '+')
	}
	bf.operations['-'] = func() {
		bf.mem[bf.memPtr]--
		bf.executed = append(bf.executed, '-')
	}
	bf.operations['>'] = func() {
		bf.memPtr++
		bf.executed = append(bf.executed, '>')
	}
	bf.operations['<'] = func() {
		if bf.memPtr > 0 {
			bf.memPtr--
			bf.executed = append(bf.executed, '<')
		}
	}
	bf.operations[','] = func() {
		bf.executed = append(bf.executed, ',')
		bf.mem[bf.memPtr], _ = bufio.NewReader(os.Stdin).ReadByte()
		bf.mem[bf.memPtr] = bf.mem[bf.memPtr]
	}
	bf.operations['.'] = func() {
		bf.executed = append(bf.executed, '.')
		bf.print(bf.mem[bf.memPtr])
	}
	bf.operations['['] = func() {
		bf.loop()
	}
	bf.operations[']'] = func() { // does nothing, just there to be a validOp
		bf.executed = append(bf.executed, ']')
	}
}

func (bf *Brainfckr) Interpret() error {
	for op, err := bf.nextOp(); err != io.EOF; op, err = bf.nextOp() {
		bf.operations[op]()
	}
	return io.EOF
}

func (bf *Brainfckr) loop() {
	var op byte = '['
	bf.code = append(bf.code, op)
	for {
		if op == '[' {
			bf.stack = bf.stack.Push(bf.codePtr)
			if bf.mem[bf.memPtr] == 0 { // skip the loop
				imbalanceCount := 1
				for imbalanceCount > 0 {
					op = bf.nextLoopOp()
					if bf.code[bf.codePtr] == ']' {
						imbalanceCount--
					} else if bf.code[bf.codePtr] == '[' {
						imbalanceCount++
					}
				}
				continue
			}
		} else if op == ']' {
			var tmpCodePtr int
			bf.stack, tmpCodePtr = bf.stack.Pop()
			if bf.mem[bf.memPtr] > 0 {
				bf.codePtr = tmpCodePtr - 1
			} else if bf.stack.IsEmpty() {
				break
			}
		} else {
			bf.operations[op]()
		}
		op = bf.nextLoopOp()
	}
	bf.code = nil
	bf.codePtr = 0
}

func (bf *Brainfckr) nextLoopOp() byte {
	var op byte
	bf.codePtr++
	if bf.codePtr == len(bf.code) {
		op, _ = bf.nextOp()
		bf.code = append(bf.code, op)
	} else {
		op = bf.code[bf.codePtr]
	}
	return op
}

func (bf *Brainfckr) nextOp() (byte, error) {
	b := make([]byte, 1)
	var err error
	for {
		_, err = bf.reader.Read(b)
		if err != nil {
			if err == io.EOF {
				bf.handleEOF()
			}
		}
		if _, validOp := bf.operations[b[0]]; validOp {
			break
		} else if b[0] == '$' {
			bf.debugOutput("\ndebug;\n")
		}
	}
	return b[0], err
}

func (bf *Brainfckr) handleEOF() {
	bf.print(0x0A)
	bf.print('E')
	bf.print('O')
	bf.print('F')
	bf.print(0x0A)
	os.Exit(0)
}

func (bf *Brainfckr) print(b byte) {
	f := bufio.NewWriter(bf.writer)
	f.WriteByte(b)
	f.Flush()
	f = nil
}

func (bf *Brainfckr) debugOutput(msg string) {
	fmt.Printf(msg)
	fmt.Printf("Stack Size %d\n", len(bf.stack))
	fmt.Printf("CodePtr %d\n", bf.codePtr)
	fmt.Printf("Code %s\n", bf.code)
	fmt.Printf("Executed %s\n", bf.executed)
	fmt.Printf("MemPtr %d\n", bf.memPtr)
	fmt.Printf("Memory Content % d\n", bf.mem[:200])
}
