package models

type AeConnection struct {
	AppId    string `mapstructure:"appId" env:"AE_APP_ID" json:"appId"`
	Sign     string `mapstructure:"sign" env:"AE_SIGN" json:"sign"`
	NonceStr string `mapstructure:"nonceStr" env:"AE_NONCE_STR" json:"nonceStr"`
	Endpoint string `mapstructure:"endpoint" env:"AE_ENDPOINT" json:"endpoint"`
}

// This object conforms to what the frontend currently expects.
type AeResponse struct {
	AeConnection
	Name string `json:"name"`
	ID   int    `json:"id"`
}
