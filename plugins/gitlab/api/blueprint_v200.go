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
	"github.com/go-playground/validator/v10"

	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

func MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*core.BlueprintScopeV200) (pp core.PipelinePlan, sc []core.Scope, err errors.Error) {
	pp = make(core.PipelinePlan, 0, 1)
	sc = make([]core.Scope, 0, 3*len(scopes))
	err = nil

	connectionHelper := helper.NewConnectionHelper(BasicRes, validator.New())

	// get the connection info for url
	connection := &models.GitlabConnection{}
	err = connectionHelper.FirstById(connection, connectionId)
	if err != nil {
		return nil, nil, err
	}

	ps := make(core.PipelineStage, 0, len(scopes))
	for _, scope := range scopes {
		var board ticket.Board
		var repo code.Repo

		id := didgen.NewDomainIdGenerator(&models.GitlabProject{}).Generate(connectionId, scope.Id)

		repo.Id = id
		repo.Name = scope.Name

		board.Id = id
		board.Name = scope.Name

		sc = append(sc, &repo)
		sc = append(sc, &board)

		ps = append(ps, &core.PipelineTask{
			Plugin: "gitlab",
			Options: map[string]interface{}{
				"name": scope.Name,
			},
		})

		ps = append(ps, &core.PipelineTask{
			Plugin: "gitextractor",
			Options: map[string]interface{}{
				"url": connection.Endpoint + scope.Name,
			},
		})
	}

	pp = append(pp, ps)

	return pp, sc, nil
}
