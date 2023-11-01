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

import type { IApiKey } from '@/types';
import { request } from '@/utils';

type ListRes = {
  count: number;
  apikeys: IApiKey[];
};

export const list = (data?: Pagination): Promise<ListRes> =>
  request('/api-keys', {
    data,
  });

type CreateForm = Pick<IApiKey, 'name' | 'expiredAt' | 'allowedPath'>;

export const create = (data: CreateForm): Promise<IApiKey> =>
  request('/api-keys', {
    method: 'POST',
    data: {
      ...data,
      type: 'devlake',
    },
  });

type RemoveRes = {
  message: string;
  success: boolean;
};

export const remove = (id: string): Promise<RemoveRes> =>
  request(`/api-keys/${id}`, {
    method: 'DELETE',
  });

export const renew = (id: ID) => request(`/api-keys/${id}`, { method: 'put' });
