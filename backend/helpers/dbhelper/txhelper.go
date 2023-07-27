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

package dbhelper

import (
	"fmt"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
)

// TxHelper is a helper for transaction management
type TxHelper[E error] struct {
	basicRes context.BasicRes
	perr     *E
	tx       dal.Transaction
}

// Begin starts a transaction
func (l *TxHelper[E]) Begin() dal.Transaction {
	if l.tx != nil {
		panic(fmt.Errorf("Begin has been called"))
	}
	l.tx = l.basicRes.GetDal().Begin()
	return l.tx
}

// LockTablesTimeout locks tables with timeout
func (l *TxHelper[E]) LockTablesTimeout(timeout time.Duration, lockTables dal.LockTables) errors.Error {
	println("timeout", timeout)
	c := make(chan errors.Error, 1)
	go func() {
		c <- l.tx.LockTables(lockTables)
		close(c)
	}()

	select {
	case err := <-c:
		if err != nil {
			panic(err)
		}
	case <-time.After(timeout):
		return errors.Timeout.New("lock tables timeout: " + fmt.Sprintf("%v", lockTables))
	}
	return nil
}

// End ends a transaction and commits it if no error was set or it will try to rollback and release locked tables
func (l *TxHelper[E]) End() {
	if l.tx == nil {
		panic("Begin was never called")
	}
	var msg string
	err := *l.perr
	if interface{}(err) == nil {
		msg = ""
		r := recover()
		if r != nil {
			msg = fmt.Sprintf("%v", r)
		}
	} else {
		msg = err.Error()
	}

	if msg == "" {
		_ = l.tx.UnlockTables()
		errors.Must(l.tx.Commit())
	} else {
		_ = l.tx.UnlockTables()
		_ = l.tx.Rollback()
	}
	l.tx = nil
}

// NewTxHelper creates a new TxHelper, the errorPointer is used to detect if any error was set
func NewTxHelper[E error](basicRes context.BasicRes, errorPointer *E) *TxHelper[E] {
	if errorPointer == nil {
		panic(fmt.Errorf("errorPointer is required"))
	}
	return &TxHelper[E]{basicRes: basicRes, perr: errorPointer}
}
