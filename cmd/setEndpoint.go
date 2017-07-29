package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tech-nico/anypoint-cli/rest"
	"github.com/spf13/viper"
	"github.com/tech-nico/anypoint-cli/utils"
	"strings"
	"log"
	"fmt"
)

var endpointType string
var endpointUri string
var endpointProxyUri string
var endpointProxyRegistration string
var endpointIsCloudHub bool
var endpointReferencesDomain bool
var endpointTimeout int
var endpointId int

// setEndpointCmd represents the endpoint command
var setEndpointCmd = &cobra.Command{
	Use:   "set",
	Short: "Set an an API version's endpoint",
	Long: `With this command you will be able to visualize the details related to an API version's endpoint and also
	to make changes to it. An API version's endpoint describes your API's connectivity.`,
	Run: func(cmd *cobra.Command, args []string) {
		endpoint := getEndpoint()
		apiMgr := rest.NewAPI(viper.GetString(utils.KEY_URI), viper.GetString(utils.KEY_TOKEN))
		error := apiMgr.SetEndpoint(endpoint)

		if error != nil {
			log.Fatalf("Error while setting endpoint : %s", error)
		}

		fmt.Println("Endpoint updated.")

	},
}

func init() {
	endpointCmd.AddCommand(setEndpointCmd)

	setEndpointCmd.Flags().StringVarP(&endpointType, "type", "t", "", "Endpoint type: http, raml or wsdl")
	setEndpointCmd.Flags().StringVarP(&endpointUri, "uri", "u", "", "Endpoint type: http, raml or wsdl")
	setEndpointCmd.Flags().StringVarP(&endpointProxyUri, "proxy-uri", "p", "", "Endpoint type: http, raml or wsdl")
	setEndpointCmd.Flags().StringVarP(&endpointProxyRegistration, "proxy-registration-uri", "r", "", "Endpoint type: http, raml or wsdl")
	setEndpointCmd.Flags().BoolVarP(&endpointIsCloudHub, "is-cloudhub", "i", false, "Endpoint type: http, raml or wsdl")
	setEndpointCmd.Flags().BoolVar(&endpointReferencesDomain, "use-domain", false, "Endpoint type: http, raml or wsdl")
	setEndpointCmd.Flags().IntVarP(&endpointTimeout, "reponse-timeout", "e", 0, "Endpoint type: http, raml or wsdl")
	setEndpointCmd.Flags().IntVarP(&endpointId, "id", "", 0, "Endpoint type: http, raml or wsdl")
}

func getEndpoint() (*rest.Endpoint) {
	endpoint := &rest.Endpoint{
		Id:                   endpointId,
		VersionID:            versionId,
		ApiID:                apiId,
		OrgID:                viper.GetString(utils.KEY_ORG_ID),
		IsCloudHub:           endpointIsCloudHub,
		ProxyRegistrationUri: endpointProxyRegistration,
		ProxyUri:             endpointProxyUri,
		ReferencesUserDomain: endpointReferencesDomain,
		ResponseTimeout:      endpointTimeout,
		Type:                 endpointType,
		Uri:                  endpointUri,
	}

	if endpoint.VersionID == 0 ||
		endpoint.ApiID == 0 ||
		endpoint.OrgID == "" ||
		endpoint.ProxyRegistrationUri == "" ||
		endpoint.ProxyUri == "" ||
		endpoint.ResponseTimeout == 0 ||
		endpoint.Type == "" ||
		endpoint.Uri == "" {
		log.Print("Warning: one or more endpoint fields have not been specified. Fields will be set to an empty value.")
	}

	if lowerType := strings.ToLower(endpoint.Type); lowerType != "" && lowerType != "http" && lowerType != "wsdl" {
		log.Fatalf("Invalid endpoint type specified %q", endpoint.Type)
	}

	return endpoint
}