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

package runner

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/migration"
	"github.com/apache/incubator-devlake/core/plugin"
	"sync"
)

var migrator plugin.Migrator

var lock sync.Mutex

// InitMigrator a Migrator singleton
func InitMigrator(basicRes context.BasicRes) (plugin.Migrator, errors.Error) {
	lock.Lock()
	defer lock.Unlock()

	if migrator != nil {
		return nil, errors.Internal.New("migrator singleton has already been initialized")
	}
	var err errors.Error
	migrator, err = migration.NewMigrator(basicRes)
	return migrator, err
}

// GetMigrator returns the shared Migrator singleton
func GetMigrator() plugin.Migrator {
	return migrator
}

/*
// RegisterMigrationScripts FIXME ...
func RegisterMigrationScripts(scripts []plugin.MigrationScript, comment string, config core.ConfigGetter, logger core.Logger) {
	for _, script := range scripts {
		if s, ok := script.(core.InjectConfigGetter); ok {
			s.SetConfigGetter(config)
		}
		if s, ok := script.(core.InjectLogger); ok {
			s.SetLogger(logger)
		}
	}
	migration.Register(scripts, comment)
}
*/
