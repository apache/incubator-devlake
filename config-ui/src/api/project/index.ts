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

import type { IProject } from '@/types';
import { request } from '@/utils';

export const list = (data: Pagination & { keyword?: string }): Promise<{ count: number; projects: IProject[] }> =>
  request('/projects', { data });

export const get = (name: string): Promise<IProject> => request(`/projects/${encodeURIComponent(name)}`);

export const checkName = (name: string) => request(`/projects/${encodeURIComponent(name)}/check`);

export const create = (data: Pick<IProject, 'name' | 'description' | 'metrics'>) =>
  request('/projects', {
    method: 'post',
    data,
  });

export const remove = (name: string) =>
  request(`/projects/${encodeURIComponent(name)}`, {
    method: 'delete',
  });

export const update = (name: string, data: Pick<IProject, 'name' | 'description' | 'metrics'>) =>
  request(`/projects/${encodeURIComponent(name)}`, {
    method: 'patch',
    data: {
      ...data,
      name: encodeURIComponent(data.name),
    },
  });
