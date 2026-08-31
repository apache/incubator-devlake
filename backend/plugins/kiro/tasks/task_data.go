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
	"fmt"

	"github.com/apache/incubator-devlake/plugins/kiro/models"
)

// DefaultWorkerCount is how many objects are fetched concurrently.
//
// Collection is bound by request count, not bandwidth: a single user can
// produce ~600 log objects on a busy day, each under a kilobyte. S3 sustains
// thousands of GETs per second per prefix, so 20 is well within limits while
// cutting a peak day from tens of minutes to a couple.
const DefaultWorkerCount = 20

// KiroOptions are the blueprint-supplied task options.
type KiroOptions struct {
	ConnectionId uint64 `json:"connectionId"`
	ScopeId      string `json:"scopeId"`
	AccountId    string `json:"accountId"`
	Year         int    `json:"year"`
	Month        *int   `json:"month"`
	// WorkerCount overrides DefaultWorkerCount when set above zero.
	WorkerCount int `json:"workerCount"`
}

// PrefixSpec is one S3 location to scan, along with the kind of file found
// there and which bucket holds it.
type PrefixSpec struct {
	Bucket   string
	Prefix   string
	FileType string
}

// KiroTaskData is shared by every subtask in a run.
type KiroTaskData struct {
	Options        *KiroOptions
	Connection     *models.KiroConnection
	S3Clients      *KiroS3Clients
	IdentityClient *KiroIdentityClient
	// Prefixes are the locations to scan, precomputed so the collector does not
	// re-derive paths.
	Prefixes []PrefixSpec
}

// WorkerCount returns the effective concurrency for this run.
func (d *KiroTaskData) WorkerCount() int {
	if d.Options != nil && d.Options.WorkerCount > 0 {
		return d.Options.WorkerCount
	}
	return DefaultWorkerCount
}

// BuildPrefixes derives the three S3 locations a scope covers.
//
// Layout confirmed against real exports:
//
//	{bucket}/{reportPrefix}/AWSLogs/{acct}/KiroLogs/user_report/{region}/{y}/{m}
//	{bucket}/{logPrefix}/AWSLogs/{acct}/KiroLogs/GenerateAssistantResponse/{region}/{y}/{m}
//	{bucket}/{logPrefix}/AWSLogs/{acct}/KiroLogs/GenerateCompletions/{region}/{y}/{m}
//
// The report path's hour segment is always 00 (reports are written at 02:00
// UTC) while log paths carry a real hour, but neither is included here: the
// prefix stops at the month so a scope lists the whole period in one sweep.
func BuildPrefixes(connection *models.KiroConnection, accountId string, timePath string) []PrefixSpec {
	region := connection.Region
	reportBase := fmt.Sprintf("%s/AWSLogs/%s/KiroLogs", connection.GetReportPrefix(), accountId)
	logBase := fmt.Sprintf("%s/AWSLogs/%s/KiroLogs", connection.GetPromptLogPrefix(), accountId)

	return []PrefixSpec{
		{
			Bucket:   connection.Bucket,
			Prefix:   fmt.Sprintf("%s/user_report/%s/%s", reportBase, region, timePath),
			FileType: models.FileTypeReport,
		},
		{
			Bucket:   connection.GetPromptLogBucket(),
			Prefix:   fmt.Sprintf("%s/GenerateAssistantResponse/%s/%s", logBase, region, timePath),
			FileType: models.FileTypeChatLog,
		},
		{
			Bucket:   connection.GetPromptLogBucket(),
			Prefix:   fmt.Sprintf("%s/GenerateCompletions/%s/%s", logBase, region, timePath),
			FileType: models.FileTypeCompletionLog,
		},
	}
}
