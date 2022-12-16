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
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/unithelper"
	"github.com/apache/incubator-devlake/mocks"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllJobs(t *testing.T) {
	const testPageSize int = 100

	var remoteData []*models.Job = []*models.Job{
		{
			Name:        "devlake",
			Color:       "blue",
			Class:       "hudson.model.FreeStyleProject",
			Base:        "",
			URL:         "https://test.nddtf.com/job/devlake/",
			Description: "",
		},
		{
			Name:        "dir-test",
			Color:       "",
			Class:       "hudson.model.FreeStyleProject",
			Base:        "",
			URL:         "https://test.nddtf.com/job/dir-test/",
			Description: "",
			Jobs:        &[]models.Job{{}},
		},
		{
			Name:        "dir-test-2",
			Color:       "",
			Class:       "hudson.model.FreeStyleProject",
			Base:        "",
			URL:         "https://test.nddtf.com/job/dir-test/job/dir-test-2/",
			Description: "",
			Jobs:        &[]models.Job{{}},
		},
		{
			Name:        "free",
			Color:       "blue",
			Class:       "hudson.model.FreeStyleProject",
			Base:        "",
			URL:         "https://test.nddtf.com/job/dir-test/job/dir-test-2/job/free/",
			Description: "",
		},
		{
			Name:        "free1",
			Color:       "blue",
			Class:       "hudson.model.FreeStyleProject",
			Base:        "",
			URL:         "https://test.nddtf.com/job/dir-test/job/dir-test-2/job/free1/",
			Description: "",
		},
	}

	var expectJobs []*models.Job = []*models.Job{
		{
			FullName:    "devlake",
			Path:        "",
			Name:        "devlake",
			Color:       "blue",
			Class:       "hudson.model.FreeStyleProject",
			Base:        "",
			URL:         "https://test.nddtf.com/job/devlake/",
			Description: "",
		},
		{
			FullName:    "dir-test/dir-test-2/free",
			Path:        "job/dir-test/job/dir-test-2/",
			Name:        "free",
			Color:       "blue",
			Class:       "hudson.model.FreeStyleProject",
			Base:        "",
			URL:         "https://test.nddtf.com/job/dir-test/job/dir-test-2/job/free/",
			Description: "",
		},
		{
			FullName:    "dir-test/dir-test-2/free1",
			Path:        "job/dir-test/job/dir-test-2/",
			Name:        "free1",
			Color:       "blue",
			Class:       "hudson.model.FreeStyleProject",
			Base:        "",
			URL:         "https://test.nddtf.com/job/dir-test/job/dir-test-2/job/free1/",
			Description: "",
		},
	}

	var expectPaths []*models.Job = []*models.Job{
		{
			FullName:    "dir-test",
			Path:        "",
			Name:        "dir-test",
			Color:       "",
			Class:       "hudson.model.FreeStyleProject",
			Base:        "",
			URL:         "https://test.nddtf.com/job/dir-test/",
			Description: "",
			Jobs:        &[]models.Job{{}},
		},
		{
			FullName:    "dir-test/dir-test-2",
			Path:        "job/dir-test/",
			Name:        "dir-test-2",
			Color:       "",
			Class:       "hudson.model.FreeStyleProject",
			Base:        "",
			URL:         "https://test.nddtf.com/job/dir-test/job/dir-test-2/",
			Description: "",
			Jobs:        &[]models.Job{{}},
		},
	}

	mockApiClient := mocks.NewApiClientGetter(t)

	var data struct {
		Jobs []json.RawMessage `json:"jobs"`
	}

	// first path jobs
	js, err1 := json.Marshal(remoteData[0])
	assert.Nil(t, err1)
	data.Jobs = append(data.Jobs, js)

	js, err1 = json.Marshal(remoteData[1])
	assert.Nil(t, err1)
	data.Jobs = append(data.Jobs, js)

	js, err1 = json.Marshal(data)
	assert.Nil(t, err1)

	res := &http.Response{}
	res.Body = io.NopCloser(bytes.NewBuffer(js))
	res.StatusCode = http.StatusOK

	mockApiClient.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(res, nil).Once()

	// second path jobs
	data.Jobs = []json.RawMessage{}

	js, err1 = json.Marshal(remoteData[2])
	assert.Nil(t, err1)
	data.Jobs = append(data.Jobs, js)

	js, err1 = json.Marshal(data)
	assert.Nil(t, err1)

	res = &http.Response{}
	res.Body = io.NopCloser(bytes.NewBuffer(js))
	res.StatusCode = http.StatusOK

	mockApiClient.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(res, nil).Once()

	// third path jobs
	data.Jobs = []json.RawMessage{}
	js, err1 = json.Marshal(remoteData[3])
	assert.Nil(t, err1)
	data.Jobs = append(data.Jobs, js)

	js, err1 = json.Marshal(remoteData[4])
	assert.Nil(t, err1)
	data.Jobs = append(data.Jobs, js)

	js, err1 = json.Marshal(data)
	assert.Nil(t, err1)

	res = &http.Response{}
	res.Body = io.NopCloser(bytes.NewBuffer(js))
	res.StatusCode = http.StatusOK

	mockApiClient.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(res, nil).Once()

	basicRes = unithelper.DummyBasicRes(func(mockDal *mocks.Dal) {})

	var jobs []*models.Job

	var paths []*models.Job

	err := GetAllJobs(mockApiClient, "", "", testPageSize, func(job *models.Job, isPath bool) errors.Error {
		if isPath {
			paths = append(paths, job)
		} else {
			jobs = append(jobs, job)
		}
		return nil
	})

	assert.Equal(t, err, nil)
	assert.Equal(t, expectJobs, jobs)
	assert.Equal(t, expectPaths, paths)
}
