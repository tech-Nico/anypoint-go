package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setEndpointCmd represents the endpoint command
var setEndpointCmd = &cobra.Command{
	Use:   "endpoint set",
	Short: "Set an an API version's endpoint",
	Long: `With this command you will be able to visualize the details related to an API version's endpoint and also
	to make changes to it. An API version's endpoint describes your API's connectivity.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("In set endpoint")
		fmt.Printf("The api-id: %s", apiId)
	},
}

func init() {
	endpointCmd.AddCommand(setEndpointCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// endpointCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
