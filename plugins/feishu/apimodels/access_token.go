package apimodels

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
