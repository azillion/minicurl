package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	netURL "net/url"
	"strings"
)

type requestType int

const (
	// GET http method enum
	GET requestType = 0
	// POST http method enum
	POST requestType = 1
)

var (
	reqURL  string
	reqType requestType
	reqData string
)

func init() {
	flag.StringVar(&reqURL, "url", "http://localhost:5400", "set url of request")
	flag.StringVar(&reqData, "data", "", "set query params of GET request or data of POST request")
	postPtr := flag.Bool("p", false, "send a post request")
	flag.Parse()

	// set request type
	if *postPtr == true {
		reqType = POST
	} else {
		reqType = GET
	}

	// parse the url
	parsedURL, err := netURL.Parse(reqURL)
	if err != nil {
		panic(err)
	}

	// add query params to reqURL if type is get and data isn't empty
	if reqType == GET && reqData != "" {
		parsedParams := parsedURL.Query()
		err := parseQueryParams(reqData, &parsedParams)
		if err != nil {
			panic(err)
		}
		parsedURL.RawQuery = parsedParams.Encode()
	}

	reqURL = parsedURL.String()
}

func parseQueryParams(s string, queryValues *netURL.Values) error {
	var jsonParams interface{}
	err := json.Unmarshal([]byte(s), &jsonParams)
	if err != nil {
		return err
	}

	m := jsonParams.(map[string]interface{})
	for k, v := range m {
		queryValues.Set(k, fmt.Sprint(v))
	}

	return nil
}

func headersToString(h http.Header) string {
	var s string
	for key, val := range h {
		// Convert each key/value pair in m to a string
		s += fmt.Sprintf("    %s: %s\n", key, val)
	}
	return s
}

func main() {
	if reqType == GET {
		resp, err := http.Get(reqURL)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		fmt.Printf("GET:\n  URL:\n    %s\n\n  Headers:\n%s\n  Body:\n%s\n", reqURL, headersToString(resp.Header), string(body))
	}

	if reqType == POST {
		resp, err := http.Post(reqURL, "application/json", strings.NewReader(reqData))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		fmt.Printf("POST:\n  URL:\n    %s\n\n  Headers:\n%s\n  Body:\n%s\n", reqURL, headersToString(resp.Header), string(body))
	}
}
