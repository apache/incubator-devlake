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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/q_dev/models"
	"github.com/apache/incubator-devlake/plugins/q_dev/tasks"

	"net/http"
)

// TestConnection 测试连接
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// 解析连接参数
	var connection models.QDevConnection
	err := api.Decode(input.Body, &connection, vld)
	if err != nil {
		return nil, err
	}

	// 测试S3连接
	_, err = tasks.NewQDevS3Client(nil, &connection)
	if err != nil {
		return nil, err
	}

	// 连接成功
	return &plugin.ApiResourceOutput{Status: http.StatusOK}, nil
}

// TestExistingConnection 测试现有连接
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.QDevConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "find connection from db")
	}
	if err := api.DecodeMapStruct(input.Body, connection, false); err != nil {
		return nil, err
	}
	// 测试连接
	_, err = tasks.NewQDevS3Client(nil, connection)
	if err != nil {
		return nil, err
	}

	// 连接成功
	return &plugin.ApiResourceOutput{Status: http.StatusOK}, nil
}
