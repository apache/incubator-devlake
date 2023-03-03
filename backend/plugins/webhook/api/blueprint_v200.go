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
	"fmt"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/webhook/models"
)

func MakeDataSourcePipelinePlanV200(connectionId uint64) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
	// get the connection info for url
	connection := &models.WebhookConnection{}
	err := connectionHelper.FirstById(connection, connectionId)
	if err != nil {
		return nil, nil, errors.Default.Wrap(err, `cannot find webhook connection`)
	}

	scopes := make([]plugin.Scope, 0)
	// add cicd_scope to scopes
	scopes = append(scopes, &devops.CicdScope{
		DomainEntity: domainlayer.DomainEntity{
			Id: fmt.Sprintf("%s:%d", "webhook", connection.ID),
		},
		Name: connection.Name,
	})

	return nil, scopes, nil
}
