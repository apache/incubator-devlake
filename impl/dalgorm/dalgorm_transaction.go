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

package dalgorm

import (
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core/dal"
)

// DalgormTransaction represents a gorm transaction which using the same underlying
// session for all queries
type DalgormTransaction struct {
	*Dalgorm
}

var _ dal.Transaction = (*DalgormTransaction)(nil)

// Rollback the transaction
func (t *DalgormTransaction) Rollback() errors.Error {
	r := t.db.Rollback()
	if r.Error != nil {
		return errors.Default.Wrap(r.Error, "failed to rollback transaction")
	}
	return nil
}

// Commit the transaction
func (t *DalgormTransaction) Commit() errors.Error {
	r := t.db.Commit()
	if r.Error != nil {
		return errors.Default.Wrap(r.Error, "failed to commit transaction")
	}
	return nil
}

func newTransaction(dalgorm *Dalgorm) *DalgormTransaction {
	return &DalgormTransaction{
		Dalgorm: NewDalgorm(dalgorm.db.Begin()),
	}
}
