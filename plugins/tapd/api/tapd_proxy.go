package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

const (
	TimeOut = 10 * time.Second
)

func Proxy(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	sourceId := input.Params["sourceId"]
	if sourceId == "" {
		return nil, fmt.Errorf("missing sourceid")
	}
	tapdSourceId, err := strconv.ParseUint(sourceId, 10, 64)
	if err != nil {
		return nil, err
	}
	tapdSource := &models.TapdSource{}
	err = db.First(tapdSource, tapdSourceId).Error
	if err != nil {
		return nil, err
	}
	encKey := cfg.GetString(core.EncodeKeyEnvStr)
	basicAuth, err := core.Decrypt(encKey, tapdSource.BasicAuthEncoded)
	if err != nil {
		return nil, err
	}
	apiClient, err := helper.NewApiClient(
		tapdSource.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", basicAuth),
		},
		30*time.Second,
		tapdSource.Proxy,
		nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := apiClient.Get(input.Params["path"], input.Query, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// verify response body is json
	var tmp interface{}
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Status: resp.StatusCode, Body: json.RawMessage(body)}, nil
}
