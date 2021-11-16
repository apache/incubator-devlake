package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/tasks"
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
	client, err := tasks.NewJiraApiClientBySourceId(jiraSourceId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(input.Params["path"], &input.Query, nil)
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
	return &core.ApiResourceOutput{Body: json.RawMessage(body)}, nil
}
