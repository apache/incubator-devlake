package tasks

import (
	"strconv"
	"time"

	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
)

func ExtractApiIssues(
	source *models.JiraSource,
	boardId uint64,
	since time.Time,
) error {
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		Table: "jira_api_issues",
		Params: JiraApiParams{
			SourceId: source.ID,
			BoardId:  boardId,
		},
		RowData: &JiraApiIssue{},
		Extractors: []helper.RawDataExtractor{
			func(row interface{}, params interface{}) (interface{}, error) {
				p := params.(*JiraApiParams)
				issue := row.(*JiraApiIssue)
				issueId, err := strconv.ParseUint(issue.Id, 10, 64)
				if err != nil {
					return nil, err
				}
				return &models.JiraIssue{
					SourceId: p.SourceId,
					IssueId:  issueId,
					// other fields
				}, nil
			},
			func(row interface{}, params interface{}) (interface{}, error) {
				p := params.(*JiraApiParams)
				issue := row.(*JiraApiIssue)
				issueId, err := strconv.ParseUint(issue.Id, 10, 64)
				if err != nil {
					return nil, err
				}
				return &models.JiraBoardIssue{
					SourceId: p.SourceId,
					BoardId:  p.BoardId,
					IssueId:  issueId,
				}, nil
			},
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
