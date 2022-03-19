package tasks

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/utils"
	"net/http"
	"time"
)

type JenkinsApiClient struct {
	core.ApiClient
}

func NewJenkinsApiClient(endpoint string, username string, password string, proxy string, ctx context.Context, scheduler *utils.WorkerScheduler, logger core.Logger) *JenkinsApiClient {
	jenkinsApiClient := &JenkinsApiClient{}
	encodedToken := utils.GetEncodedToken(username, password)

	jenkinsApiClient.Setup(
		endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", encodedToken),
		},
		50*time.Second,
		3,
		scheduler,
	)

	jenkinsApiClient.SetAfterFunction(func(res *http.Response) error {
		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("authentication failed, please check your Basic Auth Token")
		}
		return nil
	})
	if ctx != nil {
		jenkinsApiClient.SetContext(ctx)
	}
	if proxy != "" {
		err := jenkinsApiClient.SetProxy(proxy)
		if err != nil {
			panic(err)
		}
	}

	jenkinsApiClient.SetLogger(logger)
	return jenkinsApiClient
}

type JenkinsPaginationHandler func(res *http.Response) error
type JenkinsSearchPaginationHandler func(res *http.Response) (int, error)
