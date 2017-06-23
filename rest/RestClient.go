package rest

import (
	"github.com/dghubble/sling"
	"log"
	"io/ioutil"
	"net/http"
	"crypto/tls"
)

type RestClient struct {
	URI    string
	Sling  *sling.Sling
	client *http.Client
}

func NewClient(uri string) (*RestClient) {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}
	s := sling.New().
		Client(client).
		Base(uri)

	return &RestClient{
		uri,
		s,
		client,
	}
}

func (client *RestClient) AddAuthHeader(token string) (*RestClient) {
	client.Sling.Add("Authorization", "Bearer "+token)
	return client
}

func (client *RestClient) AddOrgHeader(orgId string) (*RestClient) {
	client.Sling.Add("X-ANYPNT-ORG-ID", orgId)
	return client
}

func (client *RestClient) AddEnvHeader(envId string) (*RestClient) {
	client.Sling.Add("X-ANYPNT-ENV-ID", envId)
	return client
}

func (client *RestClient) AddHeader(key, value string) (*RestClient) {
	client.Sling.Add(key, value)
	return client
}

func (client *RestClient) GET(path string) []byte {
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

func (client *RestClient) POST(body interface{}, responseObj interface{}, path string) (*http.Response, error) {

	response, err := client.
	Sling.
		Post(path).
		BodyJSON(body).
		ReceiveSuccess(responseObj)

	if err != nil {
		log.Fatal("Error while executing POST %s : %s", path, err)
	}
	return response, err
}
