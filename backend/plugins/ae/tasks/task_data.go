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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type AeOptions struct {
	ConnectionId uint64 `json:"connectionId"`
	ProjectId    int
}

type AeTaskData struct {
	Options   *AeOptions
	ApiClient *helper.ApiAsyncClient
}
type AeApiParams struct {
	ConnectionId uint64 `json:"connectionId"`
	ProjectId    int
}
