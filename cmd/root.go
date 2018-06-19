// Copyright © 2018 Peter Alexander <peter.alexander@prodatlab.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/prodatalab/cobra"
	gdax2go "github.com/prodatalab/tools/gdax2go/pkg/gdax2go"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gdax2go <dirpath>",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("gdax2go requires a directory path containing json files as its argument")
		}
		fi, err := os.Stat(args[0])
		if err != nil {
			return errors.Wrap(err, "Stat failed.. the arg must be a valid directory path?")
		}
		if !fi.IsDir() {
			return errors.New("the path provided is not a directory")
		}
		if len(args) > 1 {
			return errors.New("gdax requires a single arg of the directory path to process")
		}
		gdax2go.Val.RootPath = args[0]
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		gdax2go.Run()
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gdax2go.yaml)")
	// rootCmd.PersistentFlags().StringVar(&gdax2go.Val.RootPath, "path", "-p", "the path of the root directory")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".gdax2go" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gdax2go")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
