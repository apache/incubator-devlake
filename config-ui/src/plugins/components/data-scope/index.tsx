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

import { useState } from 'react';
import { Button, Intent } from '@blueprintjs/core';

import { Table, Dialog } from '@/components';
import { DataScopeForm } from '@/plugins';

import { BPConnectionType } from '../../types';

import * as S from './styled';

interface Props {
  connections: BPConnectionType[];
  cancelBtnProps?: {
    text?: string;
  };
  submitBtnProps?: {
    text: string;
  };
  onCancel?: () => void;
  onSubmit?: () => void;
  onChange?: (connections: BPConnectionType[]) => void;
}

export const DataScope = ({ connections, cancelBtnProps, submitBtnProps, onCancel, onSubmit, onChange }: Props) => {
  const [connection, setConnection] = useState<BPConnectionType>();

  const handleCancel = () => setConnection(undefined);

  const handleSubmit = (
    connection: BPConnectionType,
    scope: BPConnectionType['scope'],
    origin: BPConnectionType['origin'],
  ) => {
    onChange?.(
      connections.map((cs) => {
        if (cs.unique === connection.unique) {
          return {
            ...cs,
            scope,
            origin,
          };
        }
        return cs;
      }),
    );
    handleCancel();
  };

  if (connections.length === 1) {
    const [{ plugin, connectionId, ...props }] = connections;
    return (
      <DataScopeForm
        plugin={plugin}
        connectionId={connectionId}
        cancelBtnProps={cancelBtnProps}
        submitBtnProps={submitBtnProps}
        onCancel={onCancel}
        onSubmit={(scope: BPConnectionType['scope'], origin: BPConnectionType['origin']) => {
          onChange?.([
            {
              plugin,
              connectionId,
              ...props,
              scope,
              origin,
            },
          ]);
          onSubmit?.();
        }}
      />
    );
  }

  return (
    <S.Wrapper>
      <Table
        columns={[
          {
            title: 'Data Connections',
            dataIndex: ['icon', 'name'],
            key: 'connection',
            render: ({ icon, name }) => (
              <S.ConnectionColumn>
                <img src={icon} alt="" />
                <span>{name}</span>
              </S.ConnectionColumn>
            ),
          },
          {
            title: 'Data Scope',
            dataIndex: 'origin',
            key: 'scope',
            render: (scope: BPConnectionType['origin']) =>
              !scope.length ? (
                <span>No Data Scope Selected</span>
              ) : (
                <S.ScopeColumn>
                  {scope.map((sc, i) => (
                    <S.ScopeItem key={i}>
                      <span>{sc.name}</span>
                    </S.ScopeItem>
                  ))}
                </S.ScopeColumn>
              ),
          },
          {
            title: '',
            dataIndex: 'id',
            key: 'action',
            align: 'center',
            render: (_, connection) => (
              <Button
                small
                minimal
                intent={Intent.PRIMARY}
                icon="cog"
                text="Set Data Scope"
                onClick={() => setConnection(connection)}
              />
            ),
          },
        ]}
        dataSource={connections}
      />
      <S.Btns>
        <Button outlined intent={Intent.PRIMARY} text="Cancel" onClick={onCancel} {...cancelBtnProps} />
        <Button intent={Intent.PRIMARY} text="Save" onClick={onSubmit} {...submitBtnProps} />
      </S.Btns>
      {connection && (
        <Dialog isOpen title="Set Data Scope" footer={null} style={{ width: 820 }} onCancel={handleCancel}>
          <div className="connection">
            <img src={connection.icon} width={24} alt="" />
            <span>{connection.name}</span>
          </div>
          <DataScopeForm
            plugin={connection.plugin}
            connectionId={connection.connectionId}
            onCancel={handleCancel}
            onSubmit={(scope: BPConnectionType['scope'], origin: BPConnectionType['origin']) =>
              handleSubmit(connection, scope, origin)
            }
          />
        </Dialog>
      )}
    </S.Wrapper>
  );
};
