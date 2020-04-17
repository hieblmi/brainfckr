package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
		brainfckr := NewBrainfckr(file, os.Stdin, &writer)
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
	file, err := os.Open(fileName)

	if err != nil {
		t.Errorf("Could not open test file %s", fileName)
	} else {
		defer file.Close()
		brainfckr := NewBrainfckr(file, os.Stdin, resultBuffer)
		err = brainfckr.Interpret()
	}

	file, err = os.Open(resultFile)

	if err != nil {
		t.Errorf("Could not open result file %s", resultFile)
	} else {
		defer file.Close()
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

func TestPrintBrainfuck(t *testing.T) {
	fileName := "./test/printBrainfuck.bf"
	resultFile := "./test/printBrainfuck_RESULT"

	var resultBuffer *bytes.Buffer = bytes.NewBuffer([]byte{})
	file, err := os.Open(fileName)

	if err != nil {
		t.Errorf("Could not open test file %s", fileName)
	} else {
		defer file.Close()
		brainfckr := NewBrainfckr(file, os.Stdin, resultBuffer)
		err = brainfckr.Interpret()
	}

	file, err = os.Open(resultFile)
	if err != nil {
		t.Errorf("Could not open result file %s", resultFile)
	} else {
		defer file.Close()
		var expected *bufio.Reader = bufio.NewReader(file)

		a := make([]byte, 1)
		b := make([]byte, 1)
		for err != io.EOF {
			if a[0] != b[0] {
				t.Errorf("Difference detected. Calculated %x, Got %x\n", a[0], b[0])
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

func TestReverseInput(t *testing.T) {

	fileName := "./test/reverseInput.bf"
	input := []byte{1, 2, 3, 4, 5, 0}
	expectedOutput := []byte{5, 4, 3, 2, 1}
	var writer bytes.Buffer

	file, err := os.Open(fileName)
	if err != nil {
		t.Errorf("Could not open test file %s", fileName)
	} else {
		defer file.Close()
		brainfckr := NewBrainfckr(file, bytes.NewReader(input), &writer)
		err = brainfckr.Interpret()
	}

	if err != io.EOF {
		t.Errorf(fmt.Sprintf("Failed to execute Test TestOutput %s", fileName))
	} else {
		if !bytes.Equal(expectedOutput, writer.Bytes()) {
			t.Errorf("Reversed input doesn't look like expected.")
		}
	}
}

func TestWebRequestInput(t *testing.T) {
	// Setup web server
	t.Run("ServerSetup", func(t *testing.T) {
		http.HandleFunc("/bf", brainfckHandler)
		go http.ListenAndServe(":8080", nil)
	})

	input := []byte(">++++[>++++++<-]>-[[<+++++>>+<-]>-]<<[<]>>>>-\n-.<<<-.>>>-.<.<.>---.<<+++.>>>++.<<---.[>]<<.")
	expectedOutput := "brainfuck"
	resp, err := http.Post("http://localhost:8080/bf", "text/plain", bytes.NewReader(input))
	if err != nil {
		// handle error
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || string(body) != expectedOutput {
		t.Errorf("Response %s doesn't match expectedOutput %s\n", string(body), expectedOutput)
	}
	defer resp.Body.Close()
}

func brainfckHandler(w http.ResponseWriter, r *http.Request) {
	NewBrainfckr(r.Body, os.Stdin, w).Interpret()
}
