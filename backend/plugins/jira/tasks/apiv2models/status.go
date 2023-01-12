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

type Status struct {
	Description string `json:"description"`
	IconURL     string `json:"iconUrl"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Scope       *struct {
		Type    string `json:"type"`
		Project struct {
			ID string `json:"id"`
		} `json:"project"`
	} `json:"scope"`
	Self           string `json:"self"`
	StatusCategory struct {
		ColorName string `json:"colorName"`
		ID        int    `json:"id"`
		Key       string `json:"key"`
		Name      string `json:"name"`
		Self      string `json:"self"`
	} `json:"statusCategory"`
	UntranslatedName string `json:"untranslatedName"`
}
