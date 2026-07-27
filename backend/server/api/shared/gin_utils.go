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

package shared

import (
	"context"
	"net/http"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/gin-gonic/gin"
)

// restAuthKey is an unexported type used as a request-context key so it cannot
// collide with keys set by other packages.
type restAuthKey struct{}

// SetRestAuthUser stores the authenticated user in the HTTP request context.
// This is necessary because gin's HandleContext calls c.reset(), which clears
// c.Keys but leaves c.Request (and its context) intact. RestAuthentication
// calls this before rerouting so the user survives the reset.
func SetRestAuthUser(r *http.Request, user *common.User) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), restAuthKey{}, user))
}

func GetUser(c *gin.Context) (*common.User, bool) {
	userObj, exist := c.Get(common.USER)
	if exist {
		if user, ok := userObj.(*common.User); ok {
			return user, true
		}
	}
	// Fallback: RestAuthentication stores the user here before calling
	// HandleContext, which resets c.Keys but preserves c.Request.
	if user, ok := c.Request.Context().Value(restAuthKey{}).(*common.User); ok && user != nil {
		return user, true
	}
	return nil, false
}
