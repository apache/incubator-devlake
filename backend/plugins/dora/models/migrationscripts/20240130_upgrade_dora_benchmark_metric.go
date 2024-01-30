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
)

type upgradeDoraBenchmarkMetric struct{}

func (u *upgradeDoraBenchmarkMetric) Up(baseRes context.BasicRes) errors.Error {
	db := baseRes.GetDal()
	err := db.Exec("UPDATE dora_benchmarks SET low = 'Fewer than once per month', medium = 'Between once per week and per month', high = 'Between once per day and per week' WHERE id = 1")
	if err != nil {
		return err
	}

	err = db.Exec("UPDATE dora_benchmarks SET low = 'More than one month', medium = 'Between one week and one month', high = 'Between one day and one week', elite = 'Less than one day' WHERE id = 2")
	if err != nil {
		return err
	}

	err = db.Exec("UPDATE dora_benchmarks SET low = '> 15%', medium = '10%-15%', high = '5%-10%', elite = '0-5%' WHERE id = 4")
	if err != nil {
		return err
	}

	return nil
}

func (*upgradeDoraBenchmarkMetric) Version() uint64 {
	return 20240130000002
}

func (*upgradeDoraBenchmarkMetric) Name() string {
	return "upgrade dora benchmark version to 2023 dora benchmark"
}
