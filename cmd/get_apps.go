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
	"github.com/spf13/cobra"
	"log"
	"github.com/tech-nico/anypoint-cli/rest"
	"github.com/spf13/viper"
	"github.com/tech-nico/anypoint-cli/utils"
	"fmt"
)

var appName, envName string


// appCmd represents the app command
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if envName == "" {
			log.Fatalf("Please specify --environment parameter")
		}

		apiMgr := rest.NewAPI(viper.GetString(utils.KEY_URI), viper.GetString(utils.KEY_TOKEN))

		if appName != "" {
			app, err := apiMgr.ApplicationsByName(viper.GetString(utils.KEY_ORG_ID), envName, appName)
			if err != nil {
				log.Fatalf("Error when sarching for app %q : %s", appName, err)
			}

			fmt.Println("App found: %s", app)
		} else {
			apps, err := apiMgr.GetApplications(viper.GetString(utils.KEY_ORG_ID), envName)
			if err != nil {
				log.Fatalf("Error retrieving all applications: %s", err)
			}

			printApps(apps)
		}
	},
}

func init() {
	getCmd.AddCommand(appCmd)

	appCmd.Flags().StringVarP(&appName, "app-name", "a", "", "Name of the app to search for")
	appCmd.Flags().StringVarP(&envName, "environment", "e", "", "Environment name")
}

func printApps(apps []interface{}) {
	headers := []string{"NAME", "SERVER TYPE", "SERVER NAME", "STATUS", "FILE"}

	data := make([][]string, 0)
	for _, val := range apps {
		app := val.(map[string]interface{})
		artifact := app["artifact"].(map[string]interface{})
		name := fmt.Sprint(artifact["name"])
		target := app["target"].(map[string]interface{})
		serverType := fmt.Sprint(target["type"])
		if serverType == "<nil>" {
			serverType = "<Unknown>"
		}
		serverName := fmt.Sprint(target["name"])
		if serverName == "<nil>" {
			serverName = "<Unknown>"
		}
		status := fmt.Sprint(app["lastReportedStatus"])
		file := fmt.Sprint(artifact["fileName"])
		row := []string{name, serverType, serverName, status, file}
		data = append(data, row)
	}

	utils.PrintTabular(headers, data)
}

