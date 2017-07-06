package test

import (
	"testing"
	"github.com/tech-nico/anypoint-cli/rest"
)

func Test_ByNameAsString(t *testing.T) {

}

func Test_SearchAPIAsJSON(t *testing.T) {

	username, password, uri, orgId := prepTest(t)

	searchParams := &rest.SearchParameters{
		Name:      "Test",
		Limit:     10,
		Offset:    0,
		SortOrder: "",
		Filter:    rest.API_FILTER_ALL,
	}

	api := rest.NewAPIWithCredentials(uri, username, password)

	searchRes := api.SearchAPIAsJSON(orgId, searchParams)
	t.Logf("Search results : %s", searchRes)
	if searchRes == nil {
		t.Errorf("Unable to find api using criteria %s", searchParams)
	}

	if searchRes["total"].(float64) <= 0 {
		t.Errorf("API not found.")
	}

}

