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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type adddoraBenchmark2023 struct{}

type doraBenchmarkBasic struct {
	archived.Model
	Metric     string `gorm:"type:varchar(255)"`
	Low        string `gorm:"type:varchar(255)"`
	Medium     string `gorm:"type:varchar(255)"`
	High       string `gorm:"type:varchar(255)"`
	Elite      string `gorm:"type:varchar(255)"`
	DoraReport string `gorm:"type:varchar(20)"`
}

func (doraBenchmarkBasic) TableName() string {
	return "dora_benchmarks"
}

func (u *adddoraBenchmark2023) Up(baseRes context.BasicRes) errors.Error {
	db := baseRes.GetDal()
	err := db.DropTables(&doraBenchmark{})
	if err != nil {
		return err
	}
	err = migrationhelper.AutoMigrateTables(
		baseRes,
		&doraBenchmarkBasic{},
	)
	if err != nil {
		return err
	}

	// 2021 benchmarks
	doraBenchmark2021DF := &doraBenchmarkBasic{
		Model: archived.Model{
			ID: 1,
		},
		Metric:     "Deployment frequency",
		Low:        "Fewer than once per six months(low)",
		Medium:     "Between once per month and once every 6 months(medium)",
		High:       "Between once per day and once per month(high)",
		Elite:      "On-demand(elite)",
		DoraReport: "2021",
	}
	err = db.Create(doraBenchmark2021DF)
	if err != nil {
		return errors.Convert(err)
	}

	doraBenchmark2021LTC := &doraBenchmarkBasic{
		Model: archived.Model{
			ID: 2,
		},
		Metric:     "Lead time for changes",
		Low:        "More than six months(low)",
		Medium:     "Between one week and six months(medium)",
		High:       "Less than one week(high)",
		Elite:      "Less than one hour(elite)",
		DoraReport: "2021",
	}
	err = db.Create(doraBenchmark2021LTC)
	if err != nil {
		return errors.Convert(err)
	}

	doraBenchmark2021CFR := &doraBenchmarkBasic{
		Model: archived.Model{
			ID: 3,
		},
		Metric:     "Change failure rate",
		Low:        "> 30%(low)",
		Medium:     "21%-30%(medium)",
		High:       "16%-20%(high)",
		Elite:      "0-15%(elite)",
		DoraReport: "2021",
	}
	err = db.Create(doraBenchmark2021CFR)
	if err != nil {
		return errors.Convert(err)
	}

	doraBenchmark2021TTS := &doraBenchmarkBasic{
		Model: archived.Model{
			ID: 4,
		},
		Metric:     "Time to restore service",
		Low:        "More than one week(low)",
		Medium:     "Between one day and one week(medium)",
		High:       "Less than one day(high)",
		Elite:      "Less than one hour(elite)",
		DoraReport: "2021",
	}
	err = db.Create(doraBenchmark2021TTS)
	if err != nil {
		return errors.Convert(err)
	}

	// 2023 benchmarks
	doraBenchmark2023DF := &doraBenchmarkBasic{
		Model: archived.Model{
			ID: 5,
		},
		Metric:     "Deployment frequency",
		Low:        "Fewer than once per month(low)",
		Medium:     "Between once per week and once per month(medium)",
		High:       "Between once per day and once per week(high)",
		Elite:      "On-demand(elite)",
		DoraReport: "2023",
	}
	err = db.Create(doraBenchmark2023DF)
	if err != nil {
		return errors.Convert(err)
	}

	doraBenchmark2023LTC := &doraBenchmarkBasic{
		Model: archived.Model{
			ID: 6,
		},
		Metric:     "Lead time for changes",
		Low:        "More than one month(low)",
		Medium:     "Between one week and one month(medium)",
		High:       "Between one day and one week(high)",
		Elite:      "Less than one day(elite)",
		DoraReport: "2023",
	}
	err = db.Create(doraBenchmark2023LTC)
	if err != nil {
		return errors.Convert(err)
	}

	doraBenchmark2023CFR := &doraBenchmarkBasic{
		Model: archived.Model{
			ID: 7,
		},
		Metric:     "Change failure rate",
		Low:        "> 15%(low)",
		Medium:     "10%-15%(medium)",
		High:       "5%-10%(high)",
		Elite:      "0-5%(elite)",
		DoraReport: "2023",
	}
	err = db.Create(doraBenchmark2023CFR)
	if err != nil {
		return errors.Convert(err)
	}

	doraBenchmark2023FDRT := &doraBenchmarkBasic{
		Model: archived.Model{
			ID: 8,
		},
		Metric:     "Failed deployment recovery time",
		Low:        "More than one week(low)",
		Medium:     "Between one day and one week(medium)",
		High:       "Less than one day(high)",
		Elite:      "Less than one hour(elite)",
		DoraReport: "2023",
	}
	err = db.Create(doraBenchmark2023FDRT)
	if err != nil {
		return errors.Convert(err)
	}

	return nil
}

func (*adddoraBenchmark2023) Version() uint64 {
	return 20240223000003
}

func (*adddoraBenchmark2023) Name() string {
	return "add dora 2023 benchmark"
}
