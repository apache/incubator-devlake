package tasks

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/utils"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

type JiraApiIssuesResponse struct {
	JiraPagination
	Issues []JiraApiIssue `json:"issues"`
}

type JiraApiIssue struct {
	Id     string                 `json:"id"`
	Self   string                 `json:"self"`
	Key    string                 `json:"key"`
	Fields map[string]interface{} `json:"fields"`
}

type JiraApiIssueFields struct {
	Name                  string
	Summary               string
	IssueType             JiraApiIssueType
	Project               JiraApiProject
	Worklog               JiraApiWorklog
	Status                JiraApiStatus
	Creator               JiraApiUser
	Assignee              *JiraApiUser
	Priority              *JiraApiIssuePriority
	TimeTracking          JiraApiIssueTimeTracking
	AggregateTimeEstimate int64
	Parent                *JiraApiIssue
	Sprint                *JiraApiSprint
	ClosedSprints         []JiraApiSprint
	Created               core.Iso8601Time
	Updated               *core.Iso8601Time
	ResolutionDate        *core.Iso8601Time
}

type JiraApiIssuePriority struct {
	Self    string
	IconUrl string
	Name    string
	Id      string
}

type JiraApiIssueTimeTracking struct {
	OriginalEstimate        string
	RemainingEstimate       string
	OriginalEstimateSeconds int64
	RemainingEstimatSeconds int64
}

type JiraApiStatusCategory struct {
	Self      string
	Id        int
	Key       string
	ColorName string
	Name      string
}

type JiraApiStatus struct {
	Self           string
	Description    string
	IconUrl        string
	Name           string
	Id             string
	StatusCategory JiraApiStatusCategory
}

type JiraApiIssueType struct {
	Self           string
	Id             string
	Description    string
	IconUrl        string
	Name           string
	Subtask        bool
	AvatarId       int
	HierarchyLevel int
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
			logger.Error("jira collect issues:  get last sync time failed", err)
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
		logger.Error("jira collect issues: scheduler failed", err)
		return err
	}
	defer scheduler.Release()

	err = jiraApiClient.FetchPages(scheduler, fmt.Sprintf("agile/1.0/board/%v/issue", boardId), query,
		func(res *http.Response) error {
			// parse response
			jiraApiIssuesResponse := &JiraApiIssuesResponse{}
			err := core.UnmarshalResponse(res, jiraApiIssuesResponse)
			if err != nil {
				logger.Error("jira collect issues: unmarshal issue response failed", err)
				return err
			}

			// process issues
			for _, jiraApiIssue := range jiraApiIssuesResponse.Issues {
				jiraIssue, sprints, err := convertIssue(source, &jiraApiIssue)
				if err != nil {
					logger.Error("jira collect issues: convert issue failed", err)
					return err
				}
				// issue
				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(jiraIssue).Error
				if err != nil {
					logger.Error("jira collect issues: save issue failed", err)
					return err
				}

				// board / issue relationship
				lakeModels.Db.FirstOrCreate(&models.JiraBoardIssue{
					SourceId: source.ID,
					BoardId:  boardId,
					IssueId:  jiraIssue.IssueId,
				})

				// spirnt / issue relationship
				for _, sprintId := range sprints {
					err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(
						&models.JiraSprintIssue{
							SourceId:       source.ID,
							SprintId:       sprintId,
							IssueId:        jiraIssue.IssueId,
							ResolutionDate: jiraIssue.ResolutionDate,
						}).Error
					if err != nil {
						logger.Error("jira collect issues: save sprint issue relationship failed", err)
						return err
					}
				}
				err = handleWorklogs(jiraApiClient, source, &jiraApiIssue)
				if err != nil {
					logger.Error("jira collect issues: save worklogs failed", err)
					return err
				}
			}
			return nil
		})
	if err != nil {
		logger.Error("jira collect issues: fetch page failed", err)
		return err
	}
	scheduler.WaitUntilFinish()
	return nil
}

func convertIssue(source *models.JiraSource, jiraApiIssue *JiraApiIssue) (jiraIssue *models.JiraIssue, sprints []uint64, err error) {
	id, err := strconv.ParseUint(jiraApiIssue.Id, 10, 64)
	if err != nil {
		logger.Error("jira convert issue: parse issue id failed", err)
		return nil, nil, err
	}
	fields := &JiraApiIssueFields{}
	// decode known fields to a strong type struct to avoid type assertion
	err = core.DecodeMapStruct(jiraApiIssue.Fields, fields)
	if err != nil {
		logger.Error("jira convert issue: decode fields failed", err)
		return nil, nil, err
	}

	projectId, err := strconv.ParseUint(fields.Project.Id, 10, 64)
	if err != nil {
		logger.Error("jira convert issue: parse project id failed", err)
		return nil, nil, err
	}
	epicKey := ""
	if source.EpicKeyField != "" {
		epicKey, _ = jiraApiIssue.Fields[source.EpicKeyField].(string)
	}
	workload := 0.0
	if source.StoryPointField != "" {
		workload, _ = jiraApiIssue.Fields[source.StoryPointField].(float64)
	}
	jiraIssue = &models.JiraIssue{
		AllFields:          jiraApiIssue.Fields,
		SourceId:           source.ID,
		IssueId:            id,
		ProjectId:          projectId,
		Self:               jiraApiIssue.Self,
		Key:                jiraApiIssue.Key,
		Summary:            fields.Summary,
		Type:               fields.IssueType.Name,
		StatusName:         fields.Status.Name,
		StatusKey:          fields.Status.StatusCategory.Key,
		StatusCategory:     fields.Status.StatusCategory.Name,
		EpicKey:            epicKey,
		ResolutionDate:     core.Iso8601TimeToTime(fields.ResolutionDate),
		StoryPoint:         workload,
		CreatorAccountId:   fields.Creator.AccountId,
		CreatorAccountType: fields.Creator.AccountType,
		CreatorDisplayName: fields.Creator.DisplayName,
		Created:            fields.Created.ToTime(),
		Updated:            fields.Updated.ToTime(),
	}
	if fields.Assignee != nil {
		jiraIssue.AssigneeAccountId = fields.Assignee.AccountId
		jiraIssue.AssigneeAccountType = fields.Assignee.AccountType
		jiraIssue.AssigneeDisplayName = fields.Assignee.DisplayName
	}
	if fields.Priority != nil {
		priorityId, err := strconv.ParseUint(fields.Priority.Id, 10, 64)
		if err != nil {
			logger.Error("jira convert issue: parse priority id failed", err)
			return nil, nil, err
		}
		jiraIssue.PriorityId = priorityId
		jiraIssue.PriorityName = fields.Priority.Name
	}
	if fields.TimeTracking.OriginalEstimateSeconds != 0 {
		jiraIssue.OriginalEstimateMinutes = fields.TimeTracking.OriginalEstimateSeconds / 60
		jiraIssue.AggregateEstimateMinutes = fields.AggregateTimeEstimate / 60
		jiraIssue.RemainingEstimateMinutes = fields.TimeTracking.RemainingEstimatSeconds / 60
	}
	// depend on board settings, subtasks may or may not be collected.
	if fields.Parent != nil {
		parentId, err := strconv.ParseUint(fields.Parent.Id, 10, 64)
		if err != nil {
			logger.Error("jira convert issue: parse parent issue id failed", err)
			return nil, nil, err
		}
		jiraIssue.ParentId = parentId
		jiraIssue.ParentKey = fields.Parent.Key
	}
	// latest sprint
	if fields.Sprint != nil {
		// set sprint to latest sprint id/name
		jiraIssue.SprintId = fields.Sprint.Id
		jiraIssue.SprintName = fields.Sprint.Name
		sprints = append(sprints, jiraIssue.SprintId)
	}
	// closed sprint
	if fields.ClosedSprints != nil {
		for _, sprint := range fields.ClosedSprints {
			sprints = append(sprints, sprint.Id)
		}
	}
	return jiraIssue, sprints, nil
}
