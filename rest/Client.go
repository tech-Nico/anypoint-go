package rest

import (
	"github.com/dghubble/sling"
	"log"
	"io/ioutil"
	"net/http"
	"crypto/tls"
)

type Client struct {
	URI    string
	Sling  *sling.Sling
	client *http.Client
}

func NewClient(uri string) (*Client) {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}
	s := sling.New().
		Client(client).
		Base(uri)

	return &Client{
		uri,
		s,
		client,
	}
}

func (client *Client) AddAuthHeader(token string) (*Client) {
	client.Sling.Add("Authorization", "Bearer "+token)
	return client
}

func (client *Client) AddHeader(key, value string) (*Client) {
	client.Sling.Add(key, value)
	return client
}

func (client *Client) GET(path string) []byte {
	req, err := client.
	Sling.
		Get(path).Request()

	res, err := client.client.Do(req)
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
