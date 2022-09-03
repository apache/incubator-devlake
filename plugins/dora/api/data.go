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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"net/http"
)

const RAW_DEPLOYMENTS_TABLE = `dora_deplyments`

//TODO Please modify the folowing code to adapt to your plugin
/*
POST /plugins/dora/deployments
{
	TODO
}
*/
func PostDeployments(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	// TODO
	return &core.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}

const RAW_ISSUES_TABLE = `dora_issues`

//TODO Please modify the folowing code to adapt to your plugin
/*
POST /plugins/dora/issues
{
	TODO
}
*/
func PostIssues(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	// TODO
	return &core.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}

//TODO Please modify the folowing code to adapt to your plugin
/*
POST /plugins/dora/issues/:id/close
{
	TODO
}
*/
func CloseIssues(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	// TODO
	return &core.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}
