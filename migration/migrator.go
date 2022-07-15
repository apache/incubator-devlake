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
	"context"
	"fmt"
	"gorm.io/gorm"
	"sort"
	"sync"
)

var m = migrator{}

type scriptWithComment struct {
	Script
	comment string
}
type migrator struct {
	sync.Mutex
	db      *gorm.DB
	executed map[string]bool
	scripts []*scriptWithComment
	pending []*scriptWithComment
}

func Init(db *gorm.DB) {
	m.db = db
	var err error
	m.executed, err = m.getExecuted()
	if err != nil {
		panic(err)
	}
}

func (m *migrator) register(scripts []Script, comment string) {
	m.Lock()
	defer m.Unlock()
	for _, script := range scripts {
		key := fmt.Sprintf("%s:%d", script.Name(), script.Version())
		swc := &scriptWithComment{
			Script:  script,
			comment: comment,
		}
		m.scripts = append(m.scripts, swc)
		if !m.executed[key] {
			m.pending = append(m.pending, swc)
		}
	}
}

func (m *migrator) bookKeep(script *scriptWithComment) error {
	record := &MigrationHistory{
		ScriptVersion: script.Version(),
		ScriptName:    script.Name(),
		Comment:       script.comment,
	}
	return m.db.Create(record).Error
}

func (m *migrator) execute(ctx context.Context) error {
	sort.Slice(m.pending, func(i, j int) bool {
		return m.pending[i].Version() < m.pending[j].Version()
	})
	for _, script := range m.pending {
		err := script.Up(ctx, m.db)
		if err != nil {
			return err
		}
		err = m.bookKeep(script)
		if err != nil {
			return err
		}
	}
	return nil
}
func (m *migrator) getExecuted() (map[string]bool, error) {
	var err error
	versions := make(map[string]bool)
	err = m.db.Migrator().AutoMigrate(&MigrationHistory{})
	if err != nil {
		return nil, err
	}
	var records []MigrationHistory
	err = m.db.Find(&records).Error
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		versions[fmt.Sprintf("%s:%d", record.ScriptName, record.ScriptVersion)] = true
	}
	return versions, nil
}

func Register(scripts []Script, comment string) {
	m.register(scripts, comment)
}

func Execute(ctx context.Context) error {
	return m.execute(ctx)
}

func NeedConfirmation() bool {
	return len(m.executed) > 0 && len(m.pending) > 0
}
