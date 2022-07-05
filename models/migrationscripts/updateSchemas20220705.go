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

package migrationscripts

import (
	"context"
	"time"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type User struct {
	archived.DomainEntity
	Email string `gorm:"type:varchar(255)"`
	Name  string `gorm:"type:varchar(255)"`
}

type Account struct {
	archived.DomainEntity
	Email        string `gorm:"type:varchar(255)"`
	FullName     string `gorm:"type:varchar(255)"`
	UserName     string `gorm:"type:varchar(255)"`
	AvatarUrl    string `gorm:"type:varchar(255)"`
	Organization string `gorm:"type:varchar(255)"`
	CreatedDate  *time.Time
	Status       int
}

type UserAccount struct {
	UserId    string `gorm:"primaryKey;type:varchar(255)"`
	AccountId string `gorm:"primaryKey;type:varchar(255)"`
}

type Team struct {
	archived.DomainEntity
	Name         string `gorm:"type:varchar(255)"`
	Alias        string `gorm:"type:varchar(255)"`
	ParentId     string `gorm:"type:varchar(255)"`
	OrgId        string `gorm:"type:varchar(255)"`
	SortingIndex int
}

type TeamUser struct {
	TeamId string `gorm:"primaryKey;type:varchar(255)"`
	UserId string `gorm:"primaryKey;type:varchar(255)"`
}

type UpdateSchemas20220705 struct {
}

func (u *UpdateSchemas20220705) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().RenameTable("users", "accounts")
	if err != nil {
		return err
	}
	err = db.Migrator().DropColumn(&Account{}, "timezone")
	if err != nil {
		return err
	}
	err = db.Migrator().AutoMigrate(&User{})
	if err != nil {
		return err
	}
	err = db.Migrator().AutoMigrate(&Account{})
	if err != nil {
		return err
	}
	err = db.Migrator().AutoMigrate(&UserAccount{})
	if err != nil {
		return err
	}
	err = db.Migrator().AutoMigrate(&Team{})
	if err != nil {
		return err
	}
	err = db.Migrator().AutoMigrate(&TeamUser{})
	if err != nil {
		return err
	}
	return nil
}

func (*UpdateSchemas20220705) Version() uint64 {
	return 20220705141638
}

func (*UpdateSchemas20220705) Name() string {
	return "rename users to accounts, create users, user_accounts, teams, team_users"
}
