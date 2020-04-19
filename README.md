# Brainfckr
A streaming interpreter for [Brainfuck](https://en.wikipedia.org/wiki/Brainfuck) in [Golang](https://github.com/golang/go/wiki/WhyGo)
| instruction  | meaning   |
|---|---|
|		>	   | increment the data pointer (to point to the next cell to the right).		   |
|		<	   | decrement the data pointer (to point to the next cell to the left).		   |
|		+	   | increment (increase by one) the byte at the data pointer.		   |
|		-	   | decrement (decrease by one) the byte at the data pointer.		   |
|		.	   | output the byte at the data pointer.		   |
|		,	   | accept one byte of input, storing its value in the byte at the data pointer.		   |
|		[	   | if the byte at the data pointer is zero, then instead of moving the instruction pointer forward to the next command, jump it forward to the command after the matching ] command.		   |
|		]	   | 	if the byte at the data pointer is nonzero, then instead of moving the instruction pointer forward to the next command, jump it back to the command after the matching [ command.		   |
## How it works
Brainfckr performs input and output operations on Go's simplest streaming IO primitives, io.Reader and io.Writer. 
```
func NewBrainfckr(code io.Reader, input io.Reader, writer io.Writer) *Brainfckr
```
Rather than loading the complete code into memory executing it at once Brainfckr reads brainfuck code, executes it on program input and streams output on-the-fly.
While this seems more elegant it requires more IO operations for large code inputs compared to one big ```ioutil.ReadAll``` at the very start of the execution.
Another point to note is how loops are executed when code is read from a stream. Since we can't look ahead 

different streaming examples

How loops are handled in streaming
