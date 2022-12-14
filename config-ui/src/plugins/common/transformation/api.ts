/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import request from '@/components/utils/request'
import { Plugins } from '@/plugins'

type GetRulesParams = {
  page: number
  pageSize: number
}

export const getRules = (plugin: Plugins, params?: GetRulesParams) =>
  request(`/plugins/${plugin}/transformation_rules`, {
    method: 'get',
    data: params
  })

export const getDataScopeRepo = (
  plugin: Plugins,
  connectionId: ID,
  repoId: ID
) => request(`/plugins/${plugin}/connections/${connectionId}/scopes/${repoId}`)

export const updateDataScope = (
  plugin: string,
  connectionId: ID,
  payload: any
) =>
  request(`/plugins/${plugin}/connections/${connectionId}/scopes`, {
    method: 'put',
    data: payload
  })

export const createTransformation = (plugin: string, paylod: any) =>
  request(`/plugins/${plugin}/transformation_rules`, {
    method: 'POST',
    data: paylod
  })
