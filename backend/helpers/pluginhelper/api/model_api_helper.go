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
	"fmt"
	"net/http"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/helpers/utils"
)

type ModelApiHelper[M dal.Tabler] struct {
	basicRes       context.BasicRes
	dalHelper      *srvhelper.ModelSrvHelper[M]
	log            log.Logger
	modelName      string
	pkPathVarNames []string
}

func NewModelApiHelper[M dal.Tabler](
	basicRes context.BasicRes,
	dalHelper *srvhelper.ModelSrvHelper[M],
	pkPathVarNames []string, // path variable names of primary key
) *ModelApiHelper[M] {
	m := new(M)
	modelName := fmt.Sprintf("%T", m)
	return &ModelApiHelper[M]{
		basicRes:       basicRes,
		dalHelper:      dalHelper,
		log:            basicRes.GetLogger().Nested(fmt.Sprintf("%s_dal", modelName)),
		modelName:      modelName,
		pkPathVarNames: pkPathVarNames,
	}
}

func (self *ModelApiHelper[M]) Post(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	model := new(M)
	err := utils.DecodeMapStruct(input.Body, model, false)
	if err != nil {
		return nil, err
	}
	err = self.dalHelper.Create(model)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Status: http.StatusCreated,
		Body:   model,
	}, nil
}

func (self *ModelApiHelper[M]) ExtractPkValues(input *plugin.ApiResourceInput) ([]interface{}, errors.Error) {
	pkv := make([]interface{}, len(self.pkPathVarNames))
	for i, pkn := range self.pkPathVarNames {
		var ok bool
		pkv[i], ok = input.Params[pkn]
		if !ok {
			return nil, errors.BadInput.New(fmt.Sprintf("missing path variable %s", pkn))
		}
	}
	return pkv, nil
}

func (self *ModelApiHelper[M]) FindByPk(input *plugin.ApiResourceInput) (*M, errors.Error) {
	pkv, err := self.ExtractPkValues(input)
	if err != nil {
		return nil, err
	}
	return self.dalHelper.FindByPk(pkv...)
}

func (self *ModelApiHelper[M]) GetDetail(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	model, err := self.FindByPk(input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Body: model,
	}, nil
}

func (self *ModelApiHelper[M]) Patch(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	model, err := self.FindByPk(input)
	if err != nil {
		return nil, err
	}
	err = utils.DecodeMapStruct(input.Body, model, true)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, fmt.Sprintf("faled to patch %s", self.modelName))
	}
	err = self.dalHelper.Update(model)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Body: model,
	}, nil
}

func (self *ModelApiHelper[M]) Delete(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	model, err := self.FindByPk(input)
	if err != nil {
		return nil, err
	}
	err = self.dalHelper.DeleteModel(model)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Body: model,
	}, nil
}

func (self *ModelApiHelper[M]) GetAll(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	all, err := self.dalHelper.GetAll()
	return &plugin.ApiResourceOutput{
		Body: all,
	}, err
}

func (self *ModelApiHelper[M]) PutMultiple(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var req struct {
		Data []M `json:"data"`
	}
	err := utils.DecodeMapStruct(input.Body, &req, false)
	if err != nil {
		return nil, err
	}
	for i, item := range req.Data {
		err := self.dalHelper.CreateOrUpdate(&item)
		if err != nil {
			return nil, errors.BadInput.Wrap(err, fmt.Sprintf("failed to save item %d", i))
		}
	}
	return &plugin.ApiResourceOutput{
		Body: req.Data,
	}, nil
}

func ParsePagination[P any](input *plugin.ApiResourceInput, query ...dal.Clause) (*P, errors.Error) {
	pagination := new(P)
	err := utils.DecodeMapStruct(input.Query, pagination, false)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "faild to decode pagination from query string")
	}
	err = utils.DecodeMapStruct(input.Params, pagination, false)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "faild to decode pagination from path variables")
	}
	return pagination, nil
}
