package tasks

import (
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/helper"
)

const RAW_COMMIT_TABLE = "gitlab_api_commit"

type ApiCommitResponse []GitlabApiCommit
type ApiMergeRequestCommitResponse []GitlabApiCommit

type GitlabApiCommit struct {
	GitlabId       string `json:"id"`
	Title          string
	Message        string
	ProjectId      int
	ShortId        string           `json:"short_id"`
	AuthorName     string           `json:"author_name"`
	AuthorEmail    string           `json:"author_email"`
	AuthoredDate   core.Iso8601Time `json:"authored_date"`
	CommitterName  string           `json:"committer_name"`
	CommitterEmail string           `json:"committer_email"`
	CommittedDate  core.Iso8601Time `json:"committed_date"`
	WebUrl         string           `json:"web_url"`
	Stats          struct {
		Additions int
		Deletions int
		Total     int
	}
}

func CollectApiCommits(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_TABLE)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		Incremental:        false,
		UrlTemplate:        "projects/{{ .Params.ProjectId }}/repository/commits",
		Query:              GetQuery,
		Concurrency:        20,
		ResponseParser:     GetRawMessageFromResponse,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

// Convert the API response to our DB model instance
func ConvertCommit(commit *GitlabApiCommit) (*models.GitlabCommit, error) {
	gitlabCommit := &models.GitlabCommit{
		Sha:            commit.GitlabId,
		Title:          commit.Title,
		Message:        commit.Message,
		ShortId:        commit.ShortId,
		AuthorName:     commit.AuthorName,
		AuthorEmail:    commit.AuthorEmail,
		AuthoredDate:   commit.AuthoredDate.ToTime(),
		CommitterName:  commit.CommitterName,
		CommitterEmail: commit.CommitterEmail,
		CommittedDate:  commit.CommittedDate.ToTime(),
		WebUrl:         commit.WebUrl,
		Additions:      commit.Stats.Additions,
		Deletions:      commit.Stats.Deletions,
		Total:          commit.Stats.Total,
	}
	return gitlabCommit, nil
}
