package rest

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/dghubble/sling"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Auth struct {
	client *http.Client
	uri    string
	Token  string
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthToken struct {
	BearerToken string `json:"access_token,omitempty"`
}

const (
	LOGIN     string = "/accounts/login"
	ME        string = "/accounts/api/me"
	HIERARCHY string = "/accounts/api/organizations/{orgId}/hierarchy"
)

func NewAuth(uri, username, password string) *Auth {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}
	return &Auth{
		client,
		uri,
		login(client, uri, username, password),
	}
}

// Login the given user and return the bearer Token
func login(httpClient *http.Client, uri, pUsername, pPassword string) string {

	body := LoginPayload{
		Username: pUsername,
		Password: pPassword,
	}

	authToken := new(AuthToken)

	_, err := sling.
	New().
		Client(httpClient).
		Base(uri).
		Post(LOGIN).
		BodyJSON(body).
		ReceiveSuccess(authToken)

	if err != nil {
		log.Fatal("Error during login with user '", pUsername, "': ", err)
		os.Exit(1)
	} else {
		log.Print("Been able to login: ", *authToken)
	}

	return authToken.BearerToken
}

func (auth *Auth) Me() []byte {
	log.Printf("Call to %s", ME)

	return auth.httpGet(ME)
}

func (auth *Auth) Hierarchy() []byte {
	me := auth.Me()
	var data map[string]interface{}
	if err := json.Unmarshal(me, &data); err != nil {
		fmt.Printf("Error while marshalling JSON response to 'Me' endpoint: %v", err)
		os.Exit(1)
	}
	orgId := data["user"].(map[string]interface{})["organization"].(map[string]interface{})["id"].(string)
	path := strings.Replace(HIERARCHY, "{orgId}", orgId, -1)

	return auth.httpGet(path)
}

func (auth *Auth) httpGet(path string) []byte {
	req, err := sling.New().
		Client(auth.client).
		Base(auth.uri).
		Add("Authorization", "Bearer "+auth.Token).
		Get(path).Request()
	res, err := auth.client.Do(req)
	defer res.Body.Close()
	if err != nil {
		log.Fatal("Error while querying for %s: ", path, err)
	}
	// Check that the server actually sent compressed data
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error while reading response for %s : %s ", path, err)
	}
	return body
}

func (auth *Auth) FindBusinessGroup(path string) string {
	currentOrgId := ""

	groups := auth.createBusinessGroupPath(path)

	var data map[string]interface{}
	hierarchy := auth.Hierarchy()

	if err := json.Unmarshal(hierarchy, &data); err != nil {
		panic("Error while querying for hierarchy..")
	}

	subOrganizations := data["subOrganizations"].([]interface{})
	if len(groups) == 0 {
		return data["id"].(string)
	}

	for _, currGroup := range groups {
		for organization := 0; organization < len(subOrganizations); organization++ {
			jsonObject := subOrganizations[organization].(map[string]interface{})

			if jsonObject["name"].(string) == currGroup {
				currentOrgId = jsonObject["id"].(string)
				log.Printf("The matched org name is: %s", jsonObject["name"].(string))
				subOrganizations = jsonObject["subOrganizations"].([]interface{})
			}
		}
	}

	if currentOrgId == "" {
		panic("Cannot find business group " + path)
	}

	return currentOrgId
}

func (auth *Auth) createBusinessGroupPath(businessGroup string) []string {
	if businessGroup == "" {
		return make([]string, 0)
	}

	groups := []string{}
	group := ""
	pos := 0
	for ; pos < len(businessGroup)-1; pos++ {
		currChar := businessGroup[pos]
		if currChar == '\\' {
			// Double backslash maps to business group with one backslash
			if businessGroup[pos+1] == '\\' {
				group += "\\"
				pos++
				// Single backslash starts a new business group
			} else {
				groups = append(groups, group)
				group = ""
			}
			// Non backslash characters are mapped to the group
		} else {
			group += string(currChar)
		}
	}

	if pos < len(businessGroup) { // Do not end with backslash {
		group += string(businessGroup[len(businessGroup)-1])
	}
	groups = append(groups, string(group))

	return groups
}
