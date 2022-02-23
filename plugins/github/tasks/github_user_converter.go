package tasks

import (
	"context"
	"github.com/merico-dev/lake/errors"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/user"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertUsers(ctx context.Context) error {
	domainUser := &user.User{}
	githubUser := &githubModels.GithubUser{}

	cursor, err := lakeModels.Db.Model(githubUser).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	userIdGenerator := didgen.NewDomainIdGenerator(githubUser)
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.TaskCanceled
		default:
		}
		err = lakeModels.Db.ScanRows(cursor, githubUser)
		if err != nil {
			return err
		}
		domainUser.Id = userIdGenerator.Generate(githubUser.Id)
		domainUser.Name = githubUser.Login
		domainUser.AvatarUrl = githubUser.AvatarUrl
		err := lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(domainUser).Error
		if err != nil {
			return err
		}
	}
	return nil
}
