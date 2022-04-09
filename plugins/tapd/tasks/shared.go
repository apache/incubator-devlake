package tasks

import (
	"fmt"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
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

func parseIterationChangelog(taskCtx core.SubTaskContext, item *models.TapdChangelogItem) (*models.TapdChangelogItem, error) {
	data := taskCtx.GetData().(*TapdTaskData)
	db := taskCtx.GetDb()
	iterationFrom := &models.TapdIteration{}
	err := db.Model(&models.TapdIteration{}).
		Where("source_id = ? and workspace_id = ? and name = ?",
			data.Source.ID, data.Options.WorkspaceId, item.ValueBeforeParsed).Limit(1).Find(iterationFrom).Error
	if err != nil {
		return nil, err
	}
	item.IterationIdFrom = iterationFrom.ID

	iterationTo := &models.TapdIteration{}
	err = db.Model(&models.TapdIteration{}).
		Where("source_id = ? and workspace_id = ? and name = ?",
			data.Source.ID, data.Options.WorkspaceId, item.ValueAfterParsed).Limit(1).Find(iterationTo).Error
	if err != nil {
		return nil, err
	}
	item.IterationIdTo = iterationTo.ID
	return item, nil
}
