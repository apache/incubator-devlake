package models

import (
	"net/http"
	"strconv"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type ArgoTaskData struct {
	Options   ArgoOptions
	ApiClient *helper.ApiAsyncClient
}

type ArgoOptions struct {
	ConnectionId uint64
	Project      string
}

type ArgoApiParams struct {
	ConnectionId uint64
	Project      string
}

func GetTotalPagesFromResponse(res *http.Response, args *api.ApiCollectorArgs) (int, errors.Error) {
	totalPages := res.Header.Get("X-Total-Pages")
	if totalPages == "" {
		return 1, nil
	}
	pages, err := strconv.Atoi(totalPages)
	if err != nil {
		return 0, errors.Convert(err)
	}
	return pages, nil
}
