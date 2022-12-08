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
	"os"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/version"
)

// long last transaction for database locking
var lockingTx dal.Transaction

// lockDatabase prevents multiple devlake instances from sharing the same lockDatabase
// check the models.LockingHistory for the detail
func lockDatabase(db dal.Dal) errors.Error {
	// first, register the instance
	err := db.AutoMigrate(&models.LockingHistory{})
	if err != nil {
		return err
	}
	hostName, e := os.Hostname()
	if e != nil {
		return errors.Convert(e)
	}
	lockingHistory := &models.LockingHistory{
		HostName: hostName,
		Version:  version.Version,
	}
	err = db.Create(lockingHistory)
	if err != nil {
		return err
	}
	// 2. obtain the lock
	err = db.AutoMigrate(&models.LockingStub{})
	if err != nil {
		return err
	}
	lockingTx = db.Begin()
	c := make(chan error, 1)

	// This prevent multiple devlake instances from sharing the same database by locking the migration history table
	// However, it would not work if any older devlake instances were already using the database.
	go func() {
		switch db.Dialect() {
		case "mysql":
			c <- lockingTx.Exec("LOCK TABLE _devlake_locking_stub WRITE")
		case "postgres":
			c <- lockingTx.Exec("LOCK TABLE _devlake_locking_stub IN EXCLUSIVE MODE")
		}
	}()

	select {
	case err := <-c:
		if err != nil {
			return errors.Convert(err)
		}
	case <-time.After(2 * time.Second):
		return errors.Default.New("locking _devlake_locking_stub timeout, the database might be locked by another devlake instance")
	}
	// 3. update the record
	lockingHistory.Succeeded = true
	return db.Update(lockingHistory)
}
