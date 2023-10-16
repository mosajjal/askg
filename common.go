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
func ReadFromFileOrStdin(readFrom string, logger *zerolog.Logger) string {
	var fileContents string
	const READ_FROM_STDIN string = "-"
	if readFrom != READ_FROM_STDIN {
		// Read input from a side file
		fileContentsAsBytes, err := os.ReadFile(readFrom)
		if err != nil {
			logger.Fatal().Msgf("failed to load file %s: %s", readFrom, err)
		}
		fileContents = string(fileContentsAsBytes)
	} else {
		// Read input from stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				fileContents += scanner.Text()
			}
		}
	}
	return fileContents
}
