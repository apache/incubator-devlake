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

import { WebhookConnectionForm } from './form';
import { WebhookConnectionList } from './list';
import { WebhookPayloadType, WebhookFormEnmu, WebhookItemType } from './typed';
import * as API from './api';

import * as S from './styled';

export const WebhookConnection = () => {
  const [version, setVersion] = useState(0);
  const [visible, setVisible] = useState(false);
  const [formType, setFormType] = useState<WebhookFormEnmu>();
  const [initialValues, setInitialValues] = useState<any>();

  const { ready, data } = useRefreshData(API.getList, [version]);

  const handleShowModal = () => {
    setVisible(true);
  };

  const handleHideModal = () => {
    setVisible(false);
    setFormType(undefined);
  };

  const handleCreate = () => {
    setFormType(WebhookFormEnmu.add);
    setInitialValues(null);
    handleShowModal();
  };

  const handleGenerateURL = async (values: WebhookPayloadType) => {
    try {
      const { id } = await API.create(values);
      const res = await API.get(id);
      setFormType(WebhookFormEnmu.show);
      setInitialValues(res);
    } catch (err) {}
  };

  const handleDelete = async (row: WebhookItemType) => {
    const [success] = await operate(() => API.remove(row.id));
    if (success) {
      setVersion((v) => v + 1);
    }
  };

  const handleUpdate = (row: WebhookItemType) => {
    setFormType(WebhookFormEnmu.edit);
    setInitialValues(row);
    handleShowModal();
  };

  const handleSubmit = async (values: WebhookPayloadType) => {
    const [success] = await operate(() =>
      API.update(initialValues?.id, values),
    );
    if (success) {
      setVersion((v) => v + 1);
      handleHideModal();
    }
  };

  const title = useMemo(() => {
    switch (formType) {
      case WebhookFormEnmu.add:
      case WebhookFormEnmu.show:
        return 'Add Webhook';
      case WebhookFormEnmu.edit:
        return 'Edit Webhook';
    }
  }, [formType]);

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
      <WebhookConnectionList
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
        <WebhookConnectionForm
          formType={formType}
          initialValues={initialValues}
          onGenerateURL={handleGenerateURL}
          onDone={handleHideModal}
          onSubmit={handleSubmit}
        />
      </Modal>
    </S.PageContainer>
  );
};

WebhookConnection.Form = WebhookConnectionForm;
WebhookConnection.List = WebhookConnectionList;
