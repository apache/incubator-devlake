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

package tasks

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ plugin.SubTaskEntryPoint = ConvertChangelog

const RAW_CHANGELOG_TABLE = "zt_history"

var ConvertChangelogMeta = plugin.SubTaskMeta{
	Name:             "ConvertChangelog",
	EntryPoint:       ConvertChangelog,
	EnabledByDefault: true,
	Description:      "convert Zentao changelog",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

type ZentaoChangelogSelect struct {
	common.NoPKModel `json:"-"`

	models.ZentaoChangelog
	ChangelogId int64  `json:"changelogId" mapstructure:"changelogId" gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Field       string `json:"field" mapstructure:"field"`
	Old         string `json:"old" mapstructure:"old"`
	New         string `json:"new" mapstructure:"new"`
	Diff        string `json:"diff" mapstructure:"diff"`

	Account  string `json:"account" gorm:"type:varchar(100);index"`
	Avatar   string `json:"avatar" gorm:"type:varchar(255)"`
	Realname string `json:"realname" gorm:"type:varchar(100);index"`
	Role     string `json:"role" gorm:"type:varchar(100);index"`
	Dept     int64  `json:"dept" gorm:"type:BIGINT  NOT NULL;index"`

	CID  int64 `json:"cid" mapstructure:"cid" gorm:"column:cid"`
	CDID int64 `json:"cdid" mapstructure:"cdid" gorm:"column:cdid"`
	AID  int64 `json:"aid" mapstructure:"aid" gorm:"column:aid"`
}

func ConvertChangelog(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	db := taskCtx.GetDal()
	changelogIdGen := didgen.NewDomainIdGenerator(&models.ZentaoChangelogDetail{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.ZentaoAccount{})
	executionIdGen := didgen.NewDomainIdGenerator(&models.ZentaoExecution{})
	storyIdGen := didgen.NewDomainIdGenerator(&models.ZentaoStory{})
	taskIdGen := didgen.NewDomainIdGenerator(&models.ZentaoTask{})
	bugIdGen := didgen.NewDomainIdGenerator(&models.ZentaoBug{})
	cn := models.ZentaoChangelog{}.TableName()
	cdn := models.ZentaoChangelogDetail{}.TableName()
	an := models.ZentaoAccount{}.TableName()
	cursor, err := db.Cursor(
		dal.Select(fmt.Sprintf("*,%s.id cid,%s.id cdid,%s.id aid ", cn, cdn, an)),
		dal.From(&models.ZentaoChangelog{}),
		dal.Join(fmt.Sprintf("LEFT JOIN %s on %s.changelog_id = %s.id and %s.connection_id = %d", cdn, cdn, cn, cdn, data.Options.ConnectionId)),
		dal.Join(fmt.Sprintf("LEFT JOIN %s on %s.realname = %s.actor and %s.connection_id = %d", an, an, cn, an, data.Options.ConnectionId)),
		dal.Where(fmt.Sprintf(`%s.project = ? and %s.connection_id = ?`, cn, cn),
			data.Options.ProjectId,
			data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	convertor, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(ZentaoChangelogSelect{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_CHANGELOG_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			cl := inputRow.(*ZentaoChangelogSelect)
			if cl.CDID == 0 {
				return nil, nil
			}
			var issueId string
			switch cl.ObjectType {
			case "story":
				issueId = storyIdGen.Generate(data.Options.ConnectionId, cl.ObjectId)
			case "task":
				issueId = taskIdGen.Generate(data.Options.ConnectionId, cl.ObjectId)
			case "bug":
				issueId = bugIdGen.Generate(data.Options.ConnectionId, cl.ObjectId)
			}
			domainCl := &ticket.IssueChangelogs{
				DomainEntity: domainlayer.DomainEntity{
					Id: changelogIdGen.Generate(data.Options.ConnectionId, cl.CID, cl.CDID),
				},
				IssueId:           issueId,
				AuthorName:        cl.Actor,
				FieldId:           cl.Field,
				FieldName:         cl.Field,
				OriginalFromValue: cl.Old,
				OriginalToValue:   cl.New,
				FromValue:         cl.Old,
				ToValue:           cl.New,
				CreatedDate:       cl.Date,
			}
			if cl.AID != 0 {
				domainCl.AuthorId = accountIdGen.Generate(data.Options.ConnectionId, cl.AID)
			}
			if domainCl.FieldName == "assignedTo" {
				domainCl.FieldName = "assignee"
				if cl.Old != "" {
					if id := data.AccountCache.getAccountID(cl.Old); id != 0 {
						domainCl.OriginalFromValue = accountIdGen.Generate(data.Options.ConnectionId, id)
						domainCl.FromValue = accountIdGen.Generate(data.Options.ConnectionId, id)
					}
				}
				if cl.New != "" {
					if id := data.AccountCache.getAccountID(cl.New); id != 0 {
						domainCl.OriginalToValue = accountIdGen.Generate(data.Options.ConnectionId, id)
						domainCl.ToValue = accountIdGen.Generate(data.Options.ConnectionId, id)
					}
				}
			}
			if domainCl.FieldName == "execution" {
				domainCl.FieldName = "Sprint"
				if cl.Old != "" {
					oldValue, _ := strconv.ParseInt(cl.Old, 10, 64)
					if oldValue != 0 {
						domainCl.OriginalFromValue = executionIdGen.Generate(data.Options.ConnectionId, oldValue)
						domainCl.FromValue = executionIdGen.Generate(data.Options.ConnectionId, oldValue)
					}
				}
				if cl.New != "" {
					newValue, _ := strconv.ParseInt(cl.New, 10, 64)
					if newValue != 0 {
						domainCl.OriginalToValue = executionIdGen.Generate(data.Options.ConnectionId, newValue)
						domainCl.ToValue = executionIdGen.Generate(data.Options.ConnectionId, newValue)
					}
				}
			}

			return []interface{}{
				domainCl,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return convertor.Execute()
}
