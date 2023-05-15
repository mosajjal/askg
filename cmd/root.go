package cmd

import (
	"fmt"
	"os"

	bard-cli "github.com/mosajjal/bard-cli/pkg"
	"github.com/rs/zerolog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "bard-cli",
		Short: "bard-cli is awesome",
		Long:  `bard-cli is the best CLI ever!`,
		Run: func(cmd *cobra.Command, args []string) {
			bard-cli.Run()
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	// define cli arguments
	_ = rootCmd.Flags().IntP("number", "n", 7, "What is the magic number?")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	// make it required
	_ = rootCmd.MarkFlagRequired("number")

	// set up logging
	// set log level
	if l, err := zerolog.ParseLevel("debug"); err == nil {
		zerolog.SetGlobalLevel(l)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

}
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
