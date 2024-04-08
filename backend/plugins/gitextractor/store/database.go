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

package store

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"reflect"
)

const BathSize = 100

type Database struct {
	driver *helper.BatchSaveDivider
	table  string
	params string
}

func NewDatabase(basicRes context.BasicRes, repoId string) *Database {
	database := &Database{
		table:  "gitextractor",
		params: repoId,
	}
	database.driver = helper.NewBatchSaveDivider(
		basicRes,
		BathSize,
		database.table,
		database.params,
	)
	return database
}

func (d *Database) updateRawDataFields(rawData *common.RawDataOrigin) {
	rawData.RawDataTable = d.table
	rawData.RawDataParams = d.params
}

func (d *Database) RepoCommits(repoCommit *code.RepoCommit) errors.Error {
	batch, err := d.driver.ForType(reflect.TypeOf(repoCommit))
	if err != nil {
		return err
	}
	d.updateRawDataFields(&repoCommit.RawDataOrigin)
	return batch.Add(repoCommit)
}

func (d *Database) Commits(commit *code.Commit) errors.Error {
	account := &crossdomain.Account{
		DomainEntity: domainlayer.DomainEntity{Id: commit.AuthorEmail},
		Email:        commit.AuthorEmail,
		FullName:     commit.AuthorName,
		UserName:     commit.AuthorName,
	}
	accountBatch, err := d.driver.ForType(reflect.TypeOf(account))
	if err != nil {
		return err
	}
	d.updateRawDataFields(&account.RawDataOrigin)
	// Skip accounts without email, such accounts fail in PostgreSQL
	if account.Email != "" {
		err = accountBatch.Add(account)
		if err != nil {
			return err
		}
	}
	commitBatch, err := d.driver.ForType(reflect.TypeOf(commit))
	if err != nil {
		return err
	}
	d.updateRawDataFields(&account.RawDataOrigin)
	return commitBatch.Add(commit)
}

func (d *Database) Refs(ref *code.Ref) errors.Error {
	batch, err := d.driver.ForType(reflect.TypeOf(ref))
	if err != nil {
		return err
	}
	d.updateRawDataFields(&ref.RawDataOrigin)
	return batch.Add(ref)
}

func (d *Database) CommitFiles(file *code.CommitFile) errors.Error {
	batch, err := d.driver.ForType(reflect.TypeOf(file))
	if err != nil {
		return err
	}
	d.updateRawDataFields(&file.RawDataOrigin)
	return batch.Add(file)
}

func (d *Database) CommitFileComponents(commitFileComponent *code.CommitFileComponent) errors.Error {
	batch, err := d.driver.ForType(reflect.TypeOf(commitFileComponent))
	if err != nil {
		return err
	}
	d.updateRawDataFields(&commitFileComponent.RawDataOrigin)
	return batch.Add(commitFileComponent)
}

func (d *Database) RepoSnapshot(snapshotElement *code.RepoSnapshot) errors.Error {
	batch, err := d.driver.ForType(reflect.TypeOf(snapshotElement))
	if err != nil {
		return err
	}
	d.updateRawDataFields(&snapshotElement.RawDataOrigin)
	return batch.Add(snapshotElement)
}

func (d *Database) CommitLineChange(commitLineChange *code.CommitLineChange) errors.Error {
	batch, err := d.driver.ForType(reflect.TypeOf(commitLineChange))
	if err != nil {
		return err
	}
	d.updateRawDataFields(&commitLineChange.RawDataOrigin)
	return batch.Add(commitLineChange)
}

func (d *Database) CommitParents(pp []*code.CommitParent) errors.Error {
	if len(pp) == 0 {
		return nil
	}
	batch, err := d.driver.ForType(reflect.TypeOf(pp[0]))
	if err != nil {
		return err
	}
	for _, cp := range pp {
		d.updateRawDataFields(&cp.RawDataOrigin)
		err = batch.Add(cp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) Close() errors.Error {
	return d.driver.Close()
}
