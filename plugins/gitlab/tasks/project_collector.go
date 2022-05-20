package tasks

import (
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

const RAW_PROJECT_TABLE = "gitlab_api_project"

type GitlabApiProject struct {
	GitlabId          int    `json:"id"`
	Name              string `josn:"name"`
	Description       string `json:"description"`
	DefaultBranch     string `json:"default_branch"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebUrl            string `json:"web_url"`
	CreatorId         int
	Visibility        string              `json:"visibility"`
	OpenIssuesCount   int                 `json:"open_issues_count"`
	StarCount         int                 `json:"star_count"`
	ForkedFromProject *GitlabApiProject   `json:"forked_from_project"`
	CreatedAt         helper.Iso8601Time  `json:"created_at"`
	LastActivityAt    *helper.Iso8601Time `json:"last_activity_at"`
}

var CollectProjectMeta = core.SubTaskMeta{
	Name:             "collectApiProject",
	EntryPoint:       CollectApiProject,
	EnabledByDefault: true,
	Description:      "Collect project data from gitlab api",
}

func CollectApiProject(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		Incremental:        false,
		UrlTemplate:        "projects/{{ .Params.ProjectId }}",
		Query:              GetQuery,
		ResponseParser:     helper.GetRawMessageDirectFromResponse,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
