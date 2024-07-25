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
import { theme } from 'antd';
import styled from 'styled-components';

import { getPluginConfig } from '@/plugins';

const Wrapper = styled.div`
  display: flex;
  align-items: center;

  .icon {
    display: inline-flex;
    margin-right: 8px;
    width: 24px;

    & > svg {
      width: 100%;
      height: 100%;
    }
  }
`;

interface Props {
  plugin: string;
  name: string;
}

export const PluginName = ({ plugin, name }: Props) => {
  const pluginConfig = useMemo(() => getPluginConfig(plugin), [plugin]);

  const {
    token: { colorPrimary },
  } = theme.useToken();

  return (
    <Wrapper>
      <span className="icon">{pluginConfig.icon({ color: colorPrimary })}</span>
      <span>{name}</span>
    </Wrapper>
  );
};
