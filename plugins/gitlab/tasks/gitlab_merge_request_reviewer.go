package tasks

import (
	"github.com/merico-dev/lake/plugins/gitlab/models"
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

func NewReviewer(projectId int, mergeRequestId int, reviewer Reviewer) *models.GitlabReviewer {
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
	return gitlabReviewer
}
