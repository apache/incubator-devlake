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

	"github.com/apache/devlake/core/dal"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/models/domainlayer"
	"github.com/apache/devlake/core/models/domainlayer/didgen"
	"github.com/apache/devlake/core/models/domainlayer/ticket"
	"github.com/apache/devlake/core/plugin"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/incidentio/models"
)

var ConvertIncidentTypesMeta = plugin.SubTaskMeta{
	Name:             "convertIncidentTypes",
	EntryPoint:       ConvertIncidentTypes,
	EnabledByDefault: true,
	Description:      "Convert incident.io incident types",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertIncidentTypes(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*IncidentioTaskData)
	rawDataSubTaskArgs := &helper.RawDataSubTaskArgs{
		Ctx:     taskCtx,
		Options: data.Options,
		Table:   RAW_INCIDENT_TYPES_TABLE,
	}
	clauses := []dal.Clause{
		dal.Select("incident_types.*"),
		dal.From("_tool_incidentio_incident_types incident_types"),
		dal.Where("id = ? and connection_id = ?", data.Options.IncidentTypeId, data.Options.ConnectionId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.IncidentType{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			incidentType := inputRow.(*models.IncidentType)
			domainBoard := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: didgen.NewDomainIdGenerator(incidentType).Generate(incidentType.ConnectionId, incidentType.Id),
				},
				Name: incidentType.Name,
			}
			return []interface{}{
				domainBoard,
			}, nil
		},
	})
	if err != nil {
		return err
	}
	return converter.Execute()
}
