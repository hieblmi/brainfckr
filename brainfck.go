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
	// code source
	reader io.Reader
	// interpreter output
	writer io.Writer

	// code stack
	stack Stack

	// code segment holding code in loops -> ,[>+<-]>.
	// will be extended if it runs out of space by 2*MAX_FILE_CHUNK_SIZE
	code    []byte
	codePtr int

	// memory
	mem    []byte
	memPtr int

	executed []byte

	// map holding functions per operation
	operations map[byte]bfOperation
}

func NewBrainfckr(reader io.Reader, writer io.Writer) *Brainfckr {
	bf := new(Brainfckr)

	bf.opsMapSetup()
	bf.mem = make([]byte, MEM_SIZE)
	bf.memPtr = 0
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
		f := bufio.NewWriter(bf.writer)
		b := bf.mem[bf.memPtr]
		f.WriteByte(b)
		f.Flush()
	}
	bf.operations['['] = func() {
		bf.loop()
	}
	bf.operations[']'] = func() {
		bf.executed = append(bf.executed, ']')
	}
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

					if bf.codePtr == len(bf.code)-1 {
						op, _ = bf.nextOp()
						bf.code = append(bf.code, op)
					} else {
						op = bf.code[bf.codePtr+1]
					}
					bf.codePtr++

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

		bf.codePtr++
		if bf.codePtr == len(bf.code) {
			op, _ = bf.nextOp()
			bf.code = append(bf.code, op)
		} else {
			op = bf.code[bf.codePtr]
		}
		//bf.debugPrint("")
	}
	//bf.debugPrint("")
	bf.code = nil
	bf.codePtr = 0
}

func (bf *Brainfckr) nextOp() (byte, error) {
	b := make([]byte, 1)
	var err error
	for {
		_, err = bf.reader.Read(b)
		if err == io.EOF {
			os.Exit(0)
		}
		if _, validOp := bf.operations[b[0]]; validOp {
			break
		} else if b[0] == '$' {
			//bf.debugPrint("\ndebug;\n")
		}
	}
	return b[0], err
}

func (bf *Brainfckr) Interpret() error {
	for op, err := bf.nextOp(); err != io.EOF; op, err = bf.nextOp() {
		bf.operations[op]()
	}

	//bf.debugPrint("\n\nEOF\n")
	return io.EOF
}

func (bf *Brainfckr) debugPrint(msg string) {
	fmt.Printf(msg)
	fmt.Printf("Stack Size %d\n", len(bf.stack))
	fmt.Printf("CodePtr %d\n", bf.codePtr)
	fmt.Printf("Code %s\n", bf.code)
	fmt.Printf("Executed %s\n", bf.executed)
	fmt.Printf("MemPtr %d\n", bf.memPtr)
	fmt.Printf("Memory Content % d\n", bf.mem[:200])
}
