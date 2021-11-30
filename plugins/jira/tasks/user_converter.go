package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/okgen"
	"github.com/merico-dev/lake/models/domainlayer/user"
	jiraModels "github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

func ConvertUsers(sourceId uint64) error {

	var jiraUserRows []*jiraModels.JiraUser

	err := lakeModels.Db.Find(&jiraUserRows, "source_id = ?", sourceId).Error
	if err != nil {
		return err
	}

	userOriginKeyGenerator := okgen.NewOriginKeyGenerator(&jiraModels.JiraUser{})

	for _, jiraUser := range jiraUserRows {
		user := &user.User{
			DomainEntity: domainlayer.DomainEntity{
				OriginKey: userOriginKeyGenerator.Generate(jiraUser.SourceId, jiraUser.AccountId),
			},
			Name:      jiraUser.Name,
			Email:     jiraUser.Email,
			AvatarUrl: jiraUser.AvatarUrl,
			Timezone:  jiraUser.Timezone,
		}

		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(user).Error
		if err != nil {
			return err
		}

	}
	return nil
}
