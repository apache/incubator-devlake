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
	Id     int `gorm:"column:id"`
	Action int `gorm:"column:action"`
	ZentaoRemoteDbHistoryBase
}

func (ZentaoRemoteDbHistory) TableName() string {
	return "zt_history"
}

type ZentaoRemoteDbAction struct {
	Id         int       `gorm:"column:aid"`
	ObjectType string    `gorm:"column:objectType"`
	ObjectId   int       `gorm:"column:objectID"`
	Product    string    `gorm:"column:product"`
	Project    int       `gorm:"column:project"`
	Execution  int       `gorm:"column:execution"`
	Actor      string    `gorm:"column:actor"`
	Action     string    `gorm:"column:action"`
	Date       time.Time `gorm:"column:date"`
	Comment    string    `gorm:"column:comment"`
	Extra      string    `gorm:"column:extra"`
	Read       string    `gorm:"column:read"`
	Version    string    `gorm:"column:version"`
	Efforted   string    `gorm:"column:efforted"`
}

func (ZentaoRemoteDbAction) TableName() string {
	return "zt_action"
}

type ZentaoRemoteDbActionHistory struct {
	ZentaoRemoteDbAction
	ZentaoRemoteDbHistoryBase

	HistoryId int `gorm:"column:hid"`
}

func (ah *ZentaoRemoteDbActionHistory) Convert() *ZentaoChangelogCom {
	return &ZentaoChangelogCom{
		&ZentaoChangelog{
			Id:         int64(ah.Id),
			ObjectId:   ah.ObjectId,
			Execution:  ah.Execution,
			Actor:      ah.Actor,
			Action:     ah.Action,
			Extra:      ah.Extra,
			ObjectType: ah.ObjectType,
			Project:    ah.Project,
			Product:    ah.Product,
			Version:    ah.Version,
			Comment:    ah.Comment,
			Efforted:   ah.Efforted,
			Date:       ah.Date,
			Read:       ah.Read,
		},
		&ZentaoChangelogDetail{
			Id:          int64(ah.HistoryId),
			ChangelogId: int64(ah.Id),
			Field:       ah.Field,
			Old:         ah.Old,
			New:         ah.New,
			Diff:        ah.Diff,
		},
	}
}
