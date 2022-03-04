package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
)

func ExtractApiIssues(taskCtx core.TaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		Ctx:   taskCtx,
		Table: RAW_ISSUE_TABLE,
		Params: JiraApiParams{
			SourceId: data.Source.ID,
			BoardId:  data.Options.BoardId,
		},
		Extract: func(body json.RawMessage, params json.RawMessage) ([]interface{}, error) {
			b := &JiraApiIssuesResponse{}
			err := json.Unmarshal(body, b)
			if err != nil {
				return nil, err
			}
			p := &JiraApiParams{}
			err = json.Unmarshal(params, p)
			if err != nil {
				return nil, err
			}

			// need to extract 3 kinds of entities here
			result := make([]interface{}, 0, len(b.Issues)*3)
			for _, apiIssue := range b.Issues {
				jiraIssue, sprints, err := convertIssue(data.Source, &apiIssue)
				if err != nil {
					return nil, err
				}
				result = append(result, jiraIssue)
				result = append(result, &models.JiraBoardIssue{
					SourceId: data.Source.ID,
					BoardId:  data.Options.BoardId,
					IssueId:  jiraIssue.IssueId,
				})
				for _, sprint := range sprints {
					result = append(result, sprint)
				}
			}
			return result, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
