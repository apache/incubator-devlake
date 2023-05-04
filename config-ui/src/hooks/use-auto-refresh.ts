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

import { useState, useEffect, useMemo, useRef } from 'react';

export const useAutoRefresh = <T>(
  request: () => Promise<T>,
  deps: React.DependencyList = [],
  option?: {
    cancel?: (data?: T) => boolean;
    interval?: number;
    retryLimit?: number;
  },
) => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<T>();

  const timer = useRef<any>();
  const retryCount = useRef<number>(0);

  useEffect(() => {
    setLoading(true);
    request()
      .then((data: T) => {
        setData(data);
      })
      .finally(() => {
        setLoading(false);
      });
  }, [...deps]);

  useEffect(() => {
    timer.current = setInterval(() => {
      retryCount.current += 1;
      request().then((data) => setData(data));
    }, option?.interval ?? 5000);
    return () => clearInterval(timer.current);
  }, [...deps]);

  useEffect(() => {
    if (option?.cancel?.(data) || (option?.retryLimit && option?.retryLimit <= retryCount.current)) {
      clearInterval(timer.current);
    }
  }, [data]);

  return useMemo(
    () => ({
      loading,
      data,
    }),
    [loading, data],
  );
};
