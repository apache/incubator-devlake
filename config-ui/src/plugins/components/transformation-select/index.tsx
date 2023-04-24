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

import { useState, useEffect, useMemo } from 'react';
import { Button, Intent } from '@blueprintjs/core';

import { Dialog, PageLoading, Table, IconButton } from '@/components';
import { useRefreshData } from '@/hooks';
import { TransformationForm } from '@/plugins';

import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  connectionId: ID;
  scopeId: ID;
  transformationId?: ID;
  transformationType: MixConnection['transformationType'];
  onCancel: () => void;
  onSubmit: (tid: ID) => void;
}

export const TransformationSelect = ({
  plugin,
  connectionId,
  scopeId,
  transformationId,
  transformationType,
  onCancel,
  onSubmit,
}: Props) => {
  const [step, setStep] = useState(1);
  const [type, setType] = useState<'add' | 'edit'>('add');
  const [selectedId, setSelectedId] = useState<ID>();
  const [updatedId, setUpdatedId] = useState<ID>();

  useEffect(() => setSelectedId(transformationId), [transformationId]);

  useEffect(() => {
    setStep(transformationType === 'for-scope' ? 2 : 1);
    setType(transformationType === 'for-scope' && transformationId ? 'edit' : 'add');
    setUpdatedId(transformationType === 'for-scope' && transformationId ? transformationId : undefined);
  }, [transformationId, transformationType]);

  const { ready, data } = useRefreshData(() => API.getTransformations(plugin, connectionId), [step]);

  const title = useMemo(() => {
    switch (true) {
      case step === 1:
        return 'Associate Transformation';
      case type === 'add':
        return 'Add New Transformation';
      case type === 'edit':
        return 'Edit Transformation';
    }
  }, [step, type]);

  const handleNewTransformation = () => {
    setStep(2);
    setType('add');
  };

  const handleEditTransformation = (id: ID) => {
    setStep(2);
    setType('edit');
    setUpdatedId(id);
  };

  const handleReset = (tr?: any) => {
    if (transformationType === 'for-scope') {
      return tr ? onSubmit(tr.id) : onCancel();
    }
    setStep(1);
    setUpdatedId('');
  };

  const handleSubmit = () => !!selectedId && onSubmit(selectedId);

  return (
    <Dialog isOpen title={title} footer={null} style={{ width: 960 }} onCancel={onCancel}>
      {!ready || !data ? (
        <PageLoading />
      ) : step === 1 ? (
        <S.Wrapper>
          <S.Aciton>
            <Button intent={Intent.PRIMARY} icon="add" onClick={handleNewTransformation}>
              Add New Transformation
            </Button>
          </S.Aciton>
          <Table
            columns={[
              { title: 'Transformation', dataIndex: 'name', key: 'name' },
              {
                title: '',
                dataIndex: 'id',
                key: 'id',
                width: 100,
                align: 'right',
                render: (id) => (
                  <IconButton icon="annotation" tooltip="Edit" onClick={() => handleEditTransformation(id)} />
                ),
              },
            ]}
            dataSource={data}
            rowSelection={{
              rowKey: 'id',
              type: 'radio',
              selectedRowKeys: selectedId ? [selectedId] : [],
              onChange: (selectedRowKeys) => setSelectedId(selectedRowKeys[0]),
            }}
            noShadow
          />
          <S.Btns>
            <Button outlined intent={Intent.PRIMARY} text="Cancel" onClick={onCancel} />
            <Button disabled={!selectedId} intent={Intent.PRIMARY} text="Save" onClick={handleSubmit} />
          </S.Btns>
        </S.Wrapper>
      ) : (
        <TransformationForm
          plugin={plugin}
          connectionId={connectionId}
          scopeId={scopeId}
          id={updatedId}
          onCancel={handleReset}
        />
      )}
    </Dialog>
  );
};
