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

package core

import (
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/migration"
)

// MigrationScript upgrades database to a newer version
type MigrationScript interface {
	Up(basicRes BasicRes) errors.Error
	Version() uint64
	Name() string
}

// Migrator is responsible for making sure the registered scripts get applied to database and only once
type Migrator interface {
	Register(scripts []MigrationScript, comment string)
	Execute() errors.Error
	HasPendingScripts() bool
}

// PluginMigration is implemented by the plugin to declare all migration script that have to be applied to the database
type PluginMigration interface {
	MigrationScripts() []MigrationScript
}

// Deprcated: Migratable is implemented by the plugin to declare all migration script that have to be applied to the database
type Migratable interface {
	MigrationScripts() []migration.Script
}
