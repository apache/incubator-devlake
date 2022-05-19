package tasks

import (
	"encoding/json"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiCommitsMeta = core.SubTaskMeta{
	Name:             "extractApiCommits",
	EntryPoint:       ExtractApiCommits,
	EnabledByDefault: false,
	Description:      "Extract raw commit data into tool layer table github_commits",
}

type CommitsResponse struct {
	Sha       string `json:"sha"`
	Commit    Commit
	Url       string
	Author    *models.GithubUser
	Committer *models.GithubUser
}

type Commit struct {
	Author struct {
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

func ExtractApiCommits(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraCommits by Board
			*/
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			/*
				Table store raw data
			*/
			Table: RAW_COMMIT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			commit := &CommitsResponse{}
			err := json.Unmarshal(row.Data, commit)
			if err != nil {
				return nil, err
			}
			if commit.Sha == "" {
				return nil, nil
			}

			results := make([]interface{}, 0, 4)

			githubCommit := &models.GithubCommit{
				Sha:            commit.Sha,
				Message:        commit.Commit.Message,
				AuthorName:     commit.Commit.Author.Name,
				AuthorEmail:    commit.Commit.Author.Email,
				AuthoredDate:   commit.Commit.Author.Date.ToTime(),
				CommitterName:  commit.Commit.Committer.Name,
				CommitterEmail: commit.Commit.Committer.Email,
				CommittedDate:  commit.Commit.Committer.Date.ToTime(),
				Url:            commit.Url,
			}
			if commit.Author != nil {
				githubCommit.AuthorId = commit.Author.Id
				results = append(results, commit.Author)
			}
			if commit.Committer != nil {
				githubCommit.CommitterId = commit.Committer.Id
				results = append(results, commit.Committer)

			}

			githubRepoCommit := &models.GithubRepoCommit{
				RepoId:    data.Repo.GithubId,
				CommitSha: commit.Sha,
			}

			results = append(results, githubCommit)
			results = append(results, githubRepoCommit)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
