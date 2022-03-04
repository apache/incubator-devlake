package tasks

import (
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

type Reviewer struct {
	GitlabId       int `json:"id"`
	MergeRequestId int
	Name           string
	Username       string
	State          string
	AvatarUrl      string `json:"avatar_url"`
	WebUrl         string `json:"web_url"`
}

func CreateReviewers(projectId int, mergeRequestId int, reviewers []Reviewer) {
	for _, reviewer := range reviewers {
		gitlabReviewer := &models.GitlabReviewer{
			GitlabId:       reviewer.GitlabId,
			MergeRequestId: mergeRequestId,
			ProjectId:      projectId,
			Username:       reviewer.Username,
			Name:           reviewer.Name,
			State:          reviewer.State,
			AvatarUrl:      reviewer.AvatarUrl,
			WebUrl:         reviewer.WebUrl,
		}
		err := lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&gitlabReviewer).Error
		if err != nil {
			logger.Error("Could not upsert: ", err)
		}
	}
}
