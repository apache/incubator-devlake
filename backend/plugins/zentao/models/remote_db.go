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
	"time"
)

type ZentaoRemoteDbHistoryBase struct {
	Field string `gorm:"column:field"`
	Old   string `gorm:"column:old"`
	New   string `gorm:"column:new"`
	Diff  string `gorm:"column:diff"`
}

type ZentaoRemoteDbHistory struct {
	Id     int64 `gorm:"column:id"`
	Action int64 `gorm:"column:action"`
	ZentaoRemoteDbHistoryBase
}

func (h ZentaoRemoteDbHistory) ToChangelogDetail(connectionId uint64) *ZentaoChangelogDetail {
	return &ZentaoChangelogDetail{
		ConnectionId: connectionId,
		Id:           h.Id,
		ChangelogId:  h.Action,
		Field:        h.Field,
		Old:          h.Old,
		New:          h.New,
		Diff:         h.Diff,
	}
}

func (ZentaoRemoteDbHistory) TableName() string {
	return "zt_history"
}

type ZentaoRemoteDbAction struct {
	Id         int64     `gorm:"column:id"`
	ObjectType string    `gorm:"column:objectType"`
	ObjectId   int64     `gorm:"column:objectID"`
	Product    string    `gorm:"column:product"`
	Project    int64     `gorm:"column:project"`
	Execution  int64     `gorm:"column:execution"`
	Actor      string    `gorm:"column:actor"`
	Action     string    `gorm:"column:action"`
	Date       time.Time `gorm:"column:date"`
	Comment    string    `gorm:"column:comment"`
	Extra      string    `gorm:"column:extra"`
	Read       string    `gorm:"column:read"`
	Vision     string    `gorm:"column:vision"`
	Efforted   string    `gorm:"column:efforted"`
}

func (a ZentaoRemoteDbAction) ToChangelog(connectionId uint64) *ZentaoChangelog {
	return &ZentaoChangelog{
		ConnectionId: connectionId,
		Id:           a.Id,
		ObjectId:     a.ObjectId,
		Execution:    a.Execution,
		Actor:        a.Actor,
		Action:       a.Action,
		Extra:        a.Extra,
		ObjectType:   a.ObjectType,
		Project:      a.Project,
		Vision:       a.Vision,
		Comment:      a.Comment,
		Efforted:     a.Efforted,
		Date:         a.Date,
		Read:         a.Read,
	}
}

func (ZentaoRemoteDbAction) TableName() string {
	return "zt_action"
}

type ZentaoRemoteDbActionHistory struct {
	ZentaoRemoteDbAction
	ZentaoRemoteDbHistoryBase

	ActionId  int64 `gorm:"column:aid"`
	HistoryId int64 `gorm:"column:hid"`
}

func (ah *ZentaoRemoteDbActionHistory) Convert(connectId uint64) *ZentaoChangelogCom {
	return &ZentaoChangelogCom{
		&ZentaoChangelog{
			ConnectionId: connectId,
			Id:           ah.ActionId,
			ObjectId:     ah.ObjectId,
			Execution:    ah.Execution,
			Actor:        ah.Actor,
			Action:       ah.Action,
			Extra:        ah.Extra,
			ObjectType:   ah.ObjectType,
			Project:      ah.Project,
			Vision:       ah.Vision,
			Comment:      ah.Comment,
			Efforted:     ah.Efforted,
			Date:         ah.Date,
			Read:         ah.Read,
		},
		&ZentaoChangelogDetail{
			ConnectionId: connectId,
			Id:           ah.HistoryId,
			ChangelogId:  ah.ActionId,
			Field:        ah.Field,
			Old:          ah.Old,
			New:          ah.New,
			Diff:         ah.Diff,
		},
	}
}
