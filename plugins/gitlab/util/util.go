package util

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
)

func GetTotalByXTotal(url string) (int, error) {
	gitlabApiClient := CreateApiClient()

	// jsut get the first page of results. The response has a head that tells the total pages
	page := 0
	page_size := 1
	res, err := gitlabApiClient.Get(fmt.Sprintf(url, page_size, page), nil, nil)

	if err != nil {
		return 0, err
	}

	total := res.Header.Get("X-Total")
	totalInt, err := convertStringToInt(total)
	if err != nil {
		return 0, err
	}

	logger.Info("JON >>> totalInt", totalInt)
	return totalInt, nil
}
