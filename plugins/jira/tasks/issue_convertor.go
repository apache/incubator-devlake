package tasks

import (
	"net/url"
	"path/filepath"
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	jiraModels "github.com/merico-dev/lake/plugins/jira/models"
)

func ConvertIssues(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*JiraTaskData)

	jiraIssue := &jiraModels.JiraIssue{}
	// select all issues belongs to the board
	cursor, err := db.Model(jiraIssue).
		Select("_tool_jira_issues.*").
		Joins("left join _tool_jira_board_issues on _tool_jira_board_issues.issue_id = _tool_jira_issues.issue_id").
		Where(
			"_tool_jira_board_issues.source_id = ? AND _tool_jira_board_issues.board_id = ?",
			data.Options.SourceId,
			data.Options.BoardId,
		).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&jiraModels.JiraIssue{})
	userIdGen := didgen.NewDomainIdGenerator(&jiraModels.JiraUser{})
	boardIdGen := didgen.NewDomainIdGenerator(&jiraModels.JiraBoard{})
	boardId := boardIdGen.Generate(data.Options.SourceId, data.Options.BoardId)

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(jiraModels.JiraIssue{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				SourceId: data.Source.ID,
				BoardId:  data.Options.BoardId,
			},
			Table: RAW_ISSUE_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			jiraIssue := inputRow.(*jiraModels.JiraIssue)
			issue := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: issueIdGen.Generate(jiraIssue.SourceId, jiraIssue.IssueId),
				},
				Url:                     convertURL(jiraIssue.Self, jiraIssue.Key),
				Number:                  jiraIssue.Key,
				Title:                   jiraIssue.Summary,
				EpicKey:                 jiraIssue.EpicKey,
				Type:                    jiraIssue.StdType,
				Status:                  jiraIssue.StdStatus,
				StoryPoint:              jiraIssue.StdStoryPoint,
				OriginalEstimateMinutes: jiraIssue.OriginalEstimateMinutes,
				CreatorId:               userIdGen.Generate(data.Options.SourceId, jiraIssue.CreatorAccountId),
				ResolutionDate:          jiraIssue.ResolutionDate,
				Priority:                jiraIssue.PriorityName,
				CreatedDate:             &jiraIssue.Created,
				UpdatedDate:             &jiraIssue.Updated,
				LeadTimeMinutes:         jiraIssue.LeadTimeMinutes,
				TimeSpentMinutes:        jiraIssue.SpentMinutes,
			}
			if jiraIssue.AssigneeAccountId != "" {
				issue.AssigneeId = userIdGen.Generate(data.Options.SourceId, jiraIssue.AssigneeAccountId)
			}
			if jiraIssue.AssigneeDisplayName != "" {
				issue.AssigneeName = jiraIssue.AssigneeDisplayName
			}
			if jiraIssue.ParentId != 0 {
				issue.ParentIssueId = issueIdGen.Generate(data.Options.SourceId, jiraIssue.ParentId)
			}
			boardIssue := &ticket.BoardIssue{
				BoardId: boardId,
				IssueId: issue.Id,
			}
			return []interface{}{
				issue,
				boardIssue,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

func convertURL(api, issueKey string) string {
	u, err := url.Parse(api)
	if err != nil {
		return api
	}
	u.Path = filepath.Join("/browse", issueKey)
	return u.String()
}
