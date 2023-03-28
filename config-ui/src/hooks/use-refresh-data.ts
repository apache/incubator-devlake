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

import { useState, useRef } from 'react';
import { isEqualWith } from 'lodash';
import axios from 'axios';

export const useRefreshData = <T>(request: (signal: AbortSignal) => Promise<T>, deps: React.DependencyList = []) => {
  const [, setVersion] = useState(0);
  const ref = useRef<{
    state: 'ready' | 'pending' | 'error';
    deps?: React.DependencyList;
    data?: T;
    abortController?: AbortController;
    timer?: number;
  }>({
    state: 'pending',
  });

  if (isEqualWith(ref.current.deps, deps)) {
    return {
      data: ref.current.data,
      ready: ref.current.state === 'ready',
      pending: ref.current.state === 'pending',
      error: ref.current.state === 'error',
    };
  }

  // When the last state transition has not waited until the new request is completed
  // Reset status to pending
  ref.current.state = 'pending';
  ref.current.deps = deps;
  ref.current.data = undefined;
  clearTimeout(ref.current.timer);
  ref.current.abortController?.abort();
  ref.current.timer = window.setTimeout(() => {
    ref.current.abortController = new AbortController();
    request(ref.current.abortController.signal)
      .then((data: T) => {
        ref.current.data = data;
        ref.current.state = 'ready';
        setVersion((v) => v + 1);
      })
      .catch((err: unknown) => {
        if (axios.isCancel(err)) {
          return;
        }
        ref.current.state = 'error';
        console.error(err);
        setVersion((v) => v + 1);
      });
  }, 10);

  return {
    ready: false,
    pending: true,
  };
};
