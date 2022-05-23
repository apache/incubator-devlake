/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helper

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetURIStringPointer_WithSlash(t *testing.T) {
	baseUrl := "http://my-site.com/"
	relativePath := "/api/stuff"
	queryParams := url.Values{}
	queryParams.Set("id", "1")
	expected := "http://my-site.com/api/stuff?id=1"
	actual, err := GetURIStringPointer(baseUrl, relativePath, queryParams)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, expected, *actual)

}
func TestGetURIStringPointer_WithNoSlash(t *testing.T) {
	baseUrl := "http://my-site.com"
	relativePath := "api/stuff"
	queryParams := url.Values{}
	queryParams.Set("id", "1")
	expected := "http://my-site.com/api/stuff?id=1"
	actual, err := GetURIStringPointer(baseUrl, relativePath, queryParams)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, expected, *actual)
}
func TestGetURIStringPointer_WithRelativePath(t *testing.T) {
	baseUrl := "http://my-site.com/rest"
	relativePath := "api/stuff"
	queryParams := url.Values{}
	queryParams.Set("id", "1")
	expected := "http://my-site.com/rest/api/stuff?id=1"
	actual, err := GetURIStringPointer(baseUrl, relativePath, queryParams)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, expected, *actual)
}
func TestGetURIStringPointer_WithRelativePath2(t *testing.T) {
	baseUrl := "https://my-site.com/api/v4/"
	relativePath := "projects/stuff"
	queryParams := url.Values{}
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
	queryParams := url.Values{}
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
