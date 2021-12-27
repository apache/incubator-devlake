package tasks

import (
	"fmt"
	"math"
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

// issue types
const (
	Bug         = "Bug"
	Story       = "Story"
	Incident    = "Incident"
	Requirement = "Requirement"
)

const (
	BashSize = 100
)

var (
	UTCLocation, _ = time.LoadLocation("UTC")
)

type SprintIssueBurndownConverter struct {
	cache       map[string]map[int]*ticket.SprintIssueBurndown
	sprintIdGen *didgen.DomainIdGenerator
	issueIdGen  *didgen.DomainIdGenerator
	sprints     map[string]*models.JiraSprint
	sprintIssue map[string]*ticket.SprintIssue
}

func NewSprintIssueBurndownConverter() *SprintIssueBurndownConverter {
	return &SprintIssueBurndownConverter{
		cache:       make(map[string]map[int]*ticket.SprintIssueBurndown),
		sprintIdGen: didgen.NewDomainIdGenerator(&models.JiraSprint{}),
		issueIdGen:  didgen.NewDomainIdGenerator(&models.JiraIssue{}),
		sprints:     make(map[string]*models.JiraSprint),
		sprintIssue: make(map[string]*ticket.SprintIssue),
	}
}

func (c *SprintIssueBurndownConverter) FeedIn(sourceId uint64, cl ChangelogItemResult) {
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

func (c *SprintIssueBurndownConverter) UpdateSprintIssue() error {
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
func (c *SprintIssueBurndownConverter) Save() error {
	var err error
	for sprintId, sprint := range c.cache {
		err = lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).CreateInBatches(c.fill(sprintId, sprint), BatchSize).Error
		if err != nil {
			logger.Error("save sprint issue burndwon error:", err)
			return err
		}
	}
	return nil
}

func (c *SprintIssueBurndownConverter) fill(sprintId string, m map[int]*ticket.SprintIssueBurndown) []*ticket.SprintIssueBurndown {
	var result []*ticket.SprintIssueBurndown
	var max, min int
	min = math.MaxInt32
	for k := range m {
		if k > max {
			max = k
		}
		if k < min {
			min = k
		}
	}
	jiraSprint := c.sprints[sprintId]
	if jiraSprint != nil && jiraSprint.CompleteDate!= nil{
		max = c.getDateHour(*jiraSprint.CompleteDate)
	}
	for dateHour := min; dateHour <= max; dateHour = c.nextDateHour(dateHour) {
		if item, ok := m[dateHour]; ok {
			result = append(result, item)
		} else {
			result = append(result, c.newSprintIssueBurndown(sprintId, dateHour))
		}
	}

	// fill remaining
	var remain, remainBugs, remainRequirements, remainIncidents, remainStoryPoints int
	for _, item := range result {
		remain += item.Added
		remain -= item.Removed
		remainBugs += item.AddedBugs
		remainBugs -= item.RemovedBugs
		remainRequirements += item.AddedRequirements
		remainRequirements -= item.RemovedRequirements
		remainIncidents += item.AddedIncidents
		remainIncidents -= item.RemovedIncidents
		remainStoryPoints += item.AddedStoryPoints
		remainStoryPoints -= item.RemovedStoryPoints
		item.Remaining = remain
		item.RemainingBugs = remainBugs
		item.RemainingRequirements = remainRequirements
		item.RemainingIncidents = remainIncidents
		item.RemainingStoryPoints = remainStoryPoints
		item.RemainingOtherIssues = remain - remainBugs - remainIncidents - remainRequirements - remainStoryPoints
	}
	for p := len(result) - 1; p > -1; p-- {
		for i := 1; i < 24 && p-i > -1; i++ {
			result[p].Added += result[p-i].Added
			result[p].Removed += result[p-i].Removed
			result[p].AddedBugs += result[p-i].AddedBugs
			result[p].RemovedBugs += result[p-i].RemovedBugs
			result[p].AddedRequirements += result[p-i].AddedRequirements
			result[p].RemovedRequirements += result[p-i].RemovedRequirements
			result[p].AddedIncidents += result[p-i].AddedIncidents
			result[p].RemovedIncidents += result[p-i].RemovedIncidents
			result[p].AddedStoryPoints += result[p-i].AddedStoryPoints
			result[p].RemovedStoryPoints += result[p-i].RemovedStoryPoints
			result[p].AddedOtherIssues += result[p-i].AddedOtherIssues
			result[p].RemovedOtherIssues += result[p-i].RemovedOtherIssues
		}
	}
	return result
}
func (c *SprintIssueBurndownConverter) getJiraIssue(sourceId, issueId uint64) (*models.JiraIssue, error) {
	var issue models.JiraIssue
	err := lakeModels.Db.First(&issue, "issue_id = ? AND source_id = ?", issueId, sourceId).Error
	if err != nil {
		logger.Error("getJiraIssue error:", err)
		return nil, err
	}
	return &issue, err
}

func (c *SprintIssueBurndownConverter) parseFromTo(from, to string) (map[uint64]struct{}, map[uint64]struct{}, error) {
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

func (c *SprintIssueBurndownConverter) getDateHour(t time.Time) int {
	t = t.UTC().Add(time.Hour)
	y, m, d := t.Date()
	return y*1000000 + int(m)*10000 + d*100 + t.Hour()
}

func (c *SprintIssueBurndownConverter) dateHourEndPoints(dateHour int) (time.Time, time.Time) {
	y := dateHour / 1000000
	m := (dateHour / 10000) % 100
	d := (dateHour / 100) % 100
	h := dateHour % 100
	end := time.Date(y, time.Month(m), d, h, 0, 0, 0, UTCLocation)
	return end.Add(-24 * time.Hour), end
}

func (c *SprintIssueBurndownConverter) nextDateHour(dateHour int) int {
	y := dateHour / 1000000
	m := (dateHour / 10000) % 100
	d := (dateHour / 100) % 100
	h := dateHour % 100
	t := time.Date(y, time.Month(m), d, h, 0, 0, 0, UTCLocation).Add(time.Hour).UTC()
	y1, m1, d1 := t.Date()
	return y1*1000000 + int(m1)*10000 + d1*100 + t.Hour()
}

func (c *SprintIssueBurndownConverter) handleFrom(sourceId, sprintId uint64, cl ChangelogItemResult) error {
	domainSprintId := c.sprintIdGen.Generate(sourceId, sprintId)
	key := fmt.Sprintf("%d:%d:%d", sourceId, sprintId, cl.IssueId)
	if item, ok := c.sprintIssue[key]; ok {
		if item != nil && (item.RemovedDate == nil || item.RemovedDate != nil && item.RemovedDate.Before(cl.Created)) {
			item.RemovedDate = &cl.Created
		}
	} else {
		c.sprintIssue[key] = &ticket.SprintIssue{
			SprintId:    domainSprintId,
			IssueId:     c.issueIdGen.Generate(sourceId, cl.IssueId),
			AddedDate:   nil,
			RemovedDate: &cl.Created,
		}
	}
	dateHour := c.getDateHour(cl.Created)
	if _, ok := c.cache[domainSprintId]; !ok {
		c.cache[domainSprintId] = make(map[int]*ticket.SprintIssueBurndown)
		_, err := c.getSprint(sourceId, sprintId)
		if err != nil{
			return err
		}
	}
	if c.cache[domainSprintId][dateHour] == nil {
		c.cache[domainSprintId][dateHour] = c.newSprintIssueBurndown(domainSprintId, dateHour)
	}
	c.cache[domainSprintId][dateHour].Removed++
	jiraIssue, err := c.getJiraIssue(sourceId, cl.IssueId)
	if err != nil {
		return err
	}
	switch jiraIssue.StdType {
	case Bug:
		c.cache[domainSprintId][dateHour].RemovedBugs++
	case Incident:
		c.cache[domainSprintId][dateHour].RemovedIncidents++
	case Requirement:
		c.cache[domainSprintId][dateHour].RemovedRequirements++
	case Story:
		c.cache[domainSprintId][dateHour].RemovedStoryPoints++
	default:
		c.cache[domainSprintId][dateHour].RemovedOtherIssues++
	}
	return nil
}

func (c *SprintIssueBurndownConverter) handleTo(sourceId, sprintId uint64, cl ChangelogItemResult) error {
	domainSprintId := c.sprintIdGen.Generate(sourceId, sprintId)
	key := fmt.Sprintf("%d:%d:%d", sourceId, sprintId, cl.IssueId)
	if item, ok := c.sprintIssue[key]; ok {
		if item != nil && (item.AddedDate == nil || item.AddedDate != nil && item.AddedDate.After(cl.Created)) {
			item.AddedDate = &cl.Created
		}
	} else {
		c.sprintIssue[key] = &ticket.SprintIssue{
			SprintId:    domainSprintId,
			IssueId:     c.issueIdGen.Generate(sourceId, cl.IssueId),
			AddedDate:   &cl.Created,
			RemovedDate: nil,
		}
	}
	dateHour := c.getDateHour(cl.Created)
	if _, ok := c.cache[domainSprintId]; !ok {
		c.cache[domainSprintId] = make(map[int]*ticket.SprintIssueBurndown)
		_, err := c.getSprint(sourceId, sprintId)
		if err != nil{
			return err
		}
	}
	if c.cache[domainSprintId][dateHour] == nil {
		c.cache[domainSprintId][dateHour] = c.newSprintIssueBurndown(domainSprintId, dateHour)
	}
	c.cache[domainSprintId][dateHour].Added++
	jiraIssue, err := c.getJiraIssue(sourceId, cl.IssueId)
	if err != nil {
		return err
	}
	switch jiraIssue.StdType {
	case Bug:
		c.cache[domainSprintId][dateHour].AddedBugs++
	case Incident:
		c.cache[domainSprintId][dateHour].AddedIncidents++
	case Requirement:
		c.cache[domainSprintId][dateHour].AddedRequirements++
	case Story:
		c.cache[domainSprintId][dateHour].AddedStoryPoints++
	default:
		c.cache[domainSprintId][dateHour].AddedOtherIssues++
	}
	return nil
}

func (c *SprintIssueBurndownConverter) newSprintIssueBurndown(sprintId string, dateHour int) *ticket.SprintIssueBurndown {
	startedAt, endedAt := c.dateHourEndPoints(dateHour)
	return &ticket.SprintIssueBurndown{
		SprintId:    sprintId,
		StartedDate: startedAt,
		EndedDate:   endedAt,
		EndedHour:   dateHour,
	}
}

func (c *SprintIssueBurndownConverter) getSprint(sourceId, sprintId uint64) (*models.JiraSprint, error) {
	id := c.sprintIdGen.Generate(sourceId, sprintId)
	if value, ok := c.sprints[id]; ok{
		return value, nil
	}
	var sprint models.JiraSprint
	err := lakeModels.Db.First(&sprint, "source_id = ? AND sprint_id = ?", sourceId, sprintId).Error
	if err != nil{
		c.sprints[id] = &sprint
	}
	return &sprint, err
}
