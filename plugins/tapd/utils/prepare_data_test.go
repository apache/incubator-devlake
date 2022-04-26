package utils

import (
	"fmt"
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"github.com/merico-dev/lake/plugins/tapd/tasks"
	"github.com/merico-dev/lake/runner"
	"net/url"
	"testing"
)

func TestGetURIStringPointer_WithSlash(t *testing.T) {
	cfg := config.GetConfig()
	db, err := runner.NewGormDb(cfg, nil)
	if err != nil {
		panic(err)
	}
	taskCtx := helper.NewDefaultTaskContext(cfg, nil, db, nil, "", nil, nil)

	source := &models.TapdSource{}
	err = db.Find(source, 1).Error
	if err != nil {
	}

	tapdApiClient, err := tasks.NewTapdApiClient(taskCtx, source)

	query := url.Values{}
	query.Set("workspace_id", "37469667")
	query.Set("name", "test1")
	query.Set("business_value", fmt.Sprintf("%v", 12))
	query.Set("creator", "陈映初")

	res, err := tapdApiClient.Do("POST", "stories", query, nil, nil)
	var b []byte
	res.Body.Read(b)

	fmt.Sprintf(string(b))
}
