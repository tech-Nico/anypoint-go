package rest

import (
	"github.com/dghubble/sling"
	"log"
	"os"
	"net/http"
	"crypto/tls"
	"io/ioutil"
)

type Auth struct {
	client *http.Client
	uri string
	token string
}



type LoginPayload struct {

	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthToken struct {
	BearerToken string `json:"access_token,omitempty"`
}

const (
	LOGIN string = "/accounts/login"
	ME string = "/accounts/api/me"
)


func NewAuth(uri, username, password string) *Auth{
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

// Login the given user and return the bearer token
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

func (auth *Auth) Me() string {
	log.Printf("authtoken while calling me %s", auth.token)

	req, err := sling.New().
		Client(auth.client).
		Base(auth.uri).
		Add("Authorization", "Bearer " + auth.token).
		Get(ME).Request()

	res, err := auth.client.Do(req)
	defer res.Body.Close()


	if err != nil {
		log.Fatal("Error while querying for user's details: ", err)
	}
	// Check that the server actually sent compressed data
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error while reading response for %s : %s ", ME, err)
	}

	return string(body)
}
