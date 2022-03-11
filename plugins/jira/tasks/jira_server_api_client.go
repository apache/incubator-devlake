package tasks

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/plugins/jira/tasks/v8models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ServerVersion8 struct {
	db     *gorm.DB
	client *JiraApiClient
}

func NewServerVersion8(db *gorm.DB, client *JiraApiClient) *ServerVersion8 {
	return &ServerVersion8{db: db, client: client}
}

func (v8 *ServerVersion8) FetchPages(path string, query *url.Values, handler JiraSearchPaginationHandler) error {
	f := func(resp *http.Response) error {
		_, err := handler(resp)
		return err
	}
	return v8.client.FetchPages(path, query, f)
}
func (v8 *ServerVersion8) FetchWithoutPaginationHeaders(path string, query url.Values, handler JiraSearchPaginationHandler) error {
	return v8.client.FetchWithoutPaginationHeaders(path, query, handler)
}

func (v8 *ServerVersion8) Get(path string, handler JiraSearchPaginationHandler) error {
	resp, err := v8.client.Get(path, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = handler(resp)
	return err
}

func (v8 *ServerVersion8) CollectBoard(jiraApiClient *JiraApiClient, source *models.JiraSource, boardId uint64) error {
	var transformer v8models.Board
	err := v8.Get(fmt.Sprintf("agile/1.0/board/%d", boardId), v8.newHandler(source.ID, transformer))
	if err != nil {
		logger.Error("collect board", err)
		return err
	}
	return nil
}

func (v8 *ServerVersion8) CollectChangelogs(
	jiraApiClient *JiraApiClient,
	source *models.JiraSource,
	boardId uint64,
	rateLimitPerSecondInt int,
	ctx context.Context,
) error {
	return nil
}

func (v8 *ServerVersion8) collectWorklog(sourceId, issueId uint64) error {
	var transformer v8models.Worklog
	err := v8.FetchPages(fmt.Sprintf("api/2/issue/%d/worklog", issueId), nil, v8.newHandlerWithIssueId(sourceId, issueId, transformer))
	if err != nil {
		logger.Error("collect worklog", err)
	}
	return nil
}

func (v8 *ServerVersion8) collectRemotelinksByIssueId(source *models.JiraSource, jiraApiClient *JiraApiClient, issueId uint64) error {
	var transformer v8models.RemoteLink
	err := v8.Get(fmt.Sprintf("api/2/issue/%d/remotelink", issueId), v8.newHandlerWithIssueId(source.ID, issueId, transformer))
	if err != nil && err != ErrNotFoundResource {
		logger.Error("collect remotelink", err)
	}
	return nil
}

func (v8 *ServerVersion8) CollectIssues(
	jiraApiClient *JiraApiClient,
	source *models.JiraSource,
	boardId uint64,
	since time.Time,
	rateLimitPerSecondInt int,
	ctx context.Context,
) error {
	if since.IsZero() {
		var latestUpdated models.JiraIssue
		err := v8.db.Where("source_id = ?", source.ID).Order("updated DESC").Limit(1).Find(&latestUpdated).Error
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

	query := url.Values{}
	query.Set("jql", jql)
	query.Set("expand", "changelog")
	handler := func(resp *http.Response) (int, error) {
		return 0, v8.issueHandle(ctx, boardId, source, resp)
	}
	err := v8.FetchPages(fmt.Sprintf("agile/1.0/board/%d/issue", boardId), &query, handler)
	if err != nil {
		logger.Error("collect issue", err)
	}
	return nil
}

func (v8 *ServerVersion8) issueHandle(ctx context.Context, boardId uint64, source *models.JiraSource, resp *http.Response) error {
	blob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("issueHandle read response body", err)
		return err
	}
	defer resp.Body.Close()
	var issue v8models.Issue
	raw, err := issue.ExtractRawMessage(blob)
	if err != nil {
		logger.Error("issueHandle ExtractRawMessage", err)
		return err
	}
	issues, err := issue.Unmarshal(raw)
	if err != nil {
		logger.Error("issueHandle Unmarshal", err)
		return err
	}
	if len(issues) == 0 {
		return nil
	}
	var jiraIssues []*models.JiraIssue
	var boardIssues []*models.JiraBoardIssue
	for _, apiIssue := range issues {
		sprints, issue, needCollectWorklog, worklogs, changelogs, changelogItems := apiIssue.ExtractEntities(source.ID, source.StoryPointField)
		for _, sprintId := range sprints {
			err = v8.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(
				&models.JiraSprintIssue{
					SourceId:         source.ID,
					SprintId:         sprintId,
					IssueId:          issue.IssueId,
					ResolutionDate:   issue.ResolutionDate,
					IssueCreatedDate: &issue.Created,
				}).Error
			if err != nil {
				logger.Error("jira collect issues: save sprint issue relationship failed", err)
				return err
			}
		}
		jiraIssues = append(jiraIssues, issue)
		boardIssue := &models.JiraBoardIssue{SourceId: source.ID, BoardId: boardId, IssueId: apiIssue.ID}
		boardIssues = append(boardIssues, boardIssue)
		if needCollectWorklog {
			err = v8.collectWorklog(source.ID, issue.IssueId)
		} else {
			err = v8.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(worklogs, BatchSize).Error
		}
		if err != nil {
			logger.Error("jira collect issues: handle worklogs failed", err)
			return err
		}
		err = v8.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(changelogs, BatchSize).Error
		if err != nil {
			logger.Error("jira collect issues: save changelogs failed", err)
			return err
		}
		err = v8.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(changelogItems, BatchSize).Error
		if err != nil {
			logger.Error("jira collect issues: save changelogItems failed", err)
			return err
		}
	}
	err = v8.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(jiraIssues, BatchSize).Error
	if err != nil {
		logger.Error("jira collect issues: save jiraIssues failed", err)
		return err
	}
	err = v8.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(boardIssues, BatchSize).Error
	if err != nil {
		logger.Error("jira collect issues: save board issues failed", err)
		return err
	}
	return nil
}

func (v8 *ServerVersion8) CollectProjects(jiraApiClient *JiraApiClient, sourceId uint64) error {
	var transformer v8models.Project
	err := v8.Get("api/2/project", v8.newHandler(sourceId, transformer))
	if err != nil {
		logger.Error("jira collect projects", err)
		return err
	}
	return nil
}

func (v8 *ServerVersion8) CollectRemoteLinks(
	jiraApiClient *JiraApiClient,
	source *models.JiraSource,
	boardId uint64,
	rateLimitPerSecondInt int,
	ctx context.Context,
) error {
	return CollectRemoteLinks(jiraApiClient, source, boardId, rateLimitPerSecondInt, ctx, v8.collectRemotelinksByIssueId)
}

func (v8 *ServerVersion8) CollectSprint(jiraApiClient *JiraApiClient, source *models.JiraSource, boardId uint64) error {
	f := func(resp *http.Response) (int, error) {
		return v8.handleSprint(source.ID, boardId, resp)
	}
	err := v8.FetchWithoutPaginationHeaders(fmt.Sprintf("agile/1.0/board/%d/sprint", boardId), nil, f)
	if err != nil {
		logger.Error("jira collect sprint", err)
		return err
	}
	return nil
}

func (v8 *ServerVersion8) handleSprint(sourceId, boardId uint64, resp *http.Response) (int, error) {
	var s v8models.Sprint
	blob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("handleSprint read response body", err)
		return 0, err
	}
	defer resp.Body.Close()
	raw, err := s.ExtractRawMessage(blob)
	if err != nil {
		logger.Error("handleSprint ExtractRawMessage", err)
		return 0, err
	}
	sprints, err := s.GetJiraSprints(sourceId, raw)
	if err != nil {
		logger.Error("handleSprint GetJiraSprints", err)
		return 0, err
	}
	err = v8.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(sprints, BatchSize).Error
	if err != nil {
		logger.Error("handleSprint save sprints", err)
		return 0, err
	}
	var boardSprints []*models.JiraBoardSprint
	for _, sprint := range sprints {
		boardSprints = append(boardSprints, &models.JiraBoardSprint{
			SourceId: sourceId,
			BoardId:  boardId,
			SprintId: sprint.SprintId,
		})
	}
	err = v8.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(boardSprints, BatchSize).Error
	if err != nil {
		logger.Error("handleSprint BoardSprint", err)
		return 0, err
	}
	return len(sprints), nil
}

func (v8 *ServerVersion8) CollectUsers(jiraApiClient *JiraApiClient, sourceId uint64) error {
	return nil
}

func (v8 *ServerVersion8) newHandler(sourceId uint64, transformer v8models.Transformer) func(resp *http.Response) (int, error) {
	return func(resp *http.Response) (int, error) {
		if resp.StatusCode == http.StatusNotFound {
			return 0, ErrNotFoundResource
		}
		blob, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error("handler factory read response body", err)
			return 0, err
		}
		raw, err := transformer.ExtractRawMessage(blob)
		if err != nil {
			logger.Error("handler factory ExtractRawMessage", err)
			return 0, err
		}
		item, err := transformer.FromAPI(sourceId, raw)
		if err != nil {
			logger.Error("handler factory transformer.FromAPI", err)
			return 0, err
		}
		var l int
		reflectValue := reflect.Indirect(reflect.ValueOf(item))
		switch reflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			l = reflectValue.Len()
		}
		err = v8.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(item, BatchSize).Error
		if err != nil {
			logger.Error("handler factory DB error", err)
			return 0, err
		}
		return l, nil
	}
}

func (v8 *ServerVersion8) newHandlerWithIssueId(sourceId, issueId uint64, transformer v8models.TransformerWithIssueId) func(resp *http.Response) (int, error) {
	return func(resp *http.Response) (int, error) {
		if resp.StatusCode == http.StatusNotFound {
			return 0, ErrNotFoundResource
		}
		blob, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error("handler factory read response body", err)
			return 0, err
		}
		raw, err := transformer.ExtractRawMessage(blob)
		if err != nil {
			logger.Error("handler factory ExtractRawMessage", err)
			return 0, err
		}
		item, err := transformer.FromAPI(sourceId, issueId, raw)
		if err != nil {
			logger.Error("handler factory transformer.FromAPI", err)
			return 0, err
		}
		var l int
		reflectValue := reflect.Indirect(reflect.ValueOf(item))
		switch reflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			l = reflectValue.Len()
		}
		err = v8.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(item, BatchSize).Error
		if err != nil {
			logger.Error("handler factory DB error", err)
			return 0, err
		}
		return l, nil
	}
}
