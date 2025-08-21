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
import { PlusOutlined } from '@ant-design/icons';
import { Flex, Table, Modal, Input, Select, Button, Tag } from 'antd';
import dayjs from 'dayjs';

import API from '@/api';
import { PageHeader, Block, ExternalLink, CopyText, Message } from '@/components';
import { PATHS } from '@/config';
import { useRefreshData } from '@/hooks';
import { operator, formatTime } from '@/utils';

import * as C from './constant';
import * as S from './styled';

export const ApiKeys = () => {
  const [version, setVersion] = useState(1);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [operating, setOperating] = useState(false);
  const [modal, setModal] = useState<'create' | 'show' | 'delete'>();
  const [currentId, setCurrentId] = useState<string>();
  const [currentKey, setCurrentKey] = useState<string>('');
  const [form, setForm] = useState<{
    name: string;
    expiredAt?: string;
    allowedPath: string;
  }>({
    name: '',
    expiredAt: C.timeOptions[1].value,
    allowedPath: '.*',
  });

  const { data, ready } = useRefreshData(() => API.apiKey.list({ page, pageSize }), [version, page, pageSize]);

  const prefix = useMemo(() => `${window.location.origin}/api/rest/`, []);
  const [dataSource, total] = useMemo(() => [data?.apikeys ?? [], data?.count ?? 0], [data]);
  const hasError = useMemo(() => !form.name || !form.allowedPath, [form]);

  const timeSelectedValue = useMemo(() => {
    return C.timeOptions.find((it) => it.value === form.expiredAt || !it.value)?.value;
  }, [form.expiredAt]);

  const handleCancel = () => {
    setModal(undefined);
  };

  const handleSubmit = async () => {
    const [success, res] = await operator(() => API.apiKey.create(form), {
      setOperating,
    });

    if (success) {
      setVersion(version + 1);
      setModal('show');
      setCurrentKey(res.apiKey);
      setForm({
        name: '',
        expiredAt: C.timeOptions[1].value,
        allowedPath: '.*',
      });
    }
  };

  const handleRevoke = async () => {
    if (!currentId) return;

    const [success] = await operator(() => API.apiKey.remove(currentId));

    if (success) {
      setVersion(version + 1);
      setCurrentId(undefined);
      handleCancel();
    }
  };

  return (
    <PageHeader
      breadcrumbs={[{ name: 'API Keys', path: PATHS.APIKEYS() }]}
      description="You can generate and manage your API keys to access the DevLake API."
    >
      <Flex style={{ marginBottom: 16 }} justify="flex-end">
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModal('create')}>
          New API Key
        </Button>
      </Flex>
      <Table
        rowKey="id"
        size="middle"
        loading={!ready}
        columns={[
          {
            title: 'Key Name',
            dataIndex: 'name',
            key: 'name',
            width: 300,
          },
          {
            title: 'Expiration',
            dataIndex: 'expiredAt',
            key: 'expiredAt',
            width: 200,
            render: (val) => (
              <div>
                <span>{val ? formatTime(val, 'YYYY-MM-DD') : 'No expiration'}</span>
                {dayjs().isAfter(dayjs(val)) && <Tag style={{ marginLeft: 8 }}>Expired</Tag>}
              </div>
            ),
          },
          {
            title: 'Allowed Path',
            dataIndex: 'allowedPath',
            key: 'allowedPath',
            render: (val) => `${prefix}${val}`,
          },
          {
            title: '',
            dataIndex: 'id',
            key: 'id',
            width: 100,
            render: (id) => (
              <Button
                size="small"
                type="primary"
                danger
                onClick={() => {
                  setCurrentId(id);
                  setModal('delete');
                }}
              >
                Revoke
              </Button>
            ),
          },
        ]}
        dataSource={dataSource}
        pagination={{
          current: page,
          pageSize,
          total,
          onChange: ((newPage: number, newPageSize: number) => {
            setPage(newPage);
            if (newPageSize !== pageSize) {
              setPageSize(newPageSize);
            }
          }) as (newPage: number) => void,
        }}
      />
      {modal === 'create' && (
        <Modal
          open
          width={820}
          centered
          title="Generate a New API Key"
          okText="Generate"
          okButtonProps={{
            disabled: hasError,
            loading: operating,
          }}
          onCancel={handleCancel}
          onOk={handleSubmit}
        >
          <Block title="API Key Name" description="Give your API key a unique name to identify in the future." required>
            <Input
              style={{ width: 386 }}
              placeholder="My API Key"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
            />
          </Block>
          <Block title="Expiration" description="Set an expiration time for your API key." required>
            <Select
              style={{ width: 386 }}
              options={C.timeOptions}
              value={timeSelectedValue}
              onChange={(value) => setForm({ ...form, expiredAt: value ? value : undefined })}
            />
          </Block>
          <Block
            title="Allowed Path"
            description={
              <p>
                Enter a Regular Expression that matches the API URL(s) from the{' '}
                <ExternalLink link="/api/swagger/index.html">DevLake API docs</ExternalLink>. The default Regular
                Expression is set to all APIs.
              </p>
            }
            required
          >
            <S.InputContainer>
              <span>{prefix}</span>
              <Input
                placeholder=""
                value={form.allowedPath}
                onChange={(e) => setForm({ ...form, allowedPath: e.target.value })}
              />
            </S.InputContainer>
          </Block>
        </Modal>
      )}
      {modal === 'show' && (
        <Modal open width={820} centered title="Your API key has been generated!" footer={null} onCancel={handleCancel}>
          <div style={{ marginBottom: 16 }}>
            Please make sure to copy your API key now. You will not be able to see it again.
          </div>
          <CopyText content={currentKey} />
        </Modal>
      )}
      {modal === 'delete' && (
        <Modal
          open
          width={820}
          centered
          title="Are you sure you want to revoke this API key?"
          okText="Confirm"
          okButtonProps={{
            loading: operating,
          }}
          onCancel={handleCancel}
          onOk={handleRevoke}
        >
          <Message content="Any applications or scripts using this API key will no longer be able to access the DevLake API. You cannot undo this action." />
        </Modal>
      )}
    </PageHeader>
  );
};
