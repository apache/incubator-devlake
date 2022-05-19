package tasks

import (
	"github.com/apache/incubator-devlake/plugins/github/utils"
	"net/http"

	"github.com/apache/incubator-devlake/plugins/helper"
)

func GetTotalPagesFromResponse(res *http.Response, args *helper.ApiCollectorArgs) (int, error) {
	link := res.Header.Get("link")
	pageInfo, err := utils.GetPagingFromLinkHeader(link)
	if err != nil {
		return 0, nil
	}
	return pageInfo.Last, nil
}
