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

	"github.com/spf13/cobra"
	"os"
)

var apiId, versionId int

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

		fmt.Println("endpoint called")
	},
}

func init() {
	apiCmd.AddCommand(endpointCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// endpointCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	endpointCmd.Flags().IntVar(&apiId, "api-id", 0, "ID of the API for which the endpoint will be managed")
	endpointCmd.Flags().IntVar(&versionId, "version-id", 0, "ID of the API's version for which the endpoint will be managed")
}
