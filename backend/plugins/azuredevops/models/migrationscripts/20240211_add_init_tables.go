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

package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/azuredevops/models/migrationscripts/archived"
)

type addInitTables struct {
}

func (u *addInitTables) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.DropTables(
		"_tool_azuredevops_builds",
		"_tool_azuredevops_gitpullrequestcommits",
		"_tool_azuredevops_gitpullrequests",
		"_tool_azuredevops_gitrepositories",
		"_tool_azuredevops_gitrepositoryconfigs",
		"_tool_azuredevops_jobs",

		"_raw_azuredevops_builds",
		"_raw_azuredevops_gitpullrequestcommits",
		"_raw_azuredevops_gitpullrequests",
		"_raw_azuredevops_jobs",
	)
	if err != nil {
		return err
	}

	err = migrationhelper.AutoMigrateTables(
		basicRes,
		&archived.AzuredevopsBuild{},
		&archived.AzuredevopsTimelineRecord{},
		&archived.AzuredevopsCommit{},
		&archived.AzuredevopsConnection{},
		&archived.AzuredevopsPrCommit{},
		&archived.AzuredevopsPrLabel{},
		&archived.AzuredevopsPullRequest{},
		&archived.AzuredevopsRepo{},
		&archived.AzuredevopsRepoCommit{},
		&archived.AzuredevopsScopeConfig{},
		&archived.AzuredevopsUser{},
	)
	return err
}

func (*addInitTables) Version() uint64 {
	return 20240211000001
}

func (*addInitTables) Name() string {
	return "Initializing Azure DevOps schemas"
}
