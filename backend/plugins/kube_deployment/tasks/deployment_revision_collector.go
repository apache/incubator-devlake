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
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	// "fmt"
	// "net/http"
	// "net/url"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"gorm.io/datatypes"

	// "github.com/apache/incubator-devlake/helpers/pluginhelper/api"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const RAW_KUBE_DEPLOYMENT_REVISION_TABLE = "kube_deployment_revisions"

var _ plugin.SubTaskEntryPoint = CollectKubeDeploymentRevision

type KubeDeploymentAPIResult struct {
	Data []json.RawMessage `json:"data"`
}

// CollectKubeDeploymentRevision collect all DeploymentRevision that bot is in
func CollectKubeDeploymentRevision(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*KubeDeploymentTaskData)

	// // Print as JSON format
	jsonData, _ := json.Marshal(data)
	fmt.Println(string(jsonData))

	logger := taskCtx.GetLogger()
	logger.Info("Collecting kube deployment revision started...")
	db := taskCtx.GetDal()
	type RawData struct {
		ID        uint64 `gorm:"primaryKey"`
		Params    string `gorm:"type:varchar(255);index"`
		Data      []byte
		Url       string
		Input     datatypes.JSON
		CreatedAt time.Time
	}
	table := fmt.Sprintf("_raw_%s", RAW_KUBE_DEPLOYMENT_REVISION_TABLE)

	paramsString := ""
	args := api.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: KubeDeploymentApiParams{
			ConnectionId:   data.Options.ConnectionId,
			DeploymentName: data.Options.DeploymentName,
			Namespace:      data.Options.Namespace,
		},
		Table: RAW_KUBE_DEPLOYMENT_REVISION_TABLE,
	}

	paramsBytes, err := json.Marshal(args.Params)
	logger.Info("Parameters: %v", string(paramsBytes))

	if err != nil {
		return errors.Default.Wrap(err, "unable to serialize subtask parameters")
	}
	paramsString = string(paramsBytes)
	db.AutoMigrate(&RawData{}, dal.From(table))

	incremental := false
	if incremental {
		err = db.Delete(&RawData{}, dal.From(table), dal.Where("params = ?", args.Params))
		if err != nil {
			return errors.Default.Wrap(err, "error deleting data from collector")
		}
	}

	// Get the rollout history of the deployment
	rolloutHistory, err := data.KubeAPIClient.ClientSet.AppsV1().ReplicaSets(data.Options.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", data.Options.DeploymentName),
	})
	if err != nil {
		fmt.Printf("Failed to get rollout history: %v", err)
		return errors.Default.Wrap(err, "Failed to get rollout history")
	}

	count := len(rolloutHistory.Items)
	rows := make([]*RawData, count)

	logger.Info("Fetched rollout history: %v", rolloutHistory)

	for i, rs := range rolloutHistory.Items {
		revision := rs.Annotations["deployment.kubernetes.io/revision"]
		revisionNumber, _ := strconv.Atoi(revision)
		msg := map[string]interface{}{
			"id":                 rs.UID,
			"deployment_name":    data.Options.DeploymentName,
			"namespace":          rs.Namespace,
			"revision_number":    revisionNumber,
			"creation_timestamp": rs.CreationTimestamp,
		}
		var result json.RawMessage
		byteMsg, err := json.Marshal(msg)
		json.Unmarshal(byteMsg, &result)
		if err != nil {
			return errors.Default.Wrap(err, "Error marshalling message")
		}
		rows[i] = &RawData{
			Params: paramsString,
			Data:   result,
			Url:    "",
			Input:  nil,
		}
	}

	logger.Info("Saving data into : %v table", RAW_KUBE_DEPLOYMENT_REVISION_TABLE)
	err = db.Create(rows, dal.From(table))
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error inserting raw rows into %s", table))
	}

	return nil
	// // pageSize := 100
	// collector, err := api.NewApiCollector(api.ApiCollectorArgs{
	// 	RawDataSubTaskArgs: api.RawDataSubTaskArgs{
	// 		Ctx: taskCtx,
	// 		Params: KubeDeploymentApiParams{
	// 			ConnectionId:   data.Options.ConnectionId,
	// 			DeploymentName: data.Options.DeploymentName,
	// 			Namespace:      data.Options.Namespace,
	// 		},
	// 		Table: RAW_KUBE_DEPLOYMENT_REVISION_TABLE,
	// 	},
	// 	ApiClient:   data.ApiClient,
	// 	Incremental: false,
	// 	UrlTemplate: "/revisions?deployment_name={{ .Params.DeploymentName}}&namespace={{ .Params.Namespace}}",
	// 	// PageSize:    pageSize,
	// 	// GetNextPageCustomData: func(prevReqData *api.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
	// 	// 	res := MyPlugAPIResult{}
	// 	// 	err := api.UnmarshalResponse(prevPageResponse, &res)
	// 	// 	if err != nil {
	// 	// 		return nil, err
	// 	// 	}
	// 	// 	if res.ResponseMetadata.NextCursor == "" {
	// 	// 		return nil, api.ErrFinishCollect
	// 	// 	}
	// 	// 	return res.ResponseMetadata.NextCursor, nil
	// 	// },
	// 	Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
	// 		query := url.Values{}

	// 		// Print as JSON format
	// 		jsonData, _ := json.Marshal(reqData)
	// 		fmt.Println("__Query__")
	// 		fmt.Println(string(jsonData))
	// 		// query.Set("limit", strconv.Itoa(pageSize))
	// 		// if pageToken, ok := reqData.CustomData.(string); ok && pageToken != "" {
	// 		// 	query.Set("cursor", reqData.CustomData.(string))
	// 		// }
	// 		return query, nil
	// 	},
	// 	ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
	// 		var result []json.RawMessage
	// 		err := api.UnmarshalResponse(res, &result)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		return result, nil
	// 	},
	// })
	// if err != nil {
	// 	return err
	// }

	// return collector.Execute()
}

var CollectKubeDeploymentRevisionsMeta = plugin.SubTaskMeta{
	Name:             "collectKubeDeploymentRevisions",
	EntryPoint:       CollectKubeDeploymentRevision,
	EnabledByDefault: true,
	Description:      "Collect KubeDeploymentRevisions from api",
}
