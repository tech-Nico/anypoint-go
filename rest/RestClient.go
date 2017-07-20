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

package rest

import (
	"github.com/dghubble/sling"
	"log"
	"io/ioutil"
	"net/http"
	"crypto/tls"
	"github.com/tech-nico/anypoint-cli/utils"
	"fmt"
	"errors"
)

type httpError struct {
	statusCode int
	msg        string
}

func (e *httpError) Error() string {
	return fmt.Sprintf("HTTP Error %d - %s", e.statusCode, e.msg)
}

func NewHttpError(code int, theMsg string) error {
	return &httpError{
		statusCode: code,
		msg:        theMsg,
	}
}

type RestClient struct {
	URI    string
	Sling  *sling.Sling
	client *http.Client
}

func NewRestClient(uri string) (*RestClient) {
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

func (client *RestClient) GET(path string) ([]byte, error) {

	utils.Debug(func() {
		log.Println("REQEST")
		log.Printf("GET %s", client.URI+path)
	})
	req, err := client.Sling.Get(path).Request()
	if err != nil {
		fmt.Printf("\nError building GET request for path %s : %s\n", path, err)
		return nil, err
	}
	res, err := client.client.Do(req)
	defer res.Body.Close()

	httpErr := validateResponse(res, err, "GET", path)
	if httpErr != nil {
		utils.Debug(func() {
			fmt.Printf("\nError while performing GET to %q\n", path)
		})
		return nil, httpErr
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while reading response for %s : %s ", path, err))
	}

	utils.Debug(logResponse("GET", path, res))

	return body, nil
}

//POST - Perform an HTTP POST
func (client *RestClient) POST(body interface{}, responseObj interface{}, path string) (*http.Response, error) {

	utils.Debug(func() {
		log.Println("REQEST")
		log.Printf("POST %s%s", client.URI, path)
	})

	response, err := client.
	Sling.
		Post(path).
		BodyJSON(body).
		ReceiveSuccess(responseObj)

	utils.Debug(logResponse("POST", path, response))

	httpErr := validateResponse(response, err, "POST", path)

	return response, httpErr
}

func validateResponse(response *http.Response, err error, method, path string) error {

	if err != nil {
		return err
	}

	if response.StatusCode == 401 {
		return NewHttpError(401, "Auth token expired. Please login again")
	}

	if response.StatusCode == 404 {
		return NewHttpError(404, fmt.Sprintf("Entity %q not found", path))
	}

	if response.StatusCode >= 400 {
		return NewHttpError(response.StatusCode, fmt.Sprintf("\nError when invoking endpoint %s - %s \nHeaders; %s", path, response.Status, response.Request.Header))
	}

	return nil
}

func logResponse(method, path string, response *http.Response) (func()) {
	return func() {
		log.Printf("Request Headers: %s", response.Request.Header)
		log.Printf("RESPONSE")
		log.Printf("%s %s : %s", method, path, response.Status)
		log.Printf("Response Headers: %s", response.Header)
	}

}