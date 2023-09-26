package tasks

import (
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
)

var ConvertUsersMeta = plugin.SubTaskMeta{
	Name:             "convertUsers",
	EntryPoint:       ConvertUsers,
	EnabledByDefault: true,
	Description:      "convert clickup tasks",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func ConvertUsers(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*ClickupTaskData)

	clickUpUser := &models.ClickUpUser{}
	clauses := []dal.Clause{
		dal.Select("_tool_clickup_user.*"),
		dal.From(clickUpUser),
		dal.Where(
			"_tool_clickup_user.connection_id = ?",
			data.Options.ConnectionId,
		),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	accountIdGen := didgen.NewDomainIdGenerator(&models.ClickUpUser{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ClickupApiParams{
				TeamId: data.TeamId,
			},
			Table: RAW_USER_TABLE,
		},
		InputRowType: reflect.TypeOf(models.ClickUpUser{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			user := inputRow.(*models.ClickUpUser)
			u := &crossdomain.Account{
				DomainEntity: domainlayer.DomainEntity{
					Id: accountIdGen.Generate(data.Options.ConnectionId, user.AccountId),
				},
				FullName:  user.Username,
				UserName:  user.Username,
				Email:     user.Email,
				AvatarUrl: user.ProfilePictureUrl,
			}
			return []interface{}{u}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
