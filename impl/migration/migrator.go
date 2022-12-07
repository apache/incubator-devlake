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
	"sort"
	"sync"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
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

func (m *migratorImpl) loadExecuted() errors.Error {
	// make sure migration_history table exists
	err := m.tx.AutoMigrate(&MigrationHistory{})
	if err != nil {
		return errors.Default.Wrap(err, "error performing migrations")
	}
	// load executed scripts into memory
	m.executed = make(map[string]bool)
	var records []MigrationHistory
	err = m.tx.All(&records)
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
		err = m.tx.Create(&MigrationHistory{
			ScriptVersion: swc.script.Version(),
			ScriptName:    swc.script.Name(),
			Comment:       swc.comment,
		})
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("failed to execute migration script %s", scriptId))
		}
		m.executed[scriptId] = true
		m.pending = m.pending[1:]
		err = m.tx.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

// HasPendingScripts returns if there is any pending migration scripts
func (m *migratorImpl) HasPendingScripts() bool {
	return len(m.executed) > 0 && len(m.pending) > 0
}

// func lockMigrationHistory(db dal.Dal) (dal.Transaction, errors.Error) {
// 	println("start lock")
// 	err := tx.Exec("")
// 	// err := tx.First(migrationHistory, dal.Lock(true, true))
// 	println("end lock", err)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return tx, nil
// }

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
		tx:       db.Begin(),
	}

	c := make(chan error, 1)

	// This prevent multiple devlake instances from sharing the same database by locking the migration history table
	// However, it would not work if any older devlake instances were already using the database.
	go func() {
		c <- m.tx.Exec("LOCK TABLE _devlake_migration_history WRITE")
	}()

	select {
	case err := <-c:
		if err != nil {
			return nil, errors.Convert(err)
		}
	case <-time.After(2 * time.Second):
		return nil, errors.Default.New("locking _devlake_migration_history timeout, the database might be locked by another devlake instance")
	}
	err := m.loadExecuted()
	if err != nil {
		return nil, err
	}
	return m, nil
}
