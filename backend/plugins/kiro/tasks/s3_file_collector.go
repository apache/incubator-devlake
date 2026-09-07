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
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/apache/devlake/core/dal"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	"github.com/apache/devlake/plugins/kiro/models"
)

var _ plugin.SubTaskEntryPoint = CollectKiroS3Files

// CollectKiroS3FilesMeta discovers which S3 objects exist for a scope.
var CollectKiroS3FilesMeta = plugin.SubTaskMeta{
	Name:             "collectKiroS3Files",
	EntryPoint:       CollectKiroS3Files,
	EnabledByDefault: true,
	Description:      "List Kiro export objects in S3 and record them for extraction",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

// CollectKiroS3Files lists every relevant object under the scope's prefixes and
// records the ones not seen before.
//
// Work is batched per listing page rather than per object. A single S3 page
// holds up to 1000 keys, so one SELECT and one INSERT per page replaces two
// round trips per file - at tens of thousands of objects a day that is the
// difference between dozens of queries and tens of thousands.
func CollectKiroS3Files(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*KiroTaskData)
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()

	taskCtx.SetProgress(0, -1)

	for _, spec := range data.Prefixes {
		client := data.S3Clients.ForFileType(spec.FileType)
		prefix := spec.Prefix
		if prefix != "" && !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}
		logger.Info("scanning s3://%s/%s for %s files", client.Bucket, prefix, spec.FileType)

		var continuationToken *string
		for {
			output, listErr := client.S3.ListObjectsV2(&s3.ListObjectsV2Input{
				Bucket:            aws.String(client.Bucket),
				Prefix:            aws.String(prefix),
				ContinuationToken: continuationToken,
			})
			if listErr != nil {
				return errors.Convert(listErr)
			}

			candidates := collectCandidates(output, client.Bucket, spec, data.Options)
			inserted, saveErr := saveNewFileMeta(db, data.Options.ConnectionId, candidates)
			if saveErr != nil {
				return saveErr
			}
			taskCtx.IncProgress(inserted)

			// IsTruncated is a pointer; dereferencing it unguarded panics on an
			// empty response.
			if output.IsTruncated == nil || !*output.IsTruncated {
				break
			}
			continuationToken = output.NextContinuationToken
		}
	}

	return nil
}

// collectCandidates turns one listing page into file metadata rows.
//
// Only .csv and .json.gz are kept. That filter also excludes the small
// extension-less objects AWS writes at the KiroLogs root as permission probes.
func collectCandidates(output *s3.ListObjectsV2Output, bucket string, spec PrefixSpec, options *KiroOptions) []*models.KiroS3FileMeta {
	candidates := make([]*models.KiroS3FileMeta, 0, len(output.Contents))
	for _, object := range output.Contents {
		if object.Key == nil {
			continue
		}
		key := *object.Key
		if !strings.HasSuffix(key, ".csv") && !strings.HasSuffix(key, ".json.gz") {
			continue
		}
		candidates = append(candidates, &models.KiroS3FileMeta{
			ConnectionId: options.ConnectionId,
			S3Path:       key,
			// Basename only. The full key lives in S3Path, which is sized for
			// it; putting a full key here would eventually overflow the column.
			FileName:  path.Base(key),
			Bucket:    bucket,
			ScopeId:   options.ScopeId,
			FileType:  spec.FileType,
			Processed: false,
		})
	}
	return candidates
}

// saveNewFileMeta inserts the rows that are not already recorded, returning how
// many were added.
//
// The existence check queries by (connection_id, s3_path), which is exactly the
// primary key. Querying on an unindexed column here would turn each page into a
// full table scan and the task would never finish at scale.
func saveNewFileMeta(db dal.Dal, connectionId uint64, candidates []*models.KiroS3FileMeta) (int, errors.Error) {
	if len(candidates) == 0 {
		return 0, nil
	}

	paths := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		paths = append(paths, candidate.S3Path)
	}

	var existingRows []models.KiroS3FileMeta
	err := db.All(&existingRows,
		dal.Select("s3_path"),
		dal.From(&models.KiroS3FileMeta{}),
		dal.Where("connection_id = ? AND s3_path IN ?", connectionId, paths),
	)
	if err != nil {
		return 0, errors.Default.Wrap(err, "failed to query existing kiro file metadata")
	}

	existing := make(map[string]struct{}, len(existingRows))
	for _, row := range existingRows {
		existing[row.S3Path] = struct{}{}
	}

	fresh := make([]*models.KiroS3FileMeta, 0, len(candidates))
	for _, candidate := range candidates {
		if _, seen := existing[candidate.S3Path]; seen {
			continue
		}
		fresh = append(fresh, candidate)
	}
	if len(fresh) == 0 {
		return 0, nil
	}

	if err := db.Create(fresh); err != nil {
		return 0, errors.Default.Wrap(err, "failed to record kiro file metadata")
	}
	return len(fresh), nil
}
