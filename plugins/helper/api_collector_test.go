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
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/apache/incubator-devlake/helpers/unithelper"
	"github.com/apache/incubator-devlake/mocks"
	"github.com/apache/incubator-devlake/plugins/helper/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchPageUndetermined(t *testing.T) {
	mockDal := new(mocks.Dal)
	mockDal.On("AutoMigrate", mock.Anything, mock.Anything).Return(nil).Once()
	mockDal.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	mockDal.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

	mockCtx := unithelper.DummySubTaskContext(mockDal)

	mockInput := new(mocks.Iterator)
	mockInput.On("HasNext").Return(true).Once()
	mockInput.On("HasNext").Return(false).Once()
	mockInput.On("Fetch").Return(nil, nil).Once()
	mockInput.On("Close").Return(nil)

	// simulate fetching all pages of jira changelogs for 1 issue id with 1 concurrency,
	// assuming api doesn't return total number of pages.
	// then, we are expecting 2 calls for GetAsync and NextTick each, otherwise, deadlock happens
	getAsyncCounter := 0
	mockApi := new(mocks.RateLimitedApiClient)
	mockApi.On("GetAsync", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		// fake records for first page, no records for second page
		body := "[1,2,3]"
		if getAsyncCounter > 0 {
			body = "[]"
		}
		getAsyncCounter += 1
		res := &http.Response{
			Request: &http.Request{
				URL: &url.URL{},
			},
			Body: ioutil.NopCloser(bytes.NewBufferString(body)),
		}
		handler := args.Get(3).(common.ApiAsyncCallback)
		handler(res)
	}).Twice()
	mockApi.On("NextTick", mock.Anything).Run(func(args mock.Arguments) {
		handler := args.Get(0).(func() error)
		assert.Nil(t, handler())
	}).Twice()
	mockApi.On("HasError").Return(false)
	mockApi.On("WaitAsync").Return(nil)
	mockApi.On("SetAfterFunction", mock.Anything).Return()
	params := struct {
		Name string
	}{Name: "testparams"}
	mockApi.On("Release").Return()

	collector, err := NewApiCollector(ApiCollectorArgs{
		RawDataSubTaskArgs: RawDataSubTaskArgs{
			Ctx:    mockCtx,
			Table:  "whatever rawtable",
			Params: params,
		},
		ApiClient:      mockApi,
		Input:          mockInput,
		UrlTemplate:    "whatever url",
		Concurrency:    1,
		PageSize:       3,
		ResponseParser: GetRawMessageArrayFromResponse,
	})

	assert.Nil(t, err)
	assert.Nil(t, collector.Execute())

	mockDal.AssertExpectations(t)
}
