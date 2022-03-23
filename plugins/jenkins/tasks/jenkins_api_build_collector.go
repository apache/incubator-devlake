package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jenkins/models"
)

const RAW_BUILD_TABLE = "jenkins_api_builds"

var CollectApiBuildsMeta = core.SubTaskMeta{
	Name:             "collectApiBuilds",
	EntryPoint:       CollectApiBuilds,
	EnabledByDefault: true,
	Description:      "Collect builds data from jenkins api",
}

type SimpleJob struct {
	Name string
}

func CollectApiBuilds(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*JenkinsTaskData)
	cursor, err := db.Model(&models.JenkinsJob{}).Rows()
	if err != nil {
		return err
	}
	iterator, err := helper.NewCursorIterator(db, cursor, reflect.TypeOf(SimpleJob{}))
	if err != nil {
		return err
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:   taskCtx,
			Table: RAW_BUILD_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Input:       iterator,
		UrlTemplate: "job/{{ .Input.Name }}/api/json",
		/*
			(Optional) Return query string for request, or you can plug them into UrlTemplate directly
		*/
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			treeValue := fmt.Sprintf(
				"allBuilds[number,timestamp,duration,estimatedDuration,displayName,result,actions[lastBuiltRevision[SHA1],mercurialRevisionNumber],changeSet[kind,revisions[revision]]]{%d,%d}",
				reqData.Pager.Skip, reqData.Pager.Skip+reqData.Pager.Size)
			query.Set("tree", treeValue)
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Builds []json.RawMessage `json:"allBuilds"`
			}
			err := helper.UnmarshalResponse(res, &data)
			if err != nil {
				return nil, err
			}
			return data.Builds, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
