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
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tech-nico/anypoint-cli/rest"
	"github.com/tech-nico/anypoint-cli/utils"
	"log"
)

var apiId, versionId int
var apiName string
var offset, limit int

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Work with APIs managed by MuleSoft Anypoint Platform",
	Long: `Manage APIs created and governed onto the MuleSoft Anypoint Platform.
	This command allows you to:
	List apis
	List versions
	Create APIs and API versions
	Configure APIs
	Deploy API proxies`,
	Run: func(cmd *cobra.Command, args []string) {
		if apiName == "" {
			cmd.Usage()
			fmt.Println("Error: Please use option --name to specify a search criteria")
			os.Exit(1)
		}

		apiClient := rest.NewAPI(viper.GetString(utils.KEY_URI), viper.GetString(utils.KEY_TOKEN))
		searchParameter := &rest.SearchParameters{
			Name:      apiName,
			Offset:    offset,
			Limit:     limit,
			SortOrder: "", //TODO
			Filter:    "", //TODO
		}

		res, err := apiClient.SearchAPIAsJSON(viper.GetString(utils.KEY_ORG_ID), searchParameter)
		if err != nil {
			log.Fatalf("Error when searching for api %s - %s", searchParameter.Name, err)
		}

		total := res["total"].(float64)
		apis := make([]interface{}, 0)

		if total > 0 {
			apis = res["apis"].([]interface{})
		}

		headers := []string{"API NAME", "VERSION NAME", "API ID", "VERSION ID", "HAS PORTAL"}
		utils.PrintObject(apis, headers, getAPISTabularData)

	},
}

func init() {
	getCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	apiCmd.PersistentFlags().IntVar(&apiId, "api-id", 0, "ID of the API for which the endpoint will be managed")
	apiCmd.PersistentFlags().IntVar(&versionId, "version-id", 0, "ID of the API's version for which the endpoint will be managed")
	apiCmd.Flags().StringVarP(&apiName, "name", "n", "", "Name of the api")
	apiCmd.Flags().IntVar(&offset, "offset", 0, "Return results starting from the specified offset")
	apiCmd.Flags().IntVar(&limit, "limit", 25, "Number of results to return. Default to 25")
}

func getAPISTabularData(apis []interface{}) [][]string {
	data := make([][]string, 0)
	for _, api := range apis {
		versions := api.(map[string]interface{})["versions"].([]interface{})
		currAPI := api.(map[string]interface{})
		for _, version := range versions {
			currVersion := version.(map[string]interface{})
			row := []string{currAPI["name"].(string),
							currVersion["name"].(string),
							fmt.Sprint(currAPI["id"]),
							fmt.Sprint(currVersion["id"]),
							fmt.Sprint(currVersion["portalId"] != nil && currVersion["portalId"].(float64) != 0)}
			data = append(data, row)
		}
	}

	return data
}

