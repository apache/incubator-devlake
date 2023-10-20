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
import { useNavigate } from 'react-router-dom';
import { Tag, Intent } from '@blueprintjs/core';

import { useAppSelector } from '@/app/hook';
import { Dialog } from '@/components';
import { selectAllConnections } from '@/features/connections';
import { PluginConfig, PluginConfigType, ConnectionList, ConnectionForm } from '@/plugins';

import * as S from './styled';

export const ConnectionHomePage = () => {
  const [type, setType] = useState<'list' | 'form'>();
  const [pluginConfig, setPluginConfig] = useState<PluginConfigType>();

  const connections = useAppSelector(selectAllConnections);

  const navigate = useNavigate();

  const plugins = useMemo(
    () =>
      PluginConfig.map((p) => ({
        ...p,
        count: connections.filter((cs) => cs.plugin === p.plugin).length,
      })),
    [connections],
  );

  const handleShowListDialog = (pluginConfig: PluginConfigType) => {
    setType('list');
    setPluginConfig(pluginConfig);
  };

  const handleShowFormDialog = () => {
    setType('form');
  };

  const handleHideDialog = () => {
    setType(undefined);
    setPluginConfig(undefined);
  };

  const handleSuccessAfter = async (plugin: string, id: ID) => {
    navigate(`/connections/${plugin}/${id}`);
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
          {plugins
            .filter((p) => p.plugin !== 'webhook')
            .map((p) => (
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
          {plugins
            .filter((p) => p.plugin === 'webhook')
            .map((p) => (
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
          <ConnectionList plugin={pluginConfig.plugin} onCreate={handleShowFormDialog} />
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
            onSuccess={(id) => handleSuccessAfter(pluginConfig.plugin, id)}
          />
        </Dialog>
      )}
    </S.Wrapper>
  );
};
