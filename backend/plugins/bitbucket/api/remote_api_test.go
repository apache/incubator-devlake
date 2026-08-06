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

package api

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

type fakeApiClient struct{ body string }

func (f *fakeApiClient) SetData(name string, data interface{})                    {}
func (f *fakeApiClient) GetData(name string) interface{}                          { return nil }
func (f *fakeApiClient) SetHeaders(headers map[string]string)                     {}
func (f *fakeApiClient) SetBeforeFunction(callback plugin.ApiClientBeforeRequest) {}
func (f *fakeApiClient) GetBeforeFunction() plugin.ApiClientBeforeRequest         { return nil }
func (f *fakeApiClient) SetAfterFunction(callback plugin.ApiClientAfterResponse)  {}
func (f *fakeApiClient) GetAfterFunction() plugin.ApiClientAfterResponse          { return nil }

func (f *fakeApiClient) Get(path string, query url.Values, headers http.Header) (*http.Response, errors.Error) {
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(f.body))}, nil
}
func (f *fakeApiClient) Post(path string, query url.Values, body interface{}, headers http.Header) (*http.Response, errors.Error) {
	return nil, nil
}

func TestListBitbucketRepos_ReturnsNextPage(t *testing.T) {
	client := &fakeApiClient{body: `{"pagelen":2,"page":1,"size":4,
		"next":"https://api.bitbucket.org/2.0/repositories/myworkspace?page=2",
		"values":[{"name":"repo-a","full_name":"myworkspace/repo-a"},
		{"name":"repo-b","full_name":"myworkspace/repo-b"}]}`}
	children, nextPage, err := listBitbucketRepos(client, "myworkspace", BitbucketRemotePagination{Page: 1, PageLen: 2})
	assert.Nil(t, err)
	assert.Len(t, children, 2)
	if assert.NotNil(t, nextPage) {
		assert.Equal(t, 2, nextPage.Page)
	}
}

func TestListBitbucketRepos_LastPageHasNoNextPage(t *testing.T) {
	client := &fakeApiClient{body: `{"pagelen":2,"page":2,"size":4,
		"values":[{"name":"repo-c","full_name":"myworkspace/repo-c"},
		{"name":"repo-d","full_name":"myworkspace/repo-d"}]}`}
	children, nextPage, err := listBitbucketRepos(client, "myworkspace", BitbucketRemotePagination{Page: 2, PageLen: 2})
	assert.Nil(t, err)
	assert.Len(t, children, 2)
	assert.Nil(t, nextPage)
}

func TestListBitbucketWorkspaces_ReturnsNextPage(t *testing.T) {
	client := &fakeApiClient{body: `{"pagelen":1,"page":1,"size":2,
		"next":"https://api.bitbucket.org/2.0/user/workspaces?page=2",
		"values":[{"workspace":{"slug":"ws-a","name":"Workspace A"}}]}`}
	children, nextPage, err := listBitbucketWorkspaces(client, BitbucketRemotePagination{Page: 1, PageLen: 1})
	assert.Nil(t, err)
	assert.Len(t, children, 1)
	if assert.NotNil(t, nextPage) {
		assert.Equal(t, 2, nextPage.Page)
	}
}
