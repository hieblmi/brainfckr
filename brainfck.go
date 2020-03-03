package main

import (
	"fmt"
	"flag"
)

func main() {
	bfFile := flag.String("f", "brainfck_code.bf", "File path containing the brainfck code.")
	flag.Parse();

	fmt.Println(*bfFile)
	fmt.Println(flag.Args())

}