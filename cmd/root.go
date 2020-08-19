/*
Copyright © 2020 Conner Peirce <connerpeirce@gmail.com>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	defaultIndexPattern = "logstash-[0-9].*"
	defaultURL          = "http://127.0.0.1:9200"

	flagIndex = "index"
	flagURL   = "url"

	pageSize = 50
)

var (
	cfgFile string

	serializedFlags = [...]string{flagIndex, flagURL}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "elktail [query]",
	Short: "utility to tail Logstash logs from an ELK stack",
	Long:  "flags marked wite (*) are persisted between runs in the config file",
	Args:  cobra.MaximumNArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flag("index").Changed {
			fmt.Println("index was set")
			// viper.WriteConfig()
			// @todo write serialized flags if they were explicitly set
		}

		return nil

	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default: \"$HOME/.elktail.yaml\")")
	rootCmd.Flags().StringP(flagIndex, "i", defaultIndexPattern, fmt.Sprintf("(*) index pattern (default: \"%s\")", defaultIndexPattern))
	rootCmd.Flags().StringP(flagURL, "u", defaultURL, fmt.Sprintf("(*) ElasticSearch URL (default: \"%s\")", defaultURL))

	for _, f := range serializedFlags {
		viper.BindPFlag(f, rootCmd.Flag(f))
	}
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

		// Search config in home directory with name ".elktail" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".elktail")
		fmt.Println(home)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
