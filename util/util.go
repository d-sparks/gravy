package gravyutil

import (
	"bufio"
	"log"
	"os"
)

// If error, log error and die.
func FatalIfErr(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

// Create a file and wrap it in a bufio.Writer, or die.
func FileWriterOrDie(filename string) bufio.Writer {
	file, err := os.Create(filename)
	FatalIfErr(err)
	return bufio.NewWriter(file)
}

// Open a file and wrap it in a bufio.Scanner with large buffer, or die.
func FileScannerOrDie(filename string) bufio.Scanner {
	file, err := os.Open(filename)
	FatalIfErr(err)
	scanner := bufio.NewScanner(file)
	scanner.Buffer([]byte{}, 1024*1024)
	return scanner
}
