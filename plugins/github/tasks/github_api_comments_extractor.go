package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	githubUtils "github.com/merico-dev/lake/plugins/github/utils"
	"github.com/merico-dev/lake/plugins/helper"
)

var ExtractApiCommentsMeta = core.SubTaskMeta{
	Name:             "extractApiComments",
	EntryPoint:       ExtractApiComments,
	EnabledByDefault: true,
	Description: "Extract raw comment data  into tool layer table github_pull_request_comments" +
		"and github_issue_comments",
}

type IssueComment struct {
	GithubId int `json:"id"`
	Body     string
	User     struct {
		Login string
	}
	IssueUrl        string           `json:"issue_url"`
	GithubCreatedAt core.Iso8601Time `json:"created_at"`
	GithubUpdatedAt core.Iso8601Time `json:"updated_at"`
}

func ExtractApiComments(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			Table: RAW_COMMENTS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			apiComment := &IssueComment{}
			err := json.Unmarshal(row.Data, apiComment)
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, 2)
			if apiComment.GithubId == 0 {
				return nil, nil
			}
			//If this is a pr, ignore
			issueINumber, err := githubUtils.GetIssueIdByIssueUrl(apiComment.IssueUrl)
			if err != nil {
				return nil, err
			}
			issue := &models.GithubIssue{}
			err = taskCtx.GetDb().Where("number = ? and repo_id = ?", issueINumber, data.Repo.GithubId).Limit(1).Find(issue).Error
			if err != nil {
				return nil, err
			}
			//if we can not find issues with issue number above, move the comments to github_pull_request_comments
			if issue.GithubId == 0 {
				pr := &models.GithubPullRequest{}
				err = taskCtx.GetDb().Where("number = ? and repo_id = ?", issueINumber, data.Repo.GithubId).Limit(1).Find(pr).Error
				if err != nil {
					return nil, err
				}
				githubPrComment := &models.GithubPullRequestComment{
					GithubId:        apiComment.GithubId,
					PullRequestId:   pr.GithubId,
					Body:            apiComment.Body,
					AuthorUsername:  apiComment.User.Login,
					GithubCreatedAt: apiComment.GithubCreatedAt.ToTime(),
					GithubUpdatedAt: apiComment.GithubUpdatedAt.ToTime(),
				}
				results = append(results, githubPrComment)
			} else {
				githubIssueComment := &models.GithubIssueComment{
					GithubId:        apiComment.GithubId,
					IssueId:         issue.GithubId,
					Body:            apiComment.Body,
					AuthorUsername:  apiComment.User.Login,
					GithubCreatedAt: apiComment.GithubCreatedAt.ToTime(),
					GithubUpdatedAt: apiComment.GithubUpdatedAt.ToTime(),
				}
				results = append(results, githubIssueComment)
			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
