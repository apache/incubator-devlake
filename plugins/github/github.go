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

package main // must be main for plugin entry point

import (
	"fmt"
	"github.com/apache/incubator-devlake/migration"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/github/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
	"github.com/apache/incubator-devlake/runner"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var _ core.PluginMeta = (*Github)(nil)
var _ core.PluginInit = (*Github)(nil)
var _ core.PluginTask = (*Github)(nil)
var _ core.PluginApi = (*Github)(nil)
var _ core.Migratable = (*Github)(nil)

type Github struct{}

func (plugin Github) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	api.Init(config, logger, db)
	return nil
}

func (plugin Github) Description() string {
	return "To collect and enrich data from GitHub"
}

func (plugin Github) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectApiRepoMeta,
		tasks.ExtractApiRepoMeta,
		tasks.CollectApiIssuesMeta,
		tasks.ExtractApiIssuesMeta,
		tasks.CollectApiPullRequestsMeta,
		tasks.ExtractApiPullRequestsMeta,
		tasks.CollectApiCommentsMeta,
		tasks.ExtractApiCommentsMeta,
		tasks.CollectApiEventsMeta,
		tasks.ExtractApiEventsMeta,
		tasks.CollectApiPullRequestCommitsMeta,
		tasks.ExtractApiPullRequestCommitsMeta,
		tasks.CollectApiPullRequestReviewsMeta,
		tasks.ExtractApiPullRequestReviewsMeta,
		tasks.CollectApiCommitsMeta,
		tasks.ExtractApiCommitsMeta,
		tasks.CollectApiCommitStatsMeta,
		tasks.ExtractApiCommitStatsMeta,
		tasks.EnrichPullRequestIssuesMeta,
		tasks.ConvertRepoMeta,
		tasks.ConvertIssuesMeta,
		tasks.ConvertCommitsMeta,
		tasks.ConvertIssueLabelsMeta,
		tasks.ConvertPullRequestCommitsMeta,
		tasks.ConvertPullRequestsMeta,
		tasks.ConvertPullRequestLabelsMeta,
		tasks.ConvertPullRequestIssuesMeta,
		tasks.ConvertUsersMeta,
		tasks.ConvertIssueCommentsMeta,
		tasks.ConvertPullRequestCommentsMeta,
	}
}

func (plugin Github) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.GithubOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	if op.Owner == "" {
		return nil, fmt.Errorf("owner is required for GitHub execution")
	}
	if op.Repo == "" {
		return nil, fmt.Errorf("repo is required for GitHub execution")
	}
	if op.PrType == "" {
		op.PrType = "type/(.*)$"
	}
	if op.PrComponent == "" {
		op.PrComponent = "component/(.*)$"
	}
	if op.PrBodyClosePattern == "" {
		op.PrBodyClosePattern = "(?mi)(fix|close|resolve|fixes|closes|resolves|fixed|closed|resolved)[\\s]*.*(((and )?(#|https:\\/\\/github.com\\/%s\\/%s\\/issues\\/)\\d+[ ]*)+)"
	}
	if op.IssueSeverity == "" {
		op.IssueSeverity = "severity/(.*)$"
	}
	if op.IssuePriority == "" {
		op.IssuePriority = "^(highest|high|medium|low)$"
	}
	if op.IssueComponent == "" {
		op.IssueComponent = "component/(.*)$"
	}
	if op.IssueTypeBug == "" {
		op.IssueTypeBug = "^(bug|failure|error)$"
	}
	if op.IssueTypeIncident == "" {
		op.IssueTypeIncident = ""
	}
	if op.IssueTypeRequirement == "" {
		op.IssueTypeRequirement = "^(feat|feature|proposal|requirement)$"
	}

	// find the needed GitHub now
	if op.ConnectionId == 0 {
		return nil, fmt.Errorf("connectionId is invalid")
	}
	connection := &models.GithubConnection{}
	err = taskCtx.GetDal().First(connection, dal.Where("id = ?", op.ConnectionId))
	if err != nil {
		return err, nil
	}

	apiClient, err := tasks.CreateApiClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}

	return &tasks.GithubTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (plugin Github) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/github"
}

func (plugin Github) MigrationScripts() []migration.Script {
	return []migration.Script{
		new(migrationscripts.InitSchemas), new(migrationscripts.UpdateSchemas20220509),
		new(migrationscripts.UpdateSchemas20220524), new(migrationscripts.UpdateSchemas20220608),
	}
}

func (plugin Github) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":    api.GetConnection,
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
		},
	}
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Github //nolint

// standalone mode for debugging
func main() {
	cmd := &cobra.Command{Use: "github"}
	connectionId := cmd.Flags().Uint64P("connection", "c", 0, "github connection id")
	owner := cmd.Flags().StringP("owner", "o", "", "github owner")
	repo := cmd.Flags().StringP("repo", "r", "", "github repo")
	_ = cmd.MarkFlagRequired("connection")
	_ = cmd.MarkFlagRequired("owner")
	_ = cmd.MarkFlagRequired("repo")

	prType := cmd.Flags().String("prType", "type/(.*)$", "pr type")
	prComponent := cmd.Flags().String("prComponent", "component/(.*)$", "pr component")
	prBodyClosePattern := cmd.Flags().String("prBodyClosePattern", "(?mi)(fix|close|resolve|fixes|closes|resolves|fixed|closed|resolved)[\\s]*.*(((and )?(#|https:\\/\\/github.com\\/%s\\/%s\\/issues\\/)\\d+[ ]*)+)", "pr body close pattern")
	issueSeverity := cmd.Flags().String("issueSeverity", "severity/(.*)$", "issue severity")
	issuePriority := cmd.Flags().String("issuePriority", "^(highest|high|medium|low)$", "issue priority")
	issueComponent := cmd.Flags().String("issueComponent", "component/(.*)$", "issue component")
	issueTypeBug := cmd.Flags().String("issueTypeBug", "^(bug|failure|error)$", "issue type bug")
	issueTypeIncident := cmd.Flags().String("issueTypeIncident", "", "issue type incident")
	issueTypeRequirement := cmd.Flags().String("issueTypeRequirement", "^(feat|feature|proposal|requirement)$", "issue type requirement")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"connectionId":         *connectionId,
			"owner":                *owner,
			"repo":                 *repo,
			"prType":               *prType,
			"prComponent":          *prComponent,
			"prBodyClosePattern":   *prBodyClosePattern,
			"issueSeverity":        *issueSeverity,
			"issuePriority":        *issuePriority,
			"issueComponent":       *issueComponent,
			"issueTypeBug":         *issueTypeBug,
			"issueTypeIncident":    *issueTypeIncident,
			"issueTypeRequirement": *issueTypeRequirement,
		})
	}
	runner.RunCmd(cmd)
}
