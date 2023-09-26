package tasks

import (
	"fmt"
	"net/url"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
)

var ConvertIssuesMeta = plugin.SubTaskMeta{
	Name:             "convertIssues",
	EntryPoint:       ConvertIssues,
	EnabledByDefault: true,
	Description:      "convert clickup tasks",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertIssues(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*ClickupTaskData)

	clickUpTask := &models.ClickUpTask{}
	clauses := []dal.Clause{
		dal.Select("_tool_clickup_task.*"),
		dal.From(clickUpTask),
		dal.Where(
			"_tool_clickup_task.connection_id = ? AND _tool_clickup_task.space_id = ?",
			data.Options.ConnectionId,
			data.Options.ScopeId,
		),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&models.ClickUpTask{})
	boardIdGen := didgen.NewDomainIdGenerator(&models.ClickUpSpace{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.ClickUpUser{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.ClickUpTask{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ClickupApiParams{
				TeamId: data.TeamId,
			},
			Table: RAW_ISSUE_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			clickUpTask := inputRow.(*models.ClickUpTask)
			fmt.Printf("%s %s\n", *clickUpTask.CustomId, clickUpTask.StatusName)
			issue := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: issueIdGen.Generate(clickUpTask.ConnectionId, clickUpTask.TaskId),
				},
				Url:            clickUpTask.Url,
				IssueKey:       stringOrEmpty(clickUpTask.CustomId),
				Title:          clickUpTask.Name,
				Type:           clickUpTask.NormalizedType,
				Status:         convertStatus(clickUpTask.StatusName, clickUpTask.StatusType),
				OriginalStatus: clickUpTask.StatusName,
				StoryPoint:     clickUpTask.Points,
				ResolutionDate: timestampToTime(clickUpTask.DateDone),
				Priority:       clickUpTask.Priority,
				CreatedDate:    timestampToTime(clickUpTask.DateCreated),
				UpdatedDate:    timestampToTime(clickUpTask.DateUpdated),
				CreatorId:      fmt.Sprintf("%d", clickUpTask.CreatorId),
				CreatorName:    string(clickUpTask.CreatorUsername),
				// LeadTimeMinutes:         int64(clickUpTask.LeadTimeMinutes),
				// TimeSpentMinutes:        clickUpTask.SpentMinutes,
				// OriginalProject:         clickUpTask.ProjectName,
			}

			if clickUpTask.FirstAssigneeId != nil {
				issue.AssigneeId = accountIdGen.Generate(clickUpTask.ConnectionId, fromInt64ToStringOrEmpty(clickUpTask.FirstAssigneeId))
				issue.AssigneeName = stringOrEmpty(clickUpTask.FirstAssigneeUsername)
			}

			if clickUpTask.DateDone != 0 {
				// todo: LeadTimeMinutes should be calculated in extractor as
				// with other implemnetations
				issue.LeadTimeMinutes = int64(
					issue.ResolutionDate.Sub(
						*issue.CreatedDate,
					).Minutes(),
				)
			}
			// if clickUpTask.CreatorAccountId != "" {
			// 	issue.CreatorId = accountIdGen.Generate(data.Options.ConnectionId, clickUpTask.CreatorAccountId)
			// }
			// if clickUpTask.CreatorDisplayName != "" {
			// 	issue.CreatorName = clickUpTask.CreatorDisplayName
			// }
			// if clickUpTask.AssigneeAccountId != "" {
			// 	issue.AssigneeId = accountIdGen.Generate(data.Options.ConnectionId, clickUpTask.AssigneeAccountId)
			// }
			// if clickUpTask.AssigneeDisplayName != "" {
			// 	issue.AssigneeName = clickUpTask.AssigneeDisplayName
			// }
			// if clickUpTask.ParentId != 0 {
			// 	issue.ParentIssueId = issueIdGen.Generate(data.Options.ConnectionId, clickUpTask.ParentId)
			// }
			// boardIssue := &ticket.BoardIssue{
			// 	BoardId: boardId,
			// 	IssueId: issue.Id,
			// }
			boardId := boardIdGen.Generate(data.Options.ConnectionId, data.Options.ScopeId)
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

func timestampToTime(int int64) *time.Time {
	// convert timestamp in milliseconds to seconds
	tm := time.Unix(int/1000, 0)
	return &tm
}

func convertStatus(statusName, statusType string) string {
	if statusType == "closed" {
		return ticket.DONE
	}
	if statusType == "done" {
		return ticket.DONE
	}
	if statusType == "open" {
		return ticket.TODO
	}

	return ticket.IN_PROGRESS
}

func fromInt64ToStringOrEmpty(int *int64) string {
	if nil == int {
		return ""
	}

	return strconv.FormatUint(uint64(*int), 10)
}

func stringOrEmpty(string *string) string {
	if string == nil {
		return ""
	}
	return *string
}

func convertURL(api, issueKey string) string {
	u, err := url.Parse(api)
	if err != nil {
		return api
	}
	before, _, _ := strings.Cut(u.Path, "/rest/agile/1.0/issue")
	u.Path = filepath.Join(before, "browse", issueKey)
	return u.String()
}
