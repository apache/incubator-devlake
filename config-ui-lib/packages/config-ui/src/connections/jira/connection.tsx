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
import { useMemo, useState } from 'react';
import type { TableColumnType } from 'antd';
import { Button, Modal, Popconfirm } from 'antd';

import { useRefreshData } from '../../hooks';
import { operate } from '../../utils/operate';

import { JiraConnectionForm } from './form';
import { JiraConnectionList } from './list';
import type { JiraPayloadType } from './typed';
import * as API from './api';
import * as S from './styled';

export const JiraConnection = () => {
  const [version, setVersion] = useState(0);
  const [visible, setVisible] = useState(false);
  const [modalType, setModalType] = useState<'add' | 'edit'>();
  const [updateObj, setUpdateObj] = useState<any>();

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

  const handleDelete = async (row: any) => {
    const [success] = await operate(() => API.remove(row.id));
    if (success) {
      setVersion((v) => v + 1);
    }
  };

  const handleUpdate = (row: any) => {
    setUpdateObj(row);
    setModalType('edit');
    handleShowModal();
  };

  const handleSubmit = async (values: JiraPayloadType) => {
    const request =
      modalType === 'add'
        ? () => API.create(values)
        : () => API.update(updateObj?.id, values);
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
      <Button type="primary" onClick={handleCreate}>
        Add Connection
      </Button>
      <JiraConnectionList
        style={{ marginTop: 12 }}
        extraColumn={extraColumn}
        loading={!ready}
        data={data}
      />
      <Modal
        visible={visible}
        title={title}
        width={600}
        footer={null}
        onCancel={handleHideModal}
      >
        <JiraConnectionForm initialValues={updateObj} onSubmit={handleSubmit} />
      </Modal>
    </S.PageContainer>
  );
};

JiraConnection.Form = JiraConnectionForm;
JiraConnection.List = JiraConnectionList;
