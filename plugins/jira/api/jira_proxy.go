package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
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
		return nil, fmt.Errorf("invalid sourceId")
	}
	jiraSource := &models.JiraSource{}
	err = lakeModels.Db.First(jiraSource, jiraSourceId).Error
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(jiraSource.Endpoint)
	if err != nil {
		return nil, err
	}
	path := input.Params["path"]
	u.Path = path
	u.RawQuery = input.Query.Encode()
	logger.Info("request to", u.String())
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+jiraSource.BasicAuthEncoded)
	client := &http.Client{Timeout: TimeOut}
	resp, err := client.Do(req)
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
