package tasks

import (
	"time"

	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

type TapdOptions struct {
	ConnectionId uint64   `json:"connectionId"`
	WorkspaceID  uint64   `json:"workspceId"`
	CompanyId    uint64   `json:"companyId"`
	Tasks        []string `json:"tasks,omitempty"`
	Since        string
}

type TapdTaskData struct {
	Options    *TapdOptions
	ApiClient  *helper.ApiAsyncClient
	Since      *time.Time
	Connection *models.TapdConnection
}
