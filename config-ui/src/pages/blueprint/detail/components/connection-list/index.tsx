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
import { Link } from 'react-router-dom';

import { useConnections } from '@/hooks';
import { getPluginConfig } from '@/plugins';

import type { BlueprintType } from '../../../types';

import * as S from './styled';

interface Props {
  path: string;
  blueprint: BlueprintType;
}

export const ConnectionList = ({ path, blueprint }: Props) => {
  const { onGet } = useConnections();

  const list = useMemo(
    () =>
      blueprint.settings?.connections
        .filter((cs) => cs.plugin !== 'webhook')
        .map((cs: any) => {
          const unique = `${cs.plugin}-${cs.connectionId}`;
          const plugin = getPluginConfig(cs.plugin);
          const connection = onGet(unique);

          return {
            unique,
            icon: plugin.icon,
            name: connection.name,
            scope: cs.scopes,
          };
        })
        .filter(Boolean),
    [blueprint],
  );

  return (
    <S.List>
      {list.map((cs) => (
        <S.Item key={cs.unique}>
          <div className="title">
            <img src={cs.icon} alt="" />
            <span>{cs.name}</span>
          </div>
          <div className="count">
            <span>{cs.scope.length} data scope</span>
          </div>
          <div className="link">
            <Link to={`${path}${cs.unique}`}>View Detail</Link>
          </div>
        </S.Item>
      ))}
    </S.List>
  );
};
