package rest

import (
	"github.com/dghubble/sling"
	"log"
	"io/ioutil"
	"net/http"
	"crypto/tls"
	"github.com/tech-nico/anypoint-cli/utils"
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

	utils.Debug(func() {
		log.Println("REQEST")
		log.Printf("POST %s", path)
	})
	req, err := client.Sling.Get(path).Request()
	if err != nil {
		log.Fatalf("Error building GET request for path %s : %s", path, err)
	}
	res, err := client.client.Do(req)
	defer res.Body.Close()

	validateResponse(res, err, "GET", path)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error while reading response for %s : %s ", path, err)
	}
	utils.Debug(logResponse("GET", path, res))

	return body
}

//POST - Perform an HTTP POST
func (client *RestClient) POST(body interface{}, responseObj interface{}, path string) (*http.Response, error) {

	utils.Debug(func() {
		log.Println("REQEST")
		log.Printf("POST %s", path)
	})

	response, err := client.
	Sling.
		Post(path).
		BodyJSON(body).
		ReceiveSuccess(responseObj)

	utils.Debug(logResponse("POST", path, response))

	validateResponse(response, err, "POST", path)

	return response, err
}

func validateResponse(response *http.Response, err error, method, path string) {

	if err != nil {
		log.Fatalf("Error while executing http POST %s : %s", path, err)
	}

	if response.StatusCode == 401 {
		log.Fatal("Auth token expired. Please login again")
	}

	if response.StatusCode >= 400 {
		log.Fatalf("\nError while performing HTTP POST %s - %s\n Headers; %s", path, response.Status, response.Request.Header)
	}

}

func logResponse(method, path string, response *http.Response) utils.DebugFunc {
	return func() {
		log.Printf("Request Headers: %s", response.Request.Header)
		log.Printf("RESPONSE")
		log.Printf("%s %s : %s", method, path, response.Status)
		log.Printf("Response Headers: %s", response.Header)
	}

}