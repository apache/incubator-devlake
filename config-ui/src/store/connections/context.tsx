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

import React, { useState, useEffect, useMemo } from 'react';

import { PageLoading } from '@/components';

import type { PluginConfigType } from '@/plugins';
import { PluginConfig, PluginType } from '@/plugins';

import type { ConnectionItemType } from './types';
import { ConnectionStatusEnum } from './types';
import * as API from './api';

export const ConnectionContext = React.createContext<{
  connections: ConnectionItemType[];
  onGet: (unique: string) => ConnectionItemType;
  onTest: (unique: string) => void;
  onRefresh: (plugin?: string) => void;
}>(undefined!);

interface Props {
  children?: React.ReactNode;
}

export const ConnectionContextProvider = ({ children, ...props }: Props) => {
  const [loading, setLoading] = useState(true);
  const [connections, setConnections] = useState<ConnectionItemType[]>([]);

  const plugins = useMemo(() => PluginConfig.filter((p) => p.type === PluginType.Connection), []);

  const queryConnection = async (plugin: string) => {
    try {
      const res = await API.getConnection(plugin);
      const { name, icon, isBeta, entities } = plugins.find((p) => p.plugin === plugin) as PluginConfigType;

      return res.map((connection) => ({
        ...connection,
        plugin,
        pluginName: name,
        icon,
        isBeta: isBeta ?? false,
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
    secretKey,
    appId,
  }: ConnectionItemType) => {
    try {
      const res = await API.testConnection(plugin, {
        endpoint,
        proxy,
        token,
        username,
        password,
        authMethod,
        secretKey,
        appId,
      });
      return res.success ? ConnectionStatusEnum.ONLINE : ConnectionStatusEnum.OFFLINE;
    } catch {
      return ConnectionStatusEnum.OFFLINE;
    }
  };

  const transformConnection = (connections: Omit<ConnectionItemType, 'unique' | 'status'>[]) => {
    return connections.map((it) => ({
      unique: `${it.plugin}-${it.id}`,
      plugin: it.plugin,
      pluginName: it.pluginName,
      id: it.id,
      name: it.name,
      status: ConnectionStatusEnum.NULL,
      icon: it.icon,
      isBeta: it.isBeta,
      entities: it.entities,
      endpoint: it.endpoint,
      proxy: it.proxy,
      token: it.token,
      username: it.username,
      password: it.password,
      authMethod: it.authMethod,
      secretKey: it.secretKey,
      appId: it.appId,
    }));
  };



  const handleGet = (unique: string) => {
    return connections.find((cs) => cs.unique === unique) as ConnectionItemType;
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

    const connection = handleGet(unique);
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

  const handleRefresh = async (plugin?: string) => {
    if (plugin) {
      const res = await queryConnection(plugin);
      setConnections([...connections.filter((cs) => cs.plugin !== plugin), ...transformConnection(res)]);
      return;
    }

    const res = await Promise.all(plugins.map((cs) => queryConnection(cs.plugin)));

    setConnections(transformConnection(res.flat()));
    setLoading(false);
  };

  useEffect(() => {
    handleRefresh();
  }, []);

  if (loading) {
    return <PageLoading />;
  }

  return (
    <ConnectionContext.Provider
      value={{
        connections,
        onGet: handleGet,
        onTest: handleTest,
        onRefresh: handleRefresh,
      }}
    >
      {children}
    </ConnectionContext.Provider>
  );
};
