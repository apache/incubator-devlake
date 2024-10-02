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

package api

import (
	"context"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

func TestExistingConnectionForTokenCheck(input *plugin.ApiResourceInput) errors.Error {
	connection, err := dsHelper.ConnApi.GetMergedConnection(input)
	if err != nil {
		return err
	}
	testConnectionResult, testConnectionErr := testExistingConnection(context.TODO(), connection.GithubConn)
	if testConnectionErr != nil {
		return testConnectionErr
	}
	for _, token := range testConnectionResult.Tokens {
		if !token.Success {
			return errors.Default.New(fmt.Sprintf("token %s failed with msg: %s", token.Token, token.Message))
		}
	}
	return nil
}
