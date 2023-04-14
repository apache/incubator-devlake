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

import type { AxiosRequestConfig } from 'axios';
import axios from 'axios';
import { history } from '@/utils/history';

import { DEVLAKE_ENDPOINT } from '@/config';
import { toast } from '@/components/toast';

const instance = axios.create({
  baseURL: DEVLAKE_ENDPOINT,
});

export type ReuqestConfig = {
  method?: AxiosRequestConfig['method'];
  data?: unknown;
  timeout?: number;
  signal?: AbortSignal;
  headers?: Record<string, string>;
};

export const request = (path: string, config?: ReuqestConfig) => {
  const { method = 'get', data, timeout, headers, signal } = config || {};
  const cancelTokenSource = axios.CancelToken.source();
  const params: any = {
    url: path,
    method,
    timeout,
    headers: { ...headers, Authorization: `Bearer ${localStorage.getItem('accessToken')}` },
    cancelToken: cancelTokenSource?.token,
  };

  if (['GET', 'get'].includes(method)) {
    params.params = data;
  } else {
    params.data = data;
  }

  instance.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response && error.response.status === 401) {
        toast.error('Please log in first');
        history.push('/login');
      }
    },
  );

  const promise = instance.request(params).then((resp) => resp.data);

  if (signal) {
    signal.addEventListener('abort', () => {
      cancelTokenSource?.cancel();
    });
  }

  return promise;
};
