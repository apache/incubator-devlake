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

import { useAppDispatch } from '@/app/hook';
import { IconButton } from '@/components';
import { testConnection } from '@/features/connections';
import { IConnection, IConnectionStatus, IPlugin } from '@/types';

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
  [`${IConnectionStatus.IDLE}`]: 'Test',
  [`${IConnectionStatus.TESTING}`]: 'Testing',
  [`${IConnectionStatus.ONLINE}`]: 'Connected',
  [`${IConnectionStatus.OFFLINE}`]: 'Disconnected',
};

interface Props {
  connection: IConnection;
}

export const ConnectionStatus = ({ connection }: Props) => {
  const { status } = connection;

  const dispatch = useAppDispatch();

  const handleTest = () => dispatch(testConnection(connection));

  return (
    <Wrapper>
      <span className={status}>{STATUS_MAP[status]}</span>
      {status !== IConnectionStatus.ONLINE && (
        <IconButton loading={status === IConnectionStatus.TESTING} icon="repeat" tooltip="Retry" onClick={handleTest} />
      )}
    </Wrapper>
  );
};
