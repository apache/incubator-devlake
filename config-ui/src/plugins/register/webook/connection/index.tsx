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
import { ButtonGroup, Button, Intent } from '@blueprintjs/core';

import { Table, ColumnType, ExternalLink, IconButton } from '@/components';
import { useConnections } from '@/hooks';
import { DOC_URL } from '@/release';

import type { WebhookItemType } from '../types';
import { WebhookCreateDialog } from '../create-dialog';
import { WebhookDeleteDialog } from '../delete-dialog';
import { WebhookViewOrEditDialog } from '../view-or-edit-dialog';

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

  const { connections, onRefresh } = useConnections({ plugin: 'webhook' });

  const handleHideDialog = () => {
    setType(undefined);
    setCurrentID(undefined);
  };

  const handleShowDialog = (t: Type, r?: WebhookItemType) => {
    setType(t);
    setCurrentID(r?.id);
  };

  const columns: ColumnType<WebhookItemType> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: 'Incoming Webhook Name',
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
          <IconButton icon="edit" tooltip="Edit" onClick={() => handleShowDialog('edit', row)} />
          <IconButton icon="trash" tooltip="Delete" onClick={() => handleShowDialog('delete', row)} />
        </S.Action>
      ),
    },
  ];

  return (
    <S.Wrapper>
      <ButtonGroup>
        <Button icon="plus" text="Add a Webhook" intent={Intent.PRIMARY} onClick={() => handleShowDialog('add')} />
      </ButtonGroup>
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
        <WebhookCreateDialog
          isOpen
          onCancel={handleHideDialog}
          onSubmitAfter={(id) => {
            onRefresh();
            onCreateAfter?.(id);
          }}
        />
      )}
      {type === 'delete' && currentID && (
        <WebhookDeleteDialog
          isOpen
          initialID={currentID}
          onCancel={handleHideDialog}
          onSubmitAfter={(id) => {
            onRefresh();
            onDeleteAfter?.(id);
          }}
        />
      )}
      {(type === 'edit' || type === 'show') && currentID && (
        <WebhookViewOrEditDialog
          type={type}
          isOpen
          initialID={currentID}
          onCancel={handleHideDialog}
          onSubmitAfter={() => onRefresh()}
        />
      )}
    </S.Wrapper>
  );
};
