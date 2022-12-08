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

package migration

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/version"
)

type scriptWithComment struct {
	script  core.MigrationScript
	comment string
}

type migratorImpl struct {
	sync.Mutex
	basicRes core.BasicRes
	executed map[string]bool
	scripts  []*scriptWithComment
	pending  []*scriptWithComment
	db       dal.Dal
	tx       dal.Transaction
}

// LockingHistory is desgned for preventing mutiple delake instances from sharing the same database which may cause
// problems like #3537, #3466. It works by the following step:
//
// 1. Each devlake insert a record to this table whie `Succeeded=false`
// 2. Then it should try to lock the LockingStub table
// 3. Update the record with `Succeeded=true` if it had obtained the lock successfully
//
// NOTE: it works IFF all devlake instances obey the principle described above, in other words, this mechanism can
// not prevent older versions from sharing the same database
type LockingHistory struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	HostName  string
	Version   string
	Succeeded bool
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (LockingHistory) TableName() string {
	return "_devlake_locking_history"
}

// LockingStub does nothing but offer a locking target
type LockingStub struct {
	Stub string
}

func (LockingStub) TableName() string {
	return "_devlake_locking_stub"
}

func (m *migratorImpl) lockDatabase() errors.Error {
	// first, register the instance
	err := m.db.AutoMigrate(&LockingHistory{})
	if err != nil {
		return err
	}
	hostName, e := os.Hostname()
	if e != nil {
		return errors.Convert(e)
	}
	lockingHistory := &LockingHistory{
		HostName: hostName,
		Version:  version.Version,
	}
	err = m.db.Create(lockingHistory)
	if err != nil {
		return err
	}
	// 2. obtain the lock
	err = m.db.AutoMigrate(&LockingStub{})
	if err != nil {
		return err
	}
	m.tx = m.db.Begin()
	c := make(chan error, 1)

	// This prevent multiple devlake instances from sharing the same database by locking the migration history table
	// However, it would not work if any older devlake instances were already using the database.
	go func() {
		switch m.db.Dialect() {
		case "mysql":
			c <- m.tx.Exec("LOCK TABLE _devlake_locking_stub WRITE")
		case "postgres":
			c <- m.tx.Exec("LOCK TABLE _devlake_locking_stub IN EXCLUSIVE MODE")
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
	return m.db.Update(lockingHistory)
}

func (m *migratorImpl) loadExecuted() errors.Error {
	// make sure migration_history table exists
	err := m.db.AutoMigrate(&MigrationHistory{})
	if err != nil {
		return errors.Default.Wrap(err, "error performing migrations")
	}
	// load executed scripts into memory
	m.executed = make(map[string]bool)
	var records []MigrationHistory
	err = m.db.All(&records)
	if err != nil {
		return errors.Default.Wrap(err, "error finding migration history records")
	}
	for _, record := range records {
		m.executed[fmt.Sprintf("%s:%d", record.ScriptName, record.ScriptVersion)] = true
	}
	return nil
}

// Register a MigrationScript to the Migrator with comment message
func (m *migratorImpl) Register(scripts []core.MigrationScript, comment string) {
	m.Lock()
	defer m.Unlock()
	for _, script := range scripts {
		key := fmt.Sprintf("%s:%d", script.Name(), script.Version())
		swc := &scriptWithComment{
			script:  script,
			comment: comment,
		}
		m.scripts = append(m.scripts, swc)
		if !m.executed[key] {
			m.pending = append(m.pending, swc)
		}
	}
}

// Execute all registered migration script in order and mark them as executed in migration_history table
func (m *migratorImpl) Execute() errors.Error {
	// sort the scripts by version
	sort.Slice(m.pending, func(i, j int) bool {
		return m.pending[i].script.Version() < m.pending[j].script.Version()
	})
	// execute them one by one
	log := m.basicRes.GetLogger().Nested("migrator")
	for _, swc := range m.pending {
		scriptId := fmt.Sprintf("%d-%s", swc.script.Version(), swc.script.Name())
		log.Info("applying migratin script %s", scriptId)
		err := swc.script.Up(m.basicRes)
		if err != nil {
			return err
		}
		err = m.db.Create(&MigrationHistory{
			ScriptVersion: swc.script.Version(),
			ScriptName:    swc.script.Name(),
			Comment:       swc.comment,
		})
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("failed to execute migration script %s", scriptId))
		}
		m.executed[scriptId] = true
		m.pending = m.pending[1:]
	}
	return nil
}

// HasPendingScripts returns if there is any pending migration scripts
func (m *migratorImpl) HasPendingScripts() bool {
	return len(m.executed) > 0 && len(m.pending) > 0
}

// NewMigrator returns a new Migrator instance, which
// implemented based on migration_history from the same database
func NewMigrator(basicRes core.BasicRes) (core.Migrator, errors.Error) {
	db := basicRes.GetDal().Session(dal.SessionConfig{
		SkipDefaultTransaction: true,
		PrepareStmt:            false,
	})
	m := &migratorImpl{
		basicRes: basicRes,
		db:       db,
	}
	err := m.lockDatabase()
	if err != nil {
		return nil, err
	}
	err = m.loadExecuted()
	if err != nil {
		return nil, err
	}
	return m, nil
}
