package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
)

const (
	TimeOut = 10 * time.Second
)

func Proxy(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	sourceId := input.Params["sourceId"]
	if sourceId == "" {
		return nil, fmt.Errorf("missing sourceid")
	}
	jiraSourceId, err := strconv.ParseUint(sourceId, 10, 64)
	if err != nil {
		return nil, err
	}
	jiraSource := &models.JiraSource{}
	err = db.First(jiraSource, jiraSourceId).Error
	if err != nil {
		return nil, err
	}
	encKey := cfg.GetString(core.EncodeKeyEnvStr)
	basicAuth, err := core.Decrypt(encKey, jiraSource.BasicAuthEncoded)
	if err != nil {
		return nil, err
	}
	apiClient, err := helper.NewApiClient(
		jiraSource.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %v", basicAuth),
		},
		30*time.Second,
		jiraSource.Proxy,
		nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := apiClient.Get(input.Params["path"], input.Query, nil)
	if err != nil {
		fmt.Println(err)
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
