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

package plugininfo

import (
	"net/http"
	"reflect"
	"sync"

	"github.com/apache/incubator-devlake/api/shared"
	"github.com/apache/incubator-devlake/models/domainlayer/domaininfo"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/schema"
)

type SubTaskMeta struct {
	Name             string
	Required         bool
	EnabledByDefault bool
	Description      string
	DomainTypes      []string
}

func CreateSubTaskMeta(subTaskMeta []core.SubTaskMeta) []SubTaskMeta {
	ret := make([]SubTaskMeta, 0, len(subTaskMeta))
	for _, meta := range subTaskMeta {
		ret = append(ret, SubTaskMeta{
			Name:             meta.Name,
			Required:         meta.Required,
			EnabledByDefault: meta.EnabledByDefault,
			Description:      meta.Description,
			DomainTypes:      meta.DomainTypes,
		})
	}
	return ret
}

type TableInfo struct {
	TableName string
	Tags      string
}

type TableInfos struct {
	Field map[string]*TableInfo
	Error *string
}

func NewTableInfos(table core.Tabler) *TableInfos {
	tableInfos := &TableInfos{
		Field: make(map[string]*TableInfo),
		Error: nil,
	}

	fieldInfos := utils.WalkFields(reflect.TypeOf(table), nil)
	schema, err := schema.Parse(table, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		errstr := err.Error()
		tableInfos.Error = &errstr
	}
	for _, field := range fieldInfos {
		dbName := ""
		if schema != nil {
			if dbfield, ok := schema.FieldsByName[field.Name]; ok {
				dbName = dbfield.DBName
			}
		}
		tableInfos.Field[field.Name] = &TableInfo{TableName: dbName, Tags: string(field.Tag)}
	}

	return tableInfos
}

type PluginInfo struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Tables      map[string]*TableInfos `json:"tables"`
	TaskMeta    []SubTaskMeta          `json:"task_mata"`
}

func NewPluginInfo() *PluginInfo {
	return &PluginInfo{
		Tables: make(map[string]*TableInfos),
	}
}

type TotalInfo struct {
	DomainInfos map[string]*TableInfos
	PluginInfos map[string]*PluginInfo
}

func NewTotalInfo() *TotalInfo {
	return &TotalInfo{
		DomainInfos: make(map[string]*TableInfos),
		PluginInfos: make(map[string]*PluginInfo),
	}
}

// @Get detail info of plugins
// @Description GET /plugininfo
// @Description RETURN SAMPLE
// @Tags framework/plugininfo
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Router /plugininfo [get]
func Get(c *gin.Context) {
	info := NewTotalInfo()

	// set the domain layer tables info.
	domaininfo := domaininfo.GetDomainTablesInfo()
	domaininfoTable := make(map[string]*TableInfos)
	for _, table := range domaininfo {
		domaininfoTable[table.TableName()] = NewTableInfos(table)
	}
	info.DomainInfos = domaininfoTable

	// plugin info
	err := core.TraversalPlugin(func(name string, plugin core.PluginMeta) error {
		infoPlugin := NewPluginInfo()
		info.PluginInfos[name] = infoPlugin

		// plugin name and description
		infoPlugin.Name = name
		infoPlugin.Description = plugin.Description()

		// if this plugin has the plugin task info
		if pt, ok := plugin.(core.PluginTask); ok {
			infoPlugin.TaskMeta = CreateSubTaskMeta(pt.SubTaskMetas())
		}

		// if this plugin has the plugin model info
		if pm, ok := plugin.(core.PluginModel); ok {
			tables := pm.GetTablesInfo()
			infoPluginTable := make(map[string]*TableInfos)
			for _, table := range tables {
				infoPluginTable[table.TableName()] = NewTableInfos(table)
			}
			infoPlugin.Tables = infoPluginTable
		}

		return nil
	})

	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
	}

	shared.ApiOutputSuccess(c, info, http.StatusOK)
}
