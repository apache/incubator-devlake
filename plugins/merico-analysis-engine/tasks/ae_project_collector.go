package tasks

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
// IE: app_id={app_id}&key={secretKey}&nonce_str={timestamp}&page={page}&per_page={page_size}
func getSign(page int, pageSize int, nonce int64) string {
	hasher := md5.New()

	appId := config.V.GetString("AE_APP_ID")
	secretKey := config.V.GetString("AE_SECRET_KEY")

	unencodedSign := fmt.Sprintf("app_id=%v&nonce_str=%v&page=%v&per_page=%v&key=%v", appId, nonce, page, pageSize, secretKey)

	logger.Info("JON >>> unencodedSign", unencodedSign)
	hasher.Write([]byte(unencodedSign))

	md5EncodedSign := strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
	return md5EncodedSign
}

func setQueryParams(page int, pageSize int, nonceStr int64) *url.Values {
	queryParams := &url.Values{}
	queryParams.Set("app_id", config.V.GetString("AE_APP_ID"))
	queryParams.Set("nonce_str", fmt.Sprintf("%v", nonceStr))
	queryParams.Set("page", fmt.Sprintf("%v", page))
	queryParams.Set("per_page", fmt.Sprintf("%v", pageSize))
	queryParams.Set("sign", getSign(page, pageSize, nonceStr))
	return queryParams
}

func CollectProject(projectId int) error {
	aeApiClient := CreateApiClient()

	nonce := time.Now().Unix()

	res, err := aeApiClient.Get(fmt.Sprintf("/projects/%v", projectId), setQueryParams(1, 100, nonce), nil)
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
