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
	"github.com/apache/incubator-devlake/mocks"
	"github.com/stretchr/testify/mock"
)

// DummyLogger FIXME ...
func DummyLogger() *mocks.Logger {
	logger := new(mocks.Logger)
	logger.On("IsLevelEnabled", mock.Anything).Return(false).Maybe()
	logger.On("Printf", mock.Anything, mock.Anything).Maybe()
	logger.On("Log", mock.Anything, mock.Anything, mock.Anything).Maybe()
	logger.On("Debug", mock.Anything, mock.Anything).Maybe()
	logger.On("Info", mock.Anything, mock.Anything).Maybe()
	logger.On("Warn", mock.Anything, mock.Anything).Maybe()
	logger.On("Error", mock.Anything, mock.Anything).Maybe()
	logger.On("Nested", mock.Anything, mock.Anything).Return(logger).Maybe()
	return logger
}
