package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mosajjal/askg/gemini"
	"github.com/rs/zerolog"
	"github.com/valyala/gorpc"

	_ "embed"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var nocolor = strings.ToLower(os.Getenv("NO_COLOR")) == "true"
var logLevel = strings.ToLower(os.Getenv("ASKG_LOGLEVEL"))
var logger = zerolog.New(os.Stderr).
	With().
	Timestamp().
	Logger().
	Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339, NoColor: nocolor})

var (
	version string = "UNKNOWN"
	commit  string = "NOT_PROVIDED"
	rootCmd        = &cobra.Command{
		Use:   "askgd",
		Short: "Gemini is awesome",
		Long: `Use Google's AI model. This is a reverse engineered API of Gemini Web.
		in order to use this, you first need to run the browser command to get the cookies from the browser.
		note that all cookies will be stored in PLAINTEXT on this machine in the ~/.askg.yaml file.
		make sure that file is only readable by the user running this daemon.
		`,
	}
)

func main() {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the askgd daemon",
		Long:  `Run the askgd daemon`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	rootCmd.AddCommand(cmd)

	browserCmd := &cobra.Command{
		Use:   "browser",
		Short: "Get cookies from browser",
		Long:  `Get cookies from browser`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: make this into a separate function
			NewCookies := getCookiesFromBrowser("")
			j, _ := yaml.Marshal(NewCookies)

			// write the cookies to the config file.
			//TODO: fix error handling
			f, err := os.Create(rootCmd.Flags().Lookup("config").Value.String())
			if err != nil {
				logger.Error().Msgf("failed to open file: %s", err)
			}
			_, err = f.Write(j)
			if err != nil {
				logger.Error().Msgf("failed to write to file: %s", err)
			}
			err = f.Close()
			if err != nil {
				logger.Error().Msgf("failed to close file: %s", err)
			}
			os.Exit(0)
		},
	}
	rootCmd.AddCommand(browserCmd)

	flags := rootCmd.Flags()
	config := flags.StringP("config", "c", "$HOME/.askg.yaml", "path to YAML configuration file")

	_ = flags.BoolP("version", "v", false, "show version info and exit")

	rootCmd.ParseFlags(os.Args)

	// construct the ~/.askg.yaml in a cross-platform way
	if !flags.Changed("config") {
		home, err := os.UserHomeDir()
		if err != nil {
			logger.Fatal().Msgf("failed to get user home directory: %s", err)
		}
		*config = filepath.Join(home, ".askg.yaml")
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

	if err := rootCmd.Execute(); err != nil {
		logger.Error().Msgf("failed to execute command: %s", err)
		return
	}

	// below is only run if the root command is run
	if cmd.CalledAs() != "" {

		// read the cookies from the config file
		cookies := &gemini.Cookies{}
		f, err := os.Open(*config)
		if err != nil {
			logger.Error().Msgf("failed to open file: %s", err)
		}
		defer f.Close()
		err = yaml.NewDecoder(f).Decode(cookies)
		if err != nil {
			logger.Fatal().Msgf("failed to unmarshal default config: %s", err)
		}

		Gemini := gemini.New(&logger, *cookies)

		// run a daemon to rotate the cookies every 30 seconds
		go func() {
			for {
				logger.Info().Msg("Rotating cookies")
				Gemini.RotateCookies()
				time.Sleep(30 * time.Second)
			}
		}()
		// run a daemon to commit the cookies to the config file every 5 minutes
		go func() {
			for {
				logger.Info().Msg("Committing cookies")
				Gemini.CommitCookies(*config)
				time.Sleep(5 * time.Minute)
			}
		}()

		// run an rpc to expose gemini ask function
		s := &gorpc.Server{
			// Accept clients on this TCP address. TODO: make it configurable
			Addr: ":12345",

			// Echo handler - just return back the message we received from the client
			Handler: func(clientAddr string, request interface{}) interface{} {
				res, err := Gemini.Ask(request.(string))
				if err != nil {
					return err.Error()
				}
				return res
			},
		}
		if err := s.Serve(); err != nil {
			log.Fatalf("Cannot start rpc server: %s", err)
		}
	}
}
