package rest

import (
	"github.com/dghubble/sling"
	"log"
	"os"
	"net/http"
	"crypto/tls"
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
	ME string = "/accounts/me"
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
	ciccio := ""
	req, err := sling.
		New().
		Client(auth.client).
		Base(auth.uri).
		Get(ME).BodyForm(ciccio).Receive()
	log.Printf("The response: ", req.)
	log.Printf("The error: ", err)
	return "me"
}
