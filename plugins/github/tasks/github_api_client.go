package tasks

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	githubUtils "github.com/merico-dev/lake/plugins/github/utils"
	"github.com/merico-dev/lake/utils"
)

type GithubApiClient struct {
	core.ApiClient
}

var githubApiClient *GithubApiClient

func CreateApiClient() *GithubApiClient {
	if githubApiClient == nil {
		githubApiClient = &GithubApiClient{}
		auth := fmt.Sprintf("Bearer %v", config.V.GetString("GITHUB_AUTH"))
		githubApiClient.Setup(
			config.V.GetString("GITHUB_ENDPOINT"),
			map[string]string{
				"Authorization": auth,
			},
			10*time.Second,
			3,
		)
	}
	return githubApiClient
}

type GithubPaginationHandler func(res *http.Response) error

// run all requests in an Ants worker pool
// conc - number of concurent requests you want to run
func (githubApiClient *GithubApiClient) FetchWithPaginationAnts(resourceUri string, pageSize int, conc int, scheduler *utils.WorkerScheduler, handler GithubPaginationHandler) error {
	url := AddPagingQueryToUrl(resourceUri)
	err := RunConcurrently(url, pageSize, conc, scheduler, handler)
	if err != nil {
		logger.Error("runConcurrently() failed", true)
	}
	return nil
}

func AddPagingQueryToUrl(url string) string {
	var urlWithPagingQuery string
	if strings.ContainsAny(url, "?") {
		urlWithPagingQuery = url + "&page=%v&per_page=%v"
	} else {
		urlWithPagingQuery = url + "?page=%v&per_page=%v"
	}
	return urlWithPagingQuery
}

// This method exists in the case where we do not know how many pages of data we have to fetch
// This loops through the data in chunks of `conc` and if there is any request in there with no data returned, we assume we are at the end of the data required to fetch
// This is needed since we do not want to make a request to get the paging details first since the rate limit for github is so low
func RunConcurrently(resourceUriFormat string, pageSize int, conc int, scheduler *utils.WorkerScheduler, handler GithubPaginationHandler) error {

	if conc == 0 {
		logger.Error("you must send a conc count to RunConcurrently()", true)
	}

	// How many requests would you like to send at once in chunks
	step := 0
	c := make(chan bool)
	for {
		for i := conc; i > 0; i-- {
			page := step*conc + i
			err := scheduler.Submit(func() error {
				url := fmt.Sprintf(resourceUriFormat, page, pageSize)
				res, err := githubApiClient.Get(url, nil, nil)
				if err != nil {
					return err
				}
				linkHeader := res.Header.Get("Link")
				paginationInfo2, getPagingErr := githubUtils.GetPagingFromLinkHeader(linkHeader)
				if getPagingErr != nil {
					logger.Info("GetPagingFromLinkHeader err: ", getPagingErr)
				}
				handlerErr := handler(res)
				if handlerErr != nil {
					return handlerErr
				}

				// only send message to channel if I'm the last page
				if page%conc == 0 {
					if paginationInfo2.Next == 1 {
						fmt.Println(page, "has no next page")
						c <- false
					} else {
						fmt.Printf("page: %v send true\n", page)
						c <- true
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		cont := <-c
		if !cont {
			break
		}
		step += 1
	}

	scheduler.WaitUntilFinish()
	return nil
}
