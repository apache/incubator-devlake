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

import React from 'react';
import type { IntentProps } from '@blueprintjs/core';
import { Position } from '@blueprintjs/core';
import { Tooltip2 } from '@blueprintjs/popover2';
import styled from 'styled-components';

const Wrapper = styled.div`
  width: 100%;

  & > .bp4-popover2-target {
    display: block;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
`;

interface Props extends IntentProps {
  content: string;
  children: React.ReactNode;
  style?: React.CSSProperties;
}

export const TextTooltip = ({ intent, content, children, style }: Props) => {
  return (
    <Wrapper style={style}>
      <Tooltip2 intent={intent} position={Position.TOP} content={content}>
        {children}
      </Tooltip2>
    </Wrapper>
  );
};
