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
	"testing"
	"time"

	"github.com/apache/incubator-devlake/plugins/argocd/models"
	"github.com/stretchr/testify/assert"
)

func TestDescribeApplicationScopeBuildsSummary(t *testing.T) {
	desc := describeApplicationScope(&models.ArgocdApplication{
		Project:       "observability",
		Namespace:     "argocd",
		DestServer:    "https://k8s.example.com",
		DestNamespace: "prod",
	})
	assert.Equal(t, "Project: observability | Namespace: argocd | Destination: https://k8s.example.com/prod", desc)
}

func TestBuildCicdScopeFromApplication(t *testing.T) {
	created := time.Now()
	scope := buildCicdScopeFromApplication(&models.ArgocdApplication{
		Name:        "test-app",
		RepoURL:     "https://git.example.com/app.git",
		CreatedDate: &created,
	}, "argocd:app:1")

	assert.Equal(t, "argocd:app:1", scope.Id)
	assert.Equal(t, "test-app", scope.Name)
	assert.Equal(t, "https://git.example.com/app.git", scope.Url)
	if assert.NotNil(t, scope.CreatedDate) {
		assert.Equal(t, created.UTC(), scope.CreatedDate.UTC())
	}
}
