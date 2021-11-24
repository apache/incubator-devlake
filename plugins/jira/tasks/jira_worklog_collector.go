package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	BatchSize = 100
)

type JiraApiWorklog struct {
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
	Worklogs   []struct {
		Author struct {
			AccountID string `json:"accountId"`
		} `json:"author"`
		UpdateAuthor struct {
			AccountID string `json:"accountId"`
		} `json:"updateAuthor"`
		Updated          core.Iso8601Time `json:"updated"`
		Started          core.Iso8601Time `json:"started"`
		TimeSpent        string           `json:"timeSpent"`
		TimeSpentSeconds int              `json:"timeSpentSeconds"`
		ID               string           `json:"id"`
		IssueID          string           `json:"issueId"`
	} `json:"worklogs"`
}

func (w *JiraApiWorklog) toJiraWorklogs(sourceId, IssueId uint64) []models.JiraWorklog {
	worklogs := make([]models.JiraWorklog, len(w.Worklogs))
	for i, item := range w.Worklogs {
		worklogs[i] = models.JiraWorklog{
			SourceId:         sourceId,
			IssueId:          IssueId,
			WorklogId:        item.ID,
			AuthorId:         item.Author.AccountID,
			UpdateAuthorId:   item.UpdateAuthor.AccountID,
			TimeSpent:        item.TimeSpent,
			TimeSpentSeconds: item.TimeSpentSeconds,
			Updated:          item.Updated.ToTime(),
			Started:          item.Started.ToTime(),
		}
	}
	return worklogs
}

func handleWorklogs(jiraApiClient *JiraApiClient, source *models.JiraSource, jiraApiIssue *JiraApiIssue) error {
	issueId, err := strconv.ParseUint(jiraApiIssue.Id, 10, 64)
	if err != nil {
		logger.Error("jira extract worklog: parse issue id failed", err)
		return err
	}
	fields := &JiraApiIssueFields{}
	err = core.DecodeMapStruct(jiraApiIssue.Fields, fields)
	if err != nil {
		logger.Error("jira extract worklog: decode fields failed", err)
		return err
	}
	var worklogs []models.JiraWorklog
	if fields.Worklog.Total == len(fields.Worklog.Worklogs) {
		worklogs = fields.Worklog.toJiraWorklogs(source.ID, issueId)
		return saveWorklogs(worklogs)
	}
	path := fmt.Sprintf("api/2/issue/%d/worklog", issueId)
	err = jiraApiClient.FetchWithoutPaginationHeaders(path, nil, func(res *http.Response) (int, error) {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return 0, err
		}
		defer res.Body.Close()
		var jiraApiWorklog JiraApiWorklog
		err = json.Unmarshal(body, &jiraApiWorklog)
		if err != nil {
			return 0, err
		}
		worklogs = jiraApiWorklog.toJiraWorklogs(source.ID, issueId)
		err = saveWorklogs(worklogs)
		if err != nil {
			return 0, err
		}
		return len(worklogs), nil
	})
	return nil
}

func saveWorklogs(worklogs []models.JiraWorklog) error {
	if len(worklogs) == 0 {
		return nil
	}
	return lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(worklogs, BatchSize).Error
}
