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

import React, { useState, useEffect, useMemo } from 'react';
import { useHistory } from 'react-router-dom';
import { ButtonGroup, Button, Intent, Position } from '@blueprintjs/core';
import { Popover2 } from '@blueprintjs/popover2';

import { Table, ColumnType, IconButton } from '@/components';
import type { ConnectionItemType } from '@/store';
import { useConnection, ConnectionStatus } from '@/store';
import { operator } from '@/utils';

import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
}

export const Connection = ({ plugin }: Props) => {
  const [operating, setOperating] = useState(false);

  const history = useHistory();

  const { connections, onTest, onRefresh } = useConnection();

  useEffect(() => {
    connections.map((cs) => onTest(cs));
  }, []);

  const handleRefresh = () => onRefresh();

  const handleCreate = () => history.push(`/connections/${plugin}/create`);

  const handleUpdate = (id: ID) => history.push(`/connections/${plugin}/${id}`);

  const handleDelete = async (id: ID) => {
    const [success] = await operator(() => API.deleteConnection(plugin, id), {
      setOperating,
    });

    if (success) {
      onRefresh();
    }
  };

  const columns = useMemo(
    () =>
      [
        {
          title: 'ID',
          dataIndex: 'id',
          key: 'id',
          width: 100,
        },
        {
          title: 'Connection Name',
          dataIndex: 'name',
          key: 'name',
        },
        {
          title: 'Endpoint',
          dataIndex: 'endpoint',
          key: 'endpoint',
          ellipsis: true,
        },
        {
          title: 'Status',
          dataIndex: 'status',
          key: 'status',
          align: 'center',
          render: (_, row) => <ConnectionStatus connection={row} onTest={onTest} />,
        },
        {
          title: '',
          dataIndex: 'id',
          key: 'action',
          width: 100,
          align: 'center',
          render: (id) => (
            <ButtonGroup>
              <IconButton icon="edit" tooltip="Edit" onClick={() => handleUpdate(id)} />
              <Popover2
                position={Position.TOP}
                content={
                  <S.DeleteConfirm>
                    <h3>Confirm deletion</h3>
                    <p>Are you sure you want to delete this item?</p>
                    <ButtonGroup>
                      <Button
                        loading={operating}
                        intent={Intent.DANGER}
                        text="Delete"
                        onClick={() => handleDelete(id)}
                      />
                    </ButtonGroup>
                  </S.DeleteConfirm>
                }
              >
                <IconButton icon="delete" tooltip="Delete" />
              </Popover2>
            </ButtonGroup>
          ),
        },
      ] as ColumnType<ConnectionItemType>,
    [],
  );

  return (
    <S.Wrapper>
      <ButtonGroup className="action">
        <Button intent={Intent.PRIMARY} icon="plus" text="New Connection" onClick={handleCreate} />
        <Button icon="refresh" text="Refresh Connections" onClick={handleRefresh} />
      </ButtonGroup>
      <Table
        columns={columns}
        dataSource={connections}
        noData={{
          text: 'There is no data connection yet. Please add a new connection.',
          btnText: 'New Connection',
          onCreate: handleCreate,
        }}
      />
    </S.Wrapper>
  );
};
