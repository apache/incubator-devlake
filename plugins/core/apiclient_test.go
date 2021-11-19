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
	expected := "http://my-site.com/api/stuff?id=1"
	actual, err := GetURIStringPointer(baseUrl, relativePath, queryParams)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, expected, *actual)
}
