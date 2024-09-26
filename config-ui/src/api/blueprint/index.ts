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

import type { IBlueprint } from '@/types';
import { request } from '@/utils';

export const list = (data: Pagination & { type: string }): Promise<{ count: number; blueprints: IBlueprint[] }> =>
  request('/blueprints', { data });

export const get = (id: ID): Promise<IBlueprint> => request(`/blueprints/${id}`);

export const create = (data: any) =>
  request('/blueprints', {
    method: 'post',
    data,
  });

export const remove = (id: ID) => request(`/blueprints/${id}`, { method: 'delete' });

export const update = (id: ID, data: Partial<IBlueprint>) => request(`/blueprints/${id}`, { method: 'patch', data });

export const pipelines = (id: ID, data?: Pagination) => request(`/blueprints/${id}/pipelines`, { data });

type TriggerQuery = {
  skipCollectors?: boolean;
  fullSync?: boolean;
};

export const trigger = (id: ID, data: TriggerQuery = { skipCollectors: false, fullSync: false }) =>
  request(`/blueprints/${id}/trigger`, { method: 'post', data });
