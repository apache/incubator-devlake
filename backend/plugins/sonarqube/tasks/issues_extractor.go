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
	"crypto/sha256"
	"encoding/json"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
	"strings"
)

var _ plugin.SubTaskEntryPoint = ExtractIssues

func ExtractIssues(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUES_TABLE)
	hashCodeBlock := sha256.New()
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			body := &IssuesResponse{}
			err := errors.Convert(json.Unmarshal(resData.Data, body))
			if err != nil {
				return nil, err
			}
			sonarqubeIssue := &models.SonarqubeIssue{
				ConnectionId: data.Options.ConnectionId,
				IssueKey:     body.Key,
				Rule:         body.Rule,
				Severity:     body.Severity,
				Component:    body.Component,
				ProjectKey:   body.Project,
				Line:         body.Line,
				Status:       body.Status,
				Message:      body.Message,
				Author:       body.Author,
				Hash:         body.Hash,
				Type:         body.Type,
				Scope:        body.Scope,
				StartLine:    body.TextRange.StartLine,
				EndLine:      body.TextRange.EndLine,
				StartOffset:  body.TextRange.StartOffset,
				EndOffset:    body.TextRange.EndOffset,
				CreationDate: body.CreationDate,
				UpdateDate:   body.UpdateDate,
			}
			sonarqubeIssue.Debt = convertTimeToMinutes(body.Debt)
			if err != nil {
				return nil, err
			}
			sonarqubeIssue.Effort = convertTimeToMinutes(body.Effort)
			if err != nil {
				return nil, err
			}
			if len(body.Tags) > 0 {
				sonarqubeIssue.Tags = strings.Join(body.Tags, ",")
			}

			results := make([]interface{}, 0)
			results = append(results, sonarqubeIssue)
			for _, v := range body.Flows {
				for _, location := range v.Locations {
					codeBlock := &models.SonarqubeIssueCodeBlock{
						ConnectionId: data.Options.ConnectionId,
						IssueKey:     sonarqubeIssue.IssueKey,
						Component:    location.Component,
						Msg:          location.Msg,
						StartLine:    location.TextRange.StartLine,
						EndLine:      location.TextRange.EndLine,
						StartOffset:  location.TextRange.StartOffset,
						EndOffset:    location.TextRange.EndOffset,
					}
					generateId(hashCodeBlock, codeBlock)
					results = append(results, codeBlock)
				}
			}
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractIssuesMeta = plugin.SubTaskMeta{
	Name:             "ExtractIssues",
	EntryPoint:       ExtractIssues,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table sonarqube_issues",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

type IssuesResponse struct {
	Key       string `json:"key"`
	Rule      string `json:"rule"`
	Severity  string `json:"severity"`
	Component string `json:"component"`
	Project   string `json:"project"`
	Line      int    `json:"line"`
	Hash      string `json:"hash"`
	TextRange struct {
		StartLine   int `json:"startLine"`
		EndLine     int `json:"endLine"`
		StartOffset int `json:"startOffset"`
		EndOffset   int `json:"endOffset"`
	} `json:"textRange"`
	Flows             []flow              `json:"flows"`
	Status            string              `json:"status"`
	Message           string              `json:"message"`
	Effort            string              `json:"effort"`
	Debt              string              `json:"debt"`
	Author            string              `json:"author"`
	Tags              []string            `json:"tags"`
	CreationDate      *helper.Iso8601Time `json:"creationDate"`
	UpdateDate        *helper.Iso8601Time `json:"updateDate"`
	Type              string              `json:"type"`
	Scope             string              `json:"scope"`
	QuickFixAvailable bool                `json:"quickFixAvailable"`
}

type flow struct {
	Locations []Location `json:"locations"`
}
type TextRange struct {
	StartLine   int `json:"startLine"`
	EndLine     int `json:"endLine"`
	StartOffset int `json:"startOffset"`
	EndOffset   int `json:"endOffset"`
}
type Location struct {
	Component string    `json:"component"`
	TextRange TextRange `json:"textRange"`
	Msg       string    `json:"msg"`
}
