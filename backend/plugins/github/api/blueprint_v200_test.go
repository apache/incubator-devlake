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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_makeDataSourcePipelinePlanV200_with_empty_scope_config(t *testing.T) {

	mockMeta := mockplugin.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/github")
	mockMeta.On("Name").Return("github").Maybe()
	err := plugin.RegisterPlugin("github", mockMeta)
	assert.Nil(t, err)

	type args struct {
		subtaskMetas []plugin.SubTaskMeta
		scopeDetails []*srvhelper.ScopeDetail[models.GithubRepo, models.GithubScopeConfig]
		connection   *models.GithubConnection
	}
	tests := []struct {
		name                      string
		args                      args
		makeSureGitExtractorExist bool
		wantError                 bool
	}{
		{
			name: "without-empty-scope-config",
			args: args{
				subtaskMetas: []plugin.SubTaskMeta{},
				scopeDetails: []*srvhelper.ScopeDetail[models.GithubRepo, models.GithubScopeConfig]{
					&srvhelper.ScopeDetail[models.GithubRepo, models.GithubScopeConfig]{
						Scope:       models.GithubRepo{},
						ScopeConfig: &models.GithubScopeConfig{},
					},
				},
				connection: &models.GithubConnection{
					BaseConnection: api.BaseConnection{
						Model: common.Model{
							ID: 1,
						},
					},
					GithubConn:    models.GithubConn{},
					EnableGraphql: false,
				},
			},
			makeSureGitExtractorExist: true,
			wantError:                 false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := makeDataSourcePipelinePlanV200(tt.args.subtaskMetas, tt.args.scopeDetails, tt.args.connection)
			if tt.wantError {
				assert.Equalf(t, nil, got1, "makeDataSourcePipelinePlanV200 want error(%v, %v, %v)", tt.args.subtaskMetas, tt.args.scopeDetails, tt.args.connection)
			}
			if tt.makeSureGitExtractorExist {
				var existGitExtractor bool
				for _, g := range got {
					for _, v := range g {
						if v.Plugin == "gitextractor" {
							existGitExtractor = true
						}
					}
				}
				if !existGitExtractor {
					t.Fatal("gitextractor not found")
				}
			}
		})
	}
}
