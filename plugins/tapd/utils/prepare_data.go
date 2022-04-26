package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
	jiraModels "github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"github.com/merico-dev/lake/plugins/tapd/tasks"
	"gorm.io/datatypes"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var _ core.SubTaskEntryPoint = PrepareTestData

var PrepareTestDataMeta = core.SubTaskMeta{
	Name:             "PrepareTestData",
	EntryPoint:       PrepareTestData,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

type Temp struct {
	// collected fields
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
	Title       string
	Description string
	WorkspaceId string
	Reporter    string
	Created     string
	Status      string
	Fixer       string
	Priority    string
	Estimate    string
	Resolved    string
}

type TapdStoryReq struct {
	Name            string
	Description     string
	WorkspaceId     string
	Creator         string
	Created         string
	Modified        string
	Status          string
	Owner           string
	Priority        string
	Effort          string
	EffortCompleted string
	IterationID     string
	Remain          string
	Size            string
	Developer       string
	Completed       string
}

func PrepareTestData(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	source := &models.TapdSource{}
	err := db.Find(source, 1).Error
	if err != nil {
	}
	data := taskCtx.GetData().(*tasks.TapdTaskData)
	cursor, err := db.Model(&jiraModels.JiraIssue{}).
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
				WorkspaceId:     strconv.FormatUint(data.Source.WorkspaceId, 10),
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
			httpPostJson("stories", jsonstr)
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
				WorkspaceId:     strconv.FormatUint(data.Source.WorkspaceId, 10),
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
			httpPostJson("tasks", jsonstr)
		case "缺陷":
			fallthrough
		case "线上事故":
			tapdIssue := TapdBugReq{
				Title:       jiraIssueWithIter.Key,
				Description: jiraIssueWithIter.Summary,
				WorkspaceId: strconv.FormatUint(data.Source.WorkspaceId, 10),
				Reporter:    jiraIssueWithIter.CreatorDisplayName,
				Created:     jiraIssueWithIter.Created.String(),
				Status:      jiraIssueWithIter.StatusKey,
				Fixer:       jiraIssueWithIter.AssigneeDisplayName,
				Priority:    strconv.FormatUint(jiraIssueWithIter.PriorityId, 10),
				Estimate:    strconv.FormatInt(jiraIssueWithIter.OriginalEstimateMinutes/60, 10),
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
			httpPostJson("bugs", jsonstr)
		}

	}
	return nil
}

func httpPostJson(path string, jsonstr []byte) {
	url := fmt.Sprintf("https://api.tapd.cn/%s", path)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonstr))
	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic eWtQWWdqQVA6MDY2RkVEQTItNzMxNS1DMUExLUZBMjEtMzhDQkE3OTVCODgx")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	statuscode := resp.StatusCode
	if statuscode != 200 {
		fmt.Println(statuscode)
	}
	hea := resp.Header
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	fmt.Println(statuscode)
	fmt.Println(hea)

}
