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

import type { IntentProps } from '@blueprintjs/core';
import { Position, Tooltip } from '@blueprintjs/core';
import styled from 'styled-components';

const Wrapper = styled.div`
  width: 100%;

  & > .bp5-popover-target {
    width: 100%;

    & > span {
      display: block;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }
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
      <Tooltip intent={intent} position={Position.TOP} content={content}>
        {children}
      </Tooltip>
    </Wrapper>
  );
};
