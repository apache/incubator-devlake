package tasks

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/merico-analysis-engine/models"
	"gorm.io/gorm/clause"
)

type ApiProjectResponse struct {
	Id           int        `json:"id"`
	GitUrl       string     `json:"git_url"`
	Priority     int        `json:"priority"`
	AECreateTime *time.Time `json:"create_time"`
	AEUpdateTime *time.Time `json:"update_time"`
}

// This function is required by AE to prevent Man-in-the-middle attacks
// You need to fully encode all query parameters and other things in order to get a
// correct sign value in the url
func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func setQueryParams() *url.Values {
	queryParams := &url.Values{}
	queryParams.Set("app_id", config.V.GetString("AE_APP_ID"))
	queryParams.Set("nonce_str", config.V.GetString("AE_NONCE_STR"))
	queryParams.Set("sign", getMD5Hash(config.V.GetString("AE_SECRET_KEY")))
	return queryParams
}

func CollectProject(projectId int) error {
	aeApiClient := CreateApiClient()

	res, err := aeApiClient.Get(fmt.Sprintf("/projects/%v", projectId), setQueryParams(), nil)
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
