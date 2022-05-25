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
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestReadder for test io data
type TestReader struct {
	Err error
}

func (r *TestReader) Read(p []byte) (n int, err error) {
	return 0, r.Err
}

func (r *TestReader) Close() error {
	return nil
}

// it is better to move some where more public.
var ErrUnitTest error = fmt.Errorf("ErrorForTest[%d]", time.Now().UnixNano())

func callback(_ *http.Response, err error) error {
	if err == nil {
		return nil
	}
	return ErrUnitTest
}

func GetConfigForTest(basepath string) *viper.Viper {
	// create the object and load the .env file
	v := viper.New()
	envfile := ".env"
	envbasefile := basepath + ".env.example"
	bytesRead, err := ioutil.ReadFile(envbasefile)
	if err != nil {
		logrus.Warn("Failed to read ["+envbasefile+"] file:", err)
	}
	err = ioutil.WriteFile(envfile, bytesRead, 0644)

	if err != nil {
		logrus.Warn("Failed to write config file ["+envfile+"] file:", err)
	}

	v.SetConfigFile(envfile)
	err = v.ReadInConfig()
	if err != nil {
		path, _ := os.Getwd()
		logrus.Warn("Now in the path [" + path + "]")
		logrus.Warn("Failed to read ["+envfile+"] file:", err)
	}
	v.AutomaticEnv()
	// This line is essential for reading and writing
	v.WatchConfig()
	return v
}

func CreateTestAsyncApiClientWithRateLimitAndCtx(t *testing.T, rateLimiter *ApiRateLimitCalculator, ctx context.Context) (*ApiAsyncClient, error) {
	// set the function of create new default taskcontext for the AsyncApiClient
	gm := gomonkey.ApplyFunc(NewDefaultTaskContext, func(
		cfg *viper.Viper,
		_ core.Logger,
		db *gorm.DB,
		_ context.Context,
		name string,
		subtasks map[string]bool,
		progress chan core.RunningProgress,
	) core.TaskContext {
		return &DefaultTaskContext{
			&defaultExecContext{
				cfg:      cfg,
				logger:   &DefaultLogger{},
				db:       db,
				ctx:      ctx,
				name:     "Test",
				data:     nil,
				progress: progress,
			},
			subtasks,
			make(map[string]*DefaultSubTaskContext),
		}
	})
	defer gm.Reset()
	taskCtx := NewDefaultTaskContext(GetConfigForTest("../../"), nil, nil, nil, "", nil, nil)

	// create ApiClient
	apiClient := &ApiClient{}
	apiClient.Setup("", nil, 3*time.Second)
	apiClient.SetContext(taskCtx.GetContext())

	return CreateAsyncApiClient(taskCtx, apiClient, rateLimiter)
}

// Create an AsyncApiClient object for test
func CreateTestAsyncApiClient(t *testing.T) (*ApiAsyncClient, error) {
	// create rate limit calculator
	rateLimiter := &ApiRateLimitCalculator{
		UserRateLimitPerHour: 36000, // ten times each seconed
	}
	return CreateTestAsyncApiClientWithRateLimitAndCtx(t, rateLimiter, context.Background())
}

// go test -gcflags=all=-l -run ^TestWaitAsync_EmptyWork
func TestWaitAsync_EmptyWork(t *testing.T) {
	asyncApiClient, _ := CreateTestAsyncApiClient(t)

	err := asyncApiClient.WaitAsync()
	assert.Equal(t, err, nil)
}

// go test -gcflags=all=-l -run ^TestWaitAsync_WithWork
func TestWaitAsync_WithWork(t *testing.T) {
	asyncApiClient, _ := CreateTestAsyncApiClient(t)

	gm_info := gomonkey.ApplyMethod(reflect.TypeOf(&DefaultLogger{}), "Info", func(_ *DefaultLogger, _ string, _ ...interface{}) {
	})
	defer gm_info.Reset()

	gm_do := gomonkey.ApplyMethod(reflect.TypeOf(&ApiClient{}), "Do", func(
		apiClient *ApiClient,
		method string,
		path string,
		query url.Values,
		body interface{},
		headers http.Header,
	) (*http.Response, error) {
		return &http.Response{
			Body:       &TestReader{Err: io.EOF},
			StatusCode: http.StatusOK,
		}, nil
	})
	defer gm_do.Reset()

	// check if the callback1 has been finished
	waitSuc := false
	callback1 := func(_ *http.Response, err error) error {
		// wait 0.5 second for wait
		time.Sleep(500 * time.Millisecond)
		waitSuc = true
		return nil
	}

	// begin to test
	err := asyncApiClient.DoAsync("", "", nil, nil, nil, callback1, 0)
	assert.Equal(t, err, nil)

	err = asyncApiClient.WaitAsync()
	assert.Equal(t, err, nil)
	assert.Equal(t, waitSuc, true)
}

// go test -gcflags=all=-l -run ^TestWaitAsync_MutiWork
func TestWaitAsync_MutiWork(t *testing.T) {
	asyncApiClient, _ := CreateTestAsyncApiClient(t)

	gm_info := gomonkey.ApplyMethod(reflect.TypeOf(&DefaultLogger{}), "Info", func(_ *DefaultLogger, _ string, _ ...interface{}) {
	})
	defer gm_info.Reset()

	gm_do := gomonkey.ApplyMethod(reflect.TypeOf(&ApiClient{}), "Do", func(
		apiClient *ApiClient,
		method string,
		path string,
		query url.Values,
		body interface{},
		headers http.Header,
	) (*http.Response, error) {
		return &http.Response{
			Body:       &TestReader{Err: io.EOF},
			StatusCode: http.StatusOK,
		}, nil
	})
	defer gm_do.Reset()

	// check if the callback2 has been finished
	finishedCount := int64(0)
	callback2 := func(_ *http.Response, err error) error {
		// wait 0.5 second for wait
		time.Sleep(500 * time.Millisecond)
		atomic.AddInt64(&finishedCount, 1)
		return nil
	}

	testCount := int64(5)

	// begin to test
	for i := int64(0); i < testCount; i++ {
		err := asyncApiClient.DoAsync("", "", nil, nil, nil, callback2, 0)
		assert.Equal(t, err, nil)
	}

	err := asyncApiClient.WaitAsync()
	assert.Equal(t, err, nil)
	assert.Equal(t, finishedCount, testCount)
}

// go test -gcflags=all=-l -run ^TestDoAsync_OnceSuceess
func TestDoAsync_OnceSuceess(t *testing.T) {
	asyncApiClient, _ := CreateTestAsyncApiClient(t)
	gm_info := gomonkey.ApplyMethod(reflect.TypeOf(&DefaultLogger{}), "Info", func(_ *DefaultLogger, _ string, _ ...interface{}) {
	})
	defer gm_info.Reset()

	gm_do := gomonkey.ApplyMethod(reflect.TypeOf(&ApiClient{}), "Do", func(
		apiClient *ApiClient,
		method string,
		path string,
		query url.Values,
		body interface{},
		headers http.Header,
	) (*http.Response, error) {
		return &http.Response{
			Body:       &TestReader{Err: io.EOF},
			StatusCode: http.StatusOK,
		}, nil
	})
	defer gm_do.Reset()

	// ready to test
	err := asyncApiClient.DoAsync("", "", nil, nil, nil, callback, 0)
	assert.Equal(t, err, nil)

	err = asyncApiClient.WaitAsync()
	assert.Equal(t, err, nil)
}

// go test -gcflags=all=-l -run ^TestDoAsync_TryAndFail
func TestDoAsync_TryAndFail(t *testing.T) {
	asyncApiClient, _ := CreateTestAsyncApiClient(t)
	gm_info := gomonkey.ApplyMethod(reflect.TypeOf(&DefaultLogger{}), "Info", func(_ *DefaultLogger, _ string, _ ...interface{}) {
	})
	defer gm_info.Reset()

	gm_do := gomonkey.ApplyMethod(reflect.TypeOf(&ApiClient{}), "Do", func(
		apiClient *ApiClient,
		method string,
		path string,
		query url.Values,
		body interface{},
		headers http.Header,
	) (*http.Response, error) {
		return &http.Response{
			Body:       &TestReader{Err: ErrUnitTest},
			StatusCode: http.StatusInternalServerError,
		}, nil
	})
	defer gm_do.Reset()

	// ready to test
	err := asyncApiClient.DoAsync("", "", nil, nil, nil, callback, 0)
	assert.Equal(t, err, nil)

	err = asyncApiClient.WaitAsync()
	// there must have err and the err must be ErrUnitTest
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), ErrUnitTest.Error())
	}
}

// go test -gcflags=all=-l -run ^TestDoAsync_TryAndSuceess
func TestDoAsync_TryAndSuceess(t *testing.T) {
	asyncApiClient, _ := CreateTestAsyncApiClient(t)
	gm_info := gomonkey.ApplyMethod(reflect.TypeOf(&DefaultLogger{}), "Info", func(_ *DefaultLogger, _ string, _ ...interface{}) {
	})
	defer gm_info.Reset()

	// counting the retry times
	times := 0
	gm_do := gomonkey.ApplyMethod(reflect.TypeOf(&ApiClient{}), "Do", func(
		apiClient *ApiClient,
		method string,
		path string,
		query url.Values,
		body interface{},
		headers http.Header,
	) (*http.Response, error) {
		times++
		switch times {
		case 1:
			return &http.Response{
				Body:       &TestReader{Err: ErrUnitTest},
				StatusCode: http.StatusInternalServerError,
			}, nil
		case 2:
			return &http.Response{
				Body:       &TestReader{Err: io.EOF},
				StatusCode: http.StatusInternalServerError,
			}, nil
		case 3:
			return &http.Response{
				Body:       &TestReader{Err: io.EOF},
				StatusCode: http.StatusBadRequest,
			}, nil
		case 4:
			return &http.Response{
				Body:       &TestReader{Err: io.EOF},
				StatusCode: http.StatusMultipleChoices,
			}, nil
		case 5:
			return &http.Response{
				Body:       &TestReader{Err: io.EOF},
				StatusCode: http.StatusOK,
			}, nil
		default:
			assert.Empty(t, TestNoRunHere)
			return &http.Response{
				Body:       &TestReader{Err: io.EOF},
				StatusCode: http.StatusOK,
			}, TestError
		}
	})
	defer gm_do.Reset()
	asyncApiClient.SetMaxRetry(5)

	// ready to test
	err := asyncApiClient.DoAsync("", "", nil, nil, nil, callback, 0)
	assert.Equal(t, err, nil)

	err = asyncApiClient.WaitAsync()
	assert.Equal(t, err, nil)
	assert.Equal(t, times, 4)
}
