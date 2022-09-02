package commands

import (
	"fmt"
	stdlog "log"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool

var rootCmd = &cobra.Command{
	Use:   "battlesnake",
	Short: "Battlesnake Command-Line Interface",
	Long:  "Tools and utilities for Battlesnake games.",
}

func Execute() {
	rootCmd.AddCommand(NewPlayCommand())

	mapCommand := NewMapCommand()
	mapCommand.AddCommand(NewMapListCommand())
	mapCommand.AddCommand(NewMapInfoCommand())

	rootCmd.AddCommand(mapCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.battlesnake.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable debug logging")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".battlesnake" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".battlesnake")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Setup logging
	log.SetStdoutOutput(os.Stderr)
	log.SetFlags(stdlog.Ltime | stdlog.Lmicroseconds)
	if verbose {
		log.SetStdoutThreshold(log.LevelDebug)
	} else {
		log.SetStdoutThreshold(log.LevelInfo)
	}
}
