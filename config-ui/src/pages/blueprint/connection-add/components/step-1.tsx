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

import { Button, Intent } from '@blueprintjs/core';

import { Table } from '@/components';
import { getPluginConfig } from '@/plugins';
import type { ConnectionItemType } from '@/store';
import { ConnectionContextProvider, ConnectionContextConsumer } from '@/store';

import { useConnectionAdd } from '../context';

import * as S from './styled';

export const Step1 = () => {
  const { filter, connection, onChangeConnection, onCancel, onNext } = useConnectionAdd();

  return (
    <ConnectionContextProvider filterBeta filter={filter}>
      <ConnectionContextConsumer>
        {({ connections }) => (
          <S.Wrapper>
            <Table
              columns={[{ title: 'Data Connection', dataIndex: 'name', key: 'name' }]}
              dataSource={connections}
              rowSelection={{
                rowKey: 'unique',
                type: 'radio',
                selectedRowKeys: connection?.unique ? [connection?.unique] : [],
                onChange: (selectedRowKeys) => {
                  const unique = selectedRowKeys[0];
                  const connection = connections.find((cs) => cs.unique === unique) as ConnectionItemType;
                  const config = getPluginConfig(connection.plugin);
                  onChangeConnection({
                    unique: connection.unique,
                    plugin: connection.plugin,
                    connectionId: connection.id,
                    name: connection.name,
                    icon: connection.icon,
                    scope: [],
                    origin: [],
                    transformationType: config.transformationType,
                  });
                },
              }}
            />
            <S.Action>
              <Button outlined intent={Intent.PRIMARY} onClick={onCancel}>
                Cancel
              </Button>
              <Button intent={Intent.PRIMARY} disabled={!connection} onClick={onNext}>
                Next Step
              </Button>
            </S.Action>
          </S.Wrapper>
        )}
      </ConnectionContextConsumer>
    </ConnectionContextProvider>
  );
};
