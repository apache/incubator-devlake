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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var PullRequestToIncidentsMeta = plugin.SubTaskMeta{
	Name:             "ConvertPullRequestToIncidents",
	EntryPoint:       ConvertPullRequestToIncidents,
	EnabledByDefault: true,
	Description:      "Connect pull_request to incident",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET, plugin.DOMAIN_TYPE_CICD, plugin.DOMAIN_TYPE_CROSS},
}

func ConvertPullRequestToIncidents(taskCtx plugin.SubTaskContext) errors.Error {
	// TODO
	return nil
}

// lint:ignore U1000
//func generateIncidentAssigneeFromPullRequest(db dal.Dal, logger log.Logger, pullRequest *code.PullRequest) ([]*ticket.IncidentAssignee, error) {
//	if pullRequest == nil {
//		return nil, goerrors.New("pull request is nil")
//	}
//	var pullRequestAssignees []*code.PullRequestAssignee
//	if err := db.All(&pullRequestAssignees, dal.Where("pull_request_id = ?", pullRequest.Id)); err != nil {
//		logger.Error(err, "Failed to fetch pull request assignees")
//		return nil, err
//	}
//
//	var incidentAssignees []*ticket.IncidentAssignee
//	for _, pullRequestAssignee := range pullRequestAssignees {
//		incidentAssignees = append(incidentAssignees, &ticket.IncidentAssignee{
//			IncidentId:   pullRequestAssignee.PullRequestId,
//			AssigneeId:   pullRequestAssignee.AssigneeId,
//			AssigneeName: pullRequestAssignee.Name,
//			NoPKModel:    common.NewNoPKModel(),
//		})
//	}
//	return incidentAssignees, nil
//}
