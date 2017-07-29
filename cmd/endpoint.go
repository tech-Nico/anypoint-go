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

	"github.com/spf13/cobra"
	"os"
	"github.com/tech-nico/anypoint-cli/rest"
	"github.com/spf13/viper"
	"github.com/tech-nico/anypoint-cli/utils"
	"encoding/json"
	"log"
)

// endpointCmd represents the endpoint command
var endpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Manages an API version's endpoint",
	Long: `With this command you will be able to visualize the details related to an API version's endpoint and also
	to make changes to it. An API version's endpoint describes your API's connectivity.`,
	Run: func(cmd *cobra.Command, args []string) {
		if apiId == 0 {
			cmd.Usage()
			fmt.Println("Error: please specify --api-id")
			os.Exit(1)
		}

		if versionId == 0 {
			cmd.Usage()
			fmt.Println("Error: please specify --version-id")
			os.Exit(1)
		}

		apiClient := rest.NewAPI(viper.GetString(utils.KEY_URI), viper.GetString(utils.KEY_TOKEN))

		format := viper.GetString(utils.KEY_FORMAT)
		switch format {
		case "list":
			res, err := apiClient.GetEndpointAsMap(viper.GetString(utils.KEY_ORG_ID), apiId, versionId)
			if err != nil {
				log.Fatalf("Error while retrieving endpoint for api %d (version %d): %s", apiId, versionId, err)
			}
			printEndpoint(res)
			break
		case "json":
			res, err := apiClient.GetEndpointAsJSONString(viper.GetString(utils.KEY_ORG_ID), apiId, versionId)
			if err != nil {
				log.Fatalf("Error while retrieving endpoint for api %d (version %d): %s", apiId, versionId, err)
			}

			b, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				fmt.Println("error:", err)
			}
			os.Stdout.Write(b)
			break
		}

	},
}

func printEndpoint(endpoint map[string]interface{}) {
	headers := []string{"ID", "TYPE", "URI", "PROXY URI", "PROXY REGISTRATION URI", "USE DOMAIN", "RESPONSE TIMEOUT"}

	data := make([][]string, 1)
	data[0] = []string{
		fmt.Sprint(endpoint["id"]),
		fmt.Sprint(endpoint["type"]),
		fmt.Sprint(endpoint["uri"]),
		fmt.Sprint(endpoint["proxyUri"]),
		fmt.Sprint(endpoint["proxyRegistrationUri"]),
		fmt.Sprint(endpoint["referencesUserDomain"]),
		fmt.Sprint(endpoint["resonseTimeout"]),

	}

	utils.PrintTabular(headers, data)

}

func init() {
	apiCmd.AddCommand(endpointCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// endpointCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
