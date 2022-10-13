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

import { useRefreshData } from '../../hooks';
import { operate } from '../../utils/operate';

import { GiteeConnectionList } from './list';
import { GiteeConnectionForm } from './form';
import type { GiteeItemType, GiteePayloadType } from './typed';
import * as API from './api';
import * as S from './styled';

export const GiteeConnection = () => {
  const [version, setVersion] = useState(0);
  const [visible, setVisible] = useState(false);
  const [modalType, setModalType] = useState<'add' | 'edit'>();
  const [updateObj, setUpdateObj] = useState<any>();
  const [operating, setOperating] = useState(false);

  const { ready, data } = useRefreshData(API.getList, [version]);

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

  const handleDelete = async (row: GiteeItemType) => {
    const [success] = await operate(() => API.remove(row.id));
    if (success) {
      setVersion((v) => v + 1);
    }
  };

  const handleUpdate = (row: GiteeItemType) => {
    setUpdateObj(row);
    setModalType('edit');
    handleShowModal();
  };

  const handleTest = async (values: GiteePayloadType) => {
    setOperating(true);
    try {
      await API.test(values);
    } catch (err) {
      message.error((err as any).response?.data.message);
    } finally {
      setOperating(false);
    }
  };

  const handleSubmit = async (values: GiteePayloadType) => {
    const request =
      modalType === 'add'
        ? () => API.create(values)
        : () => API.update(updateObj?.id, values);
    const [success] = await operate(request, { setOperating });
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
      <Button type="primary" onClick={handleCreate}>
        Add Connection
      </Button>
      <GiteeConnectionList
        style={{ marginTop: 12 }}
        extraColumn={extraColumn}
        loading={!ready}
        data={data}
      />
      <Modal
        visible={visible}
        title={title}
        footer={null}
        onCancel={handleHideModal}
      >
        <GiteeConnectionForm
          initialValues={updateObj}
          operating={operating}
          onTest={handleTest}
          onSubmit={handleSubmit}
        />
      </Modal>
    </S.PageContainer>
  );
};

GiteeConnection.Form = GiteeConnectionForm;
GiteeConnection.List = GiteeConnectionList;
