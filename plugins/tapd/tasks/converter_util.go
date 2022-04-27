package tasks

import (
	"fmt"
	"time"

	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SprintIssuesConverter struct {
	db             *gorm.DB
	data           *TapdTaskData
	logger         core.Logger
	sprintIdGen    *didgen.DomainIdGenerator
	sprints        map[string]*ticket.Sprint
	sprintIssue    map[string]*ticket.SprintIssue
	status         map[string]*ticket.IssueStatusHistory
	assignee       map[string]*ticket.IssueAssigneeHistory
	sprintsHistory map[string]*ticket.IssueSprintsHistory
}

func NewChangelogToHistoryConverter(taskCtx core.SubTaskContext) *SprintIssuesConverter {
	return &SprintIssuesConverter{
		db:             taskCtx.GetDb(),
		logger:         taskCtx.GetLogger(),
		data:           taskCtx.GetData().(*TapdTaskData),
		sprintIdGen:    didgen.NewDomainIdGenerator(&models.TapdIteration{}),
		sprints:        make(map[string]*ticket.Sprint),
		sprintIssue:    make(map[string]*ticket.SprintIssue),
		status:         make(map[string]*ticket.IssueStatusHistory),
		assignee:       make(map[string]*ticket.IssueAssigneeHistory),
		sprintsHistory: make(map[string]*ticket.IssueSprintsHistory),
	}
}

func (c *SprintIssuesConverter) FeedIn(sourceId uint64, cl models.ChangelogTmp) {
	switch cl.FieldName {
	case "status":
		err := c.handleStatus(sourceId, cl)
		if err != nil {
			return
		}
	case "current_owner":
		err := c.handleAssignee(sourceId, cl)
		if err != nil {
			return
		}
	case "iteration_id":
		err := c.handlePrevIteration(sourceId, cl.IterationIdFrom, cl)
		if err != nil {
			c.logger.Error("handle from error:", err)
			return
		}
		err = c.handleNextIteration(sourceId, cl.IterationIdTo, cl)
		if err != nil {
			c.logger.Error("handle to error:", err)
			return
		}
	default:
		return
	}
}

func (c *SprintIssuesConverter) UpdateSprintIssue() error {
	var err error
	for _, fresh := range c.sprintIssue {
		err = c.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(fresh).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *SprintIssuesConverter) handlePrevIteration(sourceId, sprintId uint64, cl models.ChangelogTmp) error {
	domainSprintId := c.sprintIdGen.Generate(sourceId, sprintId)
	if sprint, _ := c.getSprint(domainSprintId); sprint == nil {
		return nil
	}
	key := fmt.Sprintf("%d:%d:%d", sourceId, sprintId, cl.IssueId)
	if item, ok := c.sprintIssue[key]; ok {
		if item != nil && (item.RemovedDate == nil || item.RemovedDate != nil && item.RemovedDate.Before(cl.CreatedDate)) {
			item.RemovedDate = &cl.CreatedDate
			item.IsRemoved = true
		}
	} else {
		c.sprintIssue[key] = &ticket.SprintIssue{
			SprintId:    domainSprintId,
			IssueId:     IssueIdGen.Generate(sourceId, cl.IssueId),
			AddedDate:   nil,
			RemovedDate: &cl.CreatedDate,
			IsRemoved:   true,
		}
	}
	k := fmt.Sprintf("%d:%d", sprintId, cl.IssueId)
	if item := c.sprintsHistory[k]; item != nil {
		item.EndDate = &cl.CreatedDate
		err := c.db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(item).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *SprintIssuesConverter) handleNextIteration(sourceId, sprintId uint64, cl models.ChangelogTmp) error {
	domainSprintId := c.sprintIdGen.Generate(sourceId, sprintId)
	if sprint, _ := c.getSprint(domainSprintId); sprint == nil {
		return nil
	}
	key := fmt.Sprintf("%d:%d:%d", sourceId, sprintId, cl.IssueId)
	addedStage, err := c.getStage(cl.CreatedDate, domainSprintId)
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
		if item != nil && (item.AddedDate == nil || item.AddedDate != nil && item.AddedDate.After(cl.CreatedDate)) {
			item.AddedDate = &cl.CreatedDate
			item.AddedStage = addedStage
		}
	} else {
		addedStage, _ := c.getStage(cl.CreatedDate, domainSprintId)
		c.sprintIssue[key] = &ticket.SprintIssue{
			SprintId:    domainSprintId,
			IssueId:     IssueIdGen.Generate(sourceId, cl.IssueId),
			AddedDate:   &cl.CreatedDate,
			AddedStage:  addedStage,
			RemovedDate: nil,
		}
	}
	k := fmt.Sprintf("%d:%d", sprintId, cl.IssueId)
	now := time.Now()
	c.sprintsHistory[k] = &ticket.IssueSprintsHistory{
		IssueId:   IssueIdGen.Generate(sourceId, cl.IssueId),
		SprintId:  domainSprintId,
		StartDate: cl.CreatedDate,
		EndDate:   &now,
	}
	return nil
}

func (c *SprintIssuesConverter) getSprint(id string) (*ticket.Sprint, error) {
	if value, ok := c.sprints[id]; ok {
		return value, nil
	}
	var sprint ticket.Sprint
	err := c.db.First(&sprint, "id = ?", id).Error
	if err != nil {
		c.sprints[id] = &sprint
	}
	return &sprint, err
}

func (c *SprintIssuesConverter) getStage(t time.Time, sprintId string) (*string, error) {
	sprint, err := c.getSprint(sprintId)
	if err != nil {
		return nil, err
	}
	return getStage(t, sprint.StartedDate, sprint.CompletedDate), nil
}

func (c *SprintIssuesConverter) handleStatus(sourceId uint64, cl models.ChangelogTmp) error {
	var err error
	issueId := IssueIdGen.Generate(sourceId, cl.IssueId)
	if statusHistory := c.status[issueId]; statusHistory != nil {
		statusHistory.EndDate = &cl.CreatedDate
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
		OriginalStatus: cl.From,
		StartDate:      cl.CreatedDate,
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

func (c *SprintIssuesConverter) handleAssignee(sourceId uint64, cl models.ChangelogTmp) error {
	issueId := IssueIdGen.Generate(sourceId, cl.IssueId)
	if assigneeHistory := c.assignee[issueId]; assigneeHistory != nil {
		assigneeHistory.EndDate = &cl.CreatedDate
	}
	var assignee string
	if cl.To != "" {
		assignee = UserIdGen.Generate(sourceId, c.data.Options.WorkspaceID, cl.To)
	}
	now := time.Now()
	c.assignee[issueId] = &ticket.IssueAssigneeHistory{
		IssueId:   issueId,
		Assignee:  assignee,
		StartDate: cl.CreatedDate,
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
