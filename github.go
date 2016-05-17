package github

import "github.com/spothala/go-http-api/utils"

const (
	githubApiUrl = "https://api.github.com"
	version      = "0.3.0"
	userAgent    = "Go " + version
)

var accessToken string
var debug = false

type ResponseJSON struct {
	Data map[string]string `json:"data"`
}

type Releases struct {
	Releases []GithubReleaseConfig `json:""`
}

type GithubReleaseConfig struct {
	URL     string  `json:"url"`
	UL_URL  string  `json:"upload_url"`
	ID      int     `json:"id"`
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Body    string  `json:"body"`
	State   string  `json:"state"`
	DL_URL  string  `json:"browser_download_url"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	URL         string `json:"url"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	State       string `json:"state"`
	Size        int    `json:"size"`
}

/**
*	Setting Access Token for entire GITHUB API
**/
func SetAccessToken(at string) {
	accessToken = at
}

/**
*	Preparing OAuth Header to handshake with GITHUB API
**/
func prepareAuthHeader() string {
	if accessToken == "" {
		accessToken, _ = utils.ReadFromFile(utils.GetHomeDir() + "/.github_token")
	}
	return "token " + accessToken
}

func SetDebug(value bool) {
	debug = value
}
