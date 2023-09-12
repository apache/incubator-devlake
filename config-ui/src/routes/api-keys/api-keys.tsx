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
import { Button, Intent, InputGroup } from '@blueprintjs/core';
import { CopyToClipboard } from 'react-copy-to-clipboard';

import {
  PageHeader,
  Table,
  Dialog,
  FormItem,
  Selector,
  ExternalLink,
  TextTooltip,
  IconButton,
  toast,
  Buttons,
  Message,
} from '@/components';
import { useRefreshData } from '@/hooks';
import { operator, formatTime } from '@/utils';

import * as API from './api';
import * as C from './constant';
import * as S from './styled';

export const ApiKeys = () => {
  const [version, setVersion] = useState(1);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(20);
  const [operating, setOperating] = useState(false);
  const [modal, setModal] = useState<'create' | 'show' | 'delete'>();
  const [currentId, setCurrentId] = useState<string>();
  const [currentKey, setCurrentKey] = useState<string>();
  const [form, setForm] = useState<{
    name: string;
    expiredAt?: string;
    allowedPath: string;
  }>({
    name: '',
    expiredAt: C.timeOptions[1].value,
    allowedPath: '.*',
  });

  const { data, ready } = useRefreshData(() => API.getApiKeys({ page, pageSize }), [version, page, pageSize]);

  const [dataSource, total] = useMemo(() => [data?.apikeys ?? [], data?.count ?? 0], [data]);

  const timeSelectedItem = useMemo(() => {
    return C.timeOptions.find((it) => it.value === form.expiredAt || !it.value);
  }, [form.expiredAt]);

  const handleCancel = () => {
    setModal(undefined);
  };

  const handleSubmit = async () => {
    const [success, res] = await operator(() => API.createApiKey(form), {
      setOperating,
      hideToast: true,
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

    const [success] = await operator(() => API.deleteApiKey(currentId));

    if (success) {
      setVersion(version + 1);
      setCurrentId(undefined);
      handleCancel();
    }
  };

  return (
    <PageHeader
      breadcrumbs={[{ name: 'API Keys', path: '/api-keys' }]}
      extra={<Button intent={Intent.PRIMARY} icon="plus" text="New API Key" onClick={() => setModal('create')} />}
    >
      <p>You can generate and manage your API keys to access the DevLake API.</p>
      <Table
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
            render: (val) => (val ? formatTime(val, 'YYYY-MM-DD') : 'No expiration'),
          },
          {
            title: 'Allowed Path',
            dataIndex: 'allowedPath',
            key: 'allowedPath',
            render: (val) => `http://localhost:4000/api/rest${val}`,
          },
          {
            title: '',
            dataIndex: 'id',
            key: 'id',
            width: 100,
            render: (id) => (
              <Button
                small
                intent={Intent.DANGER}
                text="Revoke"
                onClick={() => {
                  setCurrentId(id);
                  setModal('delete');
                }}
              />
            ),
          },
        ]}
        dataSource={dataSource}
        pagination={{
          page,
          pageSize,
          total,
          onChange: setPage,
        }}
        noData={{
          text: 'There is no API key yet.',
        }}
      />
      {modal === 'create' && (
        <Dialog
          style={{ width: 820 }}
          isOpen
          title="Generate a New API Key"
          okLoading={operating}
          okText="Generate"
          onCancel={handleCancel}
          onOk={handleSubmit}
        >
          <FormItem label="API Key Name" subLabel="Give your API key a unique name to identify in the future." required>
            <InputGroup
              style={{ width: 386 }}
              placeholder="My API Key"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
            />
          </FormItem>
          <FormItem label="Expiration" subLabel="Set an expiration time for your API key." required>
            <div style={{ width: 386 }}>
              <Selector
                items={C.timeOptions}
                getKey={(it) => it.value}
                getName={(it) => it.label}
                selectedItem={timeSelectedItem}
                onChangeItem={(it) => setForm({ ...form, expiredAt: it.value ? it.value : undefined })}
              />
            </div>
          </FormItem>
          <FormItem
            label="Allowed Path"
            subLabel={
              <p>
                Enter a Regular Expression that matches the API URL(s) from the{' '}
                <ExternalLink link="">DevLake API docs</ExternalLink>. The default Regular Expression is set to all
                APIs.
              </p>
            }
            required
          >
            <S.InputContainer>
              <span>http://localhost:4000/api/rest</span>
              <InputGroup
                placeholder=""
                value={form.allowedPath}
                onChange={(e) => setForm({ ...form, allowedPath: e.target.value })}
              />
            </S.InputContainer>
          </FormItem>
        </Dialog>
      )}
      {modal === 'show' && (
        <Dialog
          style={{ width: 820 }}
          isOpen
          title="Your API key has been generated!"
          footer={null}
          onCancel={handleCancel}
        >
          <div>Please make sure to copy your API key now. You will not be able to see it again.</div>
          <S.KeyContainer>
            <TextTooltip style={{ width: '96%' }} content="">
              {currentKey}
            </TextTooltip>
            <CopyToClipboard text={currentKey as string} onCopy={() => toast.success('Copy successfully.')}>
              <IconButton icon="clipboard" tooltip="Copy" />
            </CopyToClipboard>
          </S.KeyContainer>
          <Buttons position="bottom" align="right">
            <Button intent={Intent.PRIMARY} text="Confirm" onClick={handleCancel} />
          </Buttons>
        </Dialog>
      )}
      {modal === 'delete' && (
        <Dialog
          style={{ width: 820 }}
          isOpen
          title="Are you sure you want to revoke this API key?"
          okLoading={operating}
          okText="Confirm"
          onCancel={handleCancel}
          onOk={handleRevoke}
        >
          <Message content="Any applications or scripts using this API key will no longer be able to access the DevLake API. You cannot undo this action." />
        </Dialog>
      )}
    </PageHeader>
  );
};
