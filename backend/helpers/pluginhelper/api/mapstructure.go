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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/helpers/utils"
	"github.com/go-playground/validator/v10"

	"github.com/mitchellh/mapstructure"
)

func DecodeMapStruct(input interface{}, result interface{}, zeroFields bool) errors.Error {
	return utils.DecodeMapStruct(input, result, zeroFields)
}

// Decode decodes `source` into `target`. Pass an optional validator to validate the target.
func Decode(source interface{}, target interface{}, vld *validator.Validate) errors.Error {
	target = models.UnwrapObject(target)
	if err := mapstructure.Decode(source, &target); err != nil {
		return errors.Default.Wrap(err, "error decoding map into target type")
	}
	if vld != nil {
		if err := vld.Struct(target); err != nil {
			return errors.Default.Wrap(err, "error validating target")
		}
	}
	return nil
}
