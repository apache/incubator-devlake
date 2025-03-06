package tasks

import (
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
}

// converte o _tools para o modelo de dominio (pacote devops)
//
//	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
func ConvertDeployments(taskCtx plugin.SubTaskContext) errors.Error {
	log.Println("Iniciando plugin de convert.")

	converter, err := api.NewStatefulDataConverter(&api.StatefulDataConverterArgs[models.Deployment]{
		SubtaskCommonArgs: &api.SubtaskCommonArgs{
			Params: taskCtx.GetData(),
			Table:  RAW_DEPLOYMENT_TABLE,
		},
		Input:   nil,
		Convert: nil,
	})

	if err != nil {
		return err
	}

	log.Println("Conversão de deployments concluída com sucesso!")
	return converter.Execute()
}
