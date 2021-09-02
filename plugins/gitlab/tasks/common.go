package tasks

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
)

type GitlabApiClient struct {
	core.ApiClient
}

var gitlabApiClient *GitlabApiClient

func CreateApiClient() *GitlabApiClient {
	if gitlabApiClient == nil {
		gitlabApiClient = &GitlabApiClient{}
		gitlabApiClient.Setup(
			config.V.GetString("GITLAB_ENDPOINT"),
			map[string]string{
				"Authorization": fmt.Sprintf("Bearer %v", config.V.GetString("GITLAB_AUTH")),
			},
			10*time.Second,
			3,
		)
	}
	return gitlabApiClient
}

type GitlabPaginationHandler func(res *http.Response) error

func (gitlabApiClient *GitlabApiClient) FetchWithPagination(resourceUri string, page string, pageSize string, handler GitlabPaginationHandler) error {

	var resourceUriFormat string
	if strings.ContainsAny(resourceUri, "?") {
		resourceUriFormat = resourceUri + "&per_page=%v&page=%v"
	} else {
		resourceUriFormat = resourceUri + "?per_page=%v&page=%v"
	}

	for page != "" {
		res, err := gitlabApiClient.Get(fmt.Sprintf(resourceUriFormat, pageSize, page), nil, nil)

		if err != nil {
			return err
		}
		page = res.Header.Get("x-next-page")
		handler(res)
	}
	return nil
}
