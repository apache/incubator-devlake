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

package services

import (
	"testing"

	"github.com/apache/incubator-devlake/core/errors"
	dalmocks "github.com/apache/incubator-devlake/mocks/core/dal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteBlueprintInTransaction(t *testing.T) {
	t.Run("uses the caller transaction for every delete", func(t *testing.T) {
		tx := dalmocks.NewTransaction(t)
		tx.On("Delete", mock.Anything, mock.Anything).Return(nil).Times(4)

		manager := &BlueprintManager{}
		err := manager.DeleteBlueprintInTransaction(tx, 42)

		assert.NoError(t, err)
	})

	t.Run("returns dependent deletion failures", func(t *testing.T) {
		tx := dalmocks.NewTransaction(t)
		expected := errors.Default.New("unable to delete blueprint")
		tx.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		tx.On("Delete", mock.Anything, mock.Anything).Return(expected).Once()

		manager := &BlueprintManager{}
		err := manager.DeleteBlueprintInTransaction(tx, 42)

		assert.ErrorIs(t, err, expected)
	})
}
