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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/kube_deployment/models"
	"github.com/apache/incubator-devlake/plugins/kube_deployment/tasks"
)

func MakeDataSourcePipelinePlanV200(subtaskMetas []plugin.SubTaskMeta, connectionId uint64, bpScopes []*plugin.BlueprintScopeV200, syncPolicy *plugin.BlueprintSyncPolicy) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
	plan := make(plugin.PipelinePlan, len(bpScopes))
	fmt.Println(len(bpScopes), "__len_scopes")
	plan, err := makeDataSourcePipelinePlanV200(subtaskMetas, plan, bpScopes, connectionId, syncPolicy)
	if err != nil {
		return nil, nil, err
	}
	scopes, err := makeScopesV200(bpScopes, connectionId)
	if err != nil {
		return nil, nil, err
	}

	return plan, scopes, nil
}

func makeDataSourcePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	plan plugin.PipelinePlan,
	bpScopes []*plugin.BlueprintScopeV200,
	connectionId uint64,
	syncPolicy *plugin.BlueprintSyncPolicy,
) (plugin.PipelinePlan, errors.Error) {
	for i, bpScope := range bpScopes {
		stage := plan[i]
		if stage == nil {
			stage = plugin.PipelineStage{}
		}
		fmt.Println("bpScope_id: ", bpScope.Id)
		fmt.Println("bpScope_name: ", bpScope.Name)
		kubeDeployment := &models.KubeDeployment{}
		var err errors.Error
		err = basicRes.GetDal().First(kubeDeployment, dal.Where("connection_id = ? and id = ?", connectionId, bpScope.Id))
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find kube deployment details %s", bpScope.Id))
		}
		fmt.Println("kubeDeployment_name: ", kubeDeployment.Name)
		op := &tasks.KubeDeploymentOptions{
			ConnectionId:   connectionId,
			Namespace:      "default", // TODO: get from kubeDeployment table
			DeploymentName: kubeDeployment.Name,
		}

		var options map[string]interface{}
		err = helper.Decode(op, &options, nil)
		if err != nil {
			return nil, err
		}

		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, bpScope.Entities)
		if err != nil {
			return nil, err
		}

		fmt.Println("subtasks: ", options)
		stage = append(stage, &plugin.PipelineTask{
			Plugin:   "kube_deployment",
			Subtasks: subtasks,
			Options:  options,
		})
		plan[i] = stage
	}

	return plan, nil
}

func makeScopesV200(bpScopes []*plugin.BlueprintScopeV200, connectionId uint64) ([]plugin.Scope, errors.Error) {
	scopes := make([]plugin.Scope, 0)
	for _, bpScope := range bpScopes {
		kubeDeployment := &models.KubeDeployment{}
		// get repo from db
		err := basicRes.GetDal().First(kubeDeployment,
			dal.Where("connection_id = ? and id = ?",
				connectionId, bpScope.Id))
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find deployment %s", bpScope.Id))
		}
		if utils.StringsContains(bpScope.Entities, plugin.DOMAIN_TYPE_CICD) {
			stProject := &devops.CicdScope{
				DomainEntity: domainlayer.DomainEntity{
					Id: didgen.NewDomainIdGenerator(&models.KubeDeployment{}).Generate(kubeDeployment.ConnectionId, kubeDeployment.Name),
				},
				Name: kubeDeployment.Name,
			}
			scopes = append(scopes, stProject)
		}
	}
	return scopes, nil
}
