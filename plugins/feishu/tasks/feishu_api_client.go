package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"os"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/utils"
	"github.com/merico-dev/lake/config"
	"github.com/faabiosr/cachego"
	"github.com/faabiosr/cachego/file"
)

type FeishuApiClient struct{
	core.ApiClient
}
var ServerUrl  = "https://open.feishu.cn"

type getRefreshRequestFunc func() *http.Request
type DefaultAccessTokenManager struct {
	Id                    string
	GetRefreshRequestFunc getRefreshRequestFunc
	Cache                 cachego.Cache
}

var getAccessTokenLock sync.Mutex

// GetAccessToken
func (m *DefaultAccessTokenManager) GetAccessToken() (accessToken string, err error) {

	cacheKey := m.getCacheKey()
	accessToken, err = m.Cache.Fetch(cacheKey)
	if accessToken != "" {
		return
	}

	getAccessTokenLock.Lock()
	defer getAccessTokenLock.Unlock()

	accessToken, err = m.Cache.Fetch(cacheKey)
	if accessToken != "" {
		return
	}

	req := m.GetRefreshRequestFunc()

	// 添加 serverUrl
	if !strings.HasPrefix(req.URL.String(), "http") {
		parse, _ := url.Parse(ServerUrl)
		req.URL.Host = parse.Host
		req.URL.Scheme = parse.Scheme
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	defer response.Body.Close()

	var result = struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`

		AppAccessToken    string `json:"app_access_token"`
		TenantAccessToken string `json:"tenant_access_token"`

		Expire int `json:"expire"`
	}{}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		err = fmt.Errorf("unmarshal error %s", string(resp))
		return
	}

	if result.AppAccessToken == "" && result.TenantAccessToken == "" {
		err = fmt.Errorf("%s", string(resp))
		return
	}

	accessToken = result.AppAccessToken
	if result.TenantAccessToken != "" {
		accessToken = result.TenantAccessToken
	}

	err = m.Cache.Save(cacheKey, accessToken, time.Duration(result.Expire)*time.Second)
	if err != nil {
		return
	}

	return
}

// getCacheKey
func (m *DefaultAccessTokenManager) getCacheKey() (key string) {
	return "access_token:" + m.Id
}


func NewFeishuApiClient(
	endpoint string,
	scheduler *utils.WorkerScheduler,
	logger core.Logger,
) *FeishuApiClient{
	feishuApiClient := &FeishuApiClient{}
	//get feishu token
	atm := &DefaultAccessTokenManager{
		Id: config.GetConfig().GetString("FEISHU_APPID"),
		Cache: file.New(os.TempDir()),
		GetRefreshRequestFunc: func() *http.Request {
			payload := `{
                "app_id":"` + config.GetConfig().GetString("FEISHU_APPID") + `",
                "app_secret":"` + config.GetConfig().GetString("FEISHU_APPSCRECT") + `"
            }`
			req, _ := http.NewRequest(http.MethodPost, ServerUrl+"/open-apis/auth/v3/tenant_access_token/internal/", strings.NewReader(payload))
			return req
		},
	}
	// request AccessToken api
	tenantAccessToken, _ := atm.GetAccessToken()

	feishuApiClient.Setup(
		endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %v", tenantAccessToken),
		},
		50*time.Second,
		3, 
		scheduler,
	)

	feishuApiClient.SetAfterFunction(func(res *http.Response) error {
		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("feishu authentication failed, please check your Bearer Auth Token")
		}
		return nil
	})

	feishuApiClient.SetLogger(logger)
	return feishuApiClient
}