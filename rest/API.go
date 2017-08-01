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

package rest

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"github.com/tech-nico/anypoint-cli/utils"
	"errors"
	"log"
)

const (
	BASE_PATH         = "/apiplatform/repository/v2"
	ORG_PATH          = "/organizations/{orgId}"
	SEARCH_API_PATH   = "/apis?ascending=false&limit={limit}&offset={offset}&query={APIName}&sort={sortOrder}"
	API_PATH          = BASE_PATH + ORG_PATH + "/apis/{apiId}"
	VERSION_PATH      = API_PATH + "/versions/{versionId}"
	API_ENDPOINT_PATH = VERSION_PATH + "/endpoint"
	APPLICATIONS      = "/armui/api/v1/applications"
	ENVIRONMENTS      = "/accounts/api/organizations/{orgId}/environments"
)

type Endpoint struct {
	Id                   int                        `json:"id"`
	OrgID                string                `json:"masterOrganizationId"`
	ApiID                int
	VersionID            int                `json:"apiVersionId"`
	Type                 string                `json:"type"`
	Uri                  string                    `json:"uri"`
	ProxyUri             string                `json:"proxyUri"`
	ProxyRegistrationUri string `json:"proxyRegistrationUri"`
	IsCloudHub           bool                `json:"isCloudHub"`
	ReferencesUserDomain bool    `json:"referencesUserDomain"`
	ResponseTimeout      int            `json:"responseTimeout"`
}

type API struct {
	client *RestClient
}

type Filters string

const (
	API_FILTER_FAVORITES     Filters = "pinned"
	API_FILTER_ACTIVE        Filters = "active"
	API_FILTER_PUBLIC_PORTAL Filters = "public"
	API_FILTER_ALL           Filters = "all"
)

type SearchParameters struct {
	Name      string
	Limit     int
	Offset    int
	SortOrder string  `default:"createdDate"`
	Filter    Filters `default:"all""`
}

func NewAPIWithCredentials(uri, username, password string) *API {

	client := NewRestClient(uri)
	token := Login(client, uri, username, password)
	client.AddAuthHeader(token)

	return &API{
		client,
	}
}

//NewAPI - Create a new API struct
func NewAPI(uri, token string) *API {
	client := NewRestClient(uri)
	client.AddAuthHeader(token)

	return &API{
		client,
	}
}

func (api *API) SearchAPIAsString(orgID string, params *SearchParameters) (string, error) {
	typ := reflect.TypeOf(*params)

	if params.SortOrder == "" {
		f, _ := typ.FieldByName("SortOrder")
		params.SortOrder = f.Tag.Get("default")
	}

	if params.Filter == "" {
		f, _ := typ.FieldByName("Filter")
		params.Filter = getSearchFilter(f.Tag.Get("default"))
	}

	path := api.getSearchURL(params, orgID)
	apis, error := api.client.GET(path)

	if error != nil {
		fmt.Errorf("Error while searching for api using parameters %v. Error: %s", params, error)
		return "", error
	}

	return string(apis), nil
}

//SearchAPIAsJSON - Search an API by name
func (api *API) SearchAPIAsJSON(orgID string, params *SearchParameters) (map[string]interface{}, error) {
	resp, err := api.SearchAPIAsString(orgID, params)

	if err != nil {
		fmt.Errorf("Error while searching for api with parameters %v : %s", params, err)
		return nil, err
	}

	apis := []byte(resp)
	var jsonObj map[string]interface{}

	if err := json.Unmarshal(apis, &jsonObj); err != nil {
		return nil, errors.New(fmt.Sprintf("Error while querying for api with name %s : %s", params.Name, err))
	}

	return jsonObj, nil

}

func (api *API) GetEndpointAsJSONString(orgId string, apiId, versionId int) (string, error) {
	var path string
	path = strings.Replace(API_ENDPOINT_PATH, "{orgId}", orgId, -1)
	path = strings.Replace(path, "{apiId}", fmt.Sprint(apiId), -1)
	path = strings.Replace(path, "{versionId}", fmt.Sprint(versionId), -1)

	endpointStr, err := api.client.GET(path)
	if err != nil {
		fmt.Errorf("Error while getting endpoint for API %d(version-id %d)", apiId, versionId)
		return "", err
	}

	return string(endpointStr), nil
}

func (api *API) GetEndpointAsMap(orgId string, apiId, versionId int) (map[string]interface{}, error) {
	resp, err := api.GetEndpointAsJSONString(orgId, apiId, versionId)

	if err != nil {
		fmt.Errorf("Error while retrieving endpoint for API %d (version-id %d)", apiId, versionId)
		return nil, err

	}

	endpoint := []byte(resp)
	var jsonObj map[string]interface{}

	if err := json.Unmarshal(endpoint, &jsonObj); err != nil {
		fmt.Errorf("Error while retrieving endpoint: %s", err)
		return nil, err
	}

	return jsonObj, nil

}

func (api *API) getSearchURL(params *SearchParameters, orgId string) string {
	replacer := strings.NewReplacer("{limit}", fmt.Sprint(params.Limit),
		"{offset}", fmt.Sprint(params.Offset),
		"{APIName}", params.Name,
		"{sortOrder}", params.SortOrder)
	path := BASE_PATH +
		strings.Replace(ORG_PATH, "{orgId}", orgId, -1) +
		replacer.Replace(SEARCH_API_PATH)
	utils.Debug(func() {
		fmt.Println("\nThe search API url:", path)
	})

	return path
}

func (api *API) SetEndpoint(endpoint *Endpoint) (error) {
	_, err := api.GetEndpointAsMap(endpoint.OrgID, endpoint.ApiID, endpoint.VersionID)
	exists := true

	if err, ok := err.(*HttpError); ok {
		if err.StatusCode == 404 {
			exists = false
		} else {
			return err
		}
	}

	var path string
	path = strings.Replace(API_ENDPOINT_PATH, "{orgId}", endpoint.OrgID, -1)
	path = strings.Replace(path, "{apiId}", fmt.Sprint(endpoint.ApiID), -1)
	path = strings.Replace(path, "{versionId}", fmt.Sprint(endpoint.VersionID), -1)

	if exists {
		respObj := &Endpoint{}
		_, err := api.client.PATCH(endpoint, path, Application_Json, respObj)

		return err

	} else {
		respObj := make(map[string]interface{})
		_, err := api.client.POST(endpoint, path, Application_Json, respObj)
		return err
	}

	return nil
}

func (api *API) FindEnvironmentByName(orgId, environment string) (map[string]interface{}, error) {

	path := strings.Replace(ENVIRONMENTS, "{orgId}", orgId, -1)

	resp, err := api.client.GET(path)

	utils.Debug(func() {
		log.Printf("FindEnvironmentByName: Get response: %s : %s", resp, err)
	})

	if err != nil {
		fmt.Errorf("Error while searching for environment %s : %s", environment, err)
		return nil, err
	}

	envs := []byte(resp)
	var jsonObj map[string]interface{}

	if err := json.Unmarshal(envs, &jsonObj); err != nil {
		return nil, errors.New(fmt.Sprintf("Error while searching for environment %s : %s", environment, err))
	}

	total := jsonObj["total"].(float64)

	if total == 0 {
		fmt.Printf("No environment found %q", environment)
		utils.Debug(func() {
			log.Printf("Environment %q not found", environment)
		})
		return nil, nil
	}

	data := jsonObj["data"].([]interface{})

	for _, elem := range data {
		if elemMap, ok := elem.(map[string]interface{}); ok {
			if strings.ToUpper(fmt.Sprint(elemMap["name"])) == strings.ToUpper(environment) {
				return elemMap, nil
			}
		}
	}

	return nil, nil
}

func (api *API) GetApplications(orgId, environment, appName string) {
	//path := APPLICATIONS
	api.FindEnvironmentByName(orgId, environment)

}

func getSearchFilter(filter string) Filters {
	switch filter {
	case string(API_FILTER_ALL):
		return API_FILTER_ALL
	case string(API_FILTER_FAVORITES):
		return API_FILTER_FAVORITES
	case string(API_FILTER_ACTIVE):
		return API_FILTER_ACTIVE
	case string(API_FILTER_PUBLIC_PORTAL):
		return API_FILTER_PUBLIC_PORTAL
	default:
		panic("Invalid filter specified: " + filter)

	}
}
