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
	"fmt"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/go-playground/validator/v10"
	"strings"
)

var vld *validator.Validate

var basicRes context.BasicRes
var dsHelper *api.DsHelper[models.GithubConnection, models.GithubRepo, models.GithubScopeConfig]
var raProxy *api.DsRemoteApiProxyHelper[models.GithubConnection]
var raScopeList *api.DsRemoteApiScopeListHelper[models.GithubConnection, models.GithubRepo, GithubRemotePagination]
var raScopeSearch *api.DsRemoteApiScopeSearchHelper[models.GithubConnection, models.GithubRepo]

func Init(br context.BasicRes, p plugin.PluginMeta) {
	basicRes = br
	vld = validator.New()
	dsHelper = api.NewDataSourceHelper[
		models.GithubConnection,
		models.GithubRepo,
		models.GithubScopeConfig,
	](
		br,
		p.Name(),
		[]string{"full_name"},
		func(c models.GithubConnection) models.GithubConnection {
			return c.Sanitize()
		},
		nil,
		nil,
		customPatch,
		nil,
		nil,
	)
	raProxy = api.NewDsRemoteApiProxyHelper[models.GithubConnection](dsHelper.ConnApi.ModelApiHelper)
	raScopeList = api.NewDsRemoteApiScopeListHelper[models.GithubConnection, models.GithubRepo, GithubRemotePagination](raProxy, listGithubRemoteScopes)
	raScopeSearch = api.NewDsRemoteApiScopeSearchHelper[models.GithubConnection, models.GithubRepo](raProxy, searchGithubRepos)
}

func customPatch(modified, existed *models.GithubConnection) (merged *models.GithubConnection) {
	// There are many kinds of update, we just update all fields simply.
	existedTokenStr := existed.Token
	existSecretKey := existed.SecretKey

	existed.AppId = modified.AppId
	existed.SecretKey = modified.SecretKey
	existed.InstallationID = modified.InstallationID
	existed.AuthMethod = modified.AuthMethod
	existed.Proxy = modified.Proxy
	existed.Endpoint = modified.Endpoint
	existed.RateLimitPerHour = modified.RateLimitPerHour

	// handle secret
	if existSecretKey == "" {
		if modified.SecretKey == "" {
			// doesn't input secret key, pass
		} else {
			// add secret key, store it
			existed.SecretKey = modified.SecretKey
		}
	} else {
		if modified.SecretKey == "" {
			// delete secret key
			existed.SecretKey = modified.SecretKey
		} else {
			// update secret key
			sanitizeSecretKey := existed.SanitizeSecret().SecretKey
			if sanitizeSecretKey == modified.SecretKey {
				// change nothing, restore it
				existed.SecretKey = existSecretKey
			} else {
				// has changed, replace it with the new secret key
				existed.SecretKey = modified.SecretKey
			}
		}
	}

	// handle tokens
	existedTokens := strings.Split(strings.TrimSpace(existedTokenStr), ",")
	existedTokenMap := make(map[string]string)          // {originalToken:sanitizedToken}
	existedSanitizedTokenMap := make(map[string]string) // {sanitizedToken:originalToken}
	for _, token := range existedTokens {
		existedTokenMap[token] = existed.SanitizeToken(token)
		existedSanitizedTokenMap[existed.SanitizeToken(token)] = token
	}
	//fmt.Printf("debug, existed, token string: %+v, token map: %+v \n", existedTokenStr, existedTokenMap)

	modifiedTokens := strings.Split(strings.TrimSpace(modified.Token), ",")
	modifiedTokenMap := make(map[string]string) // {originalToken:sanitizedToken}
	for _, token := range modifiedTokens {
		modifiedTokenMap[token] = existed.SanitizeToken(token)
	}
	//fmt.Printf("debug, modified, token string: %+v, token map: %+v \n", modified.Token, modifiedTokenMap)

	var mergedToken []string
	mergedTokenMap := make(map[string]struct{})
	for token, sanitizeToken := range existedTokenMap {
		// check token
		if token == sanitizeToken {
			// something wrong, just ignore sanitized tokens
			continue
		}
		if _, ok := modifiedTokenMap[token]; ok {
			// case 1: find in modified tokens, keep it
			if _, ok := mergedTokenMap[token]; !ok {
				mergedToken = append(mergedToken, token)
				mergedTokenMap[token] = struct{}{}
			}
		}
		// else case 2: not found, deleted or updated, remove it

		// check sanitized token
		if _, ok := modifiedTokenMap[sanitizeToken]; ok {
			// case 1: find in modified tokens, keep it
			if _, ok := mergedTokenMap[token]; !ok {
				mergedToken = append(mergedToken, token)
				mergedTokenMap[token] = struct{}{}
			}
		}
		// else: case 2: not found, deleted or updated, remove it
	}
	//fmt.Printf("debug, after checking exist tokens, merged tokens: %+v, token map: %+v \n", mergedToken, mergedTokenMap)
	for token, sanitizeToken := range modifiedTokenMap {
		// check token
		if _, ok := existedTokenMap[token]; ok {
			// find in db, modified but no update, ignore it
			if _, ok := mergedTokenMap[token]; !ok {
				mergedToken = append(mergedToken, token)
				mergedTokenMap[token] = struct{}{}
			}
		} else {
			// not found, a new token, we should keep it
			// cannot be a sanitized token
			if _, ok := existedSanitizedTokenMap[token]; !ok {
				if token != sanitizeToken {
					if _, ok := mergedTokenMap[token]; !ok {
						mergedToken = append(mergedToken, token)
						mergedTokenMap[token] = struct{}{}
					}
				}
			}
		}

		// token may be a sanitized token
		if v, ok := existedSanitizedTokenMap[token]; ok {
			// find in db, modify nothing, just keep it
			if _, ok := mergedTokenMap[v]; !ok {
				mergedToken = append(mergedToken, v)
				mergedTokenMap[v] = struct{}{}
			}
		} else {
			// unexpected
			fmt.Printf("unexpected, token: %+v\n will be ignored", token)
		}
		// check sanitized token
		if v, ok := existedSanitizedTokenMap[sanitizeToken]; ok {
			// find in db, modify nothing, just keep it
			if _, ok := mergedTokenMap[v]; !ok {
				mergedToken = append(mergedToken, v)
				mergedTokenMap[v] = struct{}{}
			}
		} else {
			// a new token
			// but we should check it
			if sanitizeToken != token {
				if _, ok := mergedTokenMap[token]; !ok {
					mergedToken = append(mergedToken, token)
					mergedTokenMap[token] = struct{}{}
				}
			}
		}
	}

	//fmt.Printf("debug, after checking modified tokens, merged tokens: %+v, token map: %+v \n", mergedToken, mergedTokenMap)

	existed.Token = strings.Join(mergedToken, ",")
	//fmt.Printf("debug, merged, token string: %+v, token map: %+v \n", existed.Token, mergedToken)
	return existed
}
