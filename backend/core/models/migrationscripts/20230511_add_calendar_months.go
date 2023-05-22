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
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type addCalendarMonths struct{}

type calendarMonths struct {
	Month          string    `gorm:"primaryKey;type:varchar(20)"`
	MonthTimestamp time.Time `gorm:"type:timestamp"`
}

func (calendarMonths) TableName() string {
	return "calendar_months"
}

func (u *addCalendarMonths) Up(baseRes context.BasicRes) errors.Error {
	db := baseRes.GetDal()
	err := migrationhelper.AutoMigrateTables(
		baseRes,
		&calendarMonths{},
	)
	if err != nil {
		return err
	}
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startDate := firstDayOfMonth.AddDate(-2, 0, 0)
	endDate := firstDayOfMonth.AddDate(10, 0, 0)

	var months []calendarMonths
	for d := endDate; d.After(startDate); d = d.AddDate(0, -1, 0) {
		month := calendarMonths{
			Month:          d.Format("06/01"),
			MonthTimestamp: time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.Local),
		}
		months = append(months, month)
	}
	return db.Create(&months)
}

func (*addCalendarMonths) Version() uint64 {
	return 20230511000001
}

func (*addCalendarMonths) Name() string {
	return "add calendar months"
}
