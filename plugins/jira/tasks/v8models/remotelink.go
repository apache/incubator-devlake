package v8models

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/jira/models"
)

type RemoteLink struct {
	ID          uint64 `json:"id"`
	Self        string `json:"self"`
	GlobalID    string `json:"globalId"`
	Application struct {
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"application"`
	Relationship string `json:"relationship"`
	Object       struct {
		URL     string `json:"url"`
		Title   string `json:"title"`
		Summary string `json:"summary"`
		Icon    struct {
			URL16X16 string `json:"url16x16"`
			Title    string `json:"title"`
		} `json:"icon"`
		Status struct {
			Resolved bool `json:"resolved"`
			Icon     struct {
				URL16X16 string `json:"url16x16"`
				Title    string `json:"title"`
				Link     string `json:"link"`
			} `json:"icon"`
		} `json:"status"`
	} `json:"object"`
}

func (r RemoteLink) toToolLayer(sourceId, issueId uint64) *models.JiraRemotelink {
	return &models.JiraRemotelink{
		SourceId:     sourceId,
		RemotelinkId: r.ID,
		IssueId:      issueId,
		Self:         r.Self,
		Title:        r.Object.Title,
		Url:          r.Object.URL,
	}
}

func (RemoteLink) FromAPI(sourceId, issueId uint64, raw json.RawMessage) (interface{}, error) {
	var vv []RemoteLink
	err := json.Unmarshal(raw, &vv)
	if err != nil {
		return nil, err
	}
	list := make([]*models.JiraRemotelink, len(vv))
	for i, item := range vv {
		list[i] = item.toToolLayer(sourceId, issueId)
	}
	return list, nil
}

func (RemoteLink) ExtractRawMessage(blob []byte) (json.RawMessage, error) {
	return blob, nil
}
