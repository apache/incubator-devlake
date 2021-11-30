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
	"github.com/merico-dev/lake/plugins/core"
)

type AEApiClient struct {
	core.ApiClient
}

// WARNING!! HERE BE DRAGONS!!!
// Changing this code can easily break the ae authentication
func GetSign(page int, pageSize int, nonce int64) string {
	hasher := md5.New()

	appId := config.V.GetString("AE_APP_ID")
	secretKey := config.V.GetString("AE_SECRET_KEY")

	unencodedSign := fmt.Sprintf("app_id=%v&nonce_str=%v&page=%v&per_page=%v&key=%v", appId, nonce, page, pageSize, secretKey)

	_, err := hasher.Write([]byte(unencodedSign))
	if err != nil {
		return ""
	}

	md5EncodedSign := strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
	return md5EncodedSign
}

// WARNING!! HERE BE DRAGONS!!!
// Changing this code can easily break the ae authentication
func SetQueryParams(page int, pageSize int) *url.Values {
	nonce := time.Now().Unix()

	queryParams := &url.Values{}
	queryParams.Set("app_id", config.V.GetString("AE_APP_ID"))
	queryParams.Set("nonce_str", fmt.Sprintf("%v", nonce))
	queryParams.Set("page", fmt.Sprintf("%v", page))
	queryParams.Set("per_page", fmt.Sprintf("%v", pageSize))
	queryParams.Set("sign", GetSign(page, pageSize, nonce))
	return queryParams
}

func CreateApiClient() *AEApiClient {
	aeApiClient := &AEApiClient{}
	aeApiClient.Setup(
		config.V.GetString("AE_ENDPOINT"),
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %v", config.V.GetString("AE_AUTH")),
		},
		10*time.Second,
		3,
	)
	return aeApiClient
}

type AEPaginationHandler func(res *http.Response) error

// fetch paginated without ANTS worker pool
func (aeApiClient *AEApiClient) FetchWithPagination(path string, queryParams *url.Values, pageSize int, handler AEPaginationHandler) error {
	if queryParams == nil {
		queryParams = &url.Values{}
	}

	currentPage := 1

	// Loop until all pages are requested
	for {
		res, err := aeApiClient.Get(path, SetQueryParams(currentPage, pageSize), nil)
		if err != nil {
			return err
		}

		handlerErr := handler(res)
		if handlerErr != nil {
			return handlerErr
		}

		currentPage += 1
		if res.Header.Get("x-ae-has-next-page") == "false" {
			break
		}
	}

	return nil
}
