package github

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/spothala/go-http-api/client"
	"github.com/spothala/go-http-api/utils"
)

const (
	githubApiUrl = "https://api.github.com"
	version      = "0.3.0"
	userAgent    = "Go " + version
)

type Client struct {
	token string
	debug bool
}

func New(token string, debug bool) *Client {
	c := &Client{}
	c.token = token
	c.debug = debug
	return c
}

/**
* Simple Wrapper on go-http/client which returns map strcuture
**/
func (api *Client) Github(method string, header map[string][]string, endURL string, body io.Reader) (outpu []byte, httpCode int) {
	// API URL
	apiURL := ""
	if strings.Contains(endURL, "https://") {
		apiURL = endURL
	} else {
		apiURL = githubApiUrl + endURL
	}
	// HEADER VALUES
	if header == nil {
		header = url.Values{}
	}
	header["Authorization"] = []string{"token " + api.token}
	// ISSUING HTTP REQUEST
	apiResponse, httpCode := client.ProcessRequest(method, header, apiURL, body)
	if api.debug {
		fmt.Println(apiURL)
		fmt.Println(utils.ReturnPrettyPrintJson(apiResponse))
	}
	return apiResponse, httpCode
}

func checkMapConversion(apiResponse []byte, found bool) int {
	if !found {
		log.Println("Problem with the JSON Conversion to interface type. Output from API:\n " + string(apiResponse))
		return 0
	}
	return 1
}
