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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/tap"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/pagerduty/models"
)

const RAW_INCIDENTS_TABLE = "pagerduty_incidents"

var _ core.SubTaskEntryPoint = CollectIncidents

func CollectIncidents(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*PagerDutyTaskData)
	collector, err := tap.NewTapCollector(
		&tap.CollectorArgs[tap.SingerTapStream]{
			RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
				Ctx:   taskCtx,
				Table: RAW_INCIDENTS_TABLE,
				Params: models.PagerDutyParams{
					Stream:       models.IncidentStream,
					ConnectionId: data.Options.ConnectionId,
				},
			},
			TapClient:    data.Client,
			TapConfig:    data.Config,
			ConnectionId: data.Options.ConnectionId, // Seems to be an inconsequential field
			StreamName:   models.IncidentStream,
		},
	)
	if err != nil {
		return err
	}
	return collector.Execute()
}

var CollectIncidentsMeta = core.SubTaskMeta{
	Name:             "collectIncidents",
	EntryPoint:       CollectIncidents,
	EnabledByDefault: true,
	Description:      "Collect PagerDuty incidents",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}
