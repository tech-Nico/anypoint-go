package rest

import (
	"testing"
	"os"
)

func TestLogin(t *testing.T) {
	username := os.Getenv("TEST_AP_USERNAME")
	password := os.Getenv("TEST_AP_PASSWORD")
	uri := os.Getenv("TEST-AP-URI")

	if username == "" {
		t.Fatal("Unable to read TEST_AP_USERNAME env variable")
	}

	if password == "" {
		t.Fatal("Unable to read TEST_AP_PASSWORD env variable")
	}

	if uri == "" {
		t.Fatal("Unable to read TEST_AP_URI env variable")
	}

	client := NewClient(uri)
	token := login(client, uri, username, password)

	if token == "" {
		t.Errorf("Login returned a nil token")
	}

}
