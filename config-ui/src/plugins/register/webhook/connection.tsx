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

import { useAppSelector } from '@/app/hook';
import { Buttons, Table, ColumnType, ExternalLink, IconButton } from '@/components';
import { selectConnections } from '@/features/connections';
import { DOC_URL } from '@/release';

import { CreateDialog, ViewDialog, EditDialog, DeleteDialog } from './components';
import * as T from './types';

import * as S from './styled';

type Type = 'add' | 'edit' | 'show' | 'delete';

interface Props {
  filterIds?: ID[];
  onCreateAfter?: (id: ID) => void;
  onDeleteAfter?: (id: ID) => void;
}

export const WebHookConnection = ({ filterIds, onCreateAfter, onDeleteAfter }: Props) => {
  const [type, setType] = useState<Type>();
  const [currentID, setCurrentID] = useState<ID>();

  const connections = useAppSelector((state) => selectConnections(state, 'webhook'));

  const handleHideDialog = () => {
    setType(undefined);
    setCurrentID(undefined);
  };

  const handleShowDialog = (t: Type, r?: T.WebhookItemType) => {
    setType(t);
    setCurrentID(r?.id);
  };

  const columns: ColumnType<T.WebhookItemType> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: 'Webhook Name',
      dataIndex: 'name',
      key: 'name',
      render: (name, row) => <span onClick={() => handleShowDialog('show', row)}>{name}</span>,
    },
    {
      title: '',
      dataIndex: '',
      key: 'action',
      align: 'center',
      render: (_, row) => (
        <S.Action>
          <IconButton icon="array" tooltip="View" onClick={() => handleShowDialog('show', row)} />
          <IconButton icon="annotation" tooltip="Edit" onClick={() => handleShowDialog('edit', row)} />
          <IconButton icon="trash" tooltip="Delete" onClick={() => handleShowDialog('delete', row)} />
        </S.Action>
      ),
    },
  ];

  return (
    <S.Wrapper>
      <Buttons position="top">
        <Button icon="plus" text="Add a Webhook" intent={Intent.PRIMARY} onClick={() => handleShowDialog('add')} />
      </Buttons>
      <Table
        columns={columns}
        dataSource={connections.filter((cs) => (filterIds ? filterIds.includes(cs.id) : true))}
        noData={{
          text: (
            <>
              There is no Webhook yet. Please add a new Webhook.{' '}
              <ExternalLink link={DOC_URL.PLUGIN.WEBHOOK.BASIS}>Learn more</ExternalLink>
            </>
          ),
          btnText: 'Add a Webhook',
          onCreate: () => handleShowDialog('add'),
        }}
      />
      {type === 'add' && (
        <CreateDialog isOpen onCancel={handleHideDialog} onSubmitAfter={(id) => onCreateAfter?.(id)} />
      )}
      {type === 'show' && currentID && <ViewDialog initialId={currentID} onCancel={handleHideDialog} />}
      {type === 'edit' && currentID && <EditDialog initialId={currentID} onCancel={handleHideDialog} />}
      {type === 'delete' && currentID && (
        <DeleteDialog initialId={currentID} onCancel={handleHideDialog} onSubmitAfter={(id) => onDeleteAfter?.(id)} />
      )}
    </S.Wrapper>
  );
};
