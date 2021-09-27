package tasks

import (
	"fmt"
	"io/ioutil"
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

func getPaginationInfoFromGitHub(resourceUriFormat string) (githubUtils.PagingInfo, int, error) {
	// just get the first page of results. The response has a head that tells the total pages
	paginationInfo := githubUtils.PagingInfo{
		First: 1,
		Last:  1,
		Next:  1,
		Prev:  1,
	}

	page := 1
	page_size := 100 // This is the maximum
	res, err := githubApiClient.Get(fmt.Sprintf(resourceUriFormat, page, page_size), nil, nil)
	fmt.Println("INFO >>> res.Status get page info", res.Status)
	if err != nil {
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Error("UnmarshalResponse failed: ", string(resBody))
			return paginationInfo, 0, err
		}
		logger.Print(string(resBody) + "\n")
		return paginationInfo, 0, err
	}

	linkHeader := res.Header.Get("Link")
	// PagingInfo object contains Next, First, Last, and Prev page number

	paginationInfo, err = githubUtils.GetPagingFromLinkHeader(linkHeader)
	if err != nil {
		logger.Info("", err)
	}
	// These are all strings
	date := res.Header.Get("Date")
	reset := res.Header.Get("X-RateLimit-Reset")
	remaining := res.Header.Get("X-RateLimit-Remaining")
	rateLimitInfo, rateLimitInfoErr := githubUtils.ConvertRateLimitInfo(date, reset, remaining)
	if rateLimitInfoErr != nil {
		fmt.Println("ERROR >>> Rate Limit Info Err: ", rateLimitInfoErr)
	}
	rateLimitPerSecond := githubUtils.GetRateLimitPerSecond(rateLimitInfo)

	return paginationInfo, rateLimitPerSecond, nil
}

// run all requests in an Ants worker pool
func (githubApiClient *GithubApiClient) FetchWithPaginationAnts(resourceUri string, pageSize int, handler GithubPaginationHandler) error {

	var resourceUriFormat string
	if strings.ContainsAny(resourceUri, "?") {
		resourceUriFormat = resourceUri + "&page=%v&per_page=%v"
	} else {
		resourceUriFormat = resourceUri + "?page=%v&per_page=%v"
	}

	// We need to get the total pages first so we can loop through all requests concurrently
	paginationInfo, rateLimitPerSecond, err := getPaginationInfoFromGitHub(resourceUriFormat)
	if err != nil {
		return err
	}
	workerNum := 50
	// set up the worker pool
	scheduler, err := utils.NewWorkerScheduler(workerNum, rateLimitPerSecond)
	if err != nil {
		return err
	}

	defer scheduler.Release()

	// if we don't get a "last" result, try going step by step using the "next" property
	if paginationInfo.Last == 0 {
		fmt.Println("INFO >>> No last page. Using step concurrency...")
		// TODO: How do we know how high we can set the conc? Is is rateLimit?
		conc := 10
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
						if paginationInfo2.Next == 0 {
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
	} else {
		fmt.Println("INFO >>> Last page found. Looping through: ", paginationInfo.Last)
		// Loop until all pages are requested
		for i := 1; i <= paginationInfo.Last; i++ {
			// we need to save the value for the request so it is not overwritten
			currentPage := i
			err1 := scheduler.Submit(func() error {
				url := fmt.Sprintf(resourceUriFormat, currentPage, pageSize)
				res, err := githubApiClient.Get(url, nil, nil)

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
	}

	scheduler.WaitUntilFinish()
	return nil
}

// fetch paginated without ANTS worker pool
func (githubApiClient *GithubApiClient) FetchWithPagination(resourceUri string, pageSize int, handler GithubPaginationHandler) error {

	var resourceUriFormat string
	if strings.ContainsAny(resourceUri, "?") {
		resourceUriFormat = resourceUri + "&page=%v&per_page=%v"
	} else {
		resourceUriFormat = resourceUri + "?page=%v&per_page=%v"
	}

	// We need to get the total pages first so we can loop through all requests concurrently
	paginationInfo, _, _ := getPaginationInfoFromGitHub(resourceUriFormat)
	// Loop until all pages are requested
	// PLEASE NOTE: Pages start at 1. Page 1 and page 0 are the same.
	logger.Info("INFO >>> Last Page: ", paginationInfo.Last)
	for i := 1; i <= paginationInfo.Last; i++ {
		// we need to save the value for the request so it is not overwritten
		currentPage := i
		url := fmt.Sprintf(resourceUriFormat, currentPage, pageSize)
		res, err := githubApiClient.Get(url, nil, nil)
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
