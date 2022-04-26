package utils

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	jiraModels "github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"github.com/merico-dev/lake/plugins/tapd/tasks"
	"strconv"
)

var _ core.SubTaskEntryPoint = PrepareSprintTestData

var PrepareSprintTestDataMeta = core.SubTaskMeta{
	Name:             "prepareSprintTestData",
	EntryPoint:       PrepareSprintTestData,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

type TapdIterationRes struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	WorkspaceId  string `json:"workspace_id"`
	Startdate    string `json:"startdate"`
	Enddate      string `json:"enddate"`
	Status       string `json:"status"`
	ReleaseID    string `json:"release_id"`
	Description  string `json:"description"`
	Creator      string `json:"creator"`
	Created      string `json:"created"`
	Modified     string `json:"modified"`
	Completed    string `json:"completed"`
	Releaseowner string `json:"releaseowner"`
	Launchdate   string `json:"launchdate"`
	Notice       string `json:"notice"`
	Releasename  string `json:"releasename"`
}

func PrepareSprintTestData(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	source := &models.TapdSource{}
	err := db.Find(source, 1).Error
	if err != nil {
	}
	data := taskCtx.GetData().(*tasks.TapdTaskData)
	jiraSprint := jiraModels.JiraSprint{}
	cursorSprint, err := db.Model(&jiraModels.JiraSprint{}).Rows()
	defer cursorSprint.Close()
	for cursorSprint.Next() {
		_ = db.ScanRows(cursorSprint, &jiraSprint)
		tapdIter := TapdIterationRes{
			Name:        jiraSprint.Name,
			WorkspaceId: strconv.FormatUint(data.Source.WorkspaceId, 10),
			Creator:     "陈映初",
		}
		//if jiraSprint.CompleteDate != nil {
		//	tapdIter.Completed = jiraSprint.CompleteDate.String()
		//}
		if jiraSprint.StartDate != nil {
			tapdIter.Startdate = jiraSprint.StartDate.String()
			tapdIter.Created = jiraSprint.StartDate.String()
		}
		if jiraSprint.EndDate != nil {
			tapdIter.Enddate = jiraSprint.EndDate.String()
		}
		if jiraSprint.State == "closed" {
			tapdIter.Status = "done"
		} else {
			tapdIter.Status = "open"
		}

		jsonstr, _ := json.Marshal(tapdIter)
		httpPostJson("iterations", jsonstr)
	}

	return nil
}
