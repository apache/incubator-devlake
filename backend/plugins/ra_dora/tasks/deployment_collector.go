package tasks

import (
	"fmt"
	"log"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/ra_dora/models"
)

const RAW_DEPLOYMENT_TABLE = "argo_api_deployments"

// Task metadata
var CollectDeploymentsMeta = plugin.SubTaskMeta{
	Name:             "collect_deployments",
	EntryPoint:       CollectApiDeployments,
	EnabledByDefault: true,
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	DependencyTables: []string{},
	ProductTables:    []string{RAW_DEPLOYMENT_TABLE},
}

// Coletor principal
func CollectApiDeployments(taskCtx plugin.SubTaskContext) errors.Error {
	log.Println("Iniciando plugin de collect.")

	data := taskCtx.GetData().(*models.ArgoTaskData)
	apiCollector, err := helper.NewStatefulApiCollector(helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: models.ArgoApiParams{
			ConnectionId: data.Options.ConnectionId,
			Project:      data.Options.Project,
		},
		Table: RAW_DEPLOYMENT_TABLE,
	})
	if err != nil {
		return err
	}

	err = apiCollector.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    100,
		UrlTemplate: "api/v1/workflows/{{ .Params.Project }}",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("offset", fmt.Sprintf("%v", reqData.Pager.Page*reqData.Pager.Size))
			return query, nil
		},
		GetTotalPages:  models.GetTotalPagesFromResponse,
		ResponseParser: models.GetRawMessageFromResponse,
	})

	if err != nil {
		return err
	}

	log.Println("Finalizado plugin de collect.")

	return apiCollector.Execute()
}
