package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/utils"
)

type ApiAccessTokenRequest struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type ApiAccessTokenResponse struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	AppAccessToken    string `json:"app_access_token"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

const AUTH_ENDPOINT = "https://open.feishu.cn"
const ENDPOINT = "https://open.feishu.cn/open-apis/vc/v1"

func NewFeishuApiClient(taskCtx core.TaskContext) (*helper.ApiAsyncClient, error) {
	// load and process cconfiguration
	appId := taskCtx.GetConfig("FEISHU_APPID")
	if appId == "" {
		return nil, fmt.Errorf("invalid FEISHU_APPID")
	}
	secretKey := taskCtx.GetConfig("FEISHU_APPSCRECT")
	if secretKey == "" {
		return nil, fmt.Errorf("invalid FEISHU_APPSCRECT")
	}
	userRateLimit, err := utils.StrToIntOr(taskCtx.GetConfig("FEISHU_API_REQUESTS_PER_HOUR"), 18000)
	if err != nil {
		return nil, err
	}
	proxy := taskCtx.GetConfig("FEISHU_PROXY")

	authApiClient, err := helper.NewApiClient(AUTH_ENDPOINT, nil, 0, proxy, taskCtx.GetContext())
	if err != nil {
		return nil, err
	}

	// request for access token
	tokenReqBody := &ApiAccessTokenRequest{
		AppId:     appId,
		AppSecret: secretKey,
	}
	tokenRes, err := authApiClient.Post("open-apis/auth/v3/tenant_access_token/internal", nil, tokenReqBody, nil)
	if err != nil {
		return nil, err
	}
	tokenResBody := &ApiAccessTokenResponse{}
	err = helper.UnmarshalResponse(tokenRes, tokenResBody)
	if err != nil {
		return nil, err
	}
	if tokenResBody.AppAccessToken == "" && tokenResBody.TenantAccessToken == "" {
		return nil, fmt.Errorf("failed to request access token")
	}
	// real request apiClient
	apiClient, err := helper.NewApiClient(ENDPOINT, nil, 0, proxy, taskCtx.GetContext())
	if err != nil {
		return nil, err
	}
	// set token
	apiClient.SetHeaders(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", tokenResBody.TenantAccessToken),
	})

	apiClient.SetAfterFunction(func(res *http.Response) error {
		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("feishu authentication failed, please check your Bearer Auth Token")
		}
		return nil
	})

	// create async api client
	asyncApiCLient, err := helper.CreateAsyncApiClient(taskCtx, apiClient, &helper.ApiRateLimitCalculator{
		UserRateLimitPerHour: userRateLimit,
	})
	if err != nil {
		return nil, err
	}

	return asyncApiCLient, nil
}
