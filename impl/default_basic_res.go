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

package impl

import (
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/spf13/viper"
)

// DefaultBasicRes offers a common implementation for the  BasisRes interface
type DefaultBasicRes struct {
	cfg    *viper.Viper
	logger core.Logger
	db     dal.Dal
}

// GetConfig returns the value of the specificed name
func (c *DefaultBasicRes) GetConfig(name string) string {
	return c.cfg.GetString(name)
}

// GetDal returns the Dal instance
func (c *DefaultBasicRes) GetDal() dal.Dal {
	return c.db
}

// GetLogger returns the Logger instance
func (c *DefaultBasicRes) GetLogger() core.Logger {
	return c.logger
}

// NestedLogger returns a new DefaultBasicRes with a new nested logger
func (c *DefaultBasicRes) NestedLogger(name string) core.BasicRes {
	return &DefaultBasicRes{
		cfg:    c.cfg,
		logger: c.logger.Nested(name),
		db:     c.db,
	}
}

// NewDefaultBasicRes creates a new DefaultBasicRes instance
func NewDefaultBasicRes(
	cfg *viper.Viper,
	logger core.Logger,
	db dal.Dal,
) *DefaultBasicRes {
	return &DefaultBasicRes{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}
