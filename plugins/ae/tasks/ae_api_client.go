package tasks

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
)

type AEApiClient struct {
	core.ApiClient
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func getSign(query url.Values, appId, secretKey, nonceStr, timestamp string) string {
	// clone query because we need to add items
	kvs := make([]string, 0, len(query)+3)
	kvs = append(kvs, fmt.Sprintf("app_id=%s", appId))
	kvs = append(kvs, fmt.Sprintf("timestamp=%s", timestamp))
	kvs = append(kvs, fmt.Sprintf("nonce_str=%s", nonceStr))
	for key, values := range query {
		for _, value := range values {
			kvs = append(kvs, fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(value)))
		}
	}

	// sort by alphabetical order
	sort.Strings(kvs)

	// generate text for signature
	querystring := fmt.Sprintf("%s&key=%s", strings.Join(kvs, "&"), url.QueryEscape(secretKey))

	// sign it
	hasher := md5.New()
	_, err := hasher.Write([]byte(querystring))
	if err != nil {
		return ""
	}
	return strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
}

func beforeRequest(req *http.Request) error {
	appId := config.V.GetString("AE_APP_ID")
	if appId == "" {
		return fmt.Errorf("invalid AE_APP_ID")
	}
	secretKey := config.V.GetString("AE_SECRET_KEY")
	if appId == "" {
		return fmt.Errorf("invalid AE_SECRET_KEY")
	}
	nonceStr := RandStringRunes(8)
	timestamp := fmt.Sprintf("%v", time.Now().Unix())
	sign := getSign(req.URL.Query(), appId, secretKey, nonceStr, timestamp)
	req.Header.Set("x-ae-app-id", appId)
	req.Header.Set("x-ae-timestamp", timestamp)
	req.Header.Set("x-ae-nonce-str", nonceStr)
	req.Header.Set("x-ae-sign", sign)
	return nil
}

func CreateApiClient() *AEApiClient {
	aeApiClient := &AEApiClient{}
	aeApiClient.Setup(
		config.V.GetString("AE_ENDPOINT"),
		nil,
		10*time.Second,
		3,
	)
	aeApiClient.SetBeforeFunction(beforeRequest)
	return aeApiClient
}

type AEPaginationHandler func(res *http.Response) error

// fetch paginated without ANTS worker pool
func (aeApiClient *AEApiClient) FetchWithPagination(path string, pageSize int, handler AEPaginationHandler) error {
	currentPage := 1

	query := &url.Values{}
	query.Set("per_page", fmt.Sprintf("%d", pageSize))
	// Loop until all pages are requested
	for {
		query.Set("page", fmt.Sprintf("%d", currentPage))
		res, err := aeApiClient.Get(path, query, nil)
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
