package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	brainfckFile := flag.String("f", "brainfck_code.bf", "File path containing brainfck code.")
	flag.Parse()

	file, err := os.Open(*brainfckFile)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Could not open file: %s, %v\n", *brainfckFile, err))
	}
	defer file.Close()

	brainfck := NewBrainfckr(file, os.Stdout)
	fmt.Println()
	err = brainfck.Interpret()
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to execute the brainfck file: %s, %v\n", *brainfckFile, err))
	}

	return nil
}
