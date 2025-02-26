package tasks

import (
	"encoding/json"
	"log"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/ra_dora/models"
)

var _ plugin.SubTaskEntryPoint = ExtractDeployments

func init() {
	RegisterSubtaskMeta(&ExtractDeploymentsMeta)
}

// Task metadata
var ExtractDeploymentsMeta = plugin.SubTaskMeta{
	Name:             "extract_deployments",
	EntryPoint:       ExtractDeployments,
	EnabledByDefault: true,
	Description:      "Extrai deployments do DevLake para análise",
}

func ExtractDeployments(taskCtx plugin.SubTaskContext) errors.Error {
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:    taskCtx,
			Params: taskCtx.GetData(),
			Table:  "ra_dora_deployments",
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var deployments models.DatabaseDeployments

			err := errors.Convert(json.Unmarshal(row.Data, &deployments))
			if err != nil {
				return nil, err
			}

			raDeployment := &models.Deployments{
				ID:           deployments.ID,
				ScopeID:      deployments.ScopeID,
				Name:         deployments.Name,
				Result:       deployments.Result,
				Status:       deployments.Status,
				Environment:  deployments.Environment,
				CreatedDate:  deployments.CreatedDate,
				StartedDate:  deployments.StartedDate,
				FinishedDate: deployments.FinishedDate,
				DurationSec:  deployments.DurationSec,
			}

			results := make([]interface{}, 0, 2)
			results = append(results, raDeployment)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	log.Println("Extração de deployments concluída com sucesso!")
	return extractor.Execute()
}
