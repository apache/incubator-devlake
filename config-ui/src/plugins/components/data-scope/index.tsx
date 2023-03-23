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

import { useState, useMemo } from 'react';
import { Button, Icon, Intent, Position, Colors } from '@blueprintjs/core';
import { Tooltip2 } from '@blueprintjs/popover2';

import { Table, Dialog } from '@/components';
import { DataScopeForm } from '@/plugins';

import * as S from './styled';

interface Props {
  connections: MixConnection[];
  cancelBtnProps?: {
    text?: string;
  };
  submitBtnProps?: {
    text?: string;
  };
  onCancel?: () => void;
  onSubmit?: () => void;
  onChange?: (connections: MixConnection[]) => void;
}

export const DataScope = ({ connections, cancelBtnProps, submitBtnProps, onCancel, onSubmit, onChange }: Props) => {
  const [connection, setConnection] = useState<MixConnection>();

  const error = useMemo(
    () => (!connections.every((cs) => cs.scope.length) ? 'No Data Scope is Selected' : ''),
    [connections],
  );

  const handleCancel = () => setConnection(undefined);

  const handleSubmit = (connection: MixConnection, scope: MixConnection['scope'], origin: MixConnection['origin']) => {
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
        submitBtnProps={{
          ...submitBtnProps,
        }}
        onCancel={onCancel}
        onSubmit={(scope: MixConnection['scope'], origin: MixConnection['origin']) => {
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
            render: (scope: MixConnection['origin']) =>
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
        <Button
          intent={Intent.PRIMARY}
          disabled={!!error}
          icon={
            error ? (
              <Tooltip2 defaultIsOpen placement={Position.TOP} content={error}>
                <Icon icon="warning-sign" color={Colors.ORANGE5} style={{ margin: 0 }} />
              </Tooltip2>
            ) : null
          }
          text="Save"
          onClick={onSubmit}
          {...submitBtnProps}
        />
      </S.Btns>
      {connection && (
        <Dialog isOpen title="Set Data Scope" footer={null} style={{ width: 820 }} onCancel={handleCancel}>
          <S.DialogTitle>
            <img src={connection.icon} alt="" />
            <span>{connection.name}</span>
          </S.DialogTitle>
          <DataScopeForm
            plugin={connection.plugin}
            connectionId={connection.connectionId}
            onCancel={handleCancel}
            onSubmit={(scope: MixConnection['scope'], origin: MixConnection['origin']) =>
              handleSubmit(connection, scope, origin)
            }
          />
        </Dialog>
      )}
    </S.Wrapper>
  );
};
