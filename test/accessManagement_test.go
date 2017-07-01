package test

import (
	"testing"
	"os"
	"encoding/json"
	"github.com/tech-nico/anypoint-cli/rest"
)

const (
	env_uri      = "TEST_AP_URI"
	env_username = "TEST_AP_USERNAME"
	env_password = "TEST_AP_PASSWORD"
	env_org_id   = "TEST_AP_ORGID"
)

func prepTest(t *testing.T) (string, string, string) {
	username := os.Getenv(env_username)
	password := os.Getenv(env_password)
	uri := os.Getenv(env_uri)
	if username == "" {
		t.Fatalf("Environment variable %s not defined ", env_username)
	}
	if password == "" {
		t.Fatalf("Environment variable %s not defined ", env_password)
	}
	if uri == "" {
		t.Fatalf("Environment variable %s not defined ", env_uri)
	}
	return username, password, uri
}

func TestLogin(t *testing.T) {
	username, password, uri := prepTest(t)

	client := rest.NewRestClient(uri)
	token := rest.Login(client, uri, username, password)

	if token == "" {
		t.Errorf("Login returned a nil token")
	}

	t.Logf("Login successful. Token: %s", token)

}

func TestMe(t *testing.T) {
	username, password, uri := prepTest(t)
	client := rest.NewAuthWithCredentials(uri, username, password)
	meData := client.Me()
	var meJson map[string]interface{}
	err := json.Unmarshal(meData, &meJson)

	if err != nil {
		t.Errorf("Error unmarshalling the response returned by Me()")
	}

	if meJson["user"] == nil {
		t.Errorf("Expected 'user' attribute in %s", meJson)
	}

	user := meJson["user"].(map[string]interface{})

	if user["id"] == nil {
		t.Errorf("No user id found")
	}

}

func TestAuth_Hierarchy(t *testing.T) {
	username, password, uri := prepTest(t)
	client := rest.NewAuthWithCredentials(uri, username, password)
	meData := client.Hierarchy()
	var vJson map[string]interface{}
	err := json.Unmarshal(meData, &vJson)

	if err != nil {
		t.Errorf("Error unmarshalling the response returned by Hierarchy()")
	}

	if vJson["id"] == nil {
		t.Fatalf("Expected 'id' attribute in %s", vJson)
	}

	t.Logf("The org ID found is %s", vJson["id"])

}