package v8models

import (
	"encoding/json"
	"gorm.io/datatypes"

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

func (r RemoteLink) toToolLayer(sourceId, issueId uint64, raw json.RawMessage) *models.JiraRemotelink {
	return &models.JiraRemotelink{
		SourceId:     sourceId,
		RemotelinkId: r.ID,
		IssueId:      issueId,
		Self:         r.Self,
		Title:        r.Object.Title,
		Url:          r.Object.URL,
		RawJson:      datatypes.JSON(raw),
	}
}

func (RemoteLink) FromAPI(sourceId, issueId uint64, raw json.RawMessage) (interface{}, error) {
	var msgs []json.RawMessage
	err := json.Unmarshal(raw, &msgs)
	if err != nil {
		return nil, err
	}
	var list []*models.JiraRemotelink
	for _, msg := range msgs {
		var remoteLink RemoteLink
		err = json.Unmarshal(msg, &remoteLink)
		if err != nil {
			return nil, err
		}
		list = append(list, remoteLink.toToolLayer(sourceId, issueId, msg))
	}
	return list, nil
}

func (RemoteLink) ExtractRawMessage(blob []byte) (json.RawMessage, error) {
	return blob, nil
}
