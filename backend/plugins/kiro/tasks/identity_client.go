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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/identitystore"

	"github.com/apache/devlake/plugins/kiro/models"
)

// IdentityStoreAPI is the subset of the Identity Store API used here.
type IdentityStoreAPI interface {
	DescribeUser(input *identitystore.DescribeUserInput) (*identitystore.DescribeUserOutput, error)
}

// KiroIdentityClient resolves user ids to human-readable display names.
//
// This is entirely optional. Identity for joining to git history comes from the
// User_Email column of the report CSV, so a missing or misconfigured Identity
// Store degrades only the display name, never the data pipeline.
type KiroIdentityClient struct {
	IdentityStore IdentityStoreAPI
	StoreId       string
}

// NewKiroIdentityClient returns nil when Identity Store is not configured,
// which callers treat as "no display names" rather than as an error.
func NewKiroIdentityClient(connection *models.KiroConnection) (*KiroIdentityClient, error) {
	if connection.IdentityStoreId == "" || connection.IdentityStoreRegion == "" {
		return nil, nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(connection.IdentityStoreRegion),
		Credentials: credentials.NewStaticCredentials(
			connection.AccessKeyId,
			connection.SecretAccessKey,
			"",
		),
	})
	if err != nil {
		return nil, err
	}

	return &KiroIdentityClient{
		IdentityStore: identitystore.New(sess),
		StoreId:       connection.IdentityStoreId,
	}, nil
}

// ResolveDisplayName looks up a display name, returning nil when it cannot be
// determined.
//
// nil rather than the raw user id: the column exists for human readability, and
// storing an id there would make it look like a resolved name.
func (client *KiroIdentityClient) ResolveDisplayName(userId string) (*string, error) {
	if client == nil || client.IdentityStore == nil || userId == "" {
		return nil, nil
	}

	result, err := client.IdentityStore.DescribeUser(&identitystore.DescribeUserInput{
		IdentityStoreId: aws.String(client.StoreId),
		UserId:          aws.String(userId),
	})
	if err != nil {
		// Surfaced for logging, but callers proceed without a display name.
		return nil, err
	}

	if result.DisplayName != nil && *result.DisplayName != "" {
		name := *result.DisplayName
		return &name, nil
	}
	return nil, nil
}
