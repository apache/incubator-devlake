package api

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/mitchellh/mapstructure"
)

type ApiMeResponse struct {
	Name string `json:"name"`
}

/*
GET /plugins/ae/test
*/
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// TODO: implement test connection
	return &core.ApiResourceOutput{Body: true}, nil
}

type AEConfig struct {
	AE_APP_ID    string `mapstructure:"AE_APP_ID"`
	AE_SIGN      string `mapstructure:"AE_SIGN"`
	AE_NONCE_STR string `mapstructure:"AE_NONCE_STR"`
	AE_ENDPOINT  string `mapstructure:"AE_ENDPOINT"`
}

// This object conforms to what the frontend currently expects.
type AEResponse struct {
	Endpoint string
	Sign     string
	Nonce    string
	AppId    string
	Name     string
	ID       int
}

/*
PUT /plugins/ae/connections/:connectionId
*/
func PutConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	aeConnection := AEConfig{}
	err := mapstructure.Decode(input.Body, &aeConnection)
	if err != nil {
		return nil, err
	}
	v := config.GetConfig()

	if aeConnection.AE_APP_ID != "" {
		v.Set("AE_SIGN", aeConnection.AE_SIGN)
	}
	if aeConnection.AE_SIGN != "" {
		v.Set("AE_SIGN", aeConnection.AE_SIGN)
	}
	if aeConnection.AE_NONCE_STR != "" {
		v.Set("AE_NONCE_STR", aeConnection.AE_NONCE_STR)
	}
	if aeConnection.AE_ENDPOINT != "" {
		v.Set("AE_ENDPOINT", aeConnection.AE_ENDPOINT)
	}

	err = v.WriteConfig()
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: "Success"}, nil
}

/*
GET /plugins/ae/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// RETURN ONLY 1 SOURCE (FROM ENV) until multi-connection is developed.
	aeConnections, err := GetConnectionFromEnv()
	response := []AEResponse{*aeConnections}
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: response}, nil
}

/*
GET /plugins/ae/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-connection is developed.)
	aeConnections, err := GetConnectionFromEnv()
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: aeConnections}, nil
}

func GetConnectionFromEnv() (*AEResponse, error) {
	v := config.GetConfig()
	var configJson AEConfig
	err := v.Unmarshal(&configJson)
	if err != nil {
		return nil, err
	}
	return &AEResponse{
		AppId:    configJson.AE_APP_ID,
		Nonce:    configJson.AE_NONCE_STR,
		Sign:     configJson.AE_SIGN,
		Endpoint: configJson.AE_ENDPOINT,
		ID:       1,
		Name:     "AE",
	}, nil
}
