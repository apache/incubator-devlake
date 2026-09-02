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

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	"github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/kiro/models"
	"github.com/apache/devlake/plugins/kiro/tasks"
)

// TestConnection validates a connection that has not been saved yet.
// @Summary test kiro connection
// @Description Test kiro connection
// @Tags plugins/kiro
// @Param body body models.KiroConn true "json body"
// @Success 200  {object} ConnectionReport
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/kiro/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connection models.KiroConnection
	// Wrapped as BadInput: a struct-tag validation failure is the caller's
	// problem, and an unwrapped Decode error surfaces as HTTP 500 - which sends
	// the user looking at server logs instead of at their own form.
	if err := api.Decode(input.Body, &connection, vld); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid connection payload")
	}
	if err := validateConnection(&connection.KiroConn); err != nil {
		return nil, errors.BadInput.Wrap(err, "connection validation failed")
	}
	if err := testConnection(&connection); err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Body:   buildConnectionReport(&connection),
		Status: http.StatusOK,
	}, nil
}

// TestExistingConnection validates a saved connection, optionally with overrides
// from the request body.
// @Summary test existing kiro connection
// @Description Test existing kiro connection
// @Tags plugins/kiro
// @Param id path int true "connection ID"
// @Success 200  {object} ConnectionReport
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/kiro/connections/{id}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.KiroConnection{}
	if err := connectionHelper.First(connection, input.Params); err != nil {
		return nil, errors.BadInput.Wrap(err, "find connection from db")
	}
	if err := api.DecodeMapStruct(input.Body, connection, false); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid connection payload")
	}
	if err := testConnection(connection); err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Body:   buildConnectionReport(connection),
		Status: http.StatusOK,
	}, nil
}

// ConnectionReport describes what the connection can actually see.
//
// A bare success/failure verdict is not enough to tell whether a configuration
// is right, because a wrong prefix and a genuinely empty period produce the same
// outcome: collection runs, finds nothing, and reports success. Returning the
// discovered accounts and per-stream object counts makes the difference visible
// before any scope is created.
type ConnectionReport struct {
	// Success and Message satisfy the shared Config UI connection-test contract.
	// The detailed fields below remain available to explain an empty export.
	Success bool   `json:"success"`
	Message string `json:"message"`

	ReportBucket    string `json:"reportBucket"`
	PromptLogBucket string `json:"promptLogBucket"`
	// Accounts are the AWS account ids found under the report prefix. An empty
	// list is the clearest sign that the bucket or prefix is wrong.
	Accounts []string `json:"accounts"`
	// Streams reports object counts for the most recent discovered period.
	Streams []tasks.StreamCount `json:"streams"`
	// Hint explains an empty result in plain terms.
	Hint string `json:"hint,omitempty"`
}

// connectionReportCountLimit caps counting so the check stays fast on a bucket
// holding hundreds of thousands of objects. The exact number does not matter for
// verifying a path - only whether it is zero.
const connectionReportCountLimit = 500

func newConnectionReport(connection *models.KiroConnection) *ConnectionReport {
	return &ConnectionReport{
		Success:         true,
		Message:         "success",
		ReportBucket:    connection.Bucket,
		PromptLogBucket: connection.GetPromptLogBucket(),
		Accounts:        []string{},
		Streams:         []tasks.StreamCount{},
	}
}

// buildConnectionReport probes the layout and summarizes what was found.
//
// Errors are folded into the report rather than returned: the connection itself
// is already known to work at this point, and a discovery failure is more useful
// shown as an empty result with a hint than as a failed request.
func buildConnectionReport(connection *models.KiroConnection) *ConnectionReport {
	report := newConnectionReport(connection)

	discovery, err := tasks.NewDiscovery(connection)
	if err != nil {
		report.Hint = "could not initialise S3 discovery: " + err.Error()
		return report
	}

	accounts, err := discovery.ListAccounts()
	if err != nil {
		report.Hint = "could not list accounts under the report prefix: " + err.Error()
		return report
	}
	report.Accounts = accounts

	if len(accounts) == 0 {
		report.Hint = fmt.Sprintf(
			"no account directories under s3://%s/%s/AWSLogs/ - check the bucket and report prefix",
			connection.Bucket, connection.GetReportPrefix())
		return report
	}

	// Probe the newest period that exists, since that is where data is most
	// likely to be and therefore the most informative check.
	accountId := accounts[len(accounts)-1]
	years, err := discovery.ListYears(accountId)
	if err != nil || len(years) == 0 {
		report.Hint = fmt.Sprintf(
			"account %s has no year directories - check the region in the connection", accountId)
		return report
	}
	year := years[len(years)-1]

	var month *int
	if months, monthErr := discovery.ListMonths(accountId, year); monthErr == nil && len(months) > 0 {
		latest := months[len(months)-1]
		month = &latest
	}

	report.Streams = discovery.CountStreams(accountId, year, month, connectionReportCountLimit)

	period := fmt.Sprintf("%04d", year)
	if month != nil {
		period = fmt.Sprintf("%04d-%02d", year, *month)
	}
	report.Hint = fmt.Sprintf("counts are for account %s, period %s", accountId, period)
	return report
}

// testConnection issues a real request against every bucket in use.
//
// Constructing a client proves nothing - the AWS SDK builds one happily from
// invalid credentials, so a test that stops there reports success for a
// connection that cannot read anything. A single-key list is the cheapest call
// that actually exercises credentials and bucket permissions.
//
// Both buckets are checked when reports and logs are separated, since they may
// carry different KMS keys and IAM conditions; a connection that can read
// reports but not logs would otherwise pass and then silently collect nothing.
func testConnection(connection *models.KiroConnection) errors.Error {
	clients, err := tasks.NewKiroS3Clients(connection)
	if err != nil {
		return err
	}

	for _, bucket := range clients.Buckets() {
		client := clients.Report
		if bucket == clients.PromptLog.Bucket {
			client = clients.PromptLog
		}
		if _, listErr := client.S3.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:  aws.String(bucket),
			MaxKeys: aws.Int64(1),
		}); listErr != nil {
			return errors.BadInput.Wrap(listErr, "cannot access s3 bucket "+bucket)
		}
	}

	return nil
}
