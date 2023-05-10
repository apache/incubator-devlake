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

import { useState, useEffect, useMemo } from 'react';

import type { PluginConfigType } from '@/plugins';
import { PluginConfig, PluginType } from '@/plugins';

import type { ConnectionItemType } from './types';
import { ConnectionStatusEnum } from './types';
import * as API from './api';

export interface UseContextValueProps {
  plugin?: string;
  filterBeta?: boolean;
  filterPlugin?: string[];
  filter?: string[];
}

export const useContextValue = ({ plugin, filterBeta, filterPlugin, filter }: UseContextValueProps) => {
  const [loading, setLoading] = useState(true);
  const [connections, setConnections] = useState<ConnectionItemType[]>([]);

  const plugins = useMemo(
    () =>
      PluginConfig.filter((p) => p.type === PluginType.Connection)
        .filter((p) => (plugin ? p.plugin === plugin : true))
        .filter((p) => (filterBeta ? !p.isBeta : true))
        .filter((p) => (filterPlugin ? !filterPlugin.includes(p.plugin) : true)),
    [plugin],
  );

  const getConnection = async (plugin: string) => {
    try {
      const res = await API.getConnection(plugin);
      const { icon, entities } = plugins.find((p) => p.plugin === plugin) as PluginConfigType;

      return res.map((connection) => ({
        ...connection,
        plugin,
        icon,
        entities,
      }));
    } catch {
      return [];
    }
  };

  const testConnection = async ({
    plugin,
    endpoint,
    proxy,
    token,
    username,
    password,
    authMethod,
  }: ConnectionItemType) => {
    try {
      const res = await API.testConnection(plugin, {
        endpoint,
        proxy,
        token,
        username,
        password,
        authMethod,
      });
      return res.success ? ConnectionStatusEnum.ONLINE : ConnectionStatusEnum.OFFLINE;
    } catch {
      return ConnectionStatusEnum.OFFLINE;
    }
  };

  const transformConnection = (connections: Omit<ConnectionItemType, 'unique' | 'status'>[]) => {
    return connections.map((it) => ({
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
    }));
  };

  const handleRefresh = async (plugin?: string) => {
    if (plugin) {
      const res = await getConnection(plugin);
      setConnections([...connections.filter((cs) => cs.plugin !== plugin), ...transformConnection(res)]);
      return;
    }

    const res = await Promise.all(plugins.map((cs) => getConnection(cs.plugin)));

    setConnections(transformConnection(res.flat()));
    setLoading(false);
  };

  const handleTest = async (unique: string) => {
    setConnections((connections) =>
      connections.map((cs) =>
        cs.unique === unique
          ? {
              ...cs,
              status: ConnectionStatusEnum.TESTING,
            }
          : cs,
      ),
    );

    console.log(connections);

    const connection = connections.find((cs) => cs.unique === unique) as ConnectionItemType;
    const status = await testConnection(connection);

    setConnections((connections) =>
      connections.map((cs) =>
        cs.unique === unique
          ? {
              ...cs,
              status,
            }
          : cs,
      ),
    );
  };

  useEffect(() => {
    handleRefresh();
  }, []);

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
