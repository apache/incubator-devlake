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

package utils

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/errors"
	"net/http"
	"os"
	"strings"
)

const (
	EncodeKeyEnvStr = "ENCRYPTION_SECRET"
	apiKeyLen       = 128
)

func GetEncodeKeyEnv(ctx context.Context) (string, bool) {
	encodeKeyEnv := strings.TrimSpace(os.Getenv(EncodeKeyEnvStr))
	if encodeKeyEnv != "" {
		return encodeKeyEnv, true
	}
	encodeKeyEnv = strings.TrimSpace(config.GetConfig().GetString(EncodeKeyEnvStr))
	if encodeKeyEnv != "" {
		return encodeKeyEnv, true
	}
	return "", false
}

func GenerateApiKey(ctx context.Context) (string, string, errors.Error) {
	randomApiKey, randomLetterErr := RandLetterBytes(apiKeyLen)
	if randomLetterErr != nil {
		return "", "", errors.Default.Wrap(randomLetterErr, "random letters")
	}
	hashedApiKey, err := GenerateApiKeyWithToken(ctx, randomApiKey)
	return randomApiKey, hashedApiKey, err
}

func GenerateApiKeyWithToken(ctx context.Context, token string) (string, errors.Error) {
	encodeKeyEnv, exist := GetEncodeKeyEnv(ctx)
	if !exist {
		err := errors.Default.New("encode key env doesn't exist")
		return "", err
	}
	h := hmac.New(sha256.New, []byte(encodeKeyEnv))
	if _, err := h.Write([]byte(token)); err != nil {
		return "", errors.Default.Wrap(err, "hmac write api key")
	}
	hashedApiKey := fmt.Sprintf("%x", h.Sum(nil))
	return hashedApiKey, nil
}

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
