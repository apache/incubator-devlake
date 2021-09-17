package tasks

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

func getPaginationInfo(resourceUriFormat string) (int, int, error) {
	// just get the first page of results. The response has a head that tells the total pages
	page := 0
	page_size := 100 // This is the maximum
	res, err := githubApiClient.Get(fmt.Sprintf(resourceUriFormat, page, page_size), nil, nil)

	if err != nil {
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Error("UnmarshalResponse failed: ", string(resBody))
			return 0, 0, err
		}
		logger.Print(string(resBody) + "\n")
		return 0, 0, err
	}

	lastPage := 1 // Assumes that there is always at least 1 page
	linkHeader := res.Header.Get("Link")

	// PagingInfo object contains Next, First, Last, and Prev page number
	var paginationInfo githubUtils.PagingInfo
	paginationInfo, err = githubUtils.GetPagingFromLinkHeader(linkHeader)
	if err != nil {
		logger.Error("Pagination Info", err)
	}

	if paginationInfo.Last != "" {
		lastPage, err = convertStringToInt(paginationInfo.Last)
		if err != nil {
			return 0, 0, err
		}
	}
	rateRemaining := res.Header.Get("X-RateLimit-Remaining")
	date, err := http.ParseTime(res.Header.Get("Date"))
	if err != nil {
		return 0, 0, err
	}
	i, err := strconv.ParseInt(res.Header.Get("X-RateLimit-Reset"), 10, 64)
	if err != nil {
		panic(err)
	}
	rateLimitResetTime := time.Unix(i, 0)

	rateLimitInt, err := strconv.Atoi(rateRemaining)
	if err != nil {
		logger.Error("Convert error: ", err)
		return 0, 0, err
	}
	rateLimitPerSecond := rateLimitInt / int(rateLimitResetTime.Unix()-date.Unix()) * 9 / 10
	return lastPage, rateLimitPerSecond, nil
}

func convertStringToInt(input string) (int, error) {
	return strconv.Atoi(input)
}

// run all requests in an Ants worker pool
func (githubApiClient *GithubApiClient) FetchWithPaginationAnts(resourceUri string, pageSize int, handler GithubPaginationHandler) error {

	var resourceUriFormat string
	if strings.ContainsAny(resourceUri, "?") {
		resourceUriFormat = resourceUri + "&per_page=%v&page=%v"
	} else {
		resourceUriFormat = resourceUri + "?per_page=%v&page=%v"
	}
	// We need to get the total pages first so we can loop through all requests concurrently
	total, rateLimitPerSecond, err := getPaginationInfo(resourceUriFormat)
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

	// not all api return x-total header, use step concurrency
	if total == -1 {
		// TODO: How do we know how high we can set the conc? Is is rateLimit?
		conc := 10
		step := 0
		c := make(chan bool)
		for {
			for i := conc; i > 0; i-- {
				page := step*conc + i
				err := scheduler.Submit(func() error {
					url := fmt.Sprintf(resourceUriFormat, pageSize, page)
					res, err := githubApiClient.Get(url, nil, nil)
					if err != nil {
						return err
					}
					handlerErr := handler(res)
					if handlerErr != nil {
						return handlerErr
					}
					_, err = strconv.ParseInt(res.Header.Get("X-Next-Page"), 10, 32)
					// only send message to channel if I'm the last page
					if page%conc == 0 {
						if err != nil {
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
		// Loop until all pages are requested
		for i := 1; (i * pageSize) <= (total + pageSize); i++ {
			// we need to save the value for the request so it is not overwritten
			currentPage := i
			err1 := scheduler.Submit(func() error {
				url := fmt.Sprintf(resourceUriFormat, pageSize, currentPage)

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
	lastPage, _, _ := getPaginationInfo(resourceUriFormat)
	// Loop until all pages are requested
	// PLEASE NOTE: Pages start at 1. Page 1 and page 0 are the same.
	fmt.Println("INFO >>> last page: ", lastPage)
	for i := 1; i <= lastPage; i++ {
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
