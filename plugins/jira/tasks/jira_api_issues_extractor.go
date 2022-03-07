package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
)

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
			body := &JiraApiIssuesResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			// need to extract 3 kinds of entities here
			results := make([]interface{}, 0, len(body.Issues)*3)
			for _, apiIssue := range body.Issues {
				jiraIssue, sprintIds, err := convertIssue(data.Source, &apiIssue)
				if err != nil {
					return nil, err
				}
				results = append(results, jiraIssue)
				results = append(results, &models.JiraBoardIssue{
					SourceId: data.Source.ID,
					BoardId:  data.Options.BoardId,
					IssueId:  jiraIssue.IssueId,
				})
				for _, sprintId := range sprintIds {
					results = append(results, &models.JiraSprintIssue{
						SourceId: data.Source.ID,
						SprintId: sprintId,
						IssueId:  jiraIssue.IssueId,
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
