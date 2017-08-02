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

	searchRes, err := api.SearchAPIAsJSON(orgId, searchParams)
	if err != nil {
		t.Fatalf("Error while searching for api %v : %v", searchParams, err)
	}
	t.Logf("Search results : %s", searchRes)
	if searchRes == nil {
		t.Errorf("Unable to find api using criteria %s", searchParams)
	}

	if searchRes["total"].(float64) <= 0 {
		t.Errorf("API not found.")
	}

}

func Test_FindEnvironments(t *testing.T) {

	searchEnv := "Production"
	appName := "ciccio"

	username, password, uri, orgId := prepTest(t)

	api := rest.NewAPIWithCredentials(uri, username, password)

	env, err := api.FindEnvironmentByName(orgId, searchEnv)

	if err != nil {
		t.Fatalf("Error while searching for environment %q: %s", searchEnv, err)
	}

	if env != nil {
		t.Logf("Environment found: %q", env["name"])
	} else {
		t.Errorf("I was expecting to find environment %q but it couldn't be found", searchEnv)
	}

	app, err := api.ApplicationsByName(orgId, searchEnv, "ciccio")

	if err != nil {
		t.Fatalf("GetApplication did not work while searching for app %q : %s", appName, err)
	}

	t.Logf("I was able to find app %q : %s", appName, app)

}

