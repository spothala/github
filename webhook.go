package github

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spothala/go-http-api/utils"
)

func printHookHeader() {
	fmt.Println("ID \t\t HOOK URL \t\t\t\t EVENTS")
	fmt.Println("-- \t\t -------- \t\t\t\t ------")
}

func printHookDetails(hookMap map[string]interface{}) {
	output := strconv.Itoa(int(hookMap["id"].(float64))) + " \t"
	output += hookMap["config"].(map[string]interface{})["url"].(string) + " \t\t"
	for _, event := range hookMap["events"].([]interface{}) {
		output += event.(string) + ","
	}
	fmt.Println(strings.TrimSuffix(output, ","))
}

func (api *Client) ListAllWebhooks(repoName string) ([]interface{}, int) {
	apiResponse, _ := api.Github("GET", nil, "/repos/"+repoName+"/hooks", nil)
	jsonResp := utils.GetJson(apiResponse)
	hooks, found := jsonResp.([]interface{})
	checkMapConversion(apiResponse, found)
	if debug {
		printHookHeader()
	}
	hookCount := 0
	for _, hook := range hooks {
		hookMap := hook.(map[string]interface{})
		if hookMap["name"].(string) == "web" {
			if debug {
				printHookDetails(hookMap)
			}
			hookCount++
		}
	}
	return hooks, hookCount
}

func (api *Client) GetWebhook(repoName string, hookId int) map[string]interface{} {
	apiResponse, _ := api.Github("GET", nil, "/repos/"+repoName+"/hooks/"+strconv.Itoa(hookId), nil)
	if debug {
		fmt.Println(utils.ReturnPrettyPrintJson(apiResponse))
	}
	return utils.GetJson(apiResponse).(map[string]interface{})
}

func (api *Client) CreateWebhook(repoName string, url string, events string) map[string]interface{} {
	jsonOut := utils.WriteJson(map[string]interface{}{"name": "web", "active": true, "events": strings.Split(events, ","),
		"config": map[string]string{"url": url, "content_type": "json"}})
	apiResponse, _ := api.Github("POST", nil, "/repos/"+repoName+"/hooks", strings.NewReader(string(jsonOut)))
	hookMap := utils.GetJson(apiResponse).(map[string]interface{})
	if debug {
		printHookHeader()
		printHookDetails(hookMap)
	}
	return hookMap
}

func (api *Client) DeleteWebhook(repoName string, hookId int) int {
	apiResponse, httpCode := api.Github("DELETE", nil, "/repos/"+repoName+"/hooks/"+strconv.Itoa(hookId), nil)
	// 204 Means No-Content Returned
	if httpCode == 204 && utils.GetJson(apiResponse) == nil && debug {
		fmt.Println("Hook with ID " + strconv.Itoa(hookId) + ", got deleted.")
	}
	return httpCode
}
