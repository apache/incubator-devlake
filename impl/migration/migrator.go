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

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
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
}

func (m *migratorImpl) loadExecuted() errors.Error {
	db := m.basicRes.GetDal()
	// make sure migration_history table exists
	err := db.AutoMigrate(&MigrationHistory{})
	if err != nil {
		return errors.Default.Wrap(err, "error performing migrations")
	}
	// load executed scripts into memory
	m.executed = make(map[string]bool)
	var records []MigrationHistory
	err = db.All(&records)
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
	db := m.basicRes.GetDal()
	for _, swc := range m.pending {
		err := swc.script.Up(m.basicRes)
		if err != nil {
			return err
		}
		err = db.Create(&MigrationHistory{
			ScriptVersion: swc.script.Version(),
			ScriptName:    swc.script.Name(),
			Comment:       swc.comment,
		})
		if err != nil {
			return errors.Convert(err)
		}
		m.executed[fmt.Sprintf("%s:%d", swc.script.Name(), swc.script.Version())] = true
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
	m := &migratorImpl{
		basicRes: basicRes,
	}
	err := m.loadExecuted()
	if err != nil {
		return nil, err
	}
	return m, nil
}
