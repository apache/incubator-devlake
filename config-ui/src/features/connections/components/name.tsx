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
import styled from 'styled-components';

import { useAppSelector } from '@/hooks';
import { getPluginConfig } from '@/plugins';

import { selectConnection, selectWebhook } from '../slice';

const Wrapper = styled.div`
  display: flex;
  align-items: center;

  span + span {
    margin-left: 4px;
  }

  .icon {
    display: inline-block;
    width: 24px;
    height: 24px;

    & > svg {
      width: 100%;
      height: 100%;
    }
  }

  .name {
    max-width: 240px;
    white-space: nowrap;
    text-overflow: ellipsis;
    overflow: hidden;
  }
`;

interface Props {
  plugin: string;
  connectionId: ID;
}

export const ConnectionName = ({ plugin, connectionId }: Props) => {
  const {
    token: { colorPrimary },
  } = theme.useToken();

  const connection = useAppSelector((state) => selectConnection(state, `${plugin}-${connectionId}`));
  const webhook = useAppSelector((state) => selectWebhook(state, connectionId));
  const config = getPluginConfig(plugin);

  const name = connection ? connection.name : webhook ? webhook.name : `${plugin}/connection/${connectionId}`;

  return (
    <Wrapper>
      <span className="icon">{config.icon({ color: colorPrimary })}</span>
      <span className="name" title={name}>
        {name}
      </span>
    </Wrapper>
  );
};
