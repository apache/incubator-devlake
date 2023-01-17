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

import { useState, useEffect, useCallback, useMemo } from 'react';

import { PluginConfig, PluginType } from '@/plugins';

import type { TransformationItemType } from './types';
import * as API from './api';

export const useContextValue = () => {
  const [loading, setLoading] = useState(false);
  const [transformations, setTransformations] = useState<TransformationItemType[]>([]);

  const allConnections = useMemo(() => PluginConfig.filter((p) => p.type === PluginType.Connection && !p.isBeta), []);

  const getTransformation = async (plugin: string) => {
    try {
      return await API.getTransformation(plugin);
    } catch {
      return [];
    }
  };

  const handleRefresh = useCallback(async () => {
    setLoading(true);

    const res = await Promise.all(allConnections.map((cs) => getTransformation(cs.plugin)));

    const resWithPlugin = res.map((ts, i) =>
      ts.map((it: any) => {
        const { plugin } = allConnections[i];

        return {
          ...it,
          plugin,
        };
      }),
    );

    setTransformations(resWithPlugin.flat());
    setLoading(false);
  }, [allConnections]);

  useEffect(() => {
    handleRefresh();
  }, []);

  return useMemo(
    () => ({
      loading,
      plugins: allConnections,
      transformations,
    }),
    [loading, allConnections, transformations],
  );
};
