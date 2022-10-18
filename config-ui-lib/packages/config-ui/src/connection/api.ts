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
import request from '../utils/request';

import { ConnectionType } from './typed';

export const getList = (type: ConnectionType) =>
  request(`/plugins/${type}/connections`);

export const test = (type: ConnectionType, payload: any) =>
  request(`/plugins/${type}/test`, { method: 'POST', data: payload });

export const create = (type: ConnectionType, payload: any) =>
  request(`/plugins/${type}/connections`, { method: 'POST', data: payload });

export const remove = (type: ConnectionType, id: number) =>
  request(`/plugins/${type}/connections/${id}`, { method: 'DELETE' });

export const update = (type: ConnectionType, id: number, payload: any) =>
  request(`/plugins/${type}/connections/${id}`, {
    method: 'PATCH',
    data: payload,
  });
