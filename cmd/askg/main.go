package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/valyala/gorpc"

	_ "embed"

	"github.com/spf13/cobra"
)

var nocolor = strings.ToLower(os.Getenv("NO_COLOR")) == "true"
var logLevel = os.Getenv("ASKG_LOGLEVEL")
var ASKGD_RPC_ADDR = os.Getenv("ASKGD_RPC_ADDR")

var logger = zerolog.New(os.Stderr).
	With().
	Timestamp().
	Logger().
	Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339, NoColor: nocolor})

var (
	version string = "UNKNOWN"
	commit  string = "NOT_PROVIDED"
)

// Execute executes the root command.
func main() {
	cmd := &cobra.Command{
		Use:   "Gemini",
		Short: "Gemini is awesome",
		Long: `Use Google's AI model. This is a reverse engineered API of Gemini Web.
		in order to use this, you first need to start askgd and listen on an RPC port (default 12345),
		then, set up the ASKGD_RPC_ADDR environment variable to the address of the RPC server.
		to change the log level, set the LOG_LEVEL environment variable to one of the following values: debug, info, warn, error
		you can also use the NO_COLOR environment variable to disable log color output.
		`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	flags := cmd.Flags()
	const DEFAULT_FILE_SEPARATOR string = "\n"
	_ = flags.StringSliceP("file", "f", nil, `Path to a file whose contents will be appended to the prompt. If path is -, read standard input. Multiple flag-value pairs are allowed, as well as multiple comma-separated values.

Example of valid syntax:

askg -f "FILE1" -f - -f "FILE2,FILE3" "PROMPT_PART_1" "PROMPT_PART_2"

The above command will result in a prompt being composed in the following order:

<PROMPT_PART_1> <PROMPT_PART_2>`+DEFAULT_FILE_SEPARATOR+`<FILE1 contents>`+DEFAULT_FILE_SEPARATOR+`<stdin contents>`+DEFAULT_FILE_SEPARATOR+`<FILE2 contents>`+DEFAULT_FILE_SEPARATOR+`<FILE3 contents>
`)
	_ = flags.StringP("file-separator", "s", DEFAULT_FILE_SEPARATOR, "Custom string separator to insert before and between files' contents, in case the -f flag is invoked.")
	_ = flags.BoolP("interactive", "i", false, "run in interactive/conversation mode. Google will remember your previous questions and answers")
	_ = flags.BoolP("version", "v", false, "show version info and exit")

	if err := cmd.Execute(); err != nil {
		logger.Error().Msgf("failed to execute command: %s", err)
		return
	}
	if flags.Changed("help") {
		return
	}
	if flags.Changed("version") {
		fmt.Printf("gemini version %s, commit %s\n", version, commit)
		return
	}

	// set up the log level
	switch l := logLevel; l {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger = logger.With().Caller().Logger()
		logger.Println("logging in debug mode")
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger = logger.With().Caller().Logger()
		logger.Println("logging in info mode")
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		logger = logger.With().Caller().Logger()
		logger.Println("logging in warn mode")
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		logger = logger.With().Caller().Logger()
		logger.Println("logging in error mode")
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger = logger.With().Caller().Logger()
		logger.Println("logging in info mode")
	}

	c := &gorpc.Client{
		// TCP address of the server. TODO: make this configurable
		Addr: "127.0.0.1:12345",
	}
	c.Start()

	// run in interactive mode
	if flags.Changed("interactive") {
		RunInteractive(func(prompt string) (string, error) {
			res, err := c.Call(sanitizeQuestion(prompt))
			return res.(string), err
		})
		return
	}

	// run in single question mode
	pQuestionBuffer := BuildPrompt(flags)
	if pQuestionBuffer.Len() == 0 {
		logger.Fatal().Msg("no question provided")
	}

	res, err := c.Call(sanitizeQuestion(pQuestionBuffer.String()))
	if err != nil {
		logger.Fatal().Msgf("failed to ask question: %s", err)
	}

	renderToMD(os.Stdout, res.(string))
}
