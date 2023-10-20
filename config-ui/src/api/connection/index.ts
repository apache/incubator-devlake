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

import { request } from '@/utils';

import * as T from './types';

export const list = (plugin: string): Promise<T.Connection[]> => request(`/plugins/${plugin}/connections`);

export const get = (plugin: string, connectionId: ID): Promise<T.Connection> =>
  request(`/plugins/${plugin}/connections/${connectionId}`);

export const create = (plugin: string, payload: T.ConnectionForm): Promise<T.Connection> =>
  request(`/plugins/${plugin}/connections`, { method: 'post', data: payload });

export const remove = (plugin: string, id: ID): Promise<T.Connection> =>
  request(`/plugins/${plugin}/connections/${id}`, { method: 'delete' });

export const update = (plugin: string, id: ID, payload: T.ConnectionForm): Promise<T.Connection> =>
  request(`/plugins/${plugin}/connections/${id}`, {
    method: 'patch',
    data: payload,
  });

export const test = (
  plugin: string,
  payload: Pick<
    T.ConnectionForm,
    'endpoint' | 'authMethod' | 'username' | 'password' | 'token' | 'appId' | 'secretKey' | 'proxy' | 'dbUrl'
  >,
): Promise<T.ConnectionTest> => request(`/plugins/${plugin}/test`, { method: 'post', data: payload });
