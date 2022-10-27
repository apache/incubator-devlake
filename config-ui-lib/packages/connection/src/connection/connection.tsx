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
import { message, TableColumnType } from 'antd';
import { Button, Modal, Popconfirm } from 'antd';

import { useRefreshData } from '../hooks';
import { operate } from '../utils/operate';

import { ConnectionList } from './list';
import { ConnectionForm } from './form';
import type { ConnectionType, ItemType, IPaylod } from './typed';
import * as API from './api';
import * as S from './styled';

interface IConnectionProps {
  type: ConnectionType;
}

export const Connection = ({ type }: IConnectionProps) => {
  const [version, setVersion] = useState(0);
  const [visible, setVisible] = useState(false);
  const [modalType, setModalType] = useState<'add' | 'edit'>();
  const [updateObj, setUpdateObj] = useState<any>();

  const { ready, data } = useRefreshData(() => API.getList(type), [version]);

  const handleRefresh = () => {
    setVersion((v) => v + 1);
  };

  const handleShowModal = () => {
    setVisible(true);
  };

  const handleHideModal = () => {
    setVisible(false);
  };

  const handleCreate = () => {
    setModalType('add');
    setUpdateObj(null);
    handleShowModal();
  };

  const handleDelete = async (row: ItemType) => {
    const [success] = await operate(() => API.remove(type, row.id));
    if (success) {
      setVersion((v) => v + 1);
    }
  };

  const handleUpdate = (row: ItemType) => {
    setUpdateObj(row);
    setModalType('edit');
    handleShowModal();
  };

  const handleTest = async (values: IPaylod) => {
    try {
      await API.test(type, values);
    } catch (err) {
      message.error((err as any).response?.data.message);
    } finally {
    }
  };

  const handleSubmit = async (values: IPaylod) => {
    const request =
      modalType === 'add'
        ? () => API.create(type, values)
        : () => API.update(type, updateObj?.id, values);
    const [success] = await operate(request);
    if (success) {
      setVersion((v) => v + 1);
      handleHideModal();
    }
  };

  const title = useMemo(() => {
    switch (modalType) {
      case 'add':
        return 'Add Connection';
      case 'edit':
        return 'Edit Connection';
    }
  }, [modalType]);

  const extraColumn: TableColumnType<any>[] = [
    {
      title: 'Operate',
      key: 'operate',
      width: 200,
      render: (row) => (
        <>
          <Button type="text" onClick={() => handleUpdate(row)}>
            Edit
          </Button>
          <Popconfirm
            title="Are you sure you want to continue? "
            onConfirm={() => handleDelete(row)}
          >
            <Button type="text" danger>
              Delete
            </Button>
          </Popconfirm>
        </>
      ),
    },
  ];

  return (
    <S.PageContainer>
      <S.BtnContainer>
        <Button type="primary" onClick={handleCreate}>
          Add Connection
        </Button>
        <Button onClick={handleRefresh}>Refresh Connection</Button>
      </S.BtnContainer>
      <ConnectionList loading={!ready} data={data} extraColumn={extraColumn} />
      <Modal
        visible={visible}
        title={title}
        footer={null}
        onCancel={handleHideModal}
      >
        <ConnectionForm
          type={type}
          initialValues={updateObj}
          onTest={handleTest}
          onSubmit={handleSubmit}
        />
      </Modal>
    </S.PageContainer>
  );
};
