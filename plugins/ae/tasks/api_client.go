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

package tasks

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func getSign(query url.Values, appId, secretKey, nonceStr, timestamp string) string {
	// clone query because we need to add items
	kvs := make([]string, 0, len(query)+3)
	kvs = append(kvs, fmt.Sprintf("app_id=%s", appId))
	kvs = append(kvs, fmt.Sprintf("timestamp=%s", timestamp))
	kvs = append(kvs, fmt.Sprintf("nonce_str=%s", nonceStr))
	for key, values := range query {
		for _, value := range values {
			kvs = append(kvs, fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(value)))
		}
	}

	// sort by alphabetical order
	sort.Strings(kvs)

	// generate text for signature
	querystring := fmt.Sprintf("%s&key=%s", strings.Join(kvs, "&"), url.QueryEscape(secretKey))

	// sign it
	hasher := md5.New()
	_, err := hasher.Write([]byte(querystring))
	if err != nil {
		return ""
	}
	return strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
}

func CreateApiClient(taskCtx core.TaskContext) (*helper.ApiAsyncClient, error) {
	// load and process cconfiguration
	endpoint := taskCtx.GetConfig("AE_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("invalid AE_ENDPOINT")
	}
	appId := taskCtx.GetConfig("AE_APP_ID")
	if appId == "" {
		return nil, fmt.Errorf("invalid AE_APP_ID")
	}
	secretKey := taskCtx.GetConfig("AE_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("invalid AE_SECRET_KEY")
	}
	proxy := taskCtx.GetConfig("AE_PROXY")

	apiClient, err := helper.NewApiClient(endpoint, nil, 0, proxy, taskCtx.GetContext())
	if err != nil {
		return nil, err
	}
	apiClient.SetBeforeFunction(func(req *http.Request) error {
		nonceStr := RandString(8)
		timestamp := fmt.Sprintf("%v", time.Now().Unix())
		sign := getSign(req.URL.Query(), appId, secretKey, nonceStr, timestamp)
		req.Header.Set("x-ae-app-id", appId)
		req.Header.Set("x-ae-timestamp", timestamp)
		req.Header.Set("x-ae-nonce-str", nonceStr)
		req.Header.Set("x-ae-sign", sign)
		return nil
	})

	// create ae api client
	asyncApiCLient, err := helper.CreateAsyncApiClient(taskCtx, apiClient, nil)
	if err != nil {
		return nil, err
	}

	return asyncApiCLient, nil
}
