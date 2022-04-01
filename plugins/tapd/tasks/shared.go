package tasks

import (
	"fmt"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/helper"
	"net/http"
	"net/url"
)

type Page struct {
	Data Data `json:"data"`
}
type Data struct {
	Count int `json:"count"`
}

var UserIdGen *didgen.DomainIdGenerator
var WorkspaceIdGen *didgen.DomainIdGenerator
var IssueIdGen *didgen.DomainIdGenerator

// res will not be used
func GetTotalPagesFromResponse(_ *http.Response, args *helper.ApiCollectorArgs) (int, error) {
	data := args.Ctx.GetData().(*TapdTaskData)
	apiClient, err := NewTapdApiPageClient(args.Ctx.TaskContext(), data.Source)
	if err != nil {
		return 0, err
	}
	query := url.Values{}
	query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))

	res, err := apiClient.Get(fmt.Sprintf("%s/count", args.UrlTemplate), query, nil)
	if err != nil {
		return 0, err
	}
	var page Page
	err = helper.UnmarshalResponse(res, &page)

	count := page.Data.Count
	totalPage := count/args.PageSize + 1

	return totalPage, err
}
