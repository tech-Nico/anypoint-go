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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tech-nico/anypoint-cli/rest"
	"github.com/tech-nico/anypoint-cli/utils"
	"encoding/json"
)

var apiName string

// appCmd represents the app command
var apiSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search an API",
	Long: `Search an API by various criteria. For example:
  api search --name "My API" --limit 10 --filter Portal
  `,
	Run: func(cmd *cobra.Command, args []string) {

		apiClient := rest.NewAPI(viper.GetString(utils.KEY_URI), viper.GetString(utils.KEY_TOKEN))
		searchParameter := &rest.SearchParameters{
			Name:      apiName,
			Offset:    0,
			Limit:     1,
			SortOrder: "",
			Filter:    "",
		}

		format := viper.GetString(utils.KEY_FORMAT)
		switch format {
		case "list":
			res := apiClient.SearchAPIAsJSON(viper.GetString(utils.KEY_ORG_ID), searchParameter)

			total := res["total"]
			if total == 0 {
				fmt.Println("No APIs match name " + apiName)
				os.Exit(0)
			}

			apis := res["apis"]
			printAPIs(apis.([]interface{}))
			break
		case "json":
			res := apiClient.SearchAPIAsJSON(viper.GetString(utils.KEY_ORG_ID), searchParameter)
			b, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				fmt.Println("error:", err)
			}
			os.Stdout.Write(b)
			break
		}

	},
}

func printAPIs(apis []interface{}) {
	headers := []string{"API Name", "Version Name", "API ID", "Version ID"}

	data := make([][]string, 0)
	for _, api := range apis {
		versions := api.(map[string]interface{})["versions"].([]interface{})
		currAPI := api.(map[string]interface{})
		for _, version := range versions {
			currVersion := version.(map[string]interface{})
			row := []string{currAPI["name"].(string),
							currVersion["name"].(string),
							fmt.Sprint(currAPI["id"]),
							fmt.Sprint(currVersion["id"])}
			data = append(data, row)
		}
	}

	utils.PrintTabular(headers, data)
}

func init() {
	apiCmd.AddCommand(apiSearchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	apiSearchCmd.Flags().StringVar(&apiName, "api-name", "", "Name of the api")
}
