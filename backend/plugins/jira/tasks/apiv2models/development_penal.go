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

package apiv2models

type DevelopmentPanel struct {
	Errors []interface{} `json:"errors"`
	Detail []struct {
		Repositories []struct {
			Name    string `json:"name"`
			Avatar  string `json:"avatar"`
			URL     string `json:"url"`
			Commits []struct {
				ID              string `json:"id"`
				DisplayID       string `json:"displayId"`
				AuthorTimestamp string `json:"authorTimestamp"`
				URL             string `json:"url"`
				Author          struct {
					Name   string `json:"name"`
					Avatar string `json:"avatar"`
				} `json:"author"`
				FileCount int    `json:"fileCount"`
				Merge     bool   `json:"merge"`
				Message   string `json:"message"`
				Files     []struct {
					Path         string `json:"path"`
					URL          string `json:"url"`
					ChangeType   string `json:"changeType"`
					LinesAdded   int    `json:"linesAdded"`
					LinesRemoved int    `json:"linesRemoved"`
				} `json:"files"`
			} `json:"commits"`
		} `json:"repositories"`
		Instance struct {
			Name           string `json:"name"`
			BaseURL        string `json:"baseUrl"`
			Type           string `json:"type"`
			ID             string `json:"id"`
			TypeName       string `json:"typeName"`
			SingleInstance bool   `json:"singleInstance"`
		} `json:"_instance"`
	} `json:"detail"`
}
