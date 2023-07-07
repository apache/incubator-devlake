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
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"strconv"
	"strings"
)

var _ plugin.MigrationScript = (*addRawParamsTableForDomainScopes)(nil)

type addRawParamsTableForDomainScopes struct{}

type baseModel struct {
	Id            string `gorm:"column:id;type:varchar(255)"`
	RawDataTable  string `gorm:"column:_raw_data_table;type:varchar(255)"`
	RawDataParams string `gorm:"column:_raw_data_params;type:varchar(255)"`
}

func (script *addRawParamsTableForDomainScopes) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal().Begin()
	defer func() {
		if r := recover(); r != nil {
			err := db.Rollback()
			if err != nil {
				basicRes.GetLogger().Error(err, "error rolling back transaction")
			}
		}
	}()
	err := script.updateDomainScopeTables(db,
		&archived.Repo{},
		&archived.Board{},
		&archived.CicdScope{},
		&archived.CqProject{},
	)
	if err != nil {
		return err
	}
	err = db.Commit()
	return err
}

func (*addRawParamsTableForDomainScopes) Version() uint64 {
	return 20230630000001
}

func (*addRawParamsTableForDomainScopes) Name() string {
	return "populated _raw_data columns for domain scopes"
}

func (script *addRawParamsTableForDomainScopes) updateDomainScopeTables(db dal.Dal, models ...dal.Tabler) errors.Error {
	for _, model := range models {
		if err := script.updateDomainScopeTable(db, model.TableName()); err != nil {
			return err
		}
	}
	return nil
}

func (script *addRawParamsTableForDomainScopes) updateDomainScopeTable(db dal.Dal, tableName string) errors.Error {
	var scopes []*baseModel
	err := db.All(&scopes, dal.From(tableName))
	if err != nil || len(scopes) == 0 {
		return err
	}
	for _, scope := range scopes {
		IdParts := strings.Split(scope.Id, ":")
		if len(IdParts) < 4 {
			return errors.Default.New(fmt.Sprintf("malformed domain ID found in %s", tableName))
		}
		derivedPlugin := IdParts[0]
		derivedConnectionId := errors.Must1(strconv.ParseUint(IdParts[2], 10, 64))
		if _, err = plugin.GetPlugin(derivedPlugin); err != nil {
			return errors.Default.New(fmt.Sprintf("could not infer the plugin in context from the domainId: %s in table: %s",
				derivedPlugin, tableName))
		}
		scope.RawDataTable = fmt.Sprintf("_raw_%s_scopes", derivedPlugin)
		scope.RawDataParams = deriveRawDataParams(derivedPlugin, derivedConnectionId, scope.Id)
	}
	err = db.Update(&scopes, dal.From(tableName))
	return err
}

func deriveRawDataParams(pluginName string, connectionId uint64, scopeId string) string {
	switch pluginName {
	case "bamboo":
		type params20230630 struct {
			ConnectionId uint64
			ProjectKey   string
		}
		return toJsonString(&params20230630{
			ConnectionId: connectionId,
			ProjectKey:   scopeId,
		})
	case "bitbucket":
		type params20230630 struct {
			ConnectionId uint64
			FullName     string
		}
		return toJsonString(&params20230630{
			ConnectionId: connectionId,
			FullName:     scopeId,
		})
	case "github":
		type params20230630 struct {
			ConnectionId uint64
			Name         string
		}
		return toJsonString(&params20230630{
			ConnectionId: connectionId,
			Name:         scopeId,
		})
	case "gitlab":
		type params20230630 struct {
			ConnectionId uint64
			ProjectId    string
		}
		return toJsonString(&params20230630{
			ConnectionId: connectionId,
			ProjectId:    scopeId,
		})
	case "jenkins":
		type params20230630 struct {
			ConnectionId uint64
			FullName     string
		}
		return toJsonString(&params20230630{
			ConnectionId: connectionId,
			FullName:     scopeId,
		})
	case "jira":
		type params20230630 struct {
			ConnectionId uint64
			BoardId      string
		}
		return toJsonString(&params20230630{
			ConnectionId: connectionId,
			BoardId:      scopeId,
		})
	case "pagerduty":
		type params20230630 struct {
			ConnectionId uint64
			ScopeId      string
		}
		return toJsonString(&params20230630{
			ConnectionId: connectionId,
			ScopeId:      scopeId,
		})
	case "sonarqube":
		type params20230630 struct {
			ConnectionId uint64
			ProjectKey   string
		}
		return toJsonString(&params20230630{
			ConnectionId: connectionId,
			ProjectKey:   scopeId,
		})
	case "tapd":
		type params20230630 struct {
			ConnectionId uint64
			WorkspaceId  string
		}
		return toJsonString(&params20230630{
			ConnectionId: connectionId,
			WorkspaceId:  scopeId,
		})
	case "trello":
		type params20230630 struct {
			ConnectionId uint64
			BoardId      string
		}
		return toJsonString(&params20230630{
			ConnectionId: connectionId,
			BoardId:      scopeId,
		})
	case "zentao":
		type params20230630 struct {
			ConnectionId uint64
			ProjectId    string
		}
		return toJsonString(&params20230630{
			ConnectionId: connectionId,
			ProjectId:    scopeId,
		})
	default:
		panic(fmt.Sprintf("unknown plugin: %s", pluginName))
	}
}

func toJsonString(x any) string {
	return string(errors.Must1(json.Marshal(x)))
}
