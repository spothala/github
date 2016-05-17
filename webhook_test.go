package github

import (
	"strings"
	"testing"
)

const (
	URL    = "http://test-local-webhook/"
	EVENTS = "push,release"
	HOOKCT = 1
)

var hookID int

func TestCreateWebhook(t *testing.T) {
	hookMap := githubAPI.CreateWebhook(ORGREPO, URL, EVENTS)
	hookID = int(hookMap["id"].(float64))
	if hookID <= 0 {
		t.Error("Webhook is not created")
	}
	if hookMap["url"] == "" {
		t.Error("Webhook is not created")
	}
	for _, event := range hookMap["events"].([]interface{}) {
		if !strings.Contains(EVENTS, event.(string)) {
			t.Error("Webhook created but events not matched")
		}
	}
	if hookMap["config"].(map[string]interface{})["url"].(string) != URL {
		t.Error("Webhook created but URL doesnt match")
	}
}

func TestListWebhook(t *testing.T) {
	_, count := githubAPI.ListAllWebhooks(ORGREPO)
	if count != HOOKCT {
		t.Error("Didnt match the number of hooks exist")
	}
}

func TestGetWebhook(t *testing.T) {
	hookMap := githubAPI.GetWebhook(ORGREPO, hookID)
	if hookID != int(hookMap["id"].(float64)) {
		t.Error("Didnt get the requested webhook")
	}
	if hookMap["url"] == "" {
		t.Error("Didnt get the requested webhook")
	}
	for _, event := range hookMap["events"].([]interface{}) {
		if !strings.Contains(EVENTS, event.(string)) {
			t.Error("Didnt get the requested webhook")
		}
	}
	if hookMap["config"].(map[string]interface{})["url"].(string) != URL {
		t.Error("Didnt get the requested webhook")
	}
}

func TestDeleteWebhook(t *testing.T) {
	respCode := githubAPI.DeleteWebhook(ORGREPO, hookID)
	if respCode != 204 {
		t.Error("Webhook is not deleted.")
	}
}
