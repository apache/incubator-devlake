package tasks

import (
	"net/http"

	"github.com/apache/incubator-devlake/plugins/helper"
)

func GetTotalPagesFromResponse(res *http.Response, args *helper.ApiCollectorArgs) (int, error) {
	body := &JiraPagination{}
	err := helper.UnmarshalResponse(res, body)
	if err != nil {
		return 0, err
	}
	pages := body.Total / args.PageSize
	if body.Total%args.PageSize > 0 {
		pages++
	}
	return pages, nil
}
