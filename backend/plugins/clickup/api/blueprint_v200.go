package api

import (
	"fmt"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/utils"

	"github.com/apache/incubator-devlake/plugins/clickup/models"

	"github.com/apache/incubator-devlake/core/errors"

	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func MakePipelinePlanV200(subtaskMetas []plugin.SubTaskMeta, connectionId uint64, scope []*plugin.BlueprintScopeV200, syncPolicy *plugin.BlueprintSyncPolicy) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
	scopes, err := makeScopeV200(connectionId, scope)
	if err != nil {
		return nil, nil, err
	}

	plan := make(plugin.PipelinePlan, len(scope))
	plan, err = makePipelinePlanV200(subtaskMetas, plan, scope, connectionId, syncPolicy)
	if err != nil {
		return nil, nil, err
	}

	return plan, scopes, nil
}

func makeScopeV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200) ([]plugin.Scope, errors.Error) {
	sc := make([]plugin.Scope, 0, len(scopes))

	for _, scope := range scopes {
		id := didgen.NewDomainIdGenerator(&models.ClickUpSpace{}).Generate(connectionId, scope.Id)
		clickupSpace := &models.ClickUpSpace{}

		// get space from db
		err := basicRes.GetDal().First(clickupSpace,
			dal.Where(`connection_id = ? and id = ?`,
				connectionId, scope.Id))
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find space %s", scope.Id))
		}

		// add board to scopes
		if utils.StringsContains(scope.Entities, plugin.DOMAIN_TYPE_TICKET) {
			scopeTicket := ticket.NewBoard(id, clickupSpace.Name)

			sc = append(sc, scopeTicket)
		}
	}

	return sc, nil
}

func makePipelinePlanV200(subtaskMetas []plugin.SubTaskMeta, plan plugin.PipelinePlan, scopes []*plugin.BlueprintScopeV200, connectionId uint64, syncPolicy *plugin.BlueprintSyncPolicy) (plugin.PipelinePlan, errors.Error) {
	for i, scope := range scopes {
		stage := plan[i]
		if stage == nil {
			stage = plugin.PipelineStage{}
		}

		// construct task options for clickup
		options := make(map[string]interface{})
		options["connectionId"] = connectionId
		options["scopeId"] = scope.Id
		if syncPolicy.TimeAfter != nil {
			options["createdDateAfter"] = syncPolicy.TimeAfter.Format(time.RFC3339)
		}

		// construct subtasks
		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, scope.Entities)
		if err != nil {
			return nil, err
		}

		stage = append(stage, &plugin.PipelineTask{
			Plugin:   "clickup",
			Subtasks: subtasks,
			Options:  options,
		})

		plan[i] = stage
	}
	return plan, nil
}
