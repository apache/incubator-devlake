package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/merico-dev/lake/plugins/helper"
)

type GitlabApiParams struct {
	ProjectId int
}

func GetTotalPagesFromResponse(res *http.Response, args *helper.ApiCollectorArgs) (int, error) {
	total := res.Header.Get("X-Total-Pages")
	if total == "" {
		return 0, nil
	}
	totalInt, err := convertStringToInt(total)
	if err != nil {
		return 0, err
	}
	return totalInt, nil
}

func GetRawMessageFromResponse(res *http.Response) ([]json.RawMessage, error) {
	rawMessages := &[]json.RawMessage{}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("%w %s", err, res.Request.URL.String())
	}
	err = json.Unmarshal(resBody, rawMessages)
	if err != nil {
		return nil, fmt.Errorf("%w %s %s", err, res.Request.URL.String(), string(resBody))
	}

	return *rawMessages, nil
}

func GetQuery(pager *helper.Pager) (url.Values, error) {
	query := url.Values{}
	query.Set("with_stats", "true")
	query.Set("page", strconv.Itoa(pager.Page))
	query.Set("per_page", strconv.Itoa(pager.Size))
	return query, nil
}
