package tasks

import (
	"time"

	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

type JiraOptions struct {
	ConnectionId uint64   `json:"connectionId"`
	BoardId      uint64   `json:"boardId"`
	Tasks        []string `json:"tasks,omitempty"`
	Since        string
}

type JiraTaskData struct {
	Options        *JiraOptions
	ApiClient      *helper.ApiAsyncClient
	Connection     *models.JiraConnection
	Since          *time.Time
	JiraServerInfo models.JiraServerInfo
}
