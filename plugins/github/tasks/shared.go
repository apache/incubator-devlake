package tasks

import (
	"github.com/merico-dev/lake/plugins/github/utils"
	"net/http"

	"github.com/merico-dev/lake/plugins/helper"
)

func GetTotalPagesFromResponse(res *http.Response, args *helper.ApiCollectorArgs) (int, error) {
	link := res.Header.Get("link")
	pageInfo, err := utils.GetPagingFromLinkHeader(link)
	if err != nil {
		return 0, nil
	}
	return pageInfo.Last, nil
}
