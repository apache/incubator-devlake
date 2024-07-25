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

package service

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/pluginhelper"
	"github.com/apache/incubator-devlake/plugins/customize/models"
)

// Service wraps database operations
type Service struct {
	dal         dal.Dal
	nameChecker *regexp.Regexp
}

func NewService(dal dal.Dal) *Service {
	return &Service{dal: dal, nameChecker: regexp.MustCompile(`^x_[a-zA-Z0-9_]{0,50}$`)}
}

// GetFields returns all the fields of the table
func (s *Service) GetFields(table string) ([]models.CustomizedField, errors.Error) {
	// the customized fields created before v0.16.0 were not recorded in the table `_tool_customized_field`, we should take care of them
	columns, err := s.dal.GetColumns(&models.Table{Name: table}, func(columnMeta dal.ColumnMeta) bool {
		return true
	})
	if err != nil {
		return nil, errors.Default.Wrap(err, "GetColumns error")
	}
	ff, err := s.getCustomizedFields(table)
	if err != nil {
		return nil, err
	}
	fieldMap := make(map[string]models.CustomizedField)
	for _, f := range ff {
		fieldMap[f.ColumnName] = f
	}
	var result []models.CustomizedField
	for _, col := range columns {
		// original fields
		if !strings.HasPrefix(col.Name(), "x_") {
			dataType, _ := col.ColumnType()
			result = append(result, models.CustomizedField{
				TbName:     table,
				ColumnName: col.Name(),
				DataType:   dal.ColumnType(dataType),
			})
			// customized fields
		} else {
			if field, ok := fieldMap[col.Name()]; ok {
				result = append(result, field)
			} else {
				result = append(result, models.CustomizedField{
					ColumnName: col.Name(),
					DataType:   dal.Varchar,
				})
			}
		}
	}
	return result, nil
}

// checkField checks if the field exist in table
func (s *Service) checkField(table, field string) (bool, errors.Error) {
	if table == "" {
		return false, errors.Default.New("empty table name")
	}
	if !strings.HasPrefix(field, "x_") {
		return false, errors.Default.New("column name should start with `x_`")
	}
	if !s.nameChecker.MatchString(field) {
		return false, errors.Default.New("invalid column name")
	}
	fields, err := s.GetFields(table)
	if err != nil {
		return false, err
	}
	for _, fld := range fields {
		if fld.ColumnName == field {
			return true, nil
		}
	}
	return false, nil
}

// CreateField creates a new column for the table cf.TbName and creates a new record in the table `_tool_customized_fields`
func (s *Service) CreateField(cf *models.CustomizedField) errors.Error {
	exists, err := s.checkField(cf.TbName, cf.ColumnName)
	if err != nil {
		return err
	}
	if exists {
		return errors.BadInput.New(fmt.Sprintf("the column %s already exists", cf.ColumnName))
	}
	err = s.dal.Create(cf)
	if err != nil {
		return errors.Default.Wrap(err, "create customizedField")
	}
	err = s.dal.AddColumn(cf.TbName, cf.ColumnName, cf.DataType)
	if err != nil {
		return errors.Default.Wrap(err, "AddColumn error")
	}
	return nil
}

// DeleteField deletes the `field` form the `table`
func (s *Service) DeleteField(table, field string) errors.Error {
	exists, err := s.checkField(table, field)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	err = s.dal.DropColumns(table, field)
	if err != nil {
		return errors.Default.Wrap(err, "DropColumn error")
	}
	return s.dal.Delete(&models.CustomizedField{}, dal.Where("tb_name = ? AND column_name = ?", table, field))
}

// getCustomizedFields returns all the customized fields definitions of the table
func (s *Service) getCustomizedFields(table string) ([]models.CustomizedField, errors.Error) {
	var result []models.CustomizedField
	err := s.dal.All(&result, dal.Where("tb_name = ?", table))
	return result, err
}

// ImportIssue import csv file to the table `issues`, and create relations to boards
// issue could exist in multiple boards, so we should only delete an old records when it doesn't belong to another board
func (s *Service) ImportIssue(boardId string, file io.ReadCloser) errors.Error {
	err := s.dal.Delete(
		&ticket.Issue{},
		dal.Where("id IN (SELECT issue_id FROM board_issues WHERE board_id=? AND issue_id NOT IN (SELECT issue_id FROM board_issues WHERE board_id!=?))", boardId, boardId),
	)
	if err != nil {
		return err
	}

	err = s.dal.Delete(
		&ticket.IssueLabel{},
		dal.Where("issue_id IN (SELECT issue_id FROM board_issues WHERE board_id=? AND issue_id NOT IN (SELECT issue_id FROM board_issues WHERE board_id!=?))", boardId, boardId),
	)
	if err != nil {
		return err
	}

	err = s.dal.Delete(
		&ticket.BoardIssue{},
		dal.Where("board_id = ?", boardId),
	)
	if err != nil {
		return err
	}

	return s.importCSV(file, boardId, s.issueHandlerFactory(boardId))
}

// SaveBoard make sure the board exists in table `boards`
func (s *Service) SaveBoard(boardId, boardName string) errors.Error {
	return s.dal.CreateOrUpdate(&ticket.Board{
		DomainEntity: domainlayer.DomainEntity{
			Id: boardId,
		},
		Name: boardName,
		Type: "csv",
	})
}

// ImportIssueCommit imports csv file into the table `issue_commits`
func (s *Service) ImportIssueCommit(boardId string, file io.ReadCloser) errors.Error {
	err := s.dal.Delete(
		&crossdomain.IssueCommit{},
		dal.Where("issue_id IN (SELECT issue_id FROM board_issues WHERE board_id=? AND issue_id NOT IN (SELECT issue_id FROM board_issues WHERE board_id!=?))", boardId, boardId),
	)
	if err != nil {
		return err
	}
	return s.importCSV(file, boardId, s.issueCommitHandler)
}

// ImportIssueRepoCommit imports data to the table `issue_repo_commits` and `issue_commits`
func (s *Service) ImportIssueRepoCommit(boardId string, file io.ReadCloser) errors.Error {
	// delete old records of the table `issue_repo_commit` and `issue_commit`
	err := s.dal.Delete(
		&crossdomain.IssueRepoCommit{},
		dal.Where("issue_id IN (SELECT issue_id FROM board_issues WHERE board_id=? AND issue_id NOT IN (SELECT issue_id FROM board_issues WHERE board_id!=?))", boardId, boardId),
	)
	if err != nil {
		return err
	}
	err = s.dal.Delete(
		&crossdomain.IssueCommit{},
		dal.Where("issue_id IN (SELECT issue_id FROM board_issues WHERE board_id=? AND issue_id NOT IN (SELECT issue_id FROM board_issues WHERE board_id!=?))", boardId, boardId),
	)
	if err != nil {
		return err
	}
	return s.importCSV(file, boardId, s.issueRepoCommitHandler)
}

// importCSV extract records from csv file, and save them to DB using recordHandler
// the rawDataParams is used to identify the data source,
// the recordHandler is used to handle the record, it should return an error if the record is invalid
// the `created_at` and `updated_at` will be set to the current time
func (s *Service) importCSV(file io.ReadCloser, rawDataParams string, recordHandler func(map[string]interface{}) errors.Error) errors.Error {
	iterator, err := pluginhelper.NewCsvFileIteratorFromFile(file)
	if err != nil {
		return err
	}
	var hasNext bool
	var line int
	now := time.Now()
	for {
		line++
		if hasNext, err = iterator.HasNextWithError(); !hasNext {
			return errors.BadInput.Wrap(err, fmt.Sprintf("error on processing the line:%d", line))
		} else {
			record := iterator.Fetch()
			record["_raw_data_params"] = rawDataParams
			for k, v := range record {
				if v.(string) == "NULL" {
					record[k] = nil
				}
			}
			record["created_at"] = now
			record["updated_at"] = now
			err = recordHandler(record)
			if err != nil {
				return errors.BadInput.Wrap(err, fmt.Sprintf("error on processing the line:%d", line))
			}
		}
	}
}

// issueHandlerFactory returns a handler that save record into `issues`, `board_issues` and `issue_labels` table
func (s *Service) issueHandlerFactory(boardId string) func(record map[string]interface{}) errors.Error {
	return func(record map[string]interface{}) errors.Error {
		var err errors.Error
		var id string
		if record["id"] == nil {
			return errors.Default.New("record without id")
		}
		id, _ = record["id"].(string)
		if id == "" {
			return errors.Default.New("empty id")
		}
		if record["labels"] != nil {
			labels, ok := record["labels"].(string)
			if !ok {
				return errors.Default.New("labels is not string")
			}
			var issueLabels []*ticket.IssueLabel
			appearedLabels := make(map[string]struct{}) // record the labels that have appeared
			for _, label := range strings.Split(labels, ",") {
				label = strings.TrimSpace(label)
				if label == "" {
					continue
				}
				if _, appeared := appearedLabels[label]; !appeared {
					issueLabel := &ticket.IssueLabel{
						IssueId:   id,
						LabelName: label,
						NoPKModel: common.NoPKModel{
							RawDataOrigin: common.RawDataOrigin{
								RawDataParams: boardId,
							},
						},
					}
					issueLabels = append(issueLabels, issueLabel)
					appearedLabels[label] = struct{}{}
				}
			}
			if len(issueLabels) > 0 {
				err = s.dal.CreateOrUpdate(issueLabels)
				if err != nil {
					return err
				}
			}
		}
		delete(record, "labels")
		err = s.dal.CreateWithMap(&ticket.Issue{}, record)
		if err != nil {
			return err
		}
		return s.dal.CreateOrUpdate(&ticket.BoardIssue{
			BoardId: boardId,
			IssueId: id,
		})
	}
}

// issueCommitHandler save record into `issue_commits` table
func (s *Service) issueCommitHandler(record map[string]interface{}) errors.Error {
	return s.dal.CreateWithMap(&crossdomain.IssueCommit{}, record)
}

// issueRepoCommitHandlerFactory returns a handler that will populate the `issue_commits` and `issue_repo_commits` table
// ths issueCommitsFields is used to filter the fields that should be inserted into the `issue_commits` table
func (s *Service) issueRepoCommitHandler(record map[string]interface{}) errors.Error {
	err := s.dal.CreateWithMap(&crossdomain.IssueRepoCommit{}, record)
	if err != nil {
		return err
	}
	// remove fields that not in table `issue_commits`
	delete(record, "host")
	delete(record, "namespace")
	delete(record, "repo_name")
	delete(record, "repo_url")
	return s.dal.CreateWithMap(&crossdomain.IssueCommit{}, record)
}
