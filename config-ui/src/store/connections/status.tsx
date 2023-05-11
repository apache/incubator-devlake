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

import styled from 'styled-components';

import { IconButton } from '@/components';

import { ConnectionStatusEnum } from './types';

const Wrapper = styled.div`
  display: inline-flex;
  align-items: center;

  & > span.online {
    color: #4db764;
  }

  & > span.offline {
    color: #e34040;
  }
`;

const STATUS_MAP = {
  [`${ConnectionStatusEnum.NULL}`]: 'Test',
  [`${ConnectionStatusEnum.TESTING}`]: 'Testing',
  [`${ConnectionStatusEnum.ONLINE}`]: 'Connected',
  [`${ConnectionStatusEnum.OFFLINE}`]: 'Disconnected',
};

interface Props {
  status: ConnectionStatusEnum;
  unique: string;
  onTest: (unique: string) => void;
}

export const ConnectionStatus = ({ status, unique, onTest }: Props) => {
  return (
    <Wrapper>
      <span className={status}>{STATUS_MAP[status]}</span>
      {status !== ConnectionStatusEnum.ONLINE && (
        <IconButton
          loading={status === ConnectionStatusEnum.TESTING}
          icon="repeat"
          tooltip="Retry"
          onClick={() => onTest(unique)}
        />
      )}
    </Wrapper>
  );
};
