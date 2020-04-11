package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const MEM_SIZE int = 1024

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
		bf.memPtr--
		bf.executed = append(bf.executed, '<')
	}
	bf.operations[','] = func() {
		bf.executed = append(bf.executed, ',')
		bf.mem[bf.memPtr], _ = bufio.NewReader(os.Stdin).ReadByte()
		bf.mem[bf.memPtr] = bf.mem[bf.memPtr] - 48
	}
	bf.operations['.'] = func() {
		bf.executed = append(bf.executed, '.')
		/*f := bufio.NewWriter(bf.writer)
		f.WriteByte(bf.mem[bf.memPtr])
		f.Flush()*/
		//fmt.Printf("Printing bf.mem[bf.memPtr]=%d\n", bf.mem[bf.memPtr])
		//f = nil
		//bf.writer.Write(bf.mem[1:2])
		b := bf.mem[bf.memPtr]
		fmt.Printf("%c", b)
		/*if !unicode.IsSpace(rune(bf.mem[bf.memPtr])) {

		}*/
	}
	bf.operations['['] = func() {
		bf.executed = append(bf.executed, '[')
		if bf.mem[bf.memPtr] == 0 {
			bf.skipLoop()
		} else {
			bf.code = append(bf.code, '[')
			bf.stack = bf.stack.Push(0)
			bf.loop()
			//bf.debugPrint(fmt.Sprintf("#######END OF MAIN LOOP\n"))
		}
	}
	bf.operations[']'] = func() {
		bf.executed = append(bf.executed, ']')
	}
}

func (bf *Brainfckr) loop() {
	readFromCodeSegment := false
	var op byte
	for {
		if readFromCodeSegment {
			op = bf.code[bf.codePtr]
			bf.codePtr += 1
		} else {
			op, _ = bf.nextOp()
			bf.code = append(bf.code, op)
			//fmt.Printf("Added to code segment %s\n", bf.code)

		}

		if op == '[' {
			if bf.mem[bf.memPtr] > 0 {
				if readFromCodeSegment {
					bf.stack = bf.stack.Push(bf.codePtr - 1)
				} else {
					bf.stack = bf.stack.Push(len(bf.code) - 1)
				}
			} else {
				bf.skipLoop()
			}
			continue
		} else if op == ']' {
			prev := bf.codePtr
			bf.stack, bf.codePtr = bf.stack.Pop()
			if bf.mem[bf.memPtr] > 0 {
				//bf.debugPrint(fmt.Sprintf("Mempointer content end of loop %d\n", bf.mem[bf.memPtr]))
				readFromCodeSegment = true
			} else {
				//bf.debugPrint(fmt.Sprintf("End of Loop: Stack Size: %d Code Segment: %s Code Pointer: %d Prev: %d\n", len(bf.stack), bf.code, bf.codePtr, prev))
				if bf.stack.IsEmpty() {
					//fmt.Println("Loop finished, break")
					break
				}
				if prev < len(bf.code) {
					readFromCodeSegment = true
					bf.codePtr = prev
				} else {
					readFromCodeSegment = false
				}
			}
		} else {
			bf.operations[op]()
			//bf.debugPrint(fmt.Sprintf("Executed %c\n", op))
		}
	}
	//bf.debugPrint("Memory Dump")
	bf.code = nil
	bf.codePtr = 0
}

func (bf *Brainfckr) skipLoop() {
	numberOfOpenBrackets := 1
	//bf.debugPrint("SKIPPING LOOP START\n")
	for op, _ := bf.nextOp(); numberOfOpenBrackets > 0; {
		op, _ = bf.nextOp()
		if op == '[' {
			numberOfOpenBrackets++
		} else if op == ']' {
			numberOfOpenBrackets--
		}
	}
	//bf.debugPrint("SKIPPING LOOP END\n")
}

func (bf *Brainfckr) nextOp() (byte, error) {
	b := make([]byte, 1)
	var err error
	for {
		_, err = bf.reader.Read(b)
		if err == io.EOF {
			fmt.Println("EOF")
			os.Exit(0)
		}
		if _, validOp := bf.operations[b[0]]; validOp {
			break
		} else if b[0] == '$' {
			//bf.debugPrint("\ndebug;\n")
		}
		//fmt.Println("Invalid Op %c\n", b[0])
	}
	//fmt.Println("VALID Op %c %d\n", b[0], err)
	return b[0], err
}

func (bf *Brainfckr) Interpret() error {
	for op, err := bf.nextOp(); err != io.EOF; op, err = bf.nextOp() {
		bf.operations[op]()
	}

	return io.EOF
}

func (bf *Brainfckr) debugPrint(msg string) {
	fmt.Printf(msg)
	fmt.Printf("Stack Size %d\n", len(bf.stack))
	fmt.Printf("CodePtr %d\n", bf.codePtr)
	fmt.Printf("Code %s\n", bf.code)
	//fmt.Printf("Executed %s\n", bf.executed)
	fmt.Printf("MemPtr %d\n", bf.memPtr)
	fmt.Printf("Memory Content % d\n", bf.mem[:200])
}
