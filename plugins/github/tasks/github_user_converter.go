package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/user"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertUsers() error {
	user := &user.User{}
	githubUser := &githubModels.GithubUser{}

	cursor, err := lakeModels.Db.Model(githubUser).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	userIdGenerator := didgen.NewDomainIdGenerator(githubUser)
	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, githubUser)
		if err != nil {
			return err
		}
		user.Id = userIdGenerator.Generate(githubUser.Id)
		user.Name = githubUser.Login
		user.AvatarUrl = githubUser.AvatarUrl
		err := lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(user).Error
		if err != nil {
			return err
		}
	}
	return nil
}
