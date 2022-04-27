package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/helper"
	jiraModels "github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/runner"
	"github.com/mitchellh/mapstructure"
	"gorm.io/datatypes"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type TapdOptions struct {
	SourceId    uint64   `json:"sourceId"`
	WorkspaceId uint64   `json:"workspceId"`
	CompanyId   uint64   `json:"companyId"`
	Tasks       []string `json:"tasks,omitempty"`
	Since       string
}

type Temp struct {
	SourceId                 uint64 `gorm:"primaryKey"`
	IssueId                  uint64 `gorm:"primarykey"`
	ProjectId                uint64
	Self                     string `gorm:"type:varchar(255)"`
	Key                      string `gorm:"type:varchar(255)"`
	Summary                  string
	Type                     string `gorm:"type:varchar(255)"`
	EpicKey                  string `gorm:"type:varchar(255)"`
	StatusName               string `gorm:"type:varchar(255)"`
	StatusKey                string `gorm:"type:varchar(255)"`
	StoryPoint               float64
	OriginalEstimateMinutes  int64  // user input?
	AggregateEstimateMinutes int64  // sum up of all subtasks?
	RemainingEstimateMinutes int64  // could it be negative value?
	CreatorAccountId         string `gorm:"type:varchar(255)"`
	CreatorAccountType       string `gorm:"type:varchar(255)"`
	CreatorDisplayName       string `gorm:"type:varchar(255)"`
	AssigneeAccountId        string `gorm:"type:varchar(255);comment:latest assignee"`
	AssigneeAccountType      string `gorm:"type:varchar(255)"`
	AssigneeDisplayName      string `gorm:"type:varchar(255)"`
	PriorityId               uint64
	PriorityName             string `gorm:"type:varchar(255)"`
	ParentId                 uint64
	ParentKey                string `gorm:"type:varchar(255)"`
	SprintId                 uint64 // latest sprint, issue might cross multiple sprints, would be addressed by #514
	SprintName               string `gorm:"type:varchar(255)"`
	ResolutionDate           *time.Time
	Created                  time.Time
	Updated                  time.Time `gorm:"index"`
	SpentMinutes             int64
	LeadTimeMinutes          uint
	StdStoryPoint            uint
	StdType                  string `gorm:"type:varchar(255)"`
	StdStatus                string `gorm:"type:varchar(255)"`
	AllFields                datatypes.JSONMap

	// internal status tracking
	ChangelogUpdated  *time.Time
	RemotelinkUpdated *time.Time
	common.NoPKModel
	IterId uint64
}

type TapdBugReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	WorkspaceId string `json:"workspace_id"`
	Reporter    string `json:"reporter"`
	Created     string `json:"created"`
	Status      string `json:"status"`
	Fixer       string `json:"fixer"`
	Priority    string `json:"priority"`
	Estimate    string `json:"estimate"`
	Resolved    string `json:"resolved"`
	IterationID string `json:"iteration_id"`
}

type TapdStoryReq struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	WorkspaceID     string `json:"workspace_id"`
	Creator         string `json:"creator"`
	Created         string `json:"created"`
	Modified        string `json:"modified"`
	Status          string `json:"status"`
	Owner           string `json:"owner"`
	Priority        string `json:"priority"`
	Effort          string `json:"effort"`
	EffortCompleted string `json:"effort_completed"`
	Exceed          string `json:"exceed"`
	Remain          string `json:"remain"`
	IterationID     string `json:"iteration_id"`
	Size            string `json:"size"`
	Developer       string `json:"developer"`
	Completed       string `json:"completed"`
}

func PrepareTestData(options TapdOptions) error {
	cfg := config.GetConfig()
	log := logger.Global.Nested("prepare data")
	db, err := runner.NewGormDb(cfg, log)
	if err != nil {
		return err
	}
	ctx, _ := context.WithCancel(context.Background())

	scheduler, _ := helper.NewWorkerScheduler(10000, 1, time.Second*1, ctx, 3)
	cursor, _ := db.Model(&jiraModels.JiraIssue{}).
		Joins("left join _tool_jira_sprint_issues on _tool_jira_sprint_issues.issue_id = _tool_jira_issues.issue_id").
		Joins("left join _tool_jira_sprints on _tool_jira_sprint_issues.sprint_id = _tool_jira_sprints.sprint_id").
		Joins("left join _tool_tapd_iterations on _tool_tapd_iterations.name = _tool_jira_sprints.name").Select("_tool_jira_issues.*, _tool_tapd_iterations.id as iter_id").
		Rows()
	defer cursor.Close()

	jiraIssueWithIter := Temp{}

	for cursor.Next() {
		_ = db.ScanRows(cursor, &jiraIssueWithIter)
		switch jiraIssueWithIter.StdType {
		case "故事":
			fallthrough
		case "疑问":
			fallthrough
		case "风险":
			fallthrough
		case "Tech Story":
			tapdIssue := TapdStoryReq{
				Name:            jiraIssueWithIter.Key,
				Description:     jiraIssueWithIter.Summary,
				WorkspaceID:     strconv.FormatUint(options.WorkspaceId, 10),
				Creator:         jiraIssueWithIter.CreatorDisplayName,
				Created:         jiraIssueWithIter.Created.String(),
				Modified:        jiraIssueWithIter.Updated.String(),
				Status:          jiraIssueWithIter.StatusKey,
				Owner:           jiraIssueWithIter.AssigneeDisplayName,
				Size:            strconv.Itoa(int(jiraIssueWithIter.StoryPoint)),
				IterationID:     strconv.FormatUint(jiraIssueWithIter.IterId, 10),
				Priority:        strconv.FormatUint(jiraIssueWithIter.PriorityId, 10),
				Developer:       jiraIssueWithIter.AssigneeDisplayName,
				Effort:          strconv.FormatInt(jiraIssueWithIter.OriginalEstimateMinutes/60, 10),
				EffortCompleted: strconv.FormatInt(jiraIssueWithIter.SpentMinutes/60, 10),
				Remain:          strconv.FormatInt(jiraIssueWithIter.RemainingEstimateMinutes/60, 10),
			}
			if jiraIssueWithIter.ResolutionDate != nil {
				tapdIssue.Completed = jiraIssueWithIter.ResolutionDate.String()
			}
			switch jiraIssueWithIter.StdStatus {
			case "DONE":
				tapdIssue.Status = "done"
			case "TODO":
				tapdIssue.Status = "open"
			case "IN_PROGRESS":
				tapdIssue.Status = "progressing"
			}
			switch jiraIssueWithIter.PriorityId {
			case 5:
				tapdIssue.Priority = "1"
			case 4:
				tapdIssue.Priority = "2"
			case 3:
				tapdIssue.Priority = "3"
			case 2:
				tapdIssue.Priority = "4"
			case 1:
				tapdIssue.Priority = "5"
			}
			jsonstr, _ := json.Marshal(tapdIssue)
			scheduler.Submit(func() error {
				httpPostJson("stories", jsonstr)
				return nil
			})
		case "任务":
			fallthrough
		case "测试任务":
			fallthrough
		case "测试开发任务":
			fallthrough
		case "子任务":
			tapdIssue := TapdStoryReq{
				Name:            jiraIssueWithIter.Key,
				Description:     jiraIssueWithIter.Summary,
				WorkspaceID:     strconv.FormatUint(options.WorkspaceId, 10),
				Creator:         jiraIssueWithIter.CreatorDisplayName,
				Created:         jiraIssueWithIter.Created.String(),
				Modified:        jiraIssueWithIter.Updated.String(),
				Status:          jiraIssueWithIter.StatusKey,
				Owner:           jiraIssueWithIter.AssigneeDisplayName,
				Priority:        strconv.FormatUint(jiraIssueWithIter.PriorityId, 10),
				Effort:          strconv.FormatInt(jiraIssueWithIter.OriginalEstimateMinutes/60, 10),
				EffortCompleted: strconv.FormatInt(jiraIssueWithIter.SpentMinutes/60, 10),
				IterationID:     strconv.FormatUint(jiraIssueWithIter.IterId, 10),
				Remain:          strconv.FormatInt(jiraIssueWithIter.RemainingEstimateMinutes/60, 10),
			}
			if jiraIssueWithIter.ResolutionDate != nil {
				tapdIssue.Completed = jiraIssueWithIter.ResolutionDate.String()
			}
			switch jiraIssueWithIter.StdStatus {
			case "DONE":
				tapdIssue.Status = "resolved"
			case "TODO":
				tapdIssue.Status = "planning"
			case "IN_PROGRESS":
				tapdIssue.Status = "developing"
			}
			switch jiraIssueWithIter.PriorityId {
			case 5:
				tapdIssue.Priority = "1"
			case 4:
				tapdIssue.Priority = "2"
			case 3:
				tapdIssue.Priority = "3"
			case 2:
				tapdIssue.Priority = "4"
			case 1:
				tapdIssue.Priority = "5"
			}
			jsonstr, _ := json.Marshal(tapdIssue)
			scheduler.Submit(func() error {
				httpPostJson("tasks", jsonstr)
				return nil
			})
		case "缺陷":
			fallthrough
		case "线上事故":
			tapdIssue := TapdBugReq{
				Title:       jiraIssueWithIter.Key,
				Description: jiraIssueWithIter.Summary,
				WorkspaceId: strconv.FormatUint(options.WorkspaceId, 10),
				Reporter:    jiraIssueWithIter.CreatorDisplayName,
				Created:     jiraIssueWithIter.Created.String(),
				Status:      jiraIssueWithIter.StatusKey,
				Fixer:       jiraIssueWithIter.AssigneeDisplayName,
				Priority:    strconv.FormatUint(jiraIssueWithIter.PriorityId, 10),
				Estimate:    strconv.FormatInt(jiraIssueWithIter.OriginalEstimateMinutes/60, 10),
				IterationID: strconv.FormatUint(jiraIssueWithIter.IterId, 10),
			}
			if jiraIssueWithIter.ResolutionDate != nil {
				tapdIssue.Resolved = jiraIssueWithIter.ResolutionDate.String()
			}
			switch jiraIssueWithIter.PriorityId {
			case 5:
				tapdIssue.Priority = "insignificant"
			case 4:
				tapdIssue.Priority = "low"
			case 3:
				tapdIssue.Priority = "medium"
			case 2:
				tapdIssue.Priority = "high"
			case 1:
				tapdIssue.Priority = "urgent"
			}
			switch jiraIssueWithIter.StdStatus {
			case "DONE":
				tapdIssue.Status = "resolved"
			case "TODO":
				tapdIssue.Status = "planning"
			case "IN_PROGRESS":
				tapdIssue.Status = "developing"
			}
			jsonstr, _ := json.Marshal(tapdIssue)
			scheduler.Submit(func() error {
				httpPostJson("bugs", jsonstr)
				return nil
			})
		}
	}
	return nil
}

func httpPostJson(path string, jsonstr []byte) error {
	url := fmt.Sprintf("https://api.tapd.cn/%s", path)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonstr))
	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic eWtQWWdqQVA6MDY2RkVEQTItNzMxNS1DMUExLUZBMjEtMzhDQkE3OTVCODgx")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	statuscode := resp.StatusCode
	if statuscode != 200 {
		fmt.Println(statuscode)
	}
	head := resp.Header
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	fmt.Println(statuscode)
	fmt.Println(head)
	return nil

}

// standalone mode for debugging
func main() {
	options := map[string]interface{}{
		"sourceId":    1,
		"companyId":   55850509,
		"workspaceId": 62390899,
	}

	var op TapdOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		panic(err)
	}
	PrepareTestData(op)
}
