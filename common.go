package main

import (
	"bufio"
	"os"
)

// Reads input from stdin.
// Returns the input as a string.
func Read_stdin() string {
	stat, _ := os.Stdin.Stat()
	var input_as_stdin string
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input_as_stdin += scanner.Text()
		}
	}
	return input_as_stdin
}
