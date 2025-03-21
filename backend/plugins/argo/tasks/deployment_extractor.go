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
	"log"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/argo/models"
)

var _ plugin.SubTaskEntryPoint = ExtractDeployments

var ExtractDeploymentsMeta = plugin.SubTaskMeta{
	Name:             "extract_deployments",
	EntryPoint:       ExtractDeployments,
	EnabledByDefault: true,
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	DependencyTables: []string{},
	ProductTables:    []string{RAW_DEPLOYMENT_TABLE},
}

func ExtractDeployments(taskCtx plugin.SubTaskContext) errors.Error {
	log.Println("[ARGO] Iniciando plugin de extract.")

	data := taskCtx.GetData().(*ArgoTaskData)

	extractor, err := api.NewStatefulApiExtractor(&api.StatefulApiExtractorArgs[DeploymentResp]{
		SubtaskCommonArgs: &api.SubtaskCommonArgs{
			SubTaskContext: taskCtx,
			Table:          RAW_DEPLOYMENT_TABLE,
			Params: models.ArgoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
		},
		Extract: func(deploymentResp *DeploymentResp, row *api.RawData) ([]interface{}, errors.Error) {
			deployment := deploymentResp.toDeployment(data.Options.ConnectionId, data.Options.ProjectId)

			log.Println("[ARGO] Extraindo dados...")
			log.Println(deployment)
			log.Println("[ARGO] Extração finalizada.")

			return []interface{}{
				deployment,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	log.Println("[ARGO] Plugin extractor executado com sucesso!")
	return extractor.Execute()
}

func (d DeploymentResp) toDeployment(connectionId uint64, projectId string) *models.Deployment {
	return &models.Deployment{
		NoPKModel:       common.NewNoPKModel(),
		ConnectionId:    connectionId,
		ProjectId:       projectId,
		Name:            d.Metadata.Name,
		GeneratedName:   d.Metadata.GenerateName,
		Namespace:       d.Metadata.Namespace,
		UID:             d.Metadata.UID,
		ResourceVersion: d.Metadata.ResourceVersion,
		Result:          d.Status.Phase,
		Status:          GetStatus(d.Status.Phase),
		CreationDate:    d.Metadata.CreationTimestamp,
		StartedAt:       d.Status.StartedAt,
		FinishedAt:      d.Status.FinishedAt,
		RefName:         d.Spec.Arguments.Parameters[0].Value,
		CommitSha:       d.Spec.Arguments.Parameters[1].Value,
		DurationSec:     GetDuration(d.Status.StartedAt, d.Status.FinishedAt),
		RepoUrl:         GetUrl(d.Status.Nodes),
		Environment:     "production",
	}
}

type DeploymentResp struct {
	Metadata Metadata `json:"metadata"`
	Spec     Spec     `json:"spec"`
	Status   Status   `json:"status"`
}

type Metadata struct {
	Name              string            `json:"name"`
	GenerateName      string            `json:"generateName"`
	Namespace         string            `json:"namespace"`
	UID               string            `json:"uid"`
	ResourceVersion   string            `json:"resourceVersion"`
	Generation        int               `json:"generation"`
	CreationTimestamp string            `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
	ManagedFields     []ManagedField    `json:"managedFields"`
}

type ManagedField struct {
	Manager    string   `json:"manager"`
	Operation  string   `json:"operation"`
	APIVersion string   `json:"apiVersion"`
	Time       string   `json:"time"`
	FieldsType string   `json:"fieldsType"`
	FieldsV1   FieldsV1 `json:"fieldsV1"`
}

type FieldsV1 struct {
	Metadata MetadataFields `json:"f:metadata"`
	Spec     interface{}    `json:"f:spec"`
	Status   interface{}    `json:"f:status"`
}

type MetadataFields struct {
	GenerateName interface{} `json:"f:generateName"`
	Labels       Labels      `json:"f:labels"`
	Annotations  Annotations `json:"f:annotations"`
}

type Labels struct {
	EventsActionTimestamp interface{} `json:"f:events.argoproj.io/action-timestamp"`
	EventsSensor          interface{} `json:"f:events.argoproj.io/sensor"`
	EventsTrigger         interface{} `json:"f:events.argoproj.io/trigger"`
	WorkflowsCreator      interface{} `json:"f:workflows.argoproj.io/creator"`
	WorkflowsCompleted    interface{} `json:"f:workflows.argoproj.io/completed"`
	WorkflowsPhase        interface{} `json:"f:workflows.argoproj.io/phase"`
}

type Annotations struct {
	PodNameFormat interface{} `json:"f:workflows.argoproj.io/pod-name-format"`
}

type Spec struct {
	Entrypoint          string              `json:"entrypoint"`
	Arguments           Arguments           `json:"arguments"`
	WorkflowTemplateRef WorkflowTemplateRef `json:"workflowTemplateRef"`
}

type Arguments struct {
	Parameters []Parameter `json:"parameters"`
}

type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type WorkflowTemplateRef struct {
	Name string `json:"name"`
}

type Status struct {
	Phase      string          `json:"phase"`
	StartedAt  string          `json:"startedAt"`
	FinishedAt string          `json:"finishedAt"`
	Progress   string          `json:"progress"`
	Nodes      map[string]Node `json:"nodes"`
}

type Node struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	DisplayName   string      `json:"displayName"`
	Type          string      `json:"type"`
	TemplateRef   TemplateRef `json:"templateRef"`
	TemplateScope string      `json:"templateScope"`
	Phase         string      `json:"phase"`
	BoundaryID    string      `json:"boundaryID"`
	StartedAt     string      `json:"startedAt"`
	FinishedAt    string      `json:"finishedAt"`
	Progress      string      `json:"progress"`
	Inputs        Inputs      `json:"inputs"`
	Children      []string    `json:"children"`
	OutboundNodes []string    `json:"outboundNodes"`
}

type TemplateRef struct {
	Name     string `json:"name"`
	Template string `json:"template"`
}

type Inputs struct {
	Parameters []Parameter `json:"parameters"`
}

func GetDuration(started string, finished string) int64 {
	if started == "" || finished == "" {
		startedAt, err := time.Parse(time.RFC3339, started)
		if err != nil {
			log.Println("Error parsing StartedAt:", err)
		}
		finishedAt, err := time.Parse(time.RFC3339, finished)
		if err != nil {
			log.Println("Error parsing FinishedAt:", err)
		}

		return int64(finishedAt.Sub(startedAt).Seconds())
	}

	return 0
}

func GetUrl(node map[string]Node) string {
	for _, node := range node {
		if node.DisplayName == "exit-workflow" {
			for _, param := range node.Inputs.Parameters {
				if param.Name == "repository" {
					return param.Value
				}
			}
		}
	}

	return ""
}

func GetStatus(s string) string {
	if s == "succedded" {
		return "done"
	}

	return ""
}
