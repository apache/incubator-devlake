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

import { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { Tag, Intent, Button } from '@blueprintjs/core';

import { Dialog, Table } from '@/components';
import type { PluginConfigType } from '@/plugins';
import { PluginConfig, PluginType, ConnectionForm } from '@/plugins';
import { ConnectionContextProvider, useConnection, ConnectionStatus } from '@/store';

import * as S from './styled';

export const ConnectionHome = () => {
  const [type, setType] = useState<'list' | 'form'>();
  const [pluginConfig, setPluginConfig] = useState<PluginConfigType>();

  const { connections, onRefresh, onTest } = useConnection();

  const [plugins, webhook] = useMemo(
    () => [
      PluginConfig.filter((p) => p.type === PluginType.Connection && p.plugin !== 'webhook').map((p) => ({
        ...p,
        count: connections.filter((cs) => cs.plugin === p.plugin).length,
      })),
      {
        ...(PluginConfig.find((p) => p.plugin === 'webhook') as PluginConfigType),
        count: connections.filter((cs) => cs.plugin === 'webhook').length,
      },
    ],
    [],
  );

  const handleShowListDialog = (config: PluginConfigType) => {
    setType('list');
    setPluginConfig(config);
  };

  const handleShowFormDialog = () => {
    setType('form');
  };

  const handleHideDialog = () => {
    setType(undefined);
    setPluginConfig(undefined);
  };

  const handleCreateSuccess = async (unqie: string, plugin: string) => {
    onRefresh(plugin);
    setType('list');
  };

  return (
    <S.Wrapper>
      <div className="block">
        <h1>Connections</h1>
        <h5>
          Create and manage data connections from the following data sources or Webhooks to be used in syncing data in
          your Projects.
        </h5>
      </div>
      <div className="block">
        <h2>Data Connections</h2>
        <h5>
          You can create and manage data connections for the following data sources and use them in your Projects.
        </h5>
        <ul>
          {plugins.map((p) => (
            <li key={p.plugin} onClick={() => handleShowListDialog(p)}>
              <img src={p.icon} alt="" />
              <span className="name">{p.name}</span>
              <S.Count>{p.count ? `${p.count} connections` : 'No connection'}</S.Count>
              {p.isBeta && (
                <Tag intent={Intent.WARNING} round>
                  beta
                </Tag>
              )}
            </li>
          ))}
        </ul>
      </div>
      <div className="block">
        <h2>Webhooks</h2>
        <h5>
          You can use webhooks to import deployments and incidents from the unsupported data integrations to calculate
          DORA metrics, etc.
        </h5>
        <ul>
          <li onClick={() => handleShowListDialog(webhook)}>
            <img src={webhook.icon} alt="" />
            <span className="name">{webhook.name}</span>
            <S.Count>{webhook.count ? `${webhook.count} connections` : 'No connection'}</S.Count>
          </li>
        </ul>
      </div>
      {type === 'list' && pluginConfig && (
        <Dialog
          style={{ width: 820 }}
          isOpen
          title={
            <S.DialogTitle>
              <img src={pluginConfig.icon} alt="" />
              <span>Manage Connections: {pluginConfig.name}</span>
            </S.DialogTitle>
          }
          footer={null}
          onCancel={handleHideDialog}
        >
          <Table
            noShadow
            columns={[
              {
                title: 'Connection Name',
                dataIndex: 'name',
                key: 'name',
              },
              {
                title: 'Status',
                dataIndex: ['status', 'unique'],
                key: 'status',
                render: ({ status, unique }) => <ConnectionStatus status={status} unique={unique} onTest={onTest} />,
              },
              {
                title: '',
                dataIndex: ['plugin', 'id'],
                key: 'link',
                width: 100,
                render: ({ plugin, id }) => <Link to={`/connections/${plugin}/${id}`}>Details</Link>,
              },
            ]}
            dataSource={connections.filter((cs) => cs.plugin === pluginConfig.plugin)}
            noData={{
              text: 'There is no data connection yet. Please add a new connection.',
            }}
          />
          <Button
            style={{ marginTop: 16 }}
            intent={Intent.PRIMARY}
            icon="add"
            text="Create a New Connection"
            onClick={handleShowFormDialog}
          />
        </Dialog>
      )}
      {type === 'form' && pluginConfig && (
        <Dialog
          style={{ width: 820 }}
          isOpen
          title={
            <S.DialogTitle>
              <img src={pluginConfig.icon} alt="" />
              <span>Manage Connections: {pluginConfig.name}</span>
            </S.DialogTitle>
          }
          footer={null}
          onCancel={handleHideDialog}
        >
          <ConnectionForm
            plugin={pluginConfig.plugin}
            onSuccess={(unique) => handleCreateSuccess(unique, pluginConfig.plugin)}
          />
        </Dialog>
      )}
    </S.Wrapper>
  );
};

export const ConnectionHomePage = () => {
  return (
    <ConnectionContextProvider>
      <ConnectionHome />
    </ConnectionContextProvider>
  );
};
