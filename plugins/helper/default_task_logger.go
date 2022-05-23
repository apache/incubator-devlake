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

package helper

import (
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/sirupsen/logrus"
)

// bridge to current implementation at this point
// TODO: implement another TaskLogger for distributed runner/worker

func NewDefaultTaskLogger(log *logrus.Logger, prefix string, loggerPool map[string]*logrus.Logger) core.Logger {
	return NewDefaultLogger(log, prefix, loggerPool)
}
