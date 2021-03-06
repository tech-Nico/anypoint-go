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
	"github.com/spf13/viper"
	"github.com/tech-nico/anypoint-cli/rest"
	"github.com/tech-nico/anypoint-cli/utils"

	"fmt"
)

var serverSearchString string

// appCmd represents the app command
var getServersCmd = &cobra.Command{
	Use:     "servers",
	Aliases: []string{"server"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var servers []interface{}
		var err error
		if envName == "" {
			return fmt.Errorf("Please specify --environment parameter")
		}

		apiMgr := rest.NewAPI(viper.GetString(utils.KEY_URI), viper.GetString(utils.KEY_TOKEN))

		if serverSearchString != "" {
			servers, err = apiMgr.SearchServers(viper.GetString(utils.KEY_ORG_ID), envName, serverSearchString)
		} else {
			servers, err = apiMgr.GetAllServers(viper.GetString(utils.KEY_ORG_ID), envName)
		}

		if err != nil {
			return fmt.Errorf("Error in retrieving servers: %s", err)
		}
		headers := []string{"NAME", "SERVER TYPE", "RUNTIME TYPE", "STATUS", "VERSION", "AGENT", "N. APPS"}
		utils.PrintObject(servers, headers, extractList)

		return nil
	},
}

func extractList(servers []interface{}) [][]string {
	data := make([][]string, 0)
	for _, val := range servers {

		server := val.(map[string]interface{})
		details := server["details"].(map[string]interface{})
		deployments := server["deployments"].([]interface{})
		name := fmt.Sprint(server["name"])
		serverType := fmt.Sprint(server["type"])
		var runtimeType, version string

		if serverType == "CLUSTER" {
			servers := details["servers"].([]interface{})
			firstServer := servers[0].(map[string]interface{})
			firstServerDetail := firstServer["details"].(map[string]interface{})
			runtimeType = fmt.Sprint(firstServerDetail["type"])
			version = fmt.Sprint(firstServerDetail["runtimeVersion"])
		} else {
			runtimeType = fmt.Sprint(details["type"])
			version = fmt.Sprint(details["runtimeVersion"])
		}
		status := fmt.Sprint(server["status"])
		agent := fmt.Sprint(details["agentVersion"])

		numApps := fmt.Sprint(len(deployments))

		row := []string{name, serverType, runtimeType, status, version, agent, numApps}
		data = append(data, row)
	}

	return data
}

func init() {
	getCmd.AddCommand(getServersCmd)
}
