package tasks

import (
	"encoding/json"
	"testing"

	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/stretchr/testify/require"
)

func TestDetermineIssueType(t *testing.T) {
	t.Run("no custom fileds", func(t *testing.T) {
		task := Task{}
		require.Equal(t, ticket.TASK, determineIssueType(&task))
	})
	t.Run("maps to bug", func(t *testing.T) {
		task := Task{
			CustomFields: []TaskCustomField{
				{
					Name: "Type",
					Type: "drop_down",
					TypeConfig: toRawJson(TaskCustomFieldDropDown{
						Default:     0,
						Placeholder: new(string),
						Options: []TypeCustomFieldDropDownOption{
							{
								Name: "Incident",
							},
							{
								Name: "Bug",
							},
							{
								Name: "Improvement",
							},
						},
					}),
					Value: float64(1),
				},
			},
		}
		require.Equal(t, ticket.BUG, determineIssueType(&task))
	})
	t.Run("maps to incident", func(t *testing.T) {
		task := Task{
			CustomFields: []TaskCustomField{
				{
					Name: "Type",
					Type: "drop_down",
					TypeConfig: toRawJson(TaskCustomFieldDropDown{
						Default:     0,
						Placeholder: new(string),
						Options: []TypeCustomFieldDropDownOption{
							{
								Name: "Incident",
							},
							{
								Name: "Bug",
							},
							{
								Name: "Improvement",
							},
						},
					}),
					Value: float64(0),
				},
			},
		}
		require.Equal(t, ticket.INCIDENT, determineIssueType(&task))
	})
}

func toRawJson(taskCustomFieldDropDown TaskCustomFieldDropDown) *json.RawMessage {
	bytes, err := json.Marshal(taskCustomFieldDropDown)
	if err != nil {
		panic(err)
	}
	raw := json.RawMessage(bytes)
	return &raw
}
