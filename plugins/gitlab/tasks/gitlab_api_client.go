package tasks

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/utils"
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



func convertStringToInt(input string) (int, error) {
	return strconv.Atoi(input)
}

// run all requests in an Ants worker pool
func (gitlabApiClient *GitlabApiClient) FetchWithPaginationAnts(resourceUri string, pageSize string, total int, handler GitlabPaginationHandler) error {

	if total <= 0 {
		logger.Error("You failed to send a total to FetchWithPagination", total)
		return nil
	}
	
	pageSizeInt, _ := convertStringToInt(pageSize)

	var resourceUriFormat string
	if strings.ContainsAny(resourceUri, "?") {
		resourceUriFormat = resourceUri + "&per_page=%v&page=%v"
	} else {
		resourceUriFormat = resourceUri + "?per_page=%v&page=%v"
	}
	rateLimitPerSecond := 2000 / 60
	workerNum := 50
	// set up the worker pool
	scheduler, err := utils.NewWorkerScheduler(workerNum, rateLimitPerSecond)
	if err != nil {
		return err
	}

	defer scheduler.Release()

	// // We need to get the total pages first so we can loop through all requests concurrently
	// // if total was not set in the method
	// total, err = getTotal(resourceUriFormat)

	// Loop until all pages are requested
	for i := 0; (i * pageSizeInt) <= (total + pageSizeInt); i++ {
		// we need to save the value for the request so it is not overwritten
		currentPage := i
		err1 := scheduler.Submit(func() error {
			url := fmt.Sprintf(resourceUriFormat, pageSizeInt, currentPage)

			res, err := gitlabApiClient.Get(url, nil, nil)

			if err != nil {
				return err
			}

			handlerErr := handler(res)
			if handlerErr != nil {
				return handlerErr
			}
			return nil
		})

		if err1 != nil {
			return err
		}
	}

	scheduler.WaitUntilFinish()

	return nil
}

// fetch paginated without ANTS worker pool
func (gitlabApiClient *GitlabApiClient) FetchWithPagination(resourceUri string, pageSize string, total int, handler GitlabPaginationHandler) error {

	if total <= 0 {
		logger.Error("You failed to send a total to FetchWithPagination", total)
		return nil
	}

	pageSizeInt, _ := convertStringToInt(pageSize)

	var resourceUriFormat string
	if strings.ContainsAny(resourceUri, "?") {
		resourceUriFormat = resourceUri + "&per_page=%v&page=%v"
	} else {
		resourceUriFormat = resourceUri + "?per_page=%v&page=%v"
	}

	// We need to get the total pages first so we can loop through all requests concurrently
	// total, _ := getTotal(resourceUriFormat)

	// Loop until all pages are requested
	for i := 0; (i * pageSizeInt) < total; i++ {
		// we need to save the value for the request so it is not overwritten
		currentPage := i
		url := fmt.Sprintf(resourceUriFormat, pageSizeInt, currentPage)

		res, err := gitlabApiClient.Get(url, nil, nil)

		if err != nil {
			return err
		}

		handlerErr := handler(res)
		if handlerErr != nil {
			return handlerErr
		}
	}

	return nil
}
