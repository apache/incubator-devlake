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
import { Link } from 'react-router-dom';
import { PlusOutlined } from '@ant-design/icons';
import { Alert, Button } from 'antd';

import API from '@/api';
import { NoData } from '@/components';
import type { WebhookItemType } from '@/plugins/register/webhook';
import { WebhookCreateDialog, WebhookSelectorDialog, WebHookConnection } from '@/plugins/register/webhook';
import { IProject } from '@/types';
import { operator } from '@/utils';

interface Props {
  project: IProject;
  onRefresh: () => void;
}

export const WebhooksPanel = ({ project, onRefresh }: Props) => {
  const [type, setType] = useState<'selectExist' | 'create'>();
  const [operating, setOperating] = useState(false);

  const webhookIds = useMemo(
    () =>
      project.blueprint
        ? project.blueprint.connections.filter((cs) => cs.pluginName === 'webhook').map((cs: any) => cs.connectionId)
        : [],
    [project],
  );

  const handleCancel = () => {
    setType(undefined);
  };

  const handleCreate = async (id: ID) => {
    const payload = {
      ...project.blueprint,
      connections: [
        ...project.blueprint.connections,
        {
          pluginName: 'webhook',
          connectionId: id,
        },
      ],
    };

    const [success] = await operator(() => API.blueprint.update(project.blueprint.id, payload), {
      setOperating,
    });

    if (success) {
      onRefresh();
    }
  };

  const handleSelect = async (items: WebhookItemType[]) => {
    const payload = {
      ...project.blueprint,
      connections: [
        ...project.blueprint.connections,
        ...items.map((it) => ({
          pluginName: 'webhook',
          connectionId: it.id,
        })),
      ],
    };

    const [success] = await operator(() => API.blueprint.update(project.blueprint.id, payload), {
      setOperating,
    });

    if (success) {
      onRefresh();
    }
  };

  const handleDelete = async (id: ID) => {
    const payload = {
      ...project.blueprint,
      connections: project.blueprint.connections.filter(
        (cs) => !(cs.pluginName === 'webhook' && cs.connectionId === id),
      ),
    };

    const [success] = await operator(() => API.blueprint.update(project.blueprint.id, payload), {
      setOperating,
    });

    if (success) {
      onRefresh();
    }
  };

  return (
    <>
      <Alert
        style={{ marginBottom: 24 }}
        message={
          <>
            <div>
              The data pushed by Webhooks will only be calculated for DORA in the next run of the Blueprint of this
              project because DORA relies on the post-processing of "deployments," "incidents," and "pull requests"
              triggered by running the blueprint.
            </div>
            <div style={{ marginTop: 16 }}>
              To calculate DORA after receiving Webhook data immediately, you can visit the{' '}
              <b style={{ textDecoration: 'underline' }}>
                <Link to={`${window.location.pathname}?tab=status`}>Status tab</Link>
              </b>{' '}
              of the Blueprint page and click on Run Now.
            </div>
          </>
        }
      />
      {!webhookIds.length ? (
        <>
          <NoData
            text="Push `incidents` or `deployments` from your tools by incoming webhooks."
            action={
              <>
                <Button type="primary" icon={<PlusOutlined />} onClick={() => setType('create')}>
                  Add a Webhook
                </Button>
                <div style={{ margin: '8px 0' }}>or</div>
                <Button type="primary" onClick={() => setType('selectExist')}>
                  Select Existing Webhooks
                </Button>
              </>
            }
          />
          {type === 'create' && <WebhookCreateDialog open onCancel={handleCancel} onSubmitAfter={handleCreate} />}
          {type === 'selectExist' && (
            <WebhookSelectorDialog open saving={operating} onCancel={handleCancel} onSubmit={handleSelect} />
          )}
        </>
      ) : (
        <WebHookConnection filterIds={webhookIds} onCreateAfter={handleCreate} onDeleteAfter={handleDelete} />
      )}
    </>
  );
};
