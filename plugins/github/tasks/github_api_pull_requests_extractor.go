package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var _ core.SubTaskEntryPoint = ExtractApiPullRequests

type ApiPullRequestResponse []GithubApiPullRequest

type GithubApiPullRequest struct {
	GithubId int `json:"id"`
	Number   int
	State    string
	Title    string
	Body     string
	Labels   []struct {
		Name string `json:"name"`
	} `json:"labels"`
	Assignee *struct {
		Login string
		Id    int
	}
	ClosedAt        *core.Iso8601Time `json:"closed_at"`
	MergedAt        *core.Iso8601Time `json:"merged_at"`
	GithubCreatedAt core.Iso8601Time  `json:"created_at"`
	GithubUpdatedAt core.Iso8601Time  `json:"updated_at"`
	MergeCommitSha  string            `json:"merge_commit_sha"`
	Head            struct {
		Ref string
		Sha string
	}
	Base struct {
		Ref string
		Sha string
	}
}

func ExtractApiPullRequests(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			/*
				Table store raw data
			*/
			Table: RAW_PULL_REQUEST_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &ApiPullRequestResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, len(*body)*2)
			for _, apiPullRequest := range *body {
				if apiPullRequest.GithubId == 0 {
					return nil, nil
				}
				//If this is a pr, ignore

				githubPr, err := convertGithubPullRequest(&apiPullRequest, data.Repo.GithubId)
				if err != nil {
					return nil, err
				}
				results = append(results, githubPr)
				for _, label := range apiPullRequest.Labels {
					results = append(results, &models.GithubPullRequestLabel{
						PullId:    githubPr.GithubId,
						LabelName: label.Name,
					})

				}
			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
func convertGithubPullRequest(pull *GithubApiPullRequest, repoId int) (*models.GithubPullRequest, error) {
	githubPull := &models.GithubPullRequest{
		GithubId:        pull.GithubId,
		RepoId:          repoId,
		Number:          pull.Number,
		State:           pull.State,
		Title:           pull.Title,
		GithubCreatedAt: pull.GithubCreatedAt.ToTime(),
		GithubUpdatedAt: pull.GithubUpdatedAt.ToTime(),
		ClosedAt:        core.Iso8601TimeToTime(pull.ClosedAt),
		MergedAt:        core.Iso8601TimeToTime(pull.MergedAt),
		MergeCommitSha:  pull.MergeCommitSha,
		Body:            pull.Body,
		BaseRef:         pull.Base.Ref,
		BaseCommitSha:   pull.Base.Sha,
		HeadRef:         pull.Head.Ref,
		HeadCommitSha:   pull.Head.Sha,
	}
	return githubPull, nil
}
