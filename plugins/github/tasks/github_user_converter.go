package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/user"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
	"reflect"
)

func ConvertUsers(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)

	cursor, err := db.Model(&githubModels.GithubUser{}).
		Where("_raw_data_params = ?", data.Options.ParamString).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	userIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubUser{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		Ctx:          taskCtx,
		InputRowType: reflect.TypeOf(githubModels.GithubUser{}),
		Input:        cursor,
		BatchSelectors: map[reflect.Type]helper.BatchSelector{
			reflect.TypeOf(&user.User{}): {
				Query: "_raw_data_params = ?",
				Parameters: []interface{}{
					data.Options.ParamString,
				},
			},
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			githubUser := inputRow.(*githubModels.GithubUser)
			domainUser := &user.User{
				DomainEntity: domainlayer.DomainEntity{Id: userIdGen.Generate(githubUser.Id)},
				Name:         githubUser.Login,
				AvatarUrl:    githubUser.AvatarUrl,
			}
			domainUser.RawDataOrigin = githubUser.RawDataOrigin

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
