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
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
)

const RAW_JOB_BUILD_TABLE = "bamboo_api_job_build"

var _ plugin.SubTaskEntryPoint = CollectJobBuild

type SimpleJob struct {
	JobKey   string
	PlanName string
	PlanKey  string
}

func CollectJobBuild(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_JOB_BUILD_TABLE)
	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.Select("job_key, plan_name, plan_key"),
		dal.From(models.BambooJob{}.TableName()),
		dal.Where("project_key = ? and connection_id=?", data.Options.ProjectKey, data.Options.ConnectionId),
	}
	cursor, err := db.Cursor(
		clauses...,
	)
	if err != nil {
		return err
	}
	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleJob{}))
	if err != nil {
		return err
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		Input:              iterator,
		UrlTemplate:        "result/{{ .Input.JobKey }}.json",

		Query:          QueryForResult,
		GetTotalPages:  GetTotalPagesFromResult,
		ResponseParser: GetResultsResult,
	})
	if err != nil {
		return err
	}
	return collector.Execute()
}

var CollectJobBuildMeta = plugin.SubTaskMeta{
	Name:             "CollectJobBuild",
	EntryPoint:       CollectJobBuild,
	EnabledByDefault: true,
	Description:      "Collect JobBuild data from Bamboo api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}
