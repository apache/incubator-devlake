package core

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetURIStringPointer_WithSlash(t *testing.T) {
	baseUrl := "http://my-site.com/"
	relativePath := "/api/stuff"
	queryParams := &url.Values{}
	queryParams.Set("id", "1")
	expected := "http://my-site.com/api/stuff?id=1"
	actual, err := GetURIStringPointer(baseUrl, relativePath, queryParams)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, expected, *actual)

}
func TestGetURIStringPointer_WithNoSlash(t *testing.T) {
	baseUrl := "http://my-site.com"
	relativePath := "api/stuff"
	queryParams := &url.Values{}
	queryParams.Set("id", "1")
	expected := "http://my-site.com/api/stuff?id=1"
	actual, err := GetURIStringPointer(baseUrl, relativePath, queryParams)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, expected, *actual)
}
func TestGetURIStringPointer_WithRelativePath(t *testing.T) {
	baseUrl := "http://my-site.com/rest"
	relativePath := "api/stuff"
	queryParams := &url.Values{}
	queryParams.Set("id", "1")
	expected := "http://my-site.com/rest/api/stuff?id=1"
	actual, err := GetURIStringPointer(baseUrl, relativePath, queryParams)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, expected, *actual)
}
func TestGetURIStringPointer_WithRelativePath2(t *testing.T) {
	baseUrl := "https://my-site.com/api/v4/"
	relativePath := "projects/stuff"
	queryParams := &url.Values{}
	queryParams.Set("id", "1")
	expected := "https://my-site.com/api/v4/projects/stuff?id=1"
	actual, err := GetURIStringPointer(baseUrl, relativePath, queryParams)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, expected, *actual)
}

func TestGetURIStringPointer_HandlesRelativePathStartingWithSlash(t *testing.T) {
	baseUrl := "https://my-site.com/api/v4/"
	relativePath := "/user"
	expected := "https://my-site.com/api/v4/user"
	actual, err := GetURIStringPointer(baseUrl, relativePath, nil)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, expected, *actual)
}

func TestGetURIStringPointer_HandlesRelativePathStartingWithSlashWithParams(t *testing.T) {
	baseUrl := "https://my-site.com/api/v4/"
	relativePath := "/user"
	queryParams := &url.Values{}
	queryParams.Set("id", "1")
	expected := "https://my-site.com/api/v4/user?id=1"
	actual, err := GetURIStringPointer(baseUrl, relativePath, queryParams)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, expected, *actual)
}

func TestAddMissingSlashToURL_NoSlash(t *testing.T) {
	baseUrl := "http://my-site.com/rest"
	expected := "http://my-site.com/rest/"
	AddMissingSlashToURL(&baseUrl)
	assert.Equal(t, expected, baseUrl)
}

func TestAddMissingSlashToURL_WithSlash(t *testing.T) {
	baseUrl := "http://my-site.com/rest/"
	expected := "http://my-site.com/rest/"
	AddMissingSlashToURL(&baseUrl)
	assert.Equal(t, expected, baseUrl)
}

func TestRemoveStartingSlashFromPath(t *testing.T) {
	testString := "/user/api"
	expected := "user/api"
	actual := RemoveStartingSlashFromPath(testString)
	assert.Equal(t, expected, actual)
}

func TestRemoveStartingSlashFromPath_EmptyString(t *testing.T) {
	testString := ""
	expected := ""
	actual := RemoveStartingSlashFromPath(testString)
	assert.Equal(t, expected, actual)
}

func TestRemoveStartingSlashFromPath_NoStartingSlash(t *testing.T) {
	testString := "user/api"
	expected := "user/api"
	actual := RemoveStartingSlashFromPath(testString)
	assert.Equal(t, expected, actual)
}
