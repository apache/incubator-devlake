package tasks

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SprintIssuesConverter struct {
	sprintIdGen    *didgen.DomainIdGenerator
	issueIdGen     *didgen.DomainIdGenerator
	userIdGen      *didgen.DomainIdGenerator
	sprints        map[string]*ticket.Sprint
	sprintIssue    map[string]*ticket.SprintIssue
	status         map[string]*ticket.IssueStatusHistory
	assignee       map[string]*ticket.IssueAssigneeHistory
	sprintsHistory map[string]*ticket.IssueSprintsHistory
}

func NewSprintIssueConverter() *SprintIssuesConverter {
	return &SprintIssuesConverter{
		sprintIdGen:    didgen.NewDomainIdGenerator(&models.JiraSprint{}),
		issueIdGen:     didgen.NewDomainIdGenerator(&models.JiraIssue{}),
		userIdGen:      didgen.NewDomainIdGenerator(&models.JiraUser{}),
		sprints:        make(map[string]*ticket.Sprint),
		sprintIssue:    make(map[string]*ticket.SprintIssue),
		status:         make(map[string]*ticket.IssueStatusHistory),
		assignee:       make(map[string]*ticket.IssueAssigneeHistory),
		sprintsHistory: make(map[string]*ticket.IssueSprintsHistory),
	}
}

func (c *SprintIssuesConverter) FeedIn(sourceId uint64, cl ChangelogItemResult) {
	if cl.Field == "status" {
		err := c.handleStatus(sourceId, cl)
		if err != nil {
			return
		}
	}
	if cl.Field == "assignee" {
		err := c.handleAssignee(sourceId, cl)
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
		err = c.handleFrom(sourceId, sprintId, cl)
		if err != nil {
			logger.Error("handle from error:", err)
			return
		}
	}
	for sprintId := range to {
		err = c.handleTo(sourceId, sprintId, cl)
		if err != nil {
			logger.Error("handle to error:", err)
			return
		}
	}
}

func (c *SprintIssuesConverter) UpdateSprintIssue() error {
	var err error
	var flag bool
	var list []*ticket.SprintIssue
	for _, fresh := range c.sprintIssue {
		var old ticket.SprintIssue
		err = lakeModels.Db.First(&old, "sprint_id = ? AND issue_id = ?", fresh.SprintId, fresh.IssueId).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.Error("UpdateSprintIssue error:", err)
			return err
		}
		var issue ticket.Issue
		err = lakeModels.Db.First(&issue, "id = ?", fresh.IssueId).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.Error("UpdateSprintIssue error:", err)
			return err
		}
		if issue.ResolutionDate != nil {
			fresh.ResolvedStage, _ = c.getStage(*issue.ResolutionDate, fresh.IssueId)
			if fresh.ResolvedStage != old.ResolvedStage {
				flag = true
			}
		}
		if old.AddedDate == nil && fresh.AddedDate != nil || old.RemovedDate == nil && fresh.RemovedDate != nil {
			flag = true
		}
		if old.AddedDate != nil && fresh.AddedDate != nil && old.AddedDate.Before(*fresh.AddedDate) {
			fresh.AddedDate = old.AddedDate
			flag = true
		}
		if old.RemovedDate != nil && fresh.RemovedDate != nil && old.RemovedDate.After(*fresh.RemovedDate) {
			fresh.RemovedDate = old.RemovedDate
			flag = true
		}
		if fresh.AddedDate != nil && fresh.RemovedDate != nil {
			fresh.IsRemoved = fresh.AddedDate.Before(*fresh.RemovedDate)
			if fresh.IsRemoved != old.IsRemoved {
				flag = true
			}
		}
		if flag {
			list = append(list, fresh)
		}
	}
	return lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(list, BatchSize).Error
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

func (c *SprintIssuesConverter) handleFrom(sourceId, sprintId uint64, cl ChangelogItemResult) error {
	domainSprintId := c.sprintIdGen.Generate(sourceId, sprintId)
	key := fmt.Sprintf("%d:%d:%d", sourceId, sprintId, cl.IssueId)
	if item, ok := c.sprintIssue[key]; ok {
		if item != nil && (item.RemovedDate == nil || item.RemovedDate != nil && item.RemovedDate.Before(cl.Created)) {
			item.RemovedDate = &cl.Created
			item.IsRemoved = true
		}
	} else {
		c.sprintIssue[key] = &ticket.SprintIssue{
			SprintId:    domainSprintId,
			IssueId:     c.issueIdGen.Generate(sourceId, cl.IssueId),
			AddedDate:   nil,
			RemovedDate: &cl.Created,
			IsRemoved:   true,
		}
	}
	k := fmt.Sprintf("%d:%d", sprintId, cl.IssueId)
	if item := c.sprintsHistory[k]; item != nil {
		item.EndDate = &cl.Created
		err := lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(item).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *SprintIssuesConverter) handleTo(sourceId, sprintId uint64, cl ChangelogItemResult) error {
	domainSprintId := c.sprintIdGen.Generate(sourceId, sprintId)
	key := fmt.Sprintf("%d:%d:%d", sourceId, sprintId, cl.IssueId)
	addedStage, err := c.getStage(cl.Created, domainSprintId)
	if err == gorm.ErrRecordNotFound{
		return nil
	}
	if err != nil{
		return err
	}
	if addedStage == ""{
		return nil
	}
	if item, ok := c.sprintIssue[key]; ok {
		if item != nil && (item.AddedDate == nil || item.AddedDate != nil && item.AddedDate.After(cl.Created)) {
			item.AddedDate = &cl.Created
			item.AddedStage = addedStage
		}
	} else {
		addedStage, _ := c.getStage(cl.Created, domainSprintId)
		c.sprintIssue[key] = &ticket.SprintIssue{
			SprintId:    domainSprintId,
			IssueId:     c.issueIdGen.Generate(sourceId, cl.IssueId),
			AddedDate:   &cl.Created,
			AddedStage:  addedStage,
			RemovedDate: nil,
		}
	}
	k := fmt.Sprintf("%d:%d", sprintId, cl.IssueId)
	now := time.Now()
	c.sprintsHistory[k] = &ticket.IssueSprintsHistory{
		IssueId:   c.issueIdGen.Generate(sourceId, cl.IssueId),
		SprintId:  domainSprintId,
		StartDate: cl.Created,
		EndDate:   &now,
	}
	return nil
}

func (c *SprintIssuesConverter) getSprint(id string) (*ticket.Sprint, error) {
	if value, ok := c.sprints[id]; ok {
		return value, nil
	}
	var sprint ticket.Sprint
	err := lakeModels.Db.First(&sprint, "id = ?", id).Error
	if err != nil {
		c.sprints[id] = &sprint
	}
	return &sprint, err
}

func (c *SprintIssuesConverter) getStage(t time.Time, sprintId string) (string, error) {
	sprint, err := c.getSprint(sprintId)
	if err != nil {
		return "", err
	}
	return getStage(t, sprint.StartedDate, sprint.CompletedDate), nil
}

func (c *SprintIssuesConverter) handleStatus(sourceId uint64, cl ChangelogItemResult) error {
	var err error
	issueId := c.issueIdGen.Generate(sourceId, cl.IssueId)
	if statusHistory := c.status[issueId]; statusHistory != nil {
		statusHistory.EndDate = &cl.Created
		err = lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(c.status[issueId]).Error
		if err != nil{
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
	err = lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(c.status[issueId]).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *SprintIssuesConverter) handleAssignee(sourceId uint64, cl ChangelogItemResult) error {
	issueId := c.issueIdGen.Generate(sourceId, cl.IssueId)
	if assigneeHistory := c.assignee[issueId]; assigneeHistory != nil {
		assigneeHistory.EndDate = &cl.Created
	}
	var assignee string
	if cl.To != "" {
		assignee = c.userIdGen.Generate(sourceId, cl.To)
	}
	c.assignee[issueId] = &ticket.IssueAssigneeHistory{
		IssueId:   issueId,
		Assignee:  assignee,
		StartDate: cl.Created,
	}
	err := lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(c.assignee[issueId]).Error
	if err != nil {
		return err
	}
	return nil
}

func getStage(t time.Time, sprintStart, sprintComplete *time.Time) string {
	if sprintStart != nil {
		if sprintStart.After(t) {
			return ticket.BeforeSprint
		}
		if sprintStart.Equal(t) || (sprintComplete != nil && sprintComplete.Equal(t)) {
			return ticket.DuringSprint
		}
		if sprintComplete != nil && sprintStart.Before(t) && sprintComplete.After(t) {
			return ticket.DuringSprint
		}
	}
	if sprintComplete != nil && sprintComplete.Before(t) {
		return ticket.AfterSprint
	}
	return ""
}
