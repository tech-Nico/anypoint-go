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

var inputFile string

// apiCmd represents the api command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a MuleSoft entity",
	Long: `The 'set' command can be used to modify various MuleSoft component.
It accepts several subcommands representing the object that you want to set.
Specify --help for a list of objects you can set.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
	setCmd.PersistentFlags().StringVarP(&inputFile, "file", "f", "", "Name of the YAML file to set the required resource")
}
