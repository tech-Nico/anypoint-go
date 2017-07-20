package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tech-nico/anypoint-cli/rest"
	"github.com/spf13/viper"
	"github.com/tech-nico/anypoint-cli/utils"
)

// setEndpointCmd represents the endpoint command
var setEndpointCmd = &cobra.Command{
	Use:   "endpoint set",
	Short: "Set an an API version's endpoint",
	Long: `With this command you will be able to visualize the details related to an API version's endpoint and also
	to make changes to it. An API version's endpoint describes your API's connectivity.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiMgr := rest.NewAPI(viper.GetString(utils.KEY_URI), viper.GetString(utils.KEY_TOKEN))
		endpoint := apiMgr.GetEndpointAsMap(viper.GetString(utils.KEY_ORG_ID), apiId, versionId)

		if endpoint == nil {
			fmt.Println("Endpoint does not exist. Performing a POST")
		} else {
			fmt.Println("Endpoint does exist. Perfroming a PATCH")
		}
	},
}

func init() {
	endpointCmd.AddCommand(setEndpointCmd)

}
