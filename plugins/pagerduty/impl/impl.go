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

package impl

import (
	"fmt"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/tap"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/pagerduty/api"
	"github.com/apache/incubator-devlake/plugins/pagerduty/models"
	"github.com/apache/incubator-devlake/plugins/pagerduty/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/pagerduty/tasks"
)

// make sure interface is implemented
var _ core.PluginMeta = (*PagerDuty)(nil)
var _ core.PluginInit = (*PagerDuty)(nil)
var _ core.PluginTask = (*PagerDuty)(nil)
var _ core.PluginApi = (*PagerDuty)(nil)
var _ core.PluginBlueprintV100 = (*PagerDuty)(nil)
var _ core.CloseablePluginTask = (*PagerDuty)(nil)

type PagerDuty struct{}

func (plugin PagerDuty) Description() string {
	return "collect some PagerDuty data"
}

func (plugin PagerDuty) Init(basicRes core.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (plugin PagerDuty) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectIncidentsMeta,
		tasks.ExtractIncidentsMeta,
		tasks.ConvertIncidentsMeta,
	}
}

func (plugin PagerDuty) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	connection := &models.PagerDutyConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Pagerduty connection by the given connection ID")
	}
	startDate, err := parseTime("start_date", options)
	if err != nil {
		return nil, err
	}
	config := &models.PagerDutyConfig{
		Token:     connection.Token,
		Email:     "", // ignore, works without it too
		StartDate: startDate,
	}
	tapClient, err := tap.NewSingerTap(&tap.SingerTapConfig{
		TapExecutable:        models.TapExecutable,
		StreamPropertiesFile: models.StreamPropertiesFile,
	})
	if err != nil {
		return nil, err
	}
	return &tasks.PagerDutyTaskData{
		Options: op,
		Config:  config,
		Client:  tapClient,
	}, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (plugin PagerDuty) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/pagerduty"
}

func (plugin PagerDuty) MigrationScripts() []core.MigrationScript {
	return migrationscripts.All()
}

func (plugin PagerDuty) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":    api.GetConnection,
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
		},
	}
}

func (plugin PagerDuty) MakePipelinePlan(connectionId uint64, scope []*core.BlueprintScopeV100) (core.PipelinePlan, errors.Error) {
	return api.MakePipelinePlan(plugin.SubTaskMetas(), connectionId, scope)
}

func (plugin PagerDuty) Close(taskCtx core.TaskContext) errors.Error {
	_, ok := taskCtx.GetData().(*tasks.PagerDutyTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	return nil
}

func parseTime(key string, opts map[string]any) (time.Time, errors.Error) {
	var date time.Time
	dateRaw, ok := opts[key]
	if !ok {
		return date, errors.BadInput.New("time input not provided")
	}
	date, err := time.Parse("2006-01-02T15:04:05Z", dateRaw.(string))
	if err != nil {
		return date, errors.BadInput.Wrap(err, "bad type input provided")
	}
	return date, nil
}
