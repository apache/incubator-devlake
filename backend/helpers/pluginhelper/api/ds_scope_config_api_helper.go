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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/server/api/shared"
)

// DsAnyScopeConfigApiHelper
type DsAnyScopeConfigApiHelper struct {
	*AnyModelApiHelper
	*srvhelper.AnyScopeConfigSrvHelper
}

func NewDsAnyScopeConfigApiHelper(
	basicRes context.BasicRes,
	srvHelper *srvhelper.AnyScopeConfigSrvHelper,
) *DsAnyScopeConfigApiHelper {
	return &DsAnyScopeConfigApiHelper{
		AnyModelApiHelper:       NewAnyModelApiHelper(basicRes, srvHelper.AnyModelSrvHelper, []string{"connectionId", "scopeId"}, nil),
		AnyScopeConfigSrvHelper: srvHelper,
	}
}

func (scopeConfigApi *DsAnyScopeConfigApiHelper) GetAll(input *plugin.ApiResourceInput) (out *plugin.ApiResourceOutput, err errors.Error) {
	connectionId, err := extractConnectionId(input)
	if err != nil {
		return nil, err
	}
	scopeConfigs := errors.Must1(scopeConfigApi.GetAllByConnectionIdAny(connectionId))
	return &plugin.ApiResourceOutput{
		Body: scopeConfigs,
	}, nil
}

func (scopeConfigApi *DsAnyScopeConfigApiHelper) GetProjectsByScopeConfig(input *plugin.ApiResourceInput) (out *plugin.ApiResourceOutput, err errors.Error) {
	scopeConfig, err := scopeConfigApi.FindByPkAny(input)
	if err != nil {
		return nil, err
	}
	projectDetails := errors.Must1(scopeConfigApi.AnyScopeConfigSrvHelper.GetProjectsByScopeConfig(scopeConfig))
	return &plugin.ApiResourceOutput{
		Body: projectDetails,
	}, nil
}

func (scopeConfigApi *DsAnyScopeConfigApiHelper) Post(input *plugin.ApiResourceInput) (out *plugin.ApiResourceOutput, err errors.Error) {
	// fix connectionId
	connectionId, err := extractConnectionId(input)
	if err != nil {
		return nil, err
	}
	input.Body["connectionId"] = connectionId
	return scopeConfigApi.AnyModelApiHelper.Post(input)
}

func (scopeConfigApi *DsAnyScopeConfigApiHelper) Patch(input *plugin.ApiResourceInput) (out *plugin.ApiResourceOutput, err errors.Error) {
	// fix connectionId
	connectionId, err := extractConnectionId(input)
	if err != nil {
		return nil, err
	}
	input.Body["connectionId"] = connectionId
	return scopeConfigApi.AnyModelApiHelper.Patch(input)
}

func (scopeConfigApi *DsAnyScopeConfigApiHelper) Delete(input *plugin.ApiResourceInput) (out *plugin.ApiResourceOutput, err errors.Error) {
	scopeConfig, err := scopeConfigApi.FindByPkAny(input)
	if err != nil {
		return nil, err
	}
	refs, err := scopeConfigApi.DeleteScopeConfigAny(scopeConfig)
	if err != nil {
		return &plugin.ApiResourceOutput{Body: &shared.ApiBody{
			Success: false,
			Message: err.Error(),
			Data:    refs,
		}, Status: err.GetType().GetHttpCode()}, err
	}
	return &plugin.ApiResourceOutput{
		Body: scopeConfig,
	}, nil
}

type DsScopeConfigApiHelper[SC dal.Tabler] struct {
	*DsAnyScopeConfigApiHelper
	*ModelApiHelper[SC]
}

func NewDsScopeConfigApiHelper[SC dal.Tabler](
	anyScopeConfigApiHelper *DsAnyScopeConfigApiHelper,
) *DsScopeConfigApiHelper[SC] {
	return &DsScopeConfigApiHelper[SC]{
		DsAnyScopeConfigApiHelper: anyScopeConfigApiHelper,
		ModelApiHelper:            NewModelApiHelper[SC](anyScopeConfigApiHelper.AnyModelApiHelper),
	}
}
