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
)

// apiCmd represents the api command
var setApiCmd = &cobra.Command{
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
		cmd.Usage()
	},
}

func init() {
	setCmd.AddCommand(setApiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	setApiCmd.PersistentFlags().IntVar(&apiId, "api-id", 0, "ID of the API for which the endpoint will be managed")
	setApiCmd.PersistentFlags().IntVar(&versionId, "version-id", 0, "ID of the API's version for which the endpoint will be managed")
}
