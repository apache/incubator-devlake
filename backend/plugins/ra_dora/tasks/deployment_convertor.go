package tasks

import (
	"encoding/json"
	"log"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/ra_dora/models"
)

var _ plugin.SubTaskEntryPoint = ConvertDeployments

func init() {
	RegisterSubtaskMeta(&ConvertDeploymentsMeta)
}

// Task metadata
var ConvertDeploymentsMeta = plugin.SubTaskMeta{
	Name:             "convert_deployments",
	EntryPoint:       ConvertDeployments,
	EnabledByDefault: true,
	Description:      "Converte deployments do DevLake para análise de métricas DORA",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	DependencyTables: []string{RAW_DEPLOYMENT_TABLE},
	ProductTables:    []string{"dora_deployments"},
}

func ConvertDeployments(taskCtx plugin.SubTaskContext) errors.Error {
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:    taskCtx,
			Params: taskCtx.GetData(),
			Table:  RAW_DEPLOYMENT_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var rawDeployment models.Deployments

			err := errors.Convert(json.Unmarshal(row.Data, &rawDeployment))
			if err != nil {
				return nil, err
			}

			deployment := &models.Deployment{
				ID:      rawDeployment.ID,
				ScopeID: rawDeployment.ScopeID,
			}

			results := make([]interface{}, 0, 1)
			results = append(results, deployment)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	log.Println("Conversão de deployments concluída com sucesso!")
	return extractor.Execute()
}
