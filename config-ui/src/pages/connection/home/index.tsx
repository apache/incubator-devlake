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

import React, { useMemo } from 'react';
import { useHistory } from 'react-router-dom';
import { Tag, Intent } from '@blueprintjs/core';

import { PluginConfig, PluginType } from '@/plugins';

import * as S from './styled';

export const ConnectionHomePage = () => {
  const history = useHistory();

  const [connections, webhook] = useMemo(
    () => [
      PluginConfig.filter((p) => p.type === PluginType.Connection),
      PluginConfig.filter((p) => p.plugin === 'webhook'),
    ],
    [],
  );

  return (
    <S.Wrapper>
      <div className="block">
        <h2>Data Connections</h2>
        <p>Connections are available for data collection.</p>
        <ul>
          {connections.map((cs) => (
            <li key={cs.plugin} onClick={() => history.push(`/connections/${cs.plugin}`)}>
              <img src={cs.icon} alt="" />
              <span>{cs.name}</span>
              {cs.isBeta && (
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
        <p>
          You can use webhooks to import deployments and incidents from the unsupported data integrations to calculate
          DORA metrics, etc. Please note: webhooks cannot be created or managed in Blueprints.
        </p>
        <ul>
          {webhook.map((cs) => (
            <li key={cs.plugin} onClick={() => history.push(`/connections/${cs.plugin}`)}>
              <img src={cs.icon} alt="" />
              <span>{cs.name}</span>
            </li>
          ))}
        </ul>
      </div>
    </S.Wrapper>
  );
};
