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
</pre>

In line 23 the CodePtr jumps from byte 7 to byte 4 in the next operation indicating that we are reading from the code segement.

