package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = ExtractApiIssues

func ExtractApiIssues(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: JiraApiParams{
				SourceId: data.Source.ID,
				BoardId:  data.Options.BoardId,
			},
			/*
				Table store raw data
			*/
			Table: RAW_ISSUE_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var apiIssue apiv2models.Issue
			err := json.Unmarshal(row.Data, &apiIssue)
			if err != nil {
				return nil, err
			}
			var results []interface{}
			sprints, issue, _, worklogs, changelogs, changelogItems := apiIssue.ExtractEntities(data.Source.ID, data.Source.StoryPointField)
			for _, sprintId := range sprints {
				sprintIssue := &models.JiraSprintIssue{
					SourceId: data.Source.ID,
					SprintId: sprintId,
					IssueId:  issue.IssueId,
				}
				results = append(results, sprintIssue)
			}
			results = append(results, issue)
			for _, worklog := range worklogs {
				results = append(results, worklog)
			}
			for _, changelog := range changelogs {
				results = append(results, changelog)
			}
			for _, changelogItem := range changelogItems {
				results = append(results, changelogItem)
			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
