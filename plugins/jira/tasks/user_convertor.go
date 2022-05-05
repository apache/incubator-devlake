package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/user"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
	"reflect"
)

func ConvertUsers(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	connectionId := data.Connection.ID
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert user")
	cursor, err := db.Model(&models.JiraUser{}).Where("connection_id = ?", connectionId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	userIdGen := didgen.NewDomainIdGenerator(&models.JiraUser{})
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: connectionId,
				BoardId:      boardId,
			},
			Table: RAW_USERS_TABLE,
		},
		InputRowType: reflect.TypeOf(models.JiraUser{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			jiraUser := inputRow.(*models.JiraUser)
			u := &user.User{
				DomainEntity: domainlayer.DomainEntity{
					Id: userIdGen.Generate(connectionId, jiraUser.AccountId),
				},
				Name:      jiraUser.Name,
				Email:     jiraUser.Email,
				AvatarUrl: jiraUser.AvatarUrl,
				Timezone:  jiraUser.Timezone,
			}
			return []interface{}{u}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
