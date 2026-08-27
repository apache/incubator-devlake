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

import (
	"fmt"
	"strings"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"gorm.io/gorm"
)

var _ plugin.ToolLayerScope = (*KiroS3Slice)(nil)

// KiroS3Slice is a collection scope: one AWS account for one month (or one
// whole year when Month is nil).
//
// The granularity mirrors Kiro's own S3 partitioning
// ({region}/{year}/{month}/{day}/{hour}), so listing a scope touches only the
// relevant prefixes instead of scanning the bucket. It also means a month can
// be re-collected in isolation, which is how historical backfill works.
//
// Kiro requires a bucket per AWS account holding subscriptions and does not
// support cross-account buckets, so AccountId belongs in the scope.
type KiroS3Slice struct {
	common.Scope `mapstructure:",squash"`
	// Id is a URL-safe logical identifier: {accountId}_{year}[_{month}]
	Id string `json:"id" mapstructure:"id" gorm:"primaryKey;type:varchar(512)"`
	// BasePath is an optional extra path segment before AWSLogs/
	BasePath string `json:"basePath,omitempty" mapstructure:"basePath" gorm:"type:varchar(512)"`
	// AccountId is the AWS account that Kiro exports for
	AccountId string `json:"accountId" mapstructure:"accountId" gorm:"type:varchar(255);not null"`
	Year      int    `json:"year" mapstructure:"year" gorm:"not null"`
	// Month is nil to collect a whole year
	Month *int `json:"month,omitempty" mapstructure:"month"`

	Name string `json:"name" mapstructure:"name" gorm:"-"`
}

func (KiroS3Slice) TableName() string {
	return "_tool_kiro_s3_slices"
}

// BeforeSave derives Id and validates required fields before persisting.
func (s *KiroS3Slice) BeforeSave(_ *gorm.DB) error {
	return s.normalize(true)
}

// AfterFind fills the derived display name for API responses.
func (s *KiroS3Slice) AfterFind(_ *gorm.DB) error {
	return s.normalize(false)
}

// normalize trims inputs and derives Id and Name.
//
// Unlike the predecessor scope this deliberately has no path-parsing fallback: the
// scope is always constructed from its component parts, never reverse-derived
// from a raw prefix string.
func (s *KiroS3Slice) normalize(strict bool) error {
	if s == nil {
		return nil
	}

	s.BasePath = strings.Trim(strings.TrimSpace(s.BasePath), "/")
	s.AccountId = strings.TrimSpace(s.AccountId)

	if strict {
		if s.AccountId == "" {
			return fmt.Errorf("accountId is required for a Kiro S3 slice")
		}
		if s.Year <= 0 {
			return fmt.Errorf("year is required for a Kiro S3 slice")
		}
	}
	if s.Month != nil && (*s.Month < 1 || *s.Month > 12) {
		return fmt.Errorf("month must be between 1 and 12, got %d", *s.Month)
	}

	if s.Id == "" && s.AccountId != "" && s.Year > 0 {
		s.Id = s.buildId()
	}
	s.Name = s.buildName()
	return nil
}

func (s *KiroS3Slice) buildId() string {
	if s.Month != nil {
		return fmt.Sprintf("%s_%04d_%02d", s.AccountId, s.Year, *s.Month)
	}
	return fmt.Sprintf("%s_%04d", s.AccountId, s.Year)
}

func (s *KiroS3Slice) buildName() string {
	if s.AccountId == "" || s.Year <= 0 {
		return s.Id
	}
	if s.Month != nil {
		return fmt.Sprintf("%s %04d-%02d", s.AccountId, s.Year, *s.Month)
	}
	return fmt.Sprintf("%s %04d", s.AccountId, s.Year)
}

// TimePath renders the year/month portion of an S3 prefix. A nil Month yields
// just the year, which widens collection to the whole year.
func (s *KiroS3Slice) TimePath() string {
	if s.Month != nil {
		return fmt.Sprintf("%04d/%02d", s.Year, *s.Month)
	}
	return fmt.Sprintf("%04d", s.Year)
}

func (s KiroS3Slice) ScopeId() string {
	return s.Id
}

func (s KiroS3Slice) ScopeName() string {
	if s.Name != "" {
		return s.Name
	}
	return s.buildName()
}

func (s KiroS3Slice) ScopeFullName() string {
	return s.ScopeName()
}

func (s KiroS3Slice) ScopeParams() interface{} {
	return &KiroS3SliceParams{
		ConnectionId: s.ConnectionId,
		AccountId:    s.AccountId,
		Year:         s.Year,
		Month:        s.Month,
	}
}

// Sanitize returns a copy ready for JSON serialization.
func (s KiroS3Slice) Sanitize() KiroS3Slice {
	_ = s.normalize(false)
	return s
}

// KiroS3SliceParams identifies a scope in raw data records.
type KiroS3SliceParams struct {
	ConnectionId uint64
	AccountId    string
	Year         int
	Month        *int
}
