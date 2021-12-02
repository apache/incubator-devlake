package tasks

import (
	"fmt"
	"net/http"
	"time"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/ae/models"
	"github.com/merico-dev/lake/plugins/core"
	"gorm.io/gorm/clause"
)

type ApiProjectResponse struct {
	Id           int        `json:"id"`
	GitUrl       string     `json:"git_url"`
	Priority     int        `json:"priority"`
	AECreateTime *time.Time `json:"create_time"`
	AEUpdateTime *time.Time `json:"update_time"`
}

func CollectProject(projectId int) error {
	aeApiClient := CreateApiClient()

	res, err := aeApiClient.Get(fmt.Sprintf("/projects/%v", projectId), nil, nil)
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
		Id:           aeApiResponse.Id,
		GitUrl:       aeApiResponse.GitUrl,
		Priority:     aeApiResponse.Priority,
		AECreateTime: aeApiResponse.AECreateTime,
		AEUpdateTime: aeApiResponse.AEUpdateTime,
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
