package rest

import (
	"net/http"
	"strings"
	"reflect"
)

const (
	base_path = "/apiplatform/repository/v2"
	org       = "/organizations/{orgId}"
	search    = "/apis?ascending=false&limit={limit}&offset={offset}&query={APIName}&sort={sortOrder}"
)

type API struct {
	client *http.Client
	uri    string
	token  string
}

type Filters string

const (
	FAVORITES     Filters = "pinned"
	ACTIVE        Filters = "active"
	PUBLIC_PORTAL Filters = "public"
	ALL           Filters = "all"
)

type SearchParameters struct {
	limit     int
	offset    int
	name      string
	sortOrder string `default:"createdDate"`
	filter    Filters `default:"all""`
}

func NewApi(client *http.Client, uri, token string) (*API) {
	return &API{
		client,
		uri,
		token,
	}
}

//Search an API by name
func (api *API) ByName(orgId string, params SearchParameters) string {
	typ := reflect.TypeOf(params)

	if params.sortOrder == "" {
		f, _ := typ.FieldByName("sortOrder")
		params.sortOrder = f.Tag.Get("default")
	}

	if params.filter == "" {
		f, _ := typ.FieldByName("filter")
		params.filter = getSearchFilter(f.Tag.Get("default"))
	}

	path := getSearchURL(params, api, orgId)

	return path

}

func getSearchURL(params SearchParameters, api *API, orgId string) string {
	replacer := strings.NewReplacer("{limit}", string(params.limit), "{offset}", string(params.offset), "{APIName}", params.name, "{sortOrder}", params.sortOrder)
	path := api.uri +
		base_path +
		strings.Replace(org, "{orgId}", orgId, -1) +
		replacer.Replace(search)
	return path
}

func getSearchFilter(filter string) Filters {
	switch filter {
	case string(ALL):
		return ALL
	case string(FAVORITES):
		return FAVORITES
	case string(ACTIVE):
		return ACTIVE
	case string(PUBLIC_PORTAL):
		return PUBLIC_PORTAL
	default:
		panic("Invalid filter specified: " + filter)

	}
}
