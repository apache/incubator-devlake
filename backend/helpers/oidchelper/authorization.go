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

package oidchelper

import "strings"

func (c *Config) IsUserAllowed(email string) bool {
	if len(c.AllowEmails) == 0 &&
		len(c.AllowDomains) == 0 {
		return true
	}

	email = strings.ToLower(strings.TrimSpace(email))

	if _, ok := c.AllowEmails[email]; ok {
		return true
	}

	_, domain, ok := strings.Cut(email, "@")
	if ok {
		if _, ok := c.AllowDomains[domain]; ok {
			return true
		}
	}

	return false
}
