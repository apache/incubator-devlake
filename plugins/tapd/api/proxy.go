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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

func Proxy(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connectionId := input.Params["connectionId"]
	if connectionId == "" {
		return nil, fmt.Errorf("missing connectionId")
	}
	tapdConnectionId, err := strconv.ParseUint(connectionId, 10, 64)
	if err != nil {
		return nil, err
	}
	tapdConnection := &models.TapdConnection{}
	err = db.First(tapdConnection, tapdConnectionId).Error
	if err != nil {
		return nil, err
	}
	encKey := cfg.GetString(core.EncodeKeyEnvStr)
	basicAuth, err := core.Decrypt(encKey, tapdConnection.BasicAuthEncoded)
	if err != nil {
		return nil, err
	}
	apiClient, err := helper.NewApiClient(
		tapdConnection.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", basicAuth),
		},
		30*time.Second,
		"",
		nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := apiClient.Get(input.Params["path"], input.Query, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// verify response body is json
	var tmp interface{}
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Status: resp.StatusCode, Body: json.RawMessage(body)}, nil
}
