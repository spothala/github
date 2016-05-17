package github

import (
	"os"
	"testing"
)

const (
	ORGREPO = "spothala/test-api"
	TAGNAME = "test-1-0"
	TITLE   = "Release Created Automatically"
	BODY    = "This test was created by github-api test-cases"
	FILE    = "example.zip"
	CTYPE   = "application/zip"
	TOKEN   = "" // PERSONAL_ACCESS_TOKEN
)

var githubAPI = New(TOKEN, true)

func TestCreateRelease(t *testing.T) {
	relUrl := githubAPI.CreateRelease(ORGREPO, TAGNAME, TITLE, BODY)
	if relUrl["url"].(string) == "" {
		t.Error("Release is not created")
	}
}

func TestUploadAsset(t *testing.T) {
	response := githubAPI.UploadAssetToRelease(ORGREPO, TAGNAME, FILE, CTYPE)
	if response["state"].(string) != "uploaded" {
		t.Error("Uploading an asset to release " + TAGNAME + " is failed ")
	}
	os.Remove(FILE)
}

func TestDownloadAsset(t *testing.T) {
	downloaded := githubAPI.DownloadReleaseAsset(ORGREPO, TAGNAME, FILE)
	if !downloaded {
		t.Error("Downloading of file " + FILE + " failed.")
	}
}

func TestDeleteRelease(t *testing.T) {
	respCode := githubAPI.DeleteRelease(ORGREPO, TAGNAME)
	if respCode != 204 {
		t.Error("Release is not deleted.")
	}
}
