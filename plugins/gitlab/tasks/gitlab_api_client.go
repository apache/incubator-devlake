package tasks

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/merico-dev/lake/config"
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

func getTotal(resourceUriFormat string) (int, error) {
	// jsut get the first page of results. The response has a head that tells the total pages
	page := 0
	page_size := 1
	res, err := gitlabApiClient.Get(fmt.Sprintf(resourceUriFormat, page_size, page), nil, nil)

	if err != nil {
		return 0, err
	}

	total := res.Header.Get("X-Total")
	totalInt, err := convertStringToInt(total)
	if err != nil {
		return 0, err
	}
	return totalInt, nil
}

func convertStringToInt(input string) (int, error) {
	return strconv.Atoi(input)
}

// run all requests in an Ants worker pool
func (gitlabApiClient *GitlabApiClient) FetchWithPaginationAnts(resourceUri string, pageSize string, handler GitlabPaginationHandler) error {

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

	// We need to get the total pages first so we can loop through all requests concurrently
	total, err := getTotal(resourceUriFormat)

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

			handler(res)
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
func (gitlabApiClient *GitlabApiClient) FetchWithPagination(resourceUri string, pageSize string, handler GitlabPaginationHandler) error {

	pageSizeInt, _ := convertStringToInt(pageSize)

	var resourceUriFormat string
	if strings.ContainsAny(resourceUri, "?") {
		resourceUriFormat = resourceUri + "&per_page=%v&page=%v"
	} else {
		resourceUriFormat = resourceUri + "?per_page=%v&page=%v"
	}

	// We need to get the total pages first so we can loop through all requests concurrently
	total, _ := getTotal(resourceUriFormat)

	// Loop until all pages are requested
	for i := 0; (i * pageSizeInt) < total; i++ {
		// we need to save the value for the request so it is not overwritten
		currentPage := i
		url := fmt.Sprintf(resourceUriFormat, pageSizeInt, currentPage)

		res, err := gitlabApiClient.Get(url, nil, nil)

		if err != nil {
			return err
		}

		handler(res)
	}

	return nil
}
