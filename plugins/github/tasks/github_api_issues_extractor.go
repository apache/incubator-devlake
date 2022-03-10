package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var _ core.SubTaskEntryPoint = ExtractApiIssues

func ExtractApiIssues(taskCtx core.SubTaskContext) error {
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
			Table: RAW_ISSUE_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &ApiIssuesResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			// need to extract 3 kinds of entities here
			results := make([]interface{}, 0, len(*body)*2)
			for _, apiIssue := range *body {
				if apiIssue.GithubId == 0 {
					return nil, nil
				}
				//If this is a pr, ignore
				if apiIssue.PullRequest.Url != "" {
					continue
				}
				githubIssue, err := convertGithubIssue(&apiIssue, data.Repo.GithubId)
				if err != nil {
					return nil, err
				}
				results = append(results, githubIssue)
				for _, label := range apiIssue.Labels {
					results = append(results, &models.GithubIssueLabel{
						IssueId:   githubIssue.GithubId,
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
func convertGithubIssue(issue *IssuesResponse, repositoryId int) (*models.GithubIssue, error) {
	githubIssue := &models.GithubIssue{
		GithubId:        issue.GithubId,
		RepoId:          repositoryId,
		Number:          issue.Number,
		State:           issue.State,
		Title:           issue.Title,
		Body:            issue.Body,
		ClosedAt:        core.Iso8601TimeToTime(issue.ClosedAt),
		GithubCreatedAt: issue.GithubCreatedAt.ToTime(),
		GithubUpdatedAt: issue.GithubUpdatedAt.ToTime(),
	}

	if issue.Assignee != nil {
		githubIssue.AssigneeId = issue.Assignee.Id
		githubIssue.AssigneeName = issue.Assignee.Login
	}
	if issue.ClosedAt != nil {
		githubIssue.LeadTimeMinutes = uint(issue.ClosedAt.ToTime().Sub(issue.GithubCreatedAt.ToTime()).Minutes())
	}

	return githubIssue, nil
}
