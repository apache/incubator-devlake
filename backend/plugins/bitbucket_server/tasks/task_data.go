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
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
)

type BitbucketServerOptions struct {
	ConnectionId                       uint64   `json:"connectionId" mapstructure:"connectionId,omitempty"`
	Tasks                              []string `json:"tasks,omitempty" mapstructure:",omitempty"`
	FullName                           string   `json:"fullName" mapstructure:"fullName"`
	ScopeConfigId                      uint64   `json:"scopeConfigId" mapstructure:"scopeConfigId,omitempty"`
	*models.BitbucketServerScopeConfig `mapstructure:"scopeConfig,omitempty" json:"scopeConfig"`
}

type BitbucketServerTaskData struct {
	Options       *BitbucketServerOptions
	ApiClient     *api.ApiAsyncClient
	RegexEnricher *api.RegexEnricher
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*BitbucketServerOptions, errors.Error) {
	op, err := DecodeTaskOptions(options)
	if err != nil {
		return nil, err
	}
	err = ValidateTaskOptions(op)
	if err != nil {
		return nil, err
	}
	return op, nil
}

func DecodeTaskOptions(options map[string]interface{}) (*BitbucketServerOptions, errors.Error) {
	var op BitbucketServerOptions
	err := api.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func EncodeTaskOptions(op *BitbucketServerOptions) (map[string]interface{}, errors.Error) {
	var result map[string]interface{}
	err := api.Decode(op, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ValidateTaskOptions(op *BitbucketServerOptions) errors.Error {
	if op.FullName == "" {
		return errors.BadInput.New("no enough info for Bitbucket execution")
	}
	// find the needed Bitbucket now
	if op.ConnectionId == 0 {
		return errors.BadInput.New("connectionId is invalid")
	}
	return nil
}

func ignoreHTTPStatus404(res *http.Response) errors.Error {
	if res.StatusCode == http.StatusUnauthorized {
		return errors.Unauthorized.New("authentication failed, please check your AccessToken")
	}
	if res.StatusCode == http.StatusNotFound {
		return api.ErrIgnoreAndContinue
	}
	return nil
}
