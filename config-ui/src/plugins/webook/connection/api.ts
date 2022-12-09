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

export const getConnections = () => request('/plugins/webhook/connections')

export const getConnection = (id: ID) =>
  request(`/plugins/webhook/connections/${id}`)

type Paylod = {
  name: string
}

export const createConnection = (payload: Paylod) =>
  request('/plugins/webhook/connections', {
    method: 'post',
    data: payload
  })

export const updateConnection = (id: ID, payload: Paylod) =>
  request(`/plugins/webhook/connections/${id}`, {
    method: 'patch',
    data: payload
  })

export const deleteConnection = (id: ID) =>
  request(`/plugins/webhook/connections/${id}`, {
    method: 'delete'
  })
