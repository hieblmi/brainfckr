package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestSimpelOutput(t *testing.T) {

	fileName := "./test/simpleOutputTest.bf"
	expectedOutput := '7'
	var writer bytes.Buffer
	file, err := os.Open(fileName)
	defer file.Close()

	if err != nil {
		t.Errorf("Could not open test file %s", fileName)
	} else {
		brainfckr := NewBrainfckr(file, &writer)
		err = brainfckr.Interpret()
	}

	if err != io.EOF {
		t.Errorf(fmt.Sprintf("Failed to execute Test TestOutput %s", fileName))
	} else {
		result, _, _ := writer.ReadRune()
		if result != expectedOutput {
			t.Errorf("result = %c; want 7", result)
		}
	}
}

func Test99BottlesOfBeer(t *testing.T) {

	fileName := "./test/99BottlesOfBeer.bf"
	resultFile := "./test/99BottlesOfBeer_RESULT"

	var resultBuffer *bytes.Buffer = bytes.NewBuffer([]byte{})
	//writer := bufio.NewWriter(buffer)
	file, err := os.Open(fileName)
	defer file.Close()

	if err != nil {
		t.Errorf("Could not open test file %s", fileName)
	} else {
		brainfckr := NewBrainfckr(file, resultBuffer)
		err = brainfckr.Interpret()
	}

	file, err = os.Open(resultFile)
	defer file.Close()

	if err != nil {
		t.Errorf("Could not open result file %s", resultFile)
	} else {
		var expected *bufio.Reader = bufio.NewReader(file)

		a := make([]byte, 1)
		b := make([]byte, 1)
		for err != io.EOF {
			if a[0] != b[0] {
				t.Errorf("Difference in bottles of beer detected. Calculated %x, Got %x\n", a[0], b[0])
			}
			_, _ = resultBuffer.Read(a)
			_, err = expected.Read(b)
			if a[0] == 10 && b[0] == 13 { // skip CRLF from the result file
				_, err = expected.Read(b)
			}
		}

		var c, d int
		c, _ = resultBuffer.Read(a)
		d, err = expected.Read(b)
		if c != 0 && d != 0 {
			t.Errorf("Not all characters from either input or result have been processed.\n")
		}
	}

}
