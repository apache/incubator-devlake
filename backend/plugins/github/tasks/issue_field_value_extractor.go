/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiIssueFieldValuesMeta)
}

var ExtractApiIssueFieldValuesMeta = plugin.SubTaskMeta{
	Name:             "Extract Issue Field Values",
	EntryPoint:       ExtractApiIssueFieldValues,
	EnabledByDefault: false,
	Description:      "Extract raw issue field value data into the tool layer table _tool_github_issue_field_values",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	DependencyTables: []string{RAW_ISSUE_FIELD_VALUE_TABLE},
	ProductTables:    []string{models.GithubIssueFieldValue{}.TableName()},
}

type IssueFieldValueSelectOption struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type IssueFieldValueResponse struct {
	IssueFieldId       int                            `json:"issue_field_id"`
	IssueFieldName     string                         `json:"issue_field_name"`
	DataType           string                         `json:"data_type"`
	Value              json.RawMessage                `json:"value"`
	SingleSelectOption *IssueFieldValueSelectOption   `json:"single_select_option"`
	MultiSelectOptions []*IssueFieldValueSelectOption `json:"multi_select_options"`
}

func ExtractApiIssueFieldValues(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)
	extractor, err := helper.NewStatefulApiExtractor(&helper.StatefulApiExtractorArgs[IssueFieldValueResponse]{
		SubtaskCommonArgs: &helper.SubtaskCommonArgs{
			SubTaskContext: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_ISSUE_FIELD_VALUE_TABLE,
		},
		Extract: func(body *IssueFieldValueResponse, row *helper.RawData) ([]any, errors.Error) {
			if body.IssueFieldId == 0 {
				return nil, nil
			}
			issue := &SimpleIssue{}
			if err := errors.Convert(json.Unmarshal(row.Input, issue)); err != nil {
				return nil, err
			}

			fieldValue := &models.GithubIssueFieldValue{
				ConnectionId: data.Options.ConnectionId,
				IssueId:      issue.GithubId,
				FieldId:      body.IssueFieldId,
				FieldName:    body.IssueFieldName,
				DataType:     body.DataType,
				Value:        displayValue(body),
				RawValue:     string(body.Value),
			}
			if body.SingleSelectOption != nil {
				fieldValue.OptionColor = body.SingleSelectOption.Color
			}
			return []any{fieldValue}, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}

// displayValue renders the API value as queryable text. Select options are preferred over
// the raw value because the option objects carry the canonical names.
func displayValue(body *IssueFieldValueResponse) string {
	if body.SingleSelectOption != nil {
		return body.SingleSelectOption.Name
	}
	if len(body.MultiSelectOptions) > 0 {
		names := make([]string, 0, len(body.MultiSelectOptions))
		for _, option := range body.MultiSelectOptions {
			if option != nil {
				names = append(names, option.Name)
			}
		}
		return strings.Join(names, ",")
	}
	return scalarValue(body.Value)
}

// scalarValue unwraps a JSON scalar into text. Numbers are rendered without a trailing
// ".0" so an effort of 5 reads as "5" rather than "5.0" in a dashboard.
func scalarValue(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		return asString
	}
	var asNumber float64
	if err := json.Unmarshal(raw, &asNumber); err == nil {
		return strconv.FormatFloat(asNumber, 'f', -1, 64)
	}
	var asBool bool
	if err := json.Unmarshal(raw, &asBool); err == nil {
		return strconv.FormatBool(asBool)
	}
	// Arrays of plain strings can arrive without option objects on multi_select.
	var asStrings []string
	if err := json.Unmarshal(raw, &asStrings); err == nil {
		return strings.Join(asStrings, ",")
	}
	return string(raw)
}
