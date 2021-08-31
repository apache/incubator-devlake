package tasks

import (
	"fmt"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
)

func CreateApiClient() *core.ApiClient {
	return core.NewApiClient(
		config.V.GetString("GITLAB_ENDPOINT"),
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %v", config.V.GetString("GITLAB_AUTH")),
		},
		10*time.Second,
		3,
	)
}
