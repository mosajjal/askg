package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/mosajjal/bard-cli/bard"
	"github.com/rs/zerolog"

	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed config.defaults.yaml
var defaultConfig []byte

var nocolor = strings.ToLower(os.Getenv("NO_COLOR")) == "true"
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
		Use:   "bard",
		Short: "bard is awesome",
		Long: `bard is Google's AI model. This is a reverse engineered version of Bard on the web.
		in order to use this, you first need to gain access to Bard in your browser,
		and then copy the cookie "__Secure-1PSID" using developer tools. If you don't know how, follow this guide:
		https://developer.chrome.com/docs/devtools/application/cookies/
		SECURITY NOTE: NEVER share your cookies with anyone. they can be used to impersonate you and steal your data.
		`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	flags := cmd.Flags()
	config := flags.StringP("config", "c", "$HOME/.bardcli.yaml", "path to YAML configuration file")
	read_from := flags.StringP("file", "f", "", "Path to a file whose contents will be appended to the prompt. If path is -, read standard input.")
	_ = flags.Bool("defaultconfig", false, "write the default config yaml file to stdout")
	_ = flags.BoolP("interactive", "i", false, "run in interactive/conversation mode. Bard will remember your previous questions and answers")
	_ = flags.BoolP("version", "v", false, "show version info and exit")

	if err := cmd.Execute(); err != nil {
		logger.Error().Msgf("failed to execute command: %s", err)
		return
	}
	// construct the ~/.bardcli.yaml in a cross-platform way
	if !flags.Changed("config") {
		home, err := os.UserHomeDir()
		if err != nil {
			logger.Fatal().Msgf("failed to get user home directory: %s", err)
		}
		*config = filepath.Join(home, ".bardcli.yaml")
	}
	if flags.Changed("help") {
		return
	}
	if flags.Changed("version") {
		fmt.Printf("bard version %s, commit %s\n", version, commit)
		return
	}
	if flags.Changed("defaultconfig") {
		err := os.WriteFile(*config, defaultConfig, 0644)
		if err != nil {
			logger.Fatal().Msgf("failed to write default config: %s", err)
		}
		logger.Info().Msgf("wrote default config to %s", *config)
		return
	}

	k := koanf.New(".")
	// load the defaults
	if err := k.Load(rawbytes.Provider(defaultConfig), yaml.Parser()); err != nil {
		logger.Fatal().Msgf("failed to load default config: %s", err)
	}

	if err := k.Load(file.Provider(*config), yaml.Parser()); err != nil {
		logger.Fatal().Msgf("failed to load config file: %s", err)
	}

	// set up the log level
	switch l := k.String("log_level"); l {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger = logger.With().Caller().Logger()
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger = logger.With().Caller().Logger()
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// set up the bard client
	cookie1psid := k.String("cookie")
	// cookie1psid is an alias for cookie
	if cookie1psid == "" {
		cookie1psid = k.String("cookie1psid")
	}
	cookie1psidts := k.String("cookie1psidts")
	cookie1psidcc := k.String("cookie1psidcc")

	bard := bard.New(cookie1psid, cookie1psidts, cookie1psidcc, &logger)

	// run in interactive mode
	if flags.Changed("interactive") {
		RunInteractive(bard)
		return
	}

	// run in single question mode
	var input_as_args string = strings.Join(flags.Args(), " ")
	var appended_input string
	if *read_from != "" {
		appended_input = Read_from_file_or_stdin(*read_from, &logger)
	}
	var question string = strings.Join([]string{input_as_args, appended_input}, " ")
	question = strings.TrimSpace(question)
	if question == "" {
		logger.Fatal().Msg("no question provided")
	}
	answer, err := bard.Ask(sanitizeQuestion(question))
	if err != nil {
		logger.Fatal().Msgf("failed to ask question: %s", err)
	}

	renderToMD(os.Stdout, answer)

}
