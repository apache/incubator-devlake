package tasks

import (
	"encoding/json"
	"regexp"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var ExtractApiPullRequestsMeta = core.SubTaskMeta{
	Name:             "extractApiPullRequests",
	EntryPoint:       ExtractApiPullRequests,
	EnabledByDefault: true,
	Description:      "Extract raw PullRequests data into tool layer table github_pull_requests",
}

type GithubApiPullRequest struct {
	GithubId int `json:"id"`
	Number   int
	State    string
	Title    string
	Body     string
	HtmlUrl  string `json:"html_url"`
	Labels   []struct {
		Name string `json:"name"`
	} `json:"labels"`
	Assignee *struct {
		Login string
		Id    int
	}
	User *struct {
		Id    int
		Login string
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
	config := data.Options.Config
	var labelTypeRegex *regexp.Regexp
	var labelComponentRegex *regexp.Regexp
	var prType = config.GITHUB_PR_TYPE
	if prType == "" {
		prType = taskCtx.GetConfig("GITHUB_PR_TYPE")
	}
	var prComponent = config.GITHUB_PR_COMPONENT
	if prComponent == "" {
		prComponent = taskCtx.GetConfig("GITHUB_PR_COMPONENT")
	}
	if len(prType) > 0 {
		labelTypeRegex = regexp.MustCompile(prType)
	}
	if len(prComponent) > 0 {
		labelComponentRegex = regexp.MustCompile(prComponent)
	}

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
			apiPullRequest := &GithubApiPullRequest{}
			err := json.Unmarshal(row.Data, apiPullRequest)
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, 1)
			if apiPullRequest.GithubId == 0 {
				return nil, nil
			}
			//If this is a pr, ignore
			githubPr, err := convertGithubPullRequest(apiPullRequest, data.Repo.GithubId)
			if err != nil {
				return nil, err
			}
			for _, label := range apiPullRequest.Labels {
				results = append(results, &models.GithubPullRequestLabel{
					PullId:    githubPr.GithubId,
					LabelName: label.Name,
				})
				// if pr.Type has not been set and prType is set in .env, process the below
				if labelTypeRegex != nil {
					groups := labelTypeRegex.FindStringSubmatch(label.Name)
					if len(groups) > 0 {
						githubPr.Type = groups[1]
					}
				}

				// if pr.Component has not been set and prComponent is set in .env, process
				if labelComponentRegex != nil {
					groups := labelComponentRegex.FindStringSubmatch(label.Name)
					if len(groups) > 0 {
						githubPr.Component = groups[1]
					}
				}
			}
			results = append(results, githubPr)

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
		Url:             pull.HtmlUrl,
		AuthorName:      pull.User.Login,
		AuthorId:        pull.User.Id,
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
