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

package models

type ApiResponse struct {
	URL             string           `json:"url"`
	Jobs            []Job            `json:"jobs"`
	Mode            string           `json:"mode"`
	Views           []Views          `json:"views"`
	Class           string           `json:"_class"`
	NodeName        string           `json:"nodeName"`
	UseCrumbs       bool             `json:"useCrumbs"`
	Description     interface{}      `json:"description"`
	OverallLoad     OverallLoad      `json:"overallLoad"`
	PrimaryView     PrimaryView      `json:"primaryView"`
	UseSecurity     bool             `json:"useSecurity"`
	NumExecutors    int              `json:"numExecutors"`
	QuietingDown    bool             `json:"quietingDown"`
	UnlabeledLoad   UnlabeledLoad    `json:"unlabeledLoad"`
	AssignedLabels  []AssignedLabels `json:"assignedLabels"`
	SlaveAgentPort  int              `json:"slaveAgentPort"`
	NodeDescription string           `json:"nodeDescription"`
}
type Job struct {
	URL   string `json:"url"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Class string `json:"_class"`
	Jobs  *[]Job `json:"jobs"`
}
type Views struct {
	URL   string `json:"url"`
	Name  string `json:"name"`
	Class string `json:"_class"`
}
type OverallLoad struct {
}
type PrimaryView struct {
	URL   string `json:"url"`
	Name  string `json:"name"`
	Class string `json:"_class"`
}
type UnlabeledLoad struct {
	Class string `json:"_class"`
}
type AssignedLabels struct {
	Name string `json:"name"`
}

type ApiBuildResponse struct {
	Class             string    `json:"_class"`
	Number            int64     `json:"number"`
	Result            string    `json:"result"`
	Actions           []Actions `json:"actions"`
	Duration          float64   `json:"duration"`
	Timestamp         int64     `json:"timestamp"`
	DisplayName       string    `json:"displayName"`
	EstimatedDuration float64   `json:"estimatedDuration"`
	ChangeSet         ChangeSet `json:"changeSet"`
}
type LastBuiltRevision struct {
	SHA1 string `json:"SHA1"`
}

type Actions struct {
	Class                   string            `json:"_class,omitempty"`
	LastBuiltRevision       LastBuiltRevision `json:"lastBuiltRevision,omitempty"`
	MercurialRevisionNumber string            `json:"mercurialRevisionNumber"`
}
type ChangeSet struct {
	Class     string     `json:"_class"`
	Kind      string     `json:"kind"`
	Revisions []Revision `json:"revision"`
}

type Revision struct {
	Module   string
	Revision int
}
