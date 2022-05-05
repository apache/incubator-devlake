package tasks

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	BatchSize = 1000
)

type SprintIssuesConverter struct {
	db             *gorm.DB
	logger         core.Logger
	sprintIdGen    *didgen.DomainIdGenerator
	issueIdGen     *didgen.DomainIdGenerator
	userIdGen      *didgen.DomainIdGenerator
	sprints        map[string]*models.JiraSprint
	sprintIssue    map[string]*ticket.SprintIssue
	status         map[string]*ticket.IssueStatusHistory
	assignee       map[string]*ticket.IssueAssigneeHistory
	sprintsHistory map[string]*ticket.IssueSprintsHistory
	jiraIssue      map[string]*models.JiraIssue
}

func NewSprintIssueConverter(taskCtx core.SubTaskContext) (*SprintIssuesConverter, error) {
	data := taskCtx.GetData().(*JiraTaskData)
	connectionId := data.Connection.ID
	boardId := data.Options.BoardId
	converter := &SprintIssuesConverter{
		db:             taskCtx.GetDb(),
		logger:         taskCtx.GetLogger(),
		sprintIdGen:    didgen.NewDomainIdGenerator(&models.JiraSprint{}),
		issueIdGen:     didgen.NewDomainIdGenerator(&models.JiraIssue{}),
		userIdGen:      didgen.NewDomainIdGenerator(&models.JiraUser{}),
		sprints:        make(map[string]*models.JiraSprint),
		sprintIssue:    make(map[string]*ticket.SprintIssue),
		status:         make(map[string]*ticket.IssueStatusHistory),
		assignee:       make(map[string]*ticket.IssueAssigneeHistory),
		sprintsHistory: make(map[string]*ticket.IssueSprintsHistory),
		jiraIssue:      make(map[string]*models.JiraIssue),
	}
	return converter, converter.setupSprintIssue(connectionId, boardId)
}

func (c *SprintIssuesConverter) FeedIn(connectionId uint64, cl ChangelogItemResult) {
	if cl.Field == "status" {
		err := c.handleStatus(connectionId, cl)
		if err != nil {
			return
		}
	}
	if cl.Field == "assignee" {
		err := c.handleAssignee(connectionId, cl)
		if err != nil {
			return
		}
	}
	if cl.Field != "Sprint" {
		return
	}
	from, to, err := c.parseFromTo(cl.From, cl.To)
	if err != nil {
		return
	}
	for sprintId := range from {
		err = c.handleFrom(connectionId, sprintId, cl)
		if err != nil {
			c.logger.Error("handle from error:", err)
			return
		}
	}
	for sprintId := range to {
		err = c.handleTo(connectionId, sprintId, cl)
		if err != nil {
			c.logger.Error("handle to error:", err)
			return
		}
	}
}

func (c *SprintIssuesConverter) CreateSprintIssue() error {
	var err error
	cache := make([]*ticket.SprintIssue, 0, BatchSize)
	for _, item := range c.sprintIssue {
		cache = append(cache, item)
		if len(cache) == BatchSize {
			err = c.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(cache).Error
			if err != nil {
				return err
			}
			cache = make([]*ticket.SprintIssue, 0, BatchSize)
		}
	}
	if len(cache) != 0 {
		err = c.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(cache).Error
		if err == nil {
			return err
		}
	}
	return nil
}

func (c *SprintIssuesConverter) parseFromTo(from, to string) (map[uint64]struct{}, map[uint64]struct{}, error) {
	fromInts := make(map[uint64]struct{})
	toInts := make(map[uint64]struct{})
	var n uint64
	var err error
	for _, item := range strings.Split(from, ",") {
		s := strings.TrimSpace(item)
		if s == "" {
			continue
		}
		n, err = strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, nil, err
		}
		fromInts[n] = struct{}{}
	}
	for _, item := range strings.Split(to, ",") {
		s := strings.TrimSpace(item)
		if s == "" {
			continue
		}
		n, err = strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, nil, err
		}
		toInts[n] = struct{}{}
	}
	inter := make(map[uint64]struct{})
	for k := range fromInts {
		if _, ok := toInts[k]; ok {
			inter[k] = struct{}{}
			delete(toInts, k)
		}
	}
	for k := range inter {
		delete(fromInts, k)
	}
	return fromInts, toInts, nil
}

func (c *SprintIssuesConverter) handleFrom(connectionId, sprintId uint64, cl ChangelogItemResult) error {
	if sprint, _ := c.getJiraSprint(connectionId, sprintId); sprint == nil {
		return nil
	}
	key := fmt.Sprintf("%d:%d:%d", connectionId, sprintId, cl.IssueId)
	if item, ok := c.sprintIssue[key]; ok {
		if item != nil && (item.RemovedDate == nil || item.RemovedDate != nil && item.RemovedDate.Before(cl.Created)) {
			item.RemovedDate = &cl.Created
			item.IsRemoved = true
		}
	} else {
		addedStage, _ := c.getStage(cl.Created, connectionId, sprintId)
		jiraIssue, _ := c.getJiraIssue(connectionId, cl.IssueId)
		sprint, _ := c.getJiraSprint(connectionId, sprintId)
		if sprint == nil {
			return nil
		}
		sprintIssue := &ticket.SprintIssue{
			SprintId:    c.sprintIdGen.Generate(connectionId, sprintId),
			IssueId:     c.issueIdGen.Generate(connectionId, cl.IssueId),
			AddedDate:   sprint.StartDate,
			AddedStage:  addedStage,
			RemovedDate: &cl.Created,
			IsRemoved:   true,
		}
		if jiraIssue != nil {
			sprintIssue.AddedDate = &jiraIssue.Created
			sprintIssue.AddedStage = getStage(jiraIssue.Created, sprint.StartDate, sprint.CompleteDate)
		}
		c.sprintIssue[key] = sprintIssue
	}
	k := fmt.Sprintf("%d:%d", sprintId, cl.IssueId)
	if item := c.sprintsHistory[k]; item != nil {
		item.EndDate = &cl.Created
		err := c.db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(item).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *SprintIssuesConverter) handleTo(connectionId, sprintId uint64, cl ChangelogItemResult) error {
	domainSprintId := c.sprintIdGen.Generate(connectionId, sprintId)
	key := fmt.Sprintf("%d:%d:%d", connectionId, sprintId, cl.IssueId)
	addedStage, err := c.getStage(cl.Created, connectionId, sprintId)
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	if addedStage == nil {
		return nil
	}
	if item, ok := c.sprintIssue[key]; ok {
		if item != nil && (item.AddedDate == nil || item.AddedDate != nil && item.AddedDate.After(cl.Created)) {
			item.AddedDate = &cl.Created
			item.AddedStage = addedStage
		}
	} else {
		c.sprintIssue[key] = &ticket.SprintIssue{
			SprintId:    domainSprintId,
			IssueId:     c.issueIdGen.Generate(connectionId, cl.IssueId),
			AddedDate:   &cl.Created,
			AddedStage:  addedStage,
			RemovedDate: nil,
		}
	}
	k := fmt.Sprintf("%d:%d", sprintId, cl.IssueId)
	now := time.Now()
	c.sprintsHistory[k] = &ticket.IssueSprintsHistory{
		IssueId:   c.issueIdGen.Generate(connectionId, cl.IssueId),
		SprintId:  domainSprintId,
		StartDate: cl.Created,
		EndDate:   &now,
	}
	return nil
}

func (c *SprintIssuesConverter) setupSprintIssue(connectionId, boardId uint64) error {
	cursor, err := c.db.Model(&models.JiraSprintIssue{}).
		Select("_tool_jira_sprint_issues.*").
		Joins("left join _tool_jira_board_sprints on _tool_jira_board_sprints.sprint_id = _tool_jira_sprint_issues.sprint_id").
		Where("_tool_jira_board_sprints.connection_id = ? AND _tool_jira_board_sprints.board_id = ?", connectionId, boardId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	for cursor.Next() {
		var jiraSprintIssue models.JiraSprintIssue
		err = c.db.ScanRows(cursor, &jiraSprintIssue)
		if err != nil {
			return err
		}
		sprint, _ := c.getJiraSprint(connectionId, jiraSprintIssue.SprintId)
		if sprint == nil {
			continue
		}
		key := fmt.Sprintf("%d:%d:%d", connectionId, jiraSprintIssue.SprintId, jiraSprintIssue.IssueId)
		dsi := ticket.SprintIssue{
			SprintId:  c.sprintIdGen.Generate(connectionId, jiraSprintIssue.SprintId),
			IssueId:   c.issueIdGen.Generate(connectionId, jiraSprintIssue.IssueId),
			AddedDate: jiraSprintIssue.IssueCreatedDate,
		}
		if dsi.AddedDate != nil {
			dsi.AddedStage = getStage(*dsi.AddedDate, sprint.StartDate, sprint.CompleteDate)
		}
		if jiraSprintIssue.ResolutionDate != nil {
			dsi.ResolvedStage = getStage(*jiraSprintIssue.ResolutionDate, sprint.StartDate, sprint.CompleteDate)
		}
		c.sprintIssue[key] = &dsi
	}
	return nil
}
func (c *SprintIssuesConverter) getJiraSprint(connectionId, sprintId uint64) (*models.JiraSprint, error) {
	key := fmt.Sprintf("%d:%d", connectionId, sprintId)
	if value, ok := c.sprints[key]; ok {
		return value, nil
	}
	var sprint models.JiraSprint
	err := c.db.First(&sprint, "connection_id = ? AND sprint_id = ?", connectionId, sprintId).Error
	if err != nil {
		return nil, err
	}
	c.sprints[key] = &sprint
	return &sprint, nil
}

func (c *SprintIssuesConverter) getJiraIssue(connectionId, issueId uint64) (*models.JiraIssue, error) {
	key := fmt.Sprintf("%d:%d", connectionId, issueId)
	if issue, ok := c.jiraIssue[key]; ok {
		return issue, nil
	}
	var jiraIssue models.JiraIssue
	err := c.db.First(&jiraIssue, "connection_id = ? AND issue_id = ?", connectionId, issueId).Error
	if err != nil {
		return nil, err
	}
	c.jiraIssue[key] = &jiraIssue
	return &jiraIssue, nil
}

func (c *SprintIssuesConverter) getStage(t time.Time, connectionId, sprintId uint64) (*string, error) {
	sprint, err := c.getJiraSprint(sprintId, connectionId)
	if err != nil {
		return nil, err
	}
	return getStage(t, sprint.StartDate, sprint.CompleteDate), nil
}

func (c *SprintIssuesConverter) handleStatus(connectionId uint64, cl ChangelogItemResult) error {
	var err error
	issueId := c.issueIdGen.Generate(connectionId, cl.IssueId)
	if statusHistory := c.status[issueId]; statusHistory != nil {
		statusHistory.EndDate = &cl.Created
		err = c.db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(c.status[issueId]).Error
		if err != nil {
			return err
		}
	}
	now := time.Now()
	c.status[issueId] = &ticket.IssueStatusHistory{
		IssueId:        issueId,
		OriginalStatus: cl.ToString,
		StartDate:      cl.Created,
		EndDate:        &now,
	}
	err = c.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(c.status[issueId]).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *SprintIssuesConverter) handleAssignee(connectionId uint64, cl ChangelogItemResult) error {
	issueId := c.issueIdGen.Generate(connectionId, cl.IssueId)
	if assigneeHistory := c.assignee[issueId]; assigneeHistory != nil {
		assigneeHistory.EndDate = &cl.Created
	}
	var assignee string
	if cl.To != "" {
		assignee = c.userIdGen.Generate(connectionId, cl.To)
	}
	now := time.Now()
	c.assignee[issueId] = &ticket.IssueAssigneeHistory{
		IssueId:   issueId,
		Assignee:  assignee,
		StartDate: cl.Created,
		EndDate:   &now,
	}
	err := c.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(c.assignee[issueId]).Error
	if err != nil {
		return err
	}
	return nil
}

func getStage(t time.Time, sprintStart, sprintComplete *time.Time) *string {
	if sprintStart == nil {
		return &ticket.BeforeSprint
	}
	if sprintStart.After(t) {
		return &ticket.BeforeSprint
	}
	if sprintComplete == nil {
		return &ticket.DuringSprint
	}
	if sprintComplete.Before(t) {
		return &ticket.AfterSprint
	}
	return &ticket.DuringSprint
}
