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
	goerror "errors"
	"github.com/apache/incubator-devlake/plugins/helper"
	"gorm.io/gorm"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/apache/incubator-devlake/plugins/jenkins/tasks"
	"github.com/apache/incubator-devlake/utils"
	"github.com/mitchellh/mapstructure"
)

func MakeDataSourcePipelinePlanV200(subtaskMetas []core.SubTaskMeta, connectionId uint64, bpScopes []*core.BlueprintScopeV200) (core.PipelinePlan, []core.Scope, errors.Error) {
	db := BasicRes.GetDal()
	// get the connection info for url
	connection := &models.JenkinsConnection{}
	err := connectionHelper.FirstById(connection, connectionId)
	if err != nil {
		return nil, nil, err
	}

	plan := make(core.PipelinePlan, 0, len(bpScopes))
	scopes := make([]core.Scope, 0, len(bpScopes))
	for i, bpScope := range bpScopes {
		var jenkinsJob *models.JenkinsJob
		// get repo from db
		err = db.First(jenkinsJob, dal.Where(`id = ?`, bpScope.Id))
		if err != nil {
			return nil, nil, err
		}
		var transformationRule *models.JenkinsTransformationRule
		// get transformation rules from db
		err = db.First(transformationRule, dal.Where(`id = ?`, jenkinsJob.TransformationRuleId))
		if err != nil && goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		}
		var scope []core.Scope
		// make pipeline for each bpScope
		plan[i], scope, err = makeDataSourcePipelinePlanV200(subtaskMetas, bpScope, jenkinsJob, transformationRule)
		if err != nil {
			return nil, nil, err
		}
		if len(scope) > 0 {
			scopes = append(scopes, scope...)
		}

	}

	return plan, scopes, nil
}

func makeDataSourcePipelinePlanV200(
	subtaskMetas []core.SubTaskMeta,
	bpScope *core.BlueprintScopeV200,
	jenkinsJob *models.JenkinsJob,
	transformationRule *models.JenkinsTransformationRule,
) (core.PipelineStage, []core.Scope, errors.Error) {
	var err errors.Error
	var stage core.PipelineStage
	scopes := make([]core.Scope, 0)

	// construct task options for jenkins
	var options map[string]interface{}
	err = errors.Convert(mapstructure.Decode(jenkinsJob, &options))
	if err != nil {
		return nil, nil, err
	}
	// make sure task options is valid
	_, err = tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, nil, err
	}

	var transformationRuleMap map[string]interface{}
	err = errors.Convert(mapstructure.Decode(transformationRule, &transformationRuleMap))
	if err != nil {
		return nil, nil, err
	}
	options["transformationRules"] = transformationRuleMap
	subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, bpScope.Entities)
	if err != nil {
		return nil, nil, err
	}
	stage = append(stage, &core.PipelineTask{
		Plugin:   "jenkins",
		Subtasks: subtasks,
		Options:  options,
	})

	// add cicd_scope to scopes
	if utils.StringsContains(bpScope.Entities, core.DOMAIN_TYPE_CICD) {
		scopeCICD := &devops.CicdScope{
			DomainEntity: domainlayer.DomainEntity{
				Id: didgen.NewDomainIdGenerator(&models.JenkinsJob{}).Generate(jenkinsJob.ConnectionId, jenkinsJob.FullName),
			},
			Name: jenkinsJob.FullName,
		}
		scopes = append(scopes, scopeCICD)
	}

	return stage, scopes, nil
}
