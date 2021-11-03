package tasks

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/utils"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

type JiraApiIssue struct {
	Id     string                 `json:"id"`
	Self   string                 `json:"self"`
	Key    string                 `json:"key"`
	Fields map[string]interface{} `json:"fields"`
}

type JiraApiIssuesResponse struct {
	JiraPagination
	Issues []JiraApiIssue `json:"issues"`
}

func CollectIssues(
	jiraApiClient *JiraApiClient,
	source *models.JiraSource,
	boardId uint64,
	since time.Time,
	ctx context.Context,
) error {
	// user didn't specify a time range to sync, try load from database
	if since.IsZero() {
		var latestUpdated models.JiraIssue
		err := lakeModels.Db.Where("source_id = ?", source.ID).Order("updated DESC").Limit(1).Find(&latestUpdated).Error
		if err != nil {
			return err
		}
		since = latestUpdated.Updated
	}
	// build jql
	jql := "ORDER BY updated ASC"
	if !since.IsZero() {
		// prepend a time range criteria if `since` was specified, either by user or from database
		jql = fmt.Sprintf("updated >= '%v' %v", since.Format("2006/01/02 15:04"), jql)
	}

	query := &url.Values{}
	query.Set("jql", jql)

	scheduler, err := utils.NewWorkerScheduler(10, 50, ctx)
	if err != nil {
		return err
	}
	defer scheduler.Release()

	err = jiraApiClient.FetchPages(scheduler, fmt.Sprintf("/agile/1.0/board/%v/issue", boardId), query,
		func(res *http.Response) error {
			// parse response
			jiraApiIssuesResponse := &JiraApiIssuesResponse{}
			err := core.UnmarshalResponse(res, jiraApiIssuesResponse)
			if err != nil {
				logger.Error("unmarshal issue response errro", err)
				return err
			}

			// process issues
			for _, jiraApiIssue := range jiraApiIssuesResponse.Issues {
				jiraIssue, sprints, err := convertIssue(source, &jiraApiIssue)
				if err != nil {
					return err
				}
				// issue
				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(jiraIssue).Error
				if err != nil {
					return err
				}

				// board / issue relationship
				lakeModels.Db.FirstOrCreate(&models.JiraBoardIssue{
					SourceId: source.ID,
					BoardId:  boardId,
					IssueId:  jiraIssue.IssueId,
				})

				// spirnt / issue relationship
				for _, sprintId := range sprints{
					err = lakeModels.Db.FirstOrCreate(
						&models.JiraSprintIssue{
							SourceId: source.ID,
							SprintId: sprintId,
							IssueId:  jiraIssue.IssueId,
						}).Error
					if err != nil {
						logger.Error("save sprint issue relationship error", err)
						return err
					}
				}
			}
			return nil
		})
	if err != nil {
		return err
	}
	scheduler.WaitUntilFinish()
	return nil
}

func convertIssue(source *models.JiraSource, jiraApiIssue *JiraApiIssue) (jiraIssue *models.JiraIssue, sprints []uint64, err error) {
	defer func() {
		// type assertion could cause panic, this is to capture this type of error and propagate
		if r := recover(); r != nil {
			err = fmt.Errorf("jira issue converter failed: %w\n%s", r, debug.Stack())
		}
	}()
	id, err := strconv.ParseUint(jiraApiIssue.Id, 10, 64)
	if err != nil {
		return nil, nil, err
	}
	created, err := core.ConvertStringToTime(jiraApiIssue.Fields["created"].(string))
	if err != nil {
		return nil, nil, err
	}
	updated, err := core.ConvertStringToTime(jiraApiIssue.Fields["updated"].(string))
	if err != nil {
		return nil, nil, err
	}
	projectId, err := strconv.ParseUint(
		jiraApiIssue.Fields["project"].(map[string]interface{})["id"].(string), 10, 64,
	)
	if err != nil {
		return nil, nil, err
	}
	status := jiraApiIssue.Fields["status"].(map[string]interface{})
	statusName := status["name"].(string)
	statusKey := status["statusCategory"].(map[string]interface{})["key"].(string)
	statusCategory := status["statusCategory"].(map[string]interface{})["name"].(string)
	epicKey := ""
	if source.EpicKeyField != "" {
		epicKey, _ = jiraApiIssue.Fields[source.EpicKeyField].(string)
	}
	resolutionDate := sql.NullTime{}
	if rd, ok := jiraApiIssue.Fields["resolutiondate"]; ok && rd != nil {
		if resolutionDate.Time, err = core.ConvertStringToTime(rd.(string)); err == nil {
			resolutionDate.Valid = true
		}
	}
	workload := 0.0
	if source.StoryPointField != "" {
		workload, _ = jiraApiIssue.Fields[source.StoryPointField].(float64)
	}
	creator := jiraApiIssue.Fields["creator"].(map[string]interface{})
	jiraIssue = &models.JiraIssue{
		SourceId:           source.ID,
		IssueId:            id,
		ProjectId:          projectId,
		Self:               jiraApiIssue.Self,
		Key:                jiraApiIssue.Key,
		Summary:            jiraApiIssue.Fields["summary"].(string),
		Type:               jiraApiIssue.Fields["issuetype"].(map[string]interface{})["name"].(string),
		StatusName:         statusName,
		StatusKey:          statusKey,
		StatusCategory:     statusCategory,
		EpicKey:            epicKey,
		ResolutionDate:     resolutionDate,
		StoryPoint:         workload,
		CreatorAccountId:   creator["accountId"].(string),
		CreatorAccountType: creator["accountType"].(string),
		CreatorDisplayName: creator["displayName"].(string),
		Created:            created,
		Updated:            updated,
	}
	if assigneeField, ok := jiraApiIssue.Fields["assignee"]; ok && assigneeField != nil {
		assignee := assigneeField.(map[string]interface{})
		jiraIssue.AssigneeAccountId = assignee["accountId"].(string)
		jiraIssue.AssigneeAccountType = assignee["accountType"].(string)
		jiraIssue.AssigneeDisplayName = assignee["displayName"].(string)
	}
	if priorityField, ok := jiraApiIssue.Fields["priority"]; ok {
		priority := priorityField.(map[string]interface{})
		priorityId, err := strconv.ParseUint(priority["id"].(string), 10, 64)
		if err != nil {
			return nil, nil, err
		}
		jiraIssue.PriorityId = priorityId
		jiraIssue.PriorityName = priority["name"].(string)
	}
	if timetrackingField, ok := jiraApiIssue.Fields["timetracking"]; ok {
		timetracking := timetrackingField.(map[string]interface{})
		if len(timetracking) > 0 {
			if originalEstimateSeconds := timetracking["originalEstimateSeconds"]; originalEstimateSeconds != nil {
				jiraIssue.OriginalEstimateMinutes = int64(originalEstimateSeconds.(float64) / 60)
			}
			if atoe := jiraApiIssue.Fields["aggregatetimeoriginalestimate"]; atoe != nil {
				jiraIssue.AggregateEstimateMinutes = int64(atoe.(float64) / 60)
			}
			if remainingEstimateSeconds := timetracking["remainingEstimateSeconds"]; remainingEstimateSeconds != nil {
				jiraIssue.RemainingEstimateMinutes = int64(remainingEstimateSeconds.(float64) / 60)
			}
		}
	}
	// this would never be true if we collect issues by board
	if parentField, ok := jiraApiIssue.Fields["parent"]; ok {
		parent := parentField.(map[string]interface{})
		if parent != nil {
			parentId, err := strconv.ParseUint(parent["id"].(string), 10, 64)
			if err != nil {
				return nil, nil, err
			}
			jiraIssue.ParentId = parentId
			jiraIssue.ParentKey = parent["key"].(string)
		}
	}
	// latest sprint
	if sprintField, ok := jiraApiIssue.Fields["sprint"]; ok && sprintField != nil {
		if sprint := sprintField.(map[string]interface{}); ok {
			// set sprint to latest sprint id/name
			jiraIssue.SprintId = uint64(sprint["id"].(float64))
			jiraIssue.SprintName = sprint["name"].(string)
			sprints = append(sprints, jiraIssue.SprintId)
		}
	}
	// closed sprint
	if closedSprintField, ok := jiraApiIssue.Fields["closedSprints"]; ok && closedSprintField != nil {
		if clsedSprints := closedSprintField.([]interface{}); ok {
			for _, sprint := range clsedSprints {
				if s, yes := sprint.(map[string]interface{}); yes && s != nil{
					sprints = append(sprints, uint64(s["id"].(float64)))
				}
			}
		}
	}
	return jiraIssue, sprints, nil
}
