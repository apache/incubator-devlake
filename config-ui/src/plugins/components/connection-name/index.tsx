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

import { theme } from 'antd';

import { selectConnection, selectWebhook } from '@/features/connections';
import { useAppSelector } from '@/hooks';
import { getPluginConfig } from '@/plugins';

import * as S from './styled';

interface Props {
  plugin: string;
  connectionId?: ID;
  customName?: (pluginName: string) => string;
  onClick?: () => void;
}

const Name = ({ plugin, connectionId }: Required<Pick<Props, 'plugin' | 'connectionId'>>) => {
  const connection = useAppSelector((state) => selectConnection(state, `${plugin}-${connectionId}`));
  const webhook = useAppSelector((state) => selectWebhook(state, connectionId));

  return (
    <S.Name>{connection ? connection.name : webhook ? webhook.name : `${plugin}/connection/${connectionId}`}</S.Name>
  );
};

export const ConnectionName = ({ plugin, connectionId, customName, onClick }: Props) => {
  const {
    token: { colorPrimary },
  } = theme.useToken();
  const config = getPluginConfig(plugin);

  if (!connectionId) {
    return (
      <S.Wrapper onClick={onClick}>
        <S.Icon>{config.icon({ color: colorPrimary })}</S.Icon>
        <S.Name title={config.name}>{config.name}</S.Name>
      </S.Wrapper>
    );
  }

  if (customName) {
    return (
      <S.Wrapper onClick={onClick}>
        <S.Icon>{config.icon({ color: colorPrimary })}</S.Icon>
        <S.Name>{customName(config.name)}</S.Name>
      </S.Wrapper>
    );
  }

  return (
    <S.Wrapper onClick={onClick}>
      <S.Icon>{config.icon({ color: colorPrimary })}</S.Icon>
      <Name plugin={plugin} connectionId={connectionId} />
    </S.Wrapper>
  );
};
