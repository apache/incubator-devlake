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
import { PATHS } from '@/config';
import { getPluginConfig } from '@/plugins';
import { operator } from '@/utils';

import { PluginName } from '../plugin-name';
import { ScopeConfigSelect } from '../scope-config-select';
import { ScopeConfigForm } from '../scope-config-form';

const Wrapper = styled.div``;

type RelatedProjects = Array<{ name: string; blueprintId: ID; scopes: Array<{ scopeName: string }> }>;

interface Props {
  plugin: string;
  connectionId: ID;
  scopeId: ID;
  scopeName: string;
  scopeConfigId?: ID;
  scopeConfigName?: string;
  projects?: Array<{ name: string; blueprintId: ID }>;
  onSuccess: (id?: ID) => void;
}

export const ScopeConfig = ({
  plugin,
  connectionId,
  scopeId,
  scopeName,
  scopeConfigId,
  scopeConfigName,
  onSuccess,
}: Props) => {
  const [type, setType] = useState<'associate' | 'update' | 'relatedProjects' | 'duplicate'>();
  const [relatedProjects, setRelatedProjects] = useState<RelatedProjects>([]);

  const [operating, setOperating] = useState(false);

  const pluginConfig = getPluginConfig(plugin);

  const {
    token: { colorPrimary },
  } = theme.useToken();

  const [modal, contextHolder] = Modal.useModal();

  const handleHideDialog = () => setType(undefined);

  const handleCheckScopeConfig = async (id: ID) => {
    const [success, res] = await operator(() => API.scopeConfig.check(plugin, id), { hideToast: true });

    if (success) {
      const projects = (res.projects ?? []).map((it: any) => ({
        name: it.name,
        scopes: it.scopes,
      }));

      if (projects.length > 1) {
        setRelatedProjects(projects);
        setType('relatedProjects');
      } else {
        setType('update');
      }
    }
  };

  const handleRun = async (pname: string, blueprintId: ID, data?: { skipCollectors?: boolean; fullSync?: boolean }) => {
    const [success] = await operator(() => API.blueprint.trigger(blueprintId, data), {
      setOperating,
    });

    if (success) {
      window.open(PATHS.PROJECT(pname));
    }
  };

  const handleShowProjectsModal = (projects: RelatedProjects) => {
    if (!projects || !projects.length) {
      onSuccess();
    } else if (projects.length === 1) {
      const [{ name, blueprintId }] = projects;
      modal.success({
        closable: true,
        centered: true,
        width: 550,
        title: 'Scope Config Saved',
        content: 'Please re-transform data to apply the updated scope config.',
        footer: (
          <div style={{ marginTop: 20, textAlign: 'center' }}>
            <Button
              type="primary"
              loading={operating}
              onClick={() => handleRun(name, blueprintId, { skipCollectors: true })}
            >
              Re-transform now
            </Button>
          </div>
        ),
        onCancel: onSuccess,
      });
    } else {
      modal.success({
        closable: true,
        centered: true,
        width: 830,
        title: 'Scope Config Saved',
        content: (
          <>
            <div style={{ marginBottom: 16 }}>
              The listed projects are impacted. Please re-transform the data to apply the updated scope config.
            </div>
            <ul>
              {projects.map(({ name, blueprintId }: { name: string; blueprintId: ID }) => (
                <li key={name} style={{ marginBottom: 10 }}>
                  <Space>
                    <span>{name}</span>
                    <Button
                      size="small"
                      type="link"
                      loading={operating}
                      onClick={() => handleRun(name, blueprintId, { skipCollectors: true })}
                    >
                      Re-transform Data
                    </Button>
                  </Space>
                </li>
              ))}
            </ul>
          </>
        ),
        footer: null,
        onCancel: onSuccess,
      });
    }
  };

  const handleAssociate = async (trId: ID) => {
    const [success, res] = await operator(
      async () => {
        await API.scope.update(plugin, connectionId, scopeId, { scopeConfigId: trId === 'None' ? null : +trId });
        return API.scope.get(plugin, connectionId, scopeId, { blueprints: true });
      },
      {
        hideToast: true,
      },
    );

    if (success) {
      handleHideDialog();
      handleShowProjectsModal(
        (res.blueprints ?? []).map((it: any) => ({
          name: it.projectName,
          blueprintId: it.id,
          scopes: [
            {
              scopeId,
              scopeName,
            },
          ],
        })),
      );
    }
  };

  const handleUpdate = async (trId: ID) => {
    handleHideDialog();

    const [success, res] = await operator(() => API.scopeConfig.check(plugin, trId), { hideToast: true });

    if (success) {
      handleShowProjectsModal(res.projects ?? []);
    }
  };

  return (
    <Wrapper>
      {contextHolder}
      <span>{scopeConfigId ? scopeConfigName : 'N/A'}</span>
      {pluginConfig.scopeConfig && (
        <IconButton
          icon={<LinkOutlined />}
          helptip="Associate Scope Config"
          size="small"
          type="link"
          onClick={() => {
            setType('associate');
          }}
        />
      )}
      {scopeConfigId && (
        <IconButton
          icon={<EditOutlined />}
          helptip=" Edit Scope Config"
          type="link"
          size="small"
          onClick={() => handleCheckScopeConfig(scopeConfigId)}
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
              scopeConfigId={scopeConfigId}
              scopeId={scopeId}
              onCancel={handleHideDialog}
              onSubmit={handleAssociate}
            />
          ) : (
            <ScopeConfigSelect
              plugin={plugin}
              connectionId={connectionId}
              scopeConfigId={scopeConfigId}
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
            scopeConfigId={scopeConfigId}
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
            scopeConfigId={scopeConfigId}
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
          title={`Edit '${scopeConfigName}' for '${scopeName}'`}
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
