package rest

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"github.com/tech-nico/anypoint-cli/utils"
)

const (
	BASE_PATH  = "/apiplatform/repository/v2"
	ORG        = "/organizations/{orgId}"
	SEARCH_API = "/apis?ascending=false&limit={limit}&offset={offset}&query={APIName}&sort={sortOrder}"
)

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

func (api *API) SearchAPIAsString(orgID string, params *SearchParameters) string {
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
	apis := api.client.GET(path)
	return string(apis)
}

//SearchAPIAsJSON - Search an API by name
func (api *API) SearchAPIAsJSON(orgID string, params *SearchParameters) map[string]interface{} {
	apis := []byte(api.SearchAPIAsString(orgID, params))
	var jsonObj map[string]interface{}

	if err := json.Unmarshal(apis, &jsonObj); err != nil {
		fmt.Printf("Error while querying for api with name %s. : %s", params.Name, err)
		os.Exit(1)
	}

	return jsonObj

}

func (api *API) getSearchURL(params *SearchParameters, orgId string) string {
	replacer := strings.NewReplacer("{limit}", fmt.Sprint(params.Limit),
		"{offset}", fmt.Sprint(params.Offset),
		"{APIName}", params.Name,
		"{sortOrder}", params.SortOrder)
	path := BASE_PATH +
		strings.Replace(ORG, "{orgId}", orgId, -1) +
		replacer.Replace(SEARCH_API)
	utils.Debug(func() {
		fmt.Println("\nThe search API url:", path)
	})

	return path
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
