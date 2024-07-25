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

import { useState } from 'react';
import { Table } from 'antd';

import API from '@/api';
import { useRefreshData } from '@/hooks';
import { getPluginScopeId, ScopeConfig } from '@/plugins';

interface Props {
  plugin: string;
  connectionId: ID;
  scopeIds: ID[];
}

export const BlueprintConnectionDetailTable = ({ plugin, connectionId, scopeIds }: Props) => {
  const [version, setVersion] = useState(1);

  const { ready, data } = useRefreshData(async () => {
    const scopes = await Promise.all(scopeIds.map((scopeId) => API.scope.get(plugin, connectionId, scopeId)));
    return scopes.map((sc) => ({
      id: getPluginScopeId(plugin, sc.scope),
      name: sc.scope.fullName ?? sc.scope.name,
      scopeConfigId: sc.scopeConfig?.id,
      scopeConfigName: sc.scopeConfig?.name,
    }));
  }, [version]);

  return (
    <Table
      loading={!ready}
      rowKey="id"
      size="middle"
      columns={[
        {
          title: 'Data Scope',
          dataIndex: 'name',
          key: 'name',
        },
        {
          title: 'Scope Config',
          key: 'scopeConfig',
          render: (_, { id, name, scopeConfigId, scopeConfigName }) => (
            <ScopeConfig
              plugin={plugin}
              connectionId={connectionId}
              scopeId={id}
              scopeName={name}
              scopeConfigId={scopeConfigId}
              scopeConfigName={scopeConfigName}
              onSuccess={() => setVersion(version + 1)}
            />
          ),
        },
      ]}
      dataSource={data ?? []}
    />
  );
};
