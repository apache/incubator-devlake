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

import { useMemo } from 'react';
import { useHistory } from 'react-router-dom';
import { Tag, Intent } from '@blueprintjs/core';

import type { PluginConfigType } from '@/plugins';
import { PluginConfig, PluginType } from '@/plugins';
import { ConnectionContextProvider, ConnectionContextConsumer } from '@/store';

import { Count } from './count';
import * as S from './styled';

export const ConnectionHomePage = () => {
  const history = useHistory();

  const [plugins, webhook] = useMemo(
    () => [
      PluginConfig.filter((p) => p.type === PluginType.Connection && p.plugin !== 'webhook'),
      PluginConfig.find((p) => p.plugin === 'webhook') as PluginConfigType,
    ],
    [],
  );

  return (
    <ConnectionContextProvider>
      <ConnectionContextConsumer>
        {({ connections }) => (
          <S.Wrapper>
            <div className="block">
              <h1>Connections</h1>
              <h5>
                Create and manage data connections from the following data sources or Webhooks to be used in syncing
                data in your Projects.
              </h5>
            </div>
            <div className="block">
              <h2>Data Connections</h2>
              <h5>
                You can create and manage data connections for the following data sources and use them in your Projects.
              </h5>
              <ul>
                {plugins.map((p) => (
                  <li key={p.plugin} onClick={() => history.push(`/connections/${p.plugin}`)}>
                    <img src={p.icon} alt="" />
                    <span className="name">{p.name}</span>
                    <Count count={connections.filter((cs) => cs.plugin === p.plugin).length} />
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
                You can use webhooks to import deployments and incidents from the unsupported data integrations to
                calculate DORA metrics, etc.
              </h5>
              <ul>
                <li onClick={() => history.push(`/connections/${webhook.plugin}`)}>
                  <img src={webhook.icon} alt="" />
                  <span className="name">{webhook.name}</span>
                  <Count count={connections.filter((cs) => cs.plugin === 'webhook').length} />
                </li>
              </ul>
            </div>
          </S.Wrapper>
        )}
      </ConnectionContextConsumer>
    </ConnectionContextProvider>
  );
};
