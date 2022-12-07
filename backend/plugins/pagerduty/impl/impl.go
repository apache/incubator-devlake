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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/tap"
	"github.com/apache/incubator-devlake/plugins/pagerduty/api"
	"github.com/apache/incubator-devlake/plugins/pagerduty/models"
	"github.com/apache/incubator-devlake/plugins/pagerduty/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/pagerduty/tasks"
	"time"
)

// make sure interface is implemented
var _ plugin.PluginMeta = (*PagerDuty)(nil)
var _ plugin.PluginInit = (*PagerDuty)(nil)
var _ plugin.PluginTask = (*PagerDuty)(nil)
var _ plugin.PluginApi = (*PagerDuty)(nil)
var _ plugin.PluginBlueprintV100 = (*PagerDuty)(nil)
var _ plugin.CloseablePluginTask = (*PagerDuty)(nil)

type PagerDuty struct{}

func (p PagerDuty) Description() string {
	return "collect some PagerDuty data"
}

func (p PagerDuty) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (p PagerDuty) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectIncidentsMeta,
		tasks.ExtractIncidentsMeta,
		tasks.ConvertIncidentsMeta,
	}
}

func (p PagerDuty) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
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
func (p PagerDuty) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/pagerduty"
}

func (p PagerDuty) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p PagerDuty) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
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

func (p PagerDuty) MakePipelinePlan(connectionId uint64, scope []*plugin.BlueprintScopeV100) (plugin.PipelinePlan, errors.Error) {
	return api.MakePipelinePlan(p.SubTaskMetas(), connectionId, scope)
}

func (p PagerDuty) Close(taskCtx plugin.TaskContext) errors.Error {
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
