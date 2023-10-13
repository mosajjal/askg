package main

import (
	"bufio"
	"os"

	"github.com/rs/zerolog"
)

// Reads input from a side file or from standard input, but not both.
//
// Parameters:
//
// * from: Path to a file whose contents will be appended to the prompt. If path is -, read standard input.
// * logger: A zerolog logger to use for logging errors.
//
// Returns:
//
// Returns the input as a string.
func Read_from_file_or_stdin(read_from string, logger *zerolog.Logger) string {
	var file_contents string
	var from_stdin string = "-"
	if read_from != from_stdin {
		// Read input from a side file
		file_contents_as_bytes, err := os.ReadFile(read_from)
		if err != nil {
			logger.Fatal().Msgf("failed to load file %s: %s", read_from, err)
		}
		file_contents = string(file_contents_as_bytes)
	} else {
		// Read input from stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				file_contents += scanner.Text()
			}
		}
	}
	return file_contents
}
