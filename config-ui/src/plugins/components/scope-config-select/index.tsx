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
import { Button, Intent } from '@blueprintjs/core';

import { Buttons, Table, IconButton, Dialog } from '@/components';
import { useRefreshData } from '@/hooks';

import { ScopeConfigForm } from '../scope-config-form';

import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  connectionId: ID;
  onCancel?: () => void;
  onSubmit?: (trId: string) => void;
}

export const ScopeConfigSelect = ({ plugin, connectionId, onCancel, onSubmit }: Props) => {
  const [version, setVersion] = useState(1);
  const [trId, setTrId] = useState<string>();
  const [isOpen, setIsOpen] = useState(false);
  const [updatedId, setUpdatedId] = useState<ID>();

  const { ready, data } = useRefreshData(() => API.getScopeConfigs(plugin, connectionId), [version]);

  const dataSource = useMemo(() => (data ? data : []), [data]);

  const handleShowDialog = () => {
    setIsOpen(true);
  };

  const handleHideDialog = () => {
    setIsOpen(false);
    setUpdatedId(undefined);
  };

  const handleUpdate = async (id: ID) => {
    setUpdatedId(id);
    handleShowDialog();
  };

  const handleSubmit = async () => {
    handleHideDialog();
    setVersion((v) => v + 1);
  };

  return (
    <S.Wrapper>
      <Buttons position="top" align="left">
        <Button icon="add" intent={Intent.PRIMARY} text="Add New Scope Config" onClick={handleShowDialog} />
      </Buttons>
      <Table
        loading={!ready}
        columns={[
          { title: 'Name', dataIndex: 'name', key: 'name' },
          {
            title: '',
            dataIndex: 'id',
            key: 'id',
            width: 100,
            render: (id) => <IconButton icon="annotation" tooltip="Edit" onClick={() => handleUpdate(id)} />,
          },
        ]}
        dataSource={dataSource}
        rowSelection={{
          rowKey: 'id',
          type: 'radio',
          selectedRowKeys: trId ? [`${trId}`] : [],
          onChange: (selectedRowKeys) => setTrId(`${selectedRowKeys[0]}`),
        }}
        noShadow
      />
      <Buttons>
        <Button outlined intent={Intent.PRIMARY} text="Cancel" onClick={onCancel} />
        <Button disabled={!trId} intent={Intent.PRIMARY} text="Save" onClick={() => trId && onSubmit?.(trId)} />
      </Buttons>
      <Dialog
        style={{ width: 820 }}
        footer={null}
        isOpen={isOpen}
        title={!updatedId ? 'Add Scope Config' : 'Edit Scope Config'}
        onCancel={handleHideDialog}
      >
        <ScopeConfigForm
          plugin={plugin}
          connectionId={connectionId}
          showWarning={!!updatedId}
          scopeConfigId={updatedId}
          onCancel={onCancel}
          onSubmit={handleSubmit}
        />
      </Dialog>
    </S.Wrapper>
  );
};
