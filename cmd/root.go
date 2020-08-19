/*
Copyright © 2020 Conner Peirce <connerpeirce@gmail.com>

*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/conner/elktail/pkg/elktail"
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
			// fmt.Println("index was set")
			// viper.WriteConfig()
			// @todo write serialized flags if they were explicitly set
		}

		return nil

	},
	RunE: func(cmd *cobra.Command, args []string) error {
		url, err := url.Parse(viper.GetString(flagURL))
		if err != nil {
			return err
		}

		var query string
		if len(args) > 0 {
			query = args[0]
		} else {
			query = ""
		}

		t := elktail.NewTail(elktail.Config{
			IndexPattern: viper.GetString(flagIndex),
			PageSize:     pageSize,
			Query:        query,
			URL:          *url,
		})

		return t.Run(cmd.Context(), os.Stdout)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel

	// cancel the context if a signal is received
	go func() {
		s := make(chan os.Signal)
		signal.Notify(
			s,
			syscall.SIGINT,
			syscall.SIGTERM,
		)
		<-s
		cancel()
	}()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Fatalln(err)
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
			log.Fatalln(err)
		}

		// Search config in home directory with name ".elktail" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".elktail")
		// fmt.Println(home)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
