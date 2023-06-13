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

package azuredevops

import (
	"github.com/apache/incubator-devlake/core/models/common"
)

// These models correspond to Python models

type (
	Refdiff struct {
		TagsPattern string
		TagsLimit   int
		TagsOrder   string
	}

	AzureConnection struct {
		Name         string
		Token        string
		Organization string
	}

	AzureGitRepositoryConfig struct {
		common.ScopeConfig
		Refdiff           Refdiff
		DeploymentPattern string
		ProductionPattern string
	}

	AzureGitRepo struct {
		RawDataParams       string `json:"_raw_data_params"`
		Id                  string
		Name                string
		ConnectionId        uint64
		Url                 string
		RemoteUrl           string
		DefaultBranch       string
		ProjectId           string
		OrgId               string
		ParentRepositoryUrl string
		Provider            string
		// special field
		ScopeConfigId uint64
	}
)

type TestConfig struct {
	Org     string
	Project string
	Repos   []string
	Token   string
}
