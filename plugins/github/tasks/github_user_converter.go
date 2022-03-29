package tasks

import (
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/user"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var ConvertUsersMeta = core.SubTaskMeta{
	Name:             "convertUsers",
	EntryPoint:       ConvertUsers,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_users into  domain layer table users",
}

func ConvertUsers(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)

	cursor, err := db.Model(&githubModels.GithubUser{}).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	userIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubUser{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(githubModels.GithubUser{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			Table: RAW_COMMIT_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			githubUser := inputRow.(*githubModels.GithubUser)
			domainUser := &user.User{
				DomainEntity: domainlayer.DomainEntity{Id: userIdGen.Generate(githubUser.Id)},
				Name:         githubUser.Login,
				AvatarUrl:    githubUser.AvatarUrl,
			}
			return []interface{}{
				domainUser,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
