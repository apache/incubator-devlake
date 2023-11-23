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

package services

import (
	"fmt"
	"os"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/version"
)

// long last transaction for database locking
var lockingTx dal.Transaction

// lockDatabase prevents multiple devlake instances from sharing the same lockDatabase
// check the models.LockingHistory for the detail
func lockDatabase() {
	db := basicRes.GetDal()
	// first, register the instance
	errors.Must(db.AutoMigrate(&models.LockingHistory{}))
	hostName := errors.Must1(os.Hostname())
	lockingHistory := &models.LockingHistory{
		HostName: hostName,
		Version:  version.Version,
	}
	errors.Must(db.Create(lockingHistory))
	// 2. obtain the lock: using a never released transaction
	// This prevents multiple devlake instances from sharing the same database by locking the migration history table
	// However, it would not work if any older devlake instances were already using the database.
	lockingTx = db.Begin()
	c := make(chan bool, 1)
	go func() {
		errors.Must(lockingTx.DropTables(models.LockingStub{}.TableName()))
		errors.Must(lockingTx.AutoMigrate(&models.LockingStub{}))
		errors.Must(lockingTx.LockTables(dal.LockTables{{Table: "_devlake_locking_stub", Exclusive: true}}))
		lockingHistory.Succeeded = true
		errors.Must(db.Update(lockingHistory))
		c <- true
	}()

	// 3. update the record
	select {
	case <-c:
	case <-time.After(10 * time.Second):
		panic(fmt.Errorf("locking _devlake_locking_stub timeout, the database might be locked by another devlake instance"))
	}
}
