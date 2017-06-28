// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"github.com/tech-nico/anypoint-cli/utils"
)

var cfgFile string
var outputFormat string
var debug bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "anypoint-cli",
	Short: "Anypoint Platform command line client",
	Long: `Manage the Anypoint Platorm through the command line.

Anypoint-cli is made of several sub-commands, each allowing you to manage
a different entity in the Anypoint Platform. Through anypoint-cli you will
be able to manage:
- APIs
- Applications
- Users
- Access Management details`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.anypoint-cli.yaml)")
	RootCmd.PersistentFlags().StringVar(&outputFormat, "o", "json", "determines output format (json/yaml/csv)")
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Display debug information")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	setupViper()
}
func setupViper() {
	viper.SetConfigType("json")

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
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

		// Search config in home directory with name ".anypoint-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(utils.CONFIG_FILE_NAME)
	}

	viper.AutomaticEnv()
	// read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	//else {
	//		panic(err)
	//	}

	log.Printf("File used : %s", viper.ConfigFileUsed())
}
