package tasks

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

func CreateApiClient(scheduler *utils.WorkerScheduler) *GitlabApiClient {
	gitlabApiClient := &GitlabApiClient{}
	V := config.GetConfig()
	gitlabApiClient.Setup(
		V.GetString("GITLAB_ENDPOINT"),
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %v", V.GetString("GITLAB_AUTH")),
		},
		20*time.Second,
		5,
		scheduler,
	)
	return gitlabApiClient
}

type GitlabPaginationHandler func(res *http.Response) error

func (gitlabApiClient *GitlabApiClient) getTotal(path string, queryParams *url.Values) (totalInt int, rateLimitPerSecond int, err error) {
	// just get the first page of results. The response has a head that tells the total pages
	queryParams.Set("page", "0")
	queryParams.Set("per_page", "1")

	res, err := gitlabApiClient.Get(path, queryParams, nil)

	if err != nil {
		logger.Error("failed to get total page", err)
		// try to read response body if it was not nil
		if res != nil {
			resBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				logger.Error("failed to read response body", err)
				return 0, 0, err
			}
			return 0, 0, fmt.Errorf("failed to get total page: %w\n%s", err, resBody)
		}
		return 0, 0, err
	}

	totalInt = -1
	total := res.Header.Get("X-Total")
	if total != "" {
		totalInt, err = convertStringToInt(total)
		if err != nil {
			return 0, 0, err
		}
	}
	rateLimitPerSecond = 100
	rateRemaining := res.Header.Get("ratelimit-remaining")
	if rateRemaining == "" {
		return
	}
	date, err := http.ParseTime(res.Header.Get("date"))
	if err != nil {
		return 0, 0, err
	}
	rateLimitResetTime, err := http.ParseTime(res.Header.Get("ratelimit-resettime"))
	if err != nil {
		return 0, 0, err
	}
	rateLimitInt, err := strconv.Atoi(rateRemaining)
	if err != nil {
		return 0, 0, err
	}
	rateLimitPerSecond = rateLimitInt / int(rateLimitResetTime.Unix()-date.Unix()) * 9 / 10
	return
}

func convertStringToInt(input string) (int, error) {
	return strconv.Atoi(input)
}

// run all requests in an Ants worker pool
func (gitlabApiClient *GitlabApiClient) FetchWithPaginationAnts(path string, queryParams *url.Values, pageSize int, handler GitlabPaginationHandler) error {
	// We need to get the total pages first so we can loop through all requests concurrently
	if queryParams == nil {
		queryParams = &url.Values{}
	}
	total, _, err := gitlabApiClient.getTotal(path, queryParams)
	if err != nil {
		return err
	}

	// not all api return x-total header, use step concurrency
	if total == -1 {
		// TODO: How do we know how high we can set the conc? Is is rateLimit?
		conc := 30
		step := 0
		c := make(chan bool)
		errs := make([]string, 0)
		for {
			// first, we fetch 10 pages in parallel, and we want to fetch the 10th page first,
			// so we can start next 10 pages ASAP
			for i := conc; i > 0; i-- {
				page := step*conc + i
				queryCopy := url.Values{}
				for k, v := range *queryParams {
					queryCopy[k] = v
				}
				queryCopy.Set("page", strconv.Itoa(page))
				queryCopy.Set("per_page", strconv.Itoa(pageSize))
				err = gitlabApiClient.GetAsync(path, &queryCopy, func(res *http.Response) error {
					e := handler(res)
					// sth went wrong during page process
					if e != nil {
						errs = append(errs, e.Error())
						c <- false
						return e
					}
					// see if there was next page for 10th page
					if page%conc == 0 {
						_, e = strconv.ParseInt(res.Header.Get("X-Next-Page"), 10, 32)
						c <- e == nil
					}
					return nil
				})
				if err != nil {
					return err
				}
			}
			// wait for 10th page or any error
			cont := <-c
			if !cont {
				break
			}
			step += 1
		}
		if len(errs) > 0 {
			return fmt.Errorf(strings.Join(errs, "\n"))
		}
	} else {
		// Loop until all pages are requested
		for i := 1; (i * pageSize) <= (total + pageSize); i++ {
			// we need to save the value for the request so it is not overwritten
			currentPage := i
			queryCopy := url.Values{}
			for k, v := range *queryParams {
				queryCopy[k] = v

			}
			queryCopy.Set("page", strconv.Itoa(currentPage))
			queryCopy.Set("per_page", strconv.Itoa(pageSize))

			err = gitlabApiClient.GetAsync(path, &queryCopy, handler)

			if err != nil {
				return err
			}
		}
	}

	gitlabApiClient.WaitOtherGoroutines()
	return nil
}

// fetch paginated without ANTS worker pool
func (gitlabApiClient *GitlabApiClient) FetchWithPagination(path string, queryParams *url.Values, pageSize int, handler GitlabPaginationHandler) error {
	if queryParams == nil {
		queryParams = &url.Values{}
	}
	// We need to get the total pages first so we can loop through all requests concurrently
	total, _, _ := gitlabApiClient.getTotal(path, queryParams)

	// Loop until all pages are requested
	for i := 0; (i * pageSize) < total; i++ {
		// we need to save the value for the request so it is not overwritten
		currentPage := i
		queryParams.Set("per_page", strconv.Itoa(pageSize))
		queryParams.Set("page", strconv.Itoa(currentPage))
		res, err := gitlabApiClient.Get(path, queryParams, nil)

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
