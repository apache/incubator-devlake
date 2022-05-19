package tasks

import (
	"encoding/json"
	"strings"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiPullRequestCommitsMeta = core.SubTaskMeta{
	Name:             "extractApiPullRequestCommits",
	EntryPoint:       ExtractApiPullRequestCommits,
	EnabledByDefault: true,
	Description:      "Extract raw PullRequestCommits data into tool layer table github_commits",
}

type PrCommitsResponse struct {
	Sha    string `json:"sha"`
	Commit PullRequestCommit
	Url    string
}

type PullRequestCommit struct {
	Author struct {
		Id    int
		Name  string
		Email string
		Date  helper.Iso8601Time
	}
	Committer struct {
		Name  string
		Email string
		Date  helper.Iso8601Time
	}
	Message string
}

func ExtractApiPullRequestCommits(taskCtx core.SubTaskContext) error {
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
			Table: RAW_PULL_REQUEST_COMMIT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			apiPullRequestCommit := &PrCommitsResponse{}
			if strings.HasPrefix(string(row.Data), "{\"message\": \"Not Found\"") {
				return nil, nil
			}
			err := json.Unmarshal(row.Data, apiPullRequestCommit)
			if err != nil {
				return nil, err
			}
			pull := &SimplePr{}
			err = json.Unmarshal(row.Input, pull)
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, 2)

			githubCommit, err := convertPullRequestCommit(apiPullRequestCommit)
			if err != nil {
				return nil, err
			}
			results = append(results, githubCommit)

			githubPullRequestCommit := &models.GithubPullRequestCommit{
				CommitSha:     apiPullRequestCommit.Sha,
				PullRequestId: pull.GithubId,
			}
			if err != nil {
				return nil, err
			}
			results = append(results, githubPullRequestCommit)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertPullRequestCommit(prCommit *PrCommitsResponse) (*models.GithubCommit, error) {
	githubCommit := &models.GithubCommit{
		Sha:            prCommit.Sha,
		Message:        prCommit.Commit.Message,
		AuthorId:       prCommit.Commit.Author.Id,
		AuthorName:     prCommit.Commit.Author.Name,
		AuthorEmail:    prCommit.Commit.Author.Email,
		AuthoredDate:   prCommit.Commit.Author.Date.ToTime(),
		CommitterName:  prCommit.Commit.Committer.Name,
		CommitterEmail: prCommit.Commit.Committer.Email,
		CommittedDate:  prCommit.Commit.Committer.Date.ToTime(),
		Url:            prCommit.Url,
	}
	return githubCommit, nil
}
