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

import type { ConnectionItemType } from './types';
import { ConnectionStatusEnum } from './types';
import * as API from './api';

export interface UseContextValueProps {
  plugin?: string;
  filterBeta?: boolean;
  filter?: string[];
}

export const useContextValue = ({ plugin, filterBeta = false, filter }: UseContextValueProps) => {
  const [loading, setLoading] = useState(false);
  const [connections, setConnections] = useState<ConnectionItemType[]>([]);

  const allConnections = useMemo(
    () =>
      PluginConfig.filter((p) => p.type === PluginType.Connection)
        .filter((p) => (plugin ? p.plugin === plugin : true))
        .filter((p) => (filterBeta ? !p.isBeta : true)),
    [plugin],
  );

  const getConnection = async (plugin: string) => {
    try {
      return await API.getConnection(plugin);
    } catch {
      return [];
    }
  };

  const handleRefresh = useCallback(async () => {
    setLoading(true);

    const res = await Promise.all(allConnections.map((cs) => getConnection(cs.plugin)));

    const resWithPlugin = res.map((cs, i) =>
      cs.map((it: any) => {
        const { plugin, icon, entities } = allConnections[i];

        return {
          ...it,
          plugin,
          icon,
          entities,
        };
      }),
    );

    setConnections(
      resWithPlugin.flat().map((it) => ({
        unique: `${it.plugin}-${it.id}`,
        status: ConnectionStatusEnum.NULL,
        plugin: it.plugin,
        id: it.id,
        name: it.name,
        icon: it.icon,
        entities: it.entities,
        endpoint: it.endpoint,
        proxy: it.proxy,
        token: it.token,
        username: it.username,
        password: it.password,
        authMethod: it.authMethod,
      })),
    );

    setLoading(false);
  }, [allConnections]);

  useEffect(() => {
    handleRefresh();
  }, [allConnections]);

  const handleTest = useCallback(
    async (selectedConnection: ConnectionItemType) => {
      setConnections((connections) =>
        connections.map((cs) =>
          cs.unique === selectedConnection.unique
            ? {
                ...cs,
                status: ConnectionStatusEnum.TESTING,
              }
            : cs,
        ),
      );

      const { plugin, endpoint, proxy, token, username, password, authMethod } = selectedConnection;

      let status = ConnectionStatusEnum.OFFLINE;

      try {
        const res = await API.testConnection(plugin, {
          endpoint,
          proxy,
          token,
          username,
          password,
          authMethod,
        });
        status = res.success ? ConnectionStatusEnum.ONLINE : ConnectionStatusEnum.OFFLINE;
      } catch {
        status = ConnectionStatusEnum.OFFLINE;
      }

      setConnections((connections) =>
        connections.map((cs) =>
          cs.unique === selectedConnection.unique
            ? {
                ...cs,
                status,
              }
            : cs,
        ),
      );
    },
    [connections],
  );

  return useMemo(
    () => ({
      loading,
      connections: filter ? connections.filter((cs) => !filter.includes(cs.unique)) : connections,
      onRefresh: handleRefresh,
      onTest: handleTest,
    }),
    [loading, connections],
  );
};
