package test

import (
	"testing"
	"github.com/tech-nico/anypoint-cli/rest"
	"os"
)

func Test_ByNameAsJSON(t *testing.T) {

	username, password, uri := prepTest(t)
	orgId := os.Getenv(env_org_id)
	if orgId == "" {
		t.Fatalf("Env variable %s not defined", env_org_id)
	}

	searchParams := &rest.SearchParameters{
		Name:      "Flights",
		Limit:     10,
		Offset:    0,
		SortOrder: "",
		Filter:    rest.API_FILTER_ALL,
	}

	api := rest.NewAPIWithCredentials(uri, username, password)

	searchRes := api.ByNameAsJSON(orgId, searchParams)
	t.Logf("Search results : %s", searchRes)
	if searchRes == nil {
		t.Errorf("Unable to find api using criteria %s", searchParams)
	}

	if searchRes["total"].(float64) <= 0 {
		t.Errorf("API not found.")
	}

}
