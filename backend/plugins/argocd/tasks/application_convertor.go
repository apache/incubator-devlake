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
	"fmt"
	"reflect"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/argocd/models"
)

var ConvertApplicationsMeta = plugin.SubTaskMeta{
	Name:             "convertApplications",
	EntryPoint:       ConvertApplications,
	EnabledByDefault: true,
	Description:      "Convert ArgoCD applications into CICD scopes",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	DependencyTables: []string{models.ArgocdApplication{}.TableName()},
	ProductTables:    []string{devops.CicdScope{}.TableName()},
}

func ConvertApplications(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ArgocdTaskData)
	db := taskCtx.GetDal()

	cursor, err := db.Cursor(
		dal.From(&models.ArgocdApplication{}),
		dal.Where("connection_id = ?", data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	scopeIdGen := didgen.NewDomainIdGenerator(&models.ArgocdApplication{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.ArgocdApplication{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:   taskCtx,
			Table: RAW_APPLICATION_TABLE,
			Params: models.ArgocdApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.ApplicationName,
			},
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			application := inputRow.(*models.ArgocdApplication)
			scopeId := scopeIdGen.Generate(application.ConnectionId, application.Name)
			scope := buildCicdScopeFromApplication(application, scopeId)
			return []interface{}{scope}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

func buildCicdScopeFromApplication(app *models.ArgocdApplication, scopeId string) *devops.CicdScope {
	scope := &devops.CicdScope{
		DomainEntity: domainlayer.NewDomainEntity(scopeId),
		Name:         app.Name,
		Description:  describeApplicationScope(app),
		Url:          firstNonEmpty(app.RepoURL, app.DestServer),
		CreatedDate:  app.CreatedDate,
	}
	return scope
}

func describeApplicationScope(app *models.ArgocdApplication) string {
	parts := make([]string, 0, 3)
	if app.Project != "" {
		parts = append(parts, fmt.Sprintf("Project: %s", app.Project))
	}
	if app.Namespace != "" {
		parts = append(parts, fmt.Sprintf("Namespace: %s", app.Namespace))
	}
	if app.DestNamespace != "" || app.DestServer != "" {
		dest := strings.TrimSpace(fmt.Sprintf("%s/%s", app.DestServer, app.DestNamespace))
		dest = strings.Trim(dest, "/")
		if dest != "" {
			parts = append(parts, fmt.Sprintf("Destination: %s", dest))
		}
	}
	return strings.Join(parts, " | ")
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
