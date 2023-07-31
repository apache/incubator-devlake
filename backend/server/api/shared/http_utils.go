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

package shared

import (
	"encoding/base64"
	"github.com/apache/incubator-devlake/core/errors"
	"net/http"
	"strings"
)

func GetUserInfo(req *http.Request) (string, string, error) {
	if req == nil {
		return "", "", errors.Default.New("request is nil")
	}
	user := req.Header.Get("X-Forwarded-User")
	email := req.Header.Get("X-Forwarded-Email")
	if user == "" {
		// fetch with basic auth header
		user, err := GetBasicAuthUserInfo(req)
		return user, "", err
	}
	return user, email, nil
}

func GetBasicAuthUserInfo(req *http.Request) (string, error) {
	if req == nil {
		return "", errors.Default.New("request is nil")
	}
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.Default.New("Authorization is empty")
	}
	basicAuth := strings.TrimPrefix(authHeader, "Basic ")
	if basicAuth == authHeader || basicAuth == "" {
		return "", errors.Default.New("invalid basic auth")
	}
	userInfoData, err := base64.StdEncoding.DecodeString(basicAuth)
	if err != nil {
		return "", errors.Default.Wrap(err, "base64 decode")
	}
	userInfo := strings.Split(string(userInfoData), ":")
	if len(userInfo) != 2 {
		return "", errors.Default.New("invalid user info data")
	}
	return userInfo[0], nil
}
