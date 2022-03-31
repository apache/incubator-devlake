package tasks

import (
	"time"

	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

type TapdOptions struct {
	SourceId   uint64   `json:"sourceId"`
	WorkspceId uint64   `json:"workspceId"`
	CompanyId  uint64   `json:"companyId"`
	Tasks      []string `json:"tasks,omitempty"`
	Since      string
}

type TapdTaskData struct {
	Options   *TapdOptions
	ApiClient *helper.ApiAsyncClient
	Since     *time.Time
	Source    *models.TapdSource
}
