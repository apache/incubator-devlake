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

package unithelper

import (
	mockcontext "github.com/apache/incubator-devlake/mocks/core/context"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	"github.com/stretchr/testify/mock"
)

// DummyBasicRes FIXME ...
func DummyBasicRes(callback func(mockDal *mockdal.Dal)) *mockcontext.BasicRes {
	mockRes := new(mockcontext.BasicRes)
	mockLog := DummyLogger()
	mockDal := new(mockdal.Dal)

	callback(mockDal)

	mockRes.On("GetDal").Return(mockDal)
	mockRes.On("GetLogger").Return(mockLog).Maybe()
	mockRes.On("GetConfig", mock.Anything).Return("").Maybe()
	mockDal.On("AllTables").Return(nil, nil)
	return mockRes
}
