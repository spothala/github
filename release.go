package github

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spothala/go-http-api/client"
	"github.com/spothala/go-http-api/utils"
)

func (api *Client) GetReleaseByTag(repoName string, tag_name string) map[string]interface{} {
	apiResponse, _ := api.Github("GET", nil, "/repos/"+repoName+"/releases/tags/"+tag_name, nil)
	jsonResponse, found := utils.GetJson(apiResponse).(map[string]interface{})
	checkMapConversion(apiResponse, found)
	return jsonResponse
}

func (api *Client) GetReleaseById(repoName string, release_id string) map[string]interface{} {
	apiResponse, _ := api.Github("GET", nil, "/repos/"+repoName+"/releases/"+release_id, nil)
	jsonResponse, found := utils.GetJson(apiResponse).(map[string]interface{})
	checkMapConversion(apiResponse, found)
	return jsonResponse
}

func (api *Client) DownloadReleaseAsset(repoName string, tag_name string, assetName string) bool {
	asset, rUrl := api.GetAsset(repoName, tag_name, assetName)
	aName := asset["name"].(string)
	aSize := asset["size"].(float64)
	aURL := asset["url"].(string)
	if api.debug {
		fmt.Println("Downloading assets from Release: " + rUrl)
		fmt.Println("Asset Name: " + aName + "[Size: " + strconv.Itoa(int(aSize)) + "]")
	}
	// Create a File
	out, err := os.Create(aName)
	defer out.Close()
	// Create HTTP Request to Download
	httpClient, httpReq := client.CreateRequest("GET", nil, aURL, nil)
	httpReq.Header.Set("Accept", "application/octet-stream") //asset.ContentType)
	httpReq.Header.Set("Content-Length", strconv.Itoa(int(aSize)))
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	//fmt.Println("DOWNLOAD Status Code: " + resp.Status)
	if resp.StatusCode != 200 {
		client.PrintHttpResponseBody(resp)
	}
	// Downloading the Actual File
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	if api.debug {
		fmt.Println("Total Size Downloaded: " + strconv.FormatInt(n, 10))
	}
	var downloaded bool
	if int(aSize) == int(n) {
		downloaded = true
	} else {
		downloaded = false
	}
	return downloaded
}

func (api *Client) GetAsset(repoName string, tag_name string, assetName string) (map[string]interface{}, string) {
	var rConfig map[string]interface{}
	if tag_name == "" {
		rConfig = api.GetReleaseById(repoName, "latest")
	} else {
		rConfig = api.GetReleaseByTag(repoName, tag_name)
	}
	if api.debug {
		fmt.Println("Total Number of Assets: " + strconv.Itoa(len(rConfig["assets"].([]interface{}))))
	}
	var assetConfig map[string]interface{}
	for _, asset := range rConfig["assets"].([]interface{}) {
		//fmt.Println(assetName, asset.URL, asset.Name, asset.ContentType, asset.Size)
		aName := asset.(map[string]interface{})["name"].(string)
		if assetName != "" && assetName == aName {
			assetConfig = asset.(map[string]interface{})
		} else {
			assetConfig = utils.GetJson([]byte(`{"Name":"Bob","Food":"Pickle"}`)).(map[string]interface{})
		}
	}
	return assetConfig, rConfig["url"].(string)
}

func (api *Client) UploadAssetToRelease(repoName string, tag_name string,
	filePath string, contentType string) map[string]interface{} {
	var rConfig map[string]interface{}
	if tag_name == "" {
		rConfig = api.GetReleaseById(repoName, "latest")
	} else {
		rConfig = api.GetReleaseByTag(repoName, tag_name)
	}
	upload_url := rConfig["upload_url"].(string)
	if api.debug {
		fmt.Println("UPLOAD_URL: " + upload_url)
	}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}

	//Preparing Body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	header := url.Values{}
	header["Content-Type"] = []string{contentType}
	apiResponse, _ := api.Github("POST", header, strings.Split(upload_url, "{")[0]+"?name="+filePath, body)
	jsonResponse, found := utils.GetJson(apiResponse).(map[string]interface{})
	checkMapConversion(apiResponse, found)
	return jsonResponse
}

func (api *Client) CreateRelease(repoName string, tagName string, title string, body string) map[string]interface{} {
	jsonOut := utils.WriteJson(map[string]string{"tag_name": tagName, "name": title, "body": body})
	apiResponse, _ := api.Github("POST", nil, "/repos/"+repoName+"/releases", strings.NewReader(string(jsonOut)))
	jsonResponse, found := utils.GetJson(apiResponse).(map[string]interface{})
	checkMapConversion(apiResponse, found)
	if api.debug {
		fmt.Println("Status: Tag " + jsonResponse["tag_name"].(string) + " Created")
		fmt.Println("Location: " + jsonResponse["html_url"].(string))
	}
	return jsonResponse
}

func (api *Client) DeleteRelease(repoName string, tagName string) int {
	rConfig := api.GetReleaseByTag(repoName, tagName)
	jsonResp, httpCode := api.Github("DELETE", nil, "/repos/"+repoName+"/releases/"+strconv.Itoa(int(rConfig["id"].(float64))), nil)
	// 204 Means No-Content Returned
	if httpCode == 204 && jsonResp == nil && api.debug {
		fmt.Println("Tag " + tagName + " got deleted.")
	}
	return httpCode
}

func (api *Client) ListAllReleases(repoName string) []interface{} {
	apiResponse, _ := api.Github("GET", nil, "/repos/"+repoName+"/releases", nil)
	jsonResp, found := utils.GetJson(apiResponse).([]interface{})
	checkMapConversion(apiResponse, found)
	return jsonResp
}

func (api *Client) LatestRelease(repoName string) map[string]interface{} {
	apiResponse, _ := api.Github("GET", nil, "/repos/"+repoName+"/releases/latest", nil)
	release, found := utils.GetJson(apiResponse).(map[string]interface{})
	checkMapConversion(apiResponse, found)
	return release
}

func (api *Client) printReleaseInfo(release map[string]interface{}) {
	output := release["tag_name"].(string) + "\t[ "
	i := 0
	for _, asset := range release["assets"].([]interface{}) {
		i++
		if i == 1 {
			output += "Assets:"
		}
		output += " " + asset.(map[string]interface{})["name"].(string)
	}
	output += " ]"
	if api.debug {
		fmt.Println(output)
	}
}
