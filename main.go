package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {

	fileName := flag.String("f", "brainfck_code.bf", "File path containing brainfck code.")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Could not open file: %s, %v\n", *fileName, err))
	}
	defer file.Close()

	fmt.Println()
	brainfckr := NewBrainfckr(file, os.Stdout)
	err = brainfckr.Interpret()
	if err != io.EOF {
		return fmt.Errorf(fmt.Sprintf("Failed to execute the brainfck file: %s, %v\n", *fileName, err))
	}

	return nil
}
