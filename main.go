package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
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

var nocolorLog = strings.ToLower(os.Getenv("NO_COLOR")) == "true"
var logger = zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339, NoColor: nocolorLog})

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
		in order to use this, you first need to gain access to Bard in your browser, and then copy the cookie "__Secure-1PSID"
		using developer tools. If you don't know how, follow this guide:
		https://developer.chrome.com/docs/devtools/application/cookies/
		IMPORTANT NOTE: never share your cookies with anyone, as they can be used to impersonate you and steal your data.
		`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	flags := cmd.Flags()
	config := flags.StringP("config", "c", "$HOME/.bardcli.yaml", "path to YAML configuration file")
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
		err := ioutil.WriteFile(*config, defaultConfig, 0644)
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
	bard := bard.New(k.String("cookie"), &logger)

	// run in interactive mode
	if flags.Changed("interactive") {
		RunInteractive(bard)
		return
	}

	// run in single question mode
	answer, err := bard.Ask(strings.Join(flags.Args(), " "))
	if err != nil {
		logger.Fatal().Msgf("failed to ask question: %s", err)
	}

	out, err := glamour.RenderWithEnvironmentConfig(answer)
	if err != nil {
		logger.Fatal().Msgf("failed to render answer: %s", err)
	}
	fmt.Println(out)

}
