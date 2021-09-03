package tasks

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

type ApiProjectResponse struct {
	Name              string
	GitlabId          int    `json:"id"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebUrl            string `json:"web_url"`
	Visibility        string `json:"visibility"`
	OpenIssuesCount   int    `json:"open_issues_count"`
	StarCount         int    `json:"star_count"`
}

func CollectProjects(projectId int, c chan bool) error {
	gitlabApiClient := CreateApiClient()
	res, err := gitlabApiClient.Get(fmt.Sprintf("projects/%v", projectId), nil, nil)
	gitlabApiResponse := &ApiProjectResponse{}
	err = core.UnmarshalResponse(res, gitlabApiResponse)
	if err != nil {
		logger.Error("Error: ", err)
		return nil
	}
	fmt.Println("KEVIN >>> 1")
	gitlabProject := &models.GitlabProject{
		Name:              gitlabApiResponse.Name,
		GitlabId:          gitlabApiResponse.GitlabId,
		PathWithNamespace: gitlabApiResponse.PathWithNamespace,
		WebUrl:            gitlabApiResponse.WebUrl,
		Visibility:        gitlabApiResponse.Visibility,
		OpenIssuesCount:   gitlabApiResponse.OpenIssuesCount,
		StarCount:         gitlabApiResponse.StarCount,
	}
	fmt.Println("KEVIN >>> 2")
	err = lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&gitlabProject).Error
	fmt.Println("KEVIN >>> 3")
	if err != nil {
		logger.Error("Could not upsert: ", err)
	}
	fmt.Println("KEVIN >>> we have arrived at the channel!")
	c <- true
	return nil
}
