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
import { Button, Modal } from 'antd';
import styled from 'styled-components';

import API from '@/api';
import { operator } from '@/utils';

import { PluginName } from '../plugin-name';
import { ScopeConfigSelect } from '../scope-config-select';
import { ScopeConfigForm } from '../scope-config-form';

const Wrapper = styled.div``;

interface Props {
  plugin: string;
  connectionId: ID;
  scopeId: ID;
  id?: ID;
  name?: string;
  onSuccess?: () => void;
}

export const ScopeConfig = ({ plugin, connectionId, scopeId, id, name, onSuccess }: Props) => {
  const [type, setType] = useState<'associate' | 'update'>();

  const handleHideDialog = () => setType(undefined);

  const handleAssociate = async (trId: ID) => {
    const [success] = await operator(
      () => API.scope.update(plugin, connectionId, scopeId, { scopeConfigId: trId !== 'None' ? +trId : null }),
      {
        hideToast: true,
      },
    );

    if (success) {
      handleHideDialog();
      onSuccess?.();
    }
  };

  const handleUpdate = (trId: ID) => {
    handleHideDialog();
    onSuccess?.();
  };

  return (
    <Wrapper>
      <span>{id ? name : 'N/A'}</span>
      <Button
        size="small"
        type="link"
        icon={<LinkOutlined />}
        onClick={() => {
          setType('associate');
        }}
      />
      {id && (
        <Button
          size="small"
          type="link"
          icon={<EditOutlined />}
          onClick={() => {
            // TO-DO: check if the scope config is associated with any scope
            setType('update');
          }}
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
    </Wrapper>
  );
};
