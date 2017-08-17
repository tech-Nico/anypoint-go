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
	"github.com/tech-nico/anypoint-cli/rest"
	"github.com/tech-nico/anypoint-cli/utils"
	"github.com/spf13/viper"
	"log"
)

// setClusterCmd represents the setCluster command
var setClusterCmd = &cobra.Command{
	Use:     "cluster",
	Aliases: []string{"create cluster"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if inputFile == "" {
			return fmt.Errorf("Use --file to specify the YAML file containing the cluster(s) details")
		}

		if envName == "" {
			return fmt.Errorf("Please specify --environment to indicate where clusters must be created.")
		}

		clusters := rest.Clusters{}

		err := utils.OpenYAMLFile(inputFile, &clusters)

		if err != nil {
			return fmt.Errorf("Error while parsing YAML file %q . Error: %s", inputFile, err)
		}

		err = CreateMultipleClusters(clusters)

		if err != nil {
			return fmt.Errorf("Error while creating cluster(s): %s", err)
		}

		return nil

	},
}

func validateServers(allServers []interface{}, clusters rest.Clusters) error {

	for _, currCluster := range clusters.ACluster {
		for _, currServer := range currCluster.Servers {
			found := false
			for _, existingServer := range allServers {
				utils.Debug(func() {
					log.Printf("\nCheck that server %s == %s", existingServer.(map[string]interface{})["name"].(string), currServer.Name)
				})
				if existingServer.(map[string]interface{})["name"].(string) == currServer.Name {
					found = true
				}
			}

			if !found {
				return fmt.Errorf("Server %s is not a valid server in environment %s", currServer.Name, envName)
			}
		}
	}

	return nil
}

func CreateMultipleClusters(clusters rest.Clusters) error {
	api := rest.NewAPI(viper.GetString(utils.KEY_URI), viper.GetString(utils.KEY_TOKEN))
	allServers, err := api.GetAllServers(viper.GetString(utils.KEY_ORG_ID), envName)

	if err != nil {
		return err
	}

	err = validateServers(allServers, clusters)

	if err != nil {
		return err
	}
	utils.Debug(func() { log.Print("Server successfully validated!") })

	for _, server := range clusters.ACluster {
		resp, err := api.CreateCluster(viper.GetString(utils.KEY_ORG_ID), envName, &server)
		if err != nil {
			fmt.Printf("Error when creating cluster %s : %s\n", server.ClusterName, err)
		}
		utils.Debug(func() { log.Printf("Cluster created: %s", resp) })
	}
	return nil
}

func init() {
	setCmd.AddCommand(setClusterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setClusterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	setClusterCmd.Flags().StringVarP(&envName, "environment", "e", "", "Name of the environment where to create the clusters")
}
