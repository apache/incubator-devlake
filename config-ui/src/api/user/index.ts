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

import type { IUser, IUserTeams } from '@/types';
import { request } from '@/utils';

export const list = (data: Pagination & { email?: string }): Promise<{ count: number; users: IUser[] }> =>
  request('/plugins/org/users', { data });

export const get = (userId: string): Promise<IUser> => request(`/plugins/org/users/${userId}`);

export const create = (data: { users: Omit<IUser, 'id'>[] }): Promise<IUser[]> =>
  request('/plugins/org/users', {
    method: 'post',
    data,
  });

export const update = (userId: string, data: Omit<IUser, 'id'>): Promise<IUser> =>
  request(`/plugins/org/users/${userId}`, {
    method: 'put',
    data,
  });

export const remove = (userId: string) =>
  request(`/plugins/org/users/${userId}`, {
    method: 'delete',
  });

export const listTeams = (userId: string): Promise<IUserTeams> => request(`/plugins/org/users/${userId}/teams`);

export const updateTeams = (userId: string, data: { teamIds: string[] }): Promise<IUserTeams> =>
  request(`/plugins/org/users/${userId}/teams`, {
    method: 'put',
    data,
  });
