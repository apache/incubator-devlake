package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/ae/models"
	"github.com/merico-dev/lake/plugins/core"
	"gorm.io/gorm/clause"
)

type ApiProjectResponse struct {
	Name              string `josn:"name"`
	AEId              int    `json:"id"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebUrl            string `json:"web_url"`
	Visibility        string `json:"visibility"`
	OpenIssuesCount   int    `json:"open_issues_count"`
	StarCount         int    `json:"star_count"`
}

func CollectProject(projectId int) error {
	aeApiClient := CreateApiClient()
	res, err := aeApiClient.Get(fmt.Sprintf("projects/%v", projectId), nil, nil)
	if err != nil {
		logger.Error("Error: ", err)
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	aeApiResponse := &ApiProjectResponse{}
	err = core.UnmarshalResponse(res, aeApiResponse)
	if err != nil {
		logger.Error("Error: ", err)
		return err
	}
	aeProject := &models.AEProject{
		Name:              aeApiResponse.Name,
		AEId:              aeApiResponse.AEId,
		PathWithNamespace: aeApiResponse.PathWithNamespace,
		WebUrl:            aeApiResponse.WebUrl,
		Visibility:        aeApiResponse.Visibility,
		OpenIssuesCount:   aeApiResponse.OpenIssuesCount,
		StarCount:         aeApiResponse.StarCount,
	}
	err = lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&aeProject).Error
	if err != nil {
		logger.Error("Could not upsert: ", err)
		return err
	}
	return nil
}
