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
	"strconv"

	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/plugins/kiro/models"
)

// Discovery walks the report prefix to find what data actually exists.
//
// Kiro's S3 layout encodes every scope dimension as a path segment:
//
//	{reportPrefix}/AWSLogs/{accountId}/KiroLogs/user_report/{region}/{year}/{month}/
//
// so a scope never has to be typed by hand. That matters beyond convenience: a
// mistyped prefix produces exactly the same outcome as a month with no data -
// collection succeeds and finds nothing - so hand-entered paths cannot be
// verified from the result.
type Discovery struct {
	clients    *KiroS3Clients
	connection *models.KiroConnection
}

func NewDiscovery(connection *models.KiroConnection) (*Discovery, errors.Error) {
	clients, err := NewKiroS3Clients(connection)
	if err != nil {
		return nil, err
	}
	return &Discovery{clients: clients, connection: connection}, nil
}

// reportRoot is the prefix holding the per-account directories.
func (d *Discovery) reportRoot() string {
	return fmt.Sprintf("%s/AWSLogs", d.connection.GetReportPrefix())
}

// accountReportPrefix is where one account's report months live.
func (d *Discovery) accountReportPrefix(accountId string) string {
	return fmt.Sprintf("%s/AWSLogs/%s/KiroLogs/user_report/%s",
		d.connection.GetReportPrefix(), accountId, d.connection.Region)
}

// ListAccounts returns the AWS account ids that have exported data.
//
// Kiro requires a bucket per account holding subscriptions and does not support
// cross-account buckets, so in practice this is usually one entry - but reading
// it from S3 removes the chance of a typo in a 12-digit number.
func (d *Discovery) ListAccounts() ([]string, errors.Error) {
	return d.clients.Report.ListSubPrefixes(d.reportRoot())
}

// ListYears returns the years with report data for an account.
func (d *Discovery) ListYears(accountId string) ([]int, errors.Error) {
	names, err := d.clients.Report.ListSubPrefixes(d.accountReportPrefix(accountId))
	if err != nil {
		return nil, err
	}
	return parseNumericSegments(names), nil
}

// ListMonths returns the months with report data for an account and year.
func (d *Discovery) ListMonths(accountId string, year int) ([]int, errors.Error) {
	prefix := fmt.Sprintf("%s/%04d", d.accountReportPrefix(accountId), year)
	names, err := d.clients.Report.ListSubPrefixes(prefix)
	if err != nil {
		return nil, err
	}
	return parseNumericSegments(names), nil
}

// StreamCount is how many collectable objects one stream holds.
type StreamCount struct {
	FileType string `json:"fileType"`
	Bucket   string `json:"bucket"`
	Prefix   string `json:"prefix"`
	Count    int    `json:"count"`
	// AtLeast is true when counting stopped at the cap, so Count is a floor.
	AtLeast bool `json:"atLeast"`
	// Error explains why a stream could not be counted, e.g. missing
	// permission on the log bucket while the report bucket is readable.
	Error string `json:"error,omitempty"`
}

// CountStreams reports the object count for each of the three streams.
//
// This is the answer to "is my configuration right?". Reporting counts per
// stream distinguishes the three cases that otherwise look identical: a wrong
// prefix (zero everywhere), a genuinely dormant stream (zero for one type -
// inline completions stopped being produced under agentic usage), and a
// permissions gap on one bucket (an error for the log streams only).
//
// countLimit caps the work per stream; pass 0 to count everything.
func (d *Discovery) CountStreams(accountId string, year int, month *int, countLimit int) []StreamCount {
	timePath := fmt.Sprintf("%04d", year)
	if month != nil {
		timePath = fmt.Sprintf("%04d/%02d", year, *month)
	}

	specs := BuildPrefixes(d.connection, accountId, timePath)
	results := make([]StreamCount, 0, len(specs))

	for _, spec := range specs {
		client := d.clients.ForFileType(spec.FileType)
		result := StreamCount{
			FileType: spec.FileType,
			Bucket:   client.Bucket,
			Prefix:   spec.Prefix,
		}
		count, atLeast, err := client.CountObjects(spec.Prefix, countLimit)
		if err != nil {
			// Recorded rather than returned: one unreadable stream should still
			// leave the others' counts visible, since that contrast is what
			// identifies a per-bucket permission problem.
			result.Error = err.Error()
		} else {
			result.Count = count
			result.AtLeast = atLeast
		}
		results = append(results, result)
	}

	return results
}

// parseNumericSegments keeps only the segments that are numbers, in order.
//
// S3 prefixes are strings, and a non-numeric directory would otherwise surface
// as a year or month.
func parseNumericSegments(names []string) []int {
	values := make([]int, 0, len(names))
	for _, name := range names {
		value, err := strconv.Atoi(name)
		if err != nil {
			continue
		}
		values = append(values, value)
	}
	return values
}
