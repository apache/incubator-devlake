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

package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type AeAppKey helper.AppKey

// SetupAuthentication sets up the HTTP Request Authentication
func (aak *AeAppKey) SetupAuthentication(req *http.Request) errors.Error {
	nonceStr, err := utils.RandLetterBytes(8)
	if err != nil {
		return err
	}
	timestamp := fmt.Sprintf("%v", time.Now().Unix())
	sign := signRequest(req.URL.Query(), aak.AppId, aak.SecretKey, nonceStr, timestamp)
	req.Header.Set("x-ae-app-id", aak.AppId)
	req.Header.Set("x-ae-timestamp", timestamp)
	req.Header.Set("x-ae-nonce-str", nonceStr)
	req.Header.Set("x-ae-sign", sign)
	return nil
}

// AeConn holds the essential information to connect to the AE API
type AeConn struct {
	helper.RestConnection `mapstructure:",squash"`
	AeAppKey              `mapstructure:",squash"`
}

// AeConnection holds AeConn plus ID/Name for database storage
type AeConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	AeConn                `mapstructure:",squash"`
}

func (AeConnection) TableName() string {
	return "_tool_ae_connections"
}

func (connection AeConnection) Sanitize() AeConnection {
	connection.AeAppKey.SecretKey = utils.SanitizeString(connection.AeAppKey.SecretKey)
	return connection
}

func (connection *AeConnection) MergeFromRequest(target *AeConnection, body map[string]interface{}) error {
	secretKey := target.SecretKey
	if err := helper.DecodeMapStruct(body, target, true); err != nil {
		return err
	}
	modifiedSecretKey := target.SecretKey
	if modifiedSecretKey == "" || modifiedSecretKey == utils.SanitizeString(secretKey) {
		target.SecretKey = secretKey
	}
	return nil
}

func signRequest(query url.Values, appId, secretKey, nonceStr, timestamp string) string {
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
