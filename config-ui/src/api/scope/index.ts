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

export const list = (
  plugin: string,
  connectionId: ID,
  data?: T.ListQuery,
): Promise<{ count: number; scopes: T.List }> =>
  request(`/plugins/${plugin}/connections/${connectionId}/scopes`, {
    data,
  });

export const get = (plugin: string, connectionId: ID, scopeId: ID) =>
  request(`/plugins/${plugin}/connections/${connectionId}/scopes/${scopeId}`);

export const remove = (plugin: string, connectionId: ID, scopeId: ID, onlyData: boolean) =>
  request(`/plugins/${plugin}/connections/${connectionId}/scopes/${scopeId}?delete_data_only=${onlyData}`, {
    method: 'delete',
  });

export const update = (plugin: string, connectionId: ID, scopeId: ID, payload: any) =>
  request(`/plugins/${plugin}/connections/${connectionId}/scopes/${scopeId}`, {
    method: 'patch',
    data: payload,
  });

export const batch = (plugin: string, connectionId: ID, payload: any) =>
  request(`/plugins/${plugin}/connections/${connectionId}/scopes`, {
    method: 'put',
    data: payload,
  });

export const remote = (
  plugin: string,
  connectionId: ID,
  data: T.RemoteQuery,
): Promise<{ children: T.RemoteScope[]; nextPageToken: string }> =>
  request(`/plugins/${plugin}/connections/${connectionId}/remote-scopes`, {
    method: 'get',
    data,
  });

export const searchRemote = (
  plugin: string,
  connectionId: ID,
  data: T.SearchRemoteQuery,
): Promise<{ children: T.RemoteScope[]; count: number }> =>
  request(`/plugins/${plugin}/connections/${connectionId}/search-remote-scopes`, {
    method: 'get',
    data,
  });
