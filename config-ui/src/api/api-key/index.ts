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

export const list = (data?: Pagination): Promise<{ count: number; apikeys: T.Key[] }> =>
  request('/api-keys', {
    data,
  });

export const create = (data: Pick<T.Key, 'name' | 'expiredAt' | 'allowedPath'>): Promise<T.Key> =>
  request('/api-keys', {
    method: 'POST',
    data: {
      ...data,
      type: 'devlake',
    },
  });

export const remove = (id: string): Promise<void> =>
  request(`/api-keys/${id}`, {
    method: 'DELETE',
  });

export const renew = (id: ID) => request(`/api-keys/${id}`, { method: 'put' });
