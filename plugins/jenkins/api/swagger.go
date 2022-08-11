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

package api

type plan struct {
	Plugin   string   `json:"plugin" example:"jenkins"`
	Subtasks []string `json:"subtasks" example:"collectApiJobs,extractApiJobs,collectApiBuilds,extractApiBuilds,collectApiStages,extractApiStages,enrichApiBuilds,convertBuildsToCICD,convertStages,convertJobs,convertBuilds"`
	Options  struct {
		ConnectionID int `json:"connectionId" example:"1"`
	} `json:"options,omitempty"`
}
type Object struct {
	A string
	I int
}
type foo struct {
	SimpleArray []string `json:"simple_array" example:"hello,world,abc"`
	ObjectArray []Object `json:"object_array"`
}

type blueprintOutput struct {
	ID   int      `json:"id" example:"17"`
	Plan [][]plan `json:"plan"`
	blueprintInput
}
type blueprintInput struct {
	Name     string `json:"name"`
	Settings struct {
		Version     string `json:"version" example:"1.0.0"`
		Connections []struct {
			Plugin       string `json:"plugin" example:"jenkins"`
			ConnectionID int    `json:"connectionId" example:"1"`
			Scope        []struct {
				Transformation struct{} `json:"transformation"`
				Options        struct{} `json:"options"`
				Entities       []string `json:"entities" example:"CICD"`
			} `json:"scope"`
		} `json:"connections"`
	} `json:"settings"`
	CronConfig string `json:"cronConfig" example:"0 0 * * *"`
	Enable     bool   `json:"enable" example:"true"`
	Mode       string `json:"mode" example:"NORMAL"`
	IsManual   bool   `json:"isManual" example:"true"`
}
