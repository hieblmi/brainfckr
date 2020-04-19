# Brainfckr
A streaming interpreter for [Brainfuck](https://en.wikipedia.org/wiki/Brainfuck) in [Golang](https://github.com/golang/go/wiki/WhyGo)

## How it works
Brainfckr performs input and output operations on Go's simplest streaming IO primitives, io.Reader and io.Writer. 
```
func NewBrainfckr(code io.Reader, input io.Reader, writer io.Writer) *Brainfckr
```
Rather than loading the complete code into memory Brainfckr reads brainfuck code, executes it on program input and streams output on-the-fly.
While this seems more elegant it requires more IO operations for large code inputs compared to one big ```ioutil.ReadAll``` at the very start of the execution.
Another point to note is how loops are executed when code is read from a stream. 
| loop instruction  | meaning   |
|---|---|
|		[	   | if the byte at the data pointer is zero, then instead of moving the instruction pointer forward to the next command, jump it forward to the command after the matching ']' command.		   |
|		]	   | 	if the byte at the data pointer is nonzero, then instead of moving the instruction pointer forward to the next command, jump it back to the command after the matching '[' command.		   |

Since we are not looking ahead for a matching ']' at the beginning of loop Brainfckr holds looping code in memory for execution until the most outer loop sees it's loop counter at 0. To determine whether to read the next code byte from the IO stream or from the code segment in memory ```func (bf *Brainfckr) nextLoopOp() byte``` checks if the code pointer sits at the end of the code segment. If so we need to get our next instruction from the stream. Otherwise we can just read a previously stored instruction from the code segment. To furhter clarify this consider the debug output ```bf.debugOutput(string)``` for the following interleaving loops. ```++[>++[.-]<-]```

<pre>
CodePtr <b>6</b>
Code [>++[.<b>-</b>
Executed ++<++.
Memory Content [ 2  2  0  0  0  0  0  0  0  0]
CodePtr <b>7</b>
Code [>++[.-<b>]</b>
Executed ++<++.-
Memory Content [ 2  1  0  0  0  0  0  0  0  0]
CodePtr <b>4</b>
Code [>++<b>[</b>.-]
Executed ++<++.-
Memory Content [ 2  1  0  0  0  0  0  0  0  0]
CodePtr 5
Code [>++[<b>.</b>-]
Executed ++<++.-
Memory Content [ 2  1  0  0  0  0  0  0  0  0]
CodePtr 6
Code [>++[.<b>-</b>]
Executed ++<++.-.
Memory Content [ 2  1  0  0  0  0  0  0  0  0]
CodePtr 7
Code [>++[.-<b>]</b>
Executed ++<++.-.-
Memory Content [ 2  0  0  0  0  0  0  0  0  0]
.
.
.
CodePtr 8
Code [>++[.-]<
Executed ++<++.-.-
Memory Content [ 2  0  0  0  0  0  0  0  0  0]
</pre>

The codePtr jump from instruction at byte 7 to byte 4 indicates that we are reading from the code segement. Once the inner loop ran we are expanding the code segment with the next instruction from the stream ```<``` and increment the codePtr to 8.

## Examples

### Streaming web server requests
```
func TestWebRequestInput(t *testing.T) {
	// Setup web server
	t.Run("ServerSetup", func(t *testing.T) {
		http.HandleFunc("/bf", brainfckHandler)
		go http.ListenAndServe(":8080", nil)
	})

	input := []byte(">++++[>++++++<-]>-[[<+++++>>+<-]>-]<<[<]>>>>--.<<<-.>>>-.<.<.>---.<<+++.>>>++.<<---.[>]<<.")
	expectedOutput := "brainfuck"
	resp, err := http.Post("http://localhost:8080/bf", "text/plain", bytes.NewReader(input))
	if err != nil {
		// handle error
	}
	body, err := ioutil.ReadAll(resp.Body)
	if string(body) != expectedOutput {
		t.Errorf("Response %s doesn't match expectedOutput %s\n", string(body), expectedOutput)
	}
	defer resp.Body.Close()
}

func brainfckHandler(w http.ResponseWriter, r *http.Request) {
	NewBrainfckr(r.Body, os.Stdin, w).Interpret()
}
```

### 99 Bottles Of Beer
Check out testcase ```func Test99BottlesOfBeer(t *testing.T)```
Result:
```
99 bottles of beer on the wall, 99 bottles of beer.
Take one down and pass it around, 98 bottles of beer on the wall.

98 bottles of beer on the wall, 98 bottles of beer.
Take one down and pass it around, 97 bottles of beer on the wall.
.
.
.
No more bottles of beer on the wall, no more bottles of beer.
Go to the store and buy some more, 99 bottles of beer on the wall.
```

### Reversing input from a []byte 
```
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
```
