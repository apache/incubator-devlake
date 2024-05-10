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
import { LinkOutlined, EditOutlined } from '@ant-design/icons';
import { theme, Button, Modal, Flex, Space } from 'antd';
import styled from 'styled-components';

import API from '@/api';
import { IconButton, Message } from '@/components';
import { operator } from '@/utils';

import { PluginName } from '../plugin-name';
import { ScopeConfigSelect } from '../scope-config-select';
import { ScopeConfigForm } from '../scope-config-form';

const Wrapper = styled.div``;

interface Props {
  plugin: string;
  connectionId: ID;
  scopeId: ID;
  scopeName: string;
  id?: ID;
  name?: string;
  onSuccess?: (id?: ID, hideToast?: boolean) => void;
}

export const ScopeConfig = ({ plugin, connectionId, scopeId, scopeName, id, name, onSuccess }: Props) => {
  const [type, setType] = useState<'associate' | 'update' | 'relatedProjects' | 'duplicate'>();
  const [relatedProjects, setRelatedProjects] = useState<
    Array<{ name: string; scopes: Array<{ scopeId: ID; scopeName: string }> }>
  >([]);

  const {
    token: { colorPrimary },
  } = theme.useToken();

  const handleHideDialog = () => setType(undefined);

  const handleCheckScopeConfig = async () => {
    if (!id) return;

    const [success, res] = await operator(() => API.scopeConfig.check(plugin, id), { hideToast: true });

    if (success) {
      const projects = res.projects.map((it: any) => ({
        name: it.name,
        scopes: it.scopes,
      }));

      if (projects.length !== 1) {
        setRelatedProjects(projects);
        setType('relatedProjects');
      } else {
        setType('update');
      }
    }
  };

  const handleAssociate = async (trId: ID) => {
    const [success] = await operator(
      () => API.scope.update(plugin, connectionId, scopeId, { scopeConfigId: trId !== 'None' ? +trId : null }),
      {
        hideToast: true,
      },
    );

    if (success) {
      handleHideDialog();
      onSuccess?.(id, type === 'duplicate');
    }
  };

  const handleUpdate = (trId: ID) => {
    handleHideDialog();
    onSuccess?.(id);
  };

  return (
    <Wrapper>
      <span>{id ? name : 'N/A'}</span>
      <IconButton
        icon={<LinkOutlined />}
        helptip="Associate Scope Config"
        size="small"
        type="link"
        onClick={() => {
          setType('associate');
        }}
      />
      {id && (
        <IconButton
          icon={<EditOutlined />}
          helptip=" Edit Scope Config"
          type="link"
          size="small"
          onClick={handleCheckScopeConfig}
        />
      )}
      {type === 'associate' && (
        <Modal
          open
          width={960}
          centered
          footer={null}
          title={<PluginName plugin={plugin} name="Associate Scope Config" />}
          onCancel={handleHideDialog}
        >
          {plugin === 'tapd' ? (
            <ScopeConfigForm
              plugin={plugin}
              connectionId={connectionId}
              scopeConfigId={id}
              scopeId={scopeId}
              onCancel={handleHideDialog}
              onSubmit={handleAssociate}
            />
          ) : (
            <ScopeConfigSelect
              plugin={plugin}
              connectionId={connectionId}
              scopeConfigId={id}
              onCancel={handleHideDialog}
              onSubmit={handleAssociate}
            />
          )}
        </Modal>
      )}
      {type === 'update' && (
        <Modal open width={960} centered footer={null} title="Edit Scope Config" onCancel={handleHideDialog}>
          <ScopeConfigForm
            plugin={plugin}
            connectionId={connectionId}
            showWarning
            scopeConfigId={id}
            scopeId={scopeId}
            onCancel={handleHideDialog}
            onSubmit={handleUpdate}
          />
        </Modal>
      )}
      {type === 'duplicate' && (
        <Modal open width={960} centered footer={null} title="Edit Scope Config" onCancel={handleHideDialog}>
          <ScopeConfigForm
            plugin={plugin}
            connectionId={connectionId}
            showWarning
            forceCreate
            scopeConfigId={id}
            scopeId={scopeId}
            onCancel={handleHideDialog}
            onSubmit={handleAssociate}
          />
        </Modal>
      )}
      {type === 'relatedProjects' && (
        <Modal
          open
          width={830}
          centered
          footer={null}
          title={`Edit '${name}' for '${scopeName}'`}
          onCancel={handleHideDialog}
        >
          <Message content="The change will apply to all following projects:" />
          <ul style={{ margin: '15px 0 30px 30px' }}>
            {relatedProjects.map((it) => (
              <li key={it.name} style={{ color: colorPrimary }}>
                {it.name}: {it.scopes.map((sc) => sc.scopeName).join(',')}
              </li>
            ))}
          </ul>
          <Flex justify="end">
            <Space>
              <Button onClick={handleHideDialog}>Cancel</Button>
              <Button type="primary" onClick={() => setType('update')}>
                Continue
              </Button>
              <Button type="primary" onClick={() => setType('duplicate')}>
                Duplicate a scope config for {scopeName}
              </Button>
            </Space>
          </Flex>
        </Modal>
      )}
    </Wrapper>
  );
};
