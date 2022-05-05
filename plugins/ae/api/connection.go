package api

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"net/http"
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

type AeConnection struct {
	AppId    string `mapstructure:"appId" env:"AE_APP_ID" json:"appId"`
	Sign     string `mapstructure:"sign" env:"AE_SIGN" json:"sign"`
	NonceStr string `mapstructure:"nonceStr" env:"AE_NONCE_STR" json:"nonceStr"`
	Endpoint string `mapstructure:"endpoint" env:"AE_ENDPOINT" json:"endpoint"`
}

// This object conforms to what the frontend currently expects.
type AeResponse struct {
	AeConnection
	Name string
	ID   int
}

/*
PATCH /plugins/ae/connections/:connectionId
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	v := config.GetConfig()
	connection, err := helper.LoadFromConfig(v, &AeConnection{}, "env")
	if err != nil {
		return nil, err
	}
	// update from request and save to .env
	err = helper.SaveToConifgWithMap(v, connection, input.Body, "env")
	if err != nil {
		return nil, err
	}

	response := AeResponse{
		AeConnection: *connection.(*AeConnection),
		Name:         "AE",
		ID:           1,
	}
	return &core.ApiResourceOutput{Body: response, Status: http.StatusOK}, nil
}

/*
GET /plugins/ae/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// RETURN ONLY 1 SOURCE (FROM ENV) until multi-connection is developed.
	v := config.GetConfig()
	connection, err := helper.LoadFromConfig(v, &AeConnection{}, "env")
	if err != nil {
		return nil, err
	}
	response := AeResponse{
		AeConnection: *connection.(*AeConnection),
		Name:         "AE",
		ID:           1,
	}

	return &core.ApiResourceOutput{Body: []AeResponse{response}}, nil
}

/*
GET /plugins/ae/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-connection is developed.)
	v := config.GetConfig()
	connection, err := helper.LoadFromConfig(v, &AeConnection{}, "env")
	if err != nil {
		return nil, err
	}
	response := &AeResponse{
		AeConnection: *connection.(*AeConnection),
		Name:         "AE",
		ID:           1,
	}
	return &core.ApiResourceOutput{Body: response}, nil
}
