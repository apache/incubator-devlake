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
import { Icon, Colors, Position, Intent } from '@blueprintjs/core';
import { Tooltip2 } from '@blueprintjs/popover2';
import styled from 'styled-components';

import { Loading } from '@/components';

import type { ConnectionItemType } from './types';
import { ConnectionStatusEnum } from './types';

const Wrapper = styled.div`
  display: inline-flex;
  align-items: center;

  & > span.online {
    color: ${Colors.GREEN3};
  }

  & > span.offline {
    color: ${Colors.RED3};
  }

  & > span.testing {
    color: #7497f7;
  }
`;

const STATUS_MAP = {
  [`${ConnectionStatusEnum.NULL}`]: 'Init',
  [`${ConnectionStatusEnum.TESTING}`]: 'Testing',
  [`${ConnectionStatusEnum.ONLINE}`]: 'Online',
  [`${ConnectionStatusEnum.OFFLINE}`]: 'Offline',
};

interface Props {
  connection: ConnectionItemType;
  onTest: (connection: ConnectionItemType) => void;
}

export const ConnectionStatus = ({ connection, onTest }: Props) => {
  const { status } = connection;

  return (
    <Wrapper>
      {status === ConnectionStatusEnum.TESTING && <Loading size={14} style={{ marginRight: 4 }} />}
      {status === ConnectionStatusEnum.OFFLINE && (
        <Tooltip2 intent={Intent.PRIMARY} position={Position.TOP} content="Retry">
          <Icon
            size={14}
            icon="repeat"
            style={{ marginRight: 4, color: Colors.RED3, cursor: 'pointer' }}
            onClick={() => onTest(connection)}
          />
        </Tooltip2>
      )}
      <span className={status}>{STATUS_MAP[status]}</span>
    </Wrapper>
  );
};
