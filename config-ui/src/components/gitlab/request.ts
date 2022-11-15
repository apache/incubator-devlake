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

import type { AxiosRequestConfig } from 'axios'
import axios from 'axios'

const instance = axios.create({
  baseURL: '/api'
})

export type ReuqestConfig = {
  method?: AxiosRequestConfig['method']
  data?: unknown
  timeout?: number
  signal?: AbortSignal
  headers?: Record<string, string>
}

const request = (path: string, config?: ReuqestConfig) => {
  const { method = 'GET', data, timeout, headers, signal } = config || {}

  const cancelTokenSource = axios.CancelToken.source()
  const params: any = {
    url: path,
    method,
    timeout,
    headers,
    cancelToken: cancelTokenSource?.token
  }

  if (method === 'GET') {
    params.params = data
  } else {
    params.data = data
  }

  const promise = instance.request(params).then((resp) => resp.data)

  if (signal) {
    signal.addEventListener('abort', () => {
      cancelTokenSource?.cancel()
    })
  }

  return promise
}

export default request
