package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
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
// Returns the pointer to the input string.
func ReadFromFileOrStdin(readFrom string, logger *zerolog.Logger) *string {
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
	return &fileContents
}

// Builds the full prompt from the provided arguments, with the positional arguments being the first part of the prompt, followed by the named arguments -f. For example, the command
//   bard-cli -f "FILE1" -f - -f "FILE2,FILE3" "PROMPT_PART_1" "PROMPT_PART_2"
// will result in the following prompt:
//   <PROMPT_PART_1> <PROMPT_PART_2>
//   <FILE1 contents>
//   <stdin contents>
//   <FILE2 contents>
//   <FILE3 contents>
//
// Parameters:
//
// * flags: a *pflag.FlagSet pointer
//
// Returns:
//
// Returns the pointer to the buffer containing the full prompt.
func BuildPrompt(flags *pflag.FlagSet) *bytes.Buffer {
	fullPromptAsBuffer := bytes.Buffer{}
	firstPartOfPrompt := strings.Join(flags.Args(), " ")
	fullPromptAsBuffer.WriteString(firstPartOfPrompt)
	alsoReadFrom, _ := flags.GetStringSlice("file")
	if len(alsoReadFrom) == 0 {
		return &fullPromptAsBuffer
	}
	fileSeparator, _ := flags.GetString("file-separator")
	for _, pathOrStdin := range alsoReadFrom {
		pNextPartOfPrompt := ReadFromFileOrStdin(pathOrStdin, &logger)
		if *pNextPartOfPrompt == "" {
			continue
		}
		if fullPromptAsBuffer.Len() > 0 {
			fullPromptAsBuffer.WriteString(fileSeparator)
		}
		fullPromptAsBuffer.WriteString(*pNextPartOfPrompt)
	}
	return &fullPromptAsBuffer
}
