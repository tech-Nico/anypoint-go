// Copyright © 2017 Nico Balestra <functions@protonmail.com>
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

	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tech-nico/anypoint-cli/rest"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"syscall"
	"github.com/tech-nico/anypoint-cli/utils"
)

var username string
var password string
var uri string
var orgPath string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login user into the Anypoint Platform",
	Long: `Login into the Anypoint Platform providing your username and password.
	Bear in mind that if the Anypoint Platform you are trying to login onto is configured
	with an External Identity Provider, you will need to provide credentials for such IDP.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			fmt.Println("Debug mode enabled!")
			viper.Set(utils.KEY_DEBUG, true)
		}

		if uri == "" {
			uri = "https://anypoint.mulesoft.com"
		}

		if username == "" {
			errors.New("Please specify --username")
		}

		if password == "" {
			password = promptForPassword()
		}

		auth := rest.NewAuthWithCredentials(uri, username, password)
		fmt.Printf("Login successful. Got token '%s'", auth.Token)
		viper.Set(utils.KEY_TOKEN, auth.Token)
		viper.Set(utils.KEY_URI, uri)
		if orgPath != "" {
			orgId := auth.FindBusinessGroup(orgPath)
			fmt.Printf("The orgId is '%s'\n", orgId)
			viper.Set(utils.KEY_ORG_ID, orgId)
		}
		utils.WriteConfig()
		//All good
		return nil
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//loginCmd.PersistentFlags().StringVar(&username, "username", "", "Username to login to Anypoint Platform")
	//loginCmd.PersistentFlags().StringVar(&password, "password", "", "Password to login to Anypoint Platform. If not specified it will prompt for a password")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	loginCmd.Flags().StringVar(&username, "username", "", "Specify the username to login into Anypoint Platform")
	loginCmd.Flags().StringVar(&password, "password", "", "Specify the password to login into Anypoint Platform")
	loginCmd.Flags().StringVar(&uri, "uri", "", "Specify the url of the Anypoint Platform instance where you would like to login to")
	loginCmd.Flags().StringVar(&orgPath, "org", "", "Specify the path to the org you want to manage resources for. Ex: Root\\Sub Org")

}

func promptForPassword() string {

	fmt.Print("Enter password:")
	bytepassword, _ := terminal.ReadPassword(int(syscall.Stdin))

	password := string(bytepassword)

	return strings.TrimSpace(password)
}
