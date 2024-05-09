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
import { useNavigate, useParams } from 'react-router-dom';
import { Helmet } from 'react-helmet';
import { DeleteOutlined, FormOutlined } from '@ant-design/icons';
import { Flex, Table, Popconfirm, Modal, Button } from 'antd';

import API from '@/api';
import { PageLoading, PageHeader, ExternalLink } from '@/components';
import { PATHS } from '@/config';
import { useRefreshData } from '@/hooks';
import { ScopeConfig, DataScopeSelect, getPluginScopeId } from '@/plugins';
import { operator } from '@/utils';

import * as S from './styled';

const brandName = import.meta.env.DEVLAKE_BRAND_NAME ?? 'DevLake';

export const BlueprintConnectionDetailPage = () => {
  const [version, setVersion] = useState(1);
  const [open, setOpen] = useState(false);
  const [operating, setOperating] = useState(false);

  const { pname, bid, unique } = useParams() as { pname?: string; bid?: string; unique: string };
  const navigate = useNavigate();

  const [modal, contextHolder] = Modal.useModal();

  const getBlueprint = async (pname?: string, bid?: string) => {
    if (pname) {
      const res = await API.project.get(pname);
      return res.blueprint;
    }

    return API.blueprint.get(bid as any);
  };

  const [plugin, connectionId] = unique.split('-');

  const { ready, data } = useRefreshData(async () => {
    const [blueprint, connection] = await Promise.all([
      getBlueprint(pname, bid),
      API.connection.get(plugin, connectionId),
    ]);

    const scopeIds =
      blueprint.connections
        .find((cs) => cs.pluginName === plugin && cs.connectionId === +connectionId)
        ?.scopes?.map((sc: any) => sc.scopeId) ?? [];

    const scopes = await Promise.all(scopeIds.map((scopeId) => API.scope.get(plugin, connectionId, scopeId)));

    return {
      blueprint,
      connection: {
        unique,
        plugin,
        id: +connectionId,
        name: connection.name,
      },
      scopes: scopes.map((sc) => ({
        id: getPluginScopeId(plugin, sc.scope),
        name: sc.scope.fullName ?? sc.scope.name,
        scopeConfigId: sc.scopeConfig?.id,
        scopeConfigName: sc.scopeConfig?.name,
      })),
    };
  }, [version, pname, bid]);

  if (!ready || !data) {
    return <PageLoading />;
  }

  const { blueprint, connection, scopes } = data;

  const handleShowDataScope = () => setOpen(true);
  const handleHideDataScope = () => setOpen(false);

  const handleRun = async (data?: { skipCollectors?: boolean; fullSync?: boolean }) => {
    const [success] = await operator(() => API.blueprint.trigger(blueprint.id, data), {
      setOperating,
      hideToast: true,
    });

    if (success) {
      navigate(pname ? PATHS.PROJECT(pname, { tab: 'status' }) : PATHS.BLUEPRINT(blueprint.id, 'status'));
    }
  };

  const handleRemoveConnection = async () => {
    const [success] = await operator(() =>
      API.blueprint.update(blueprint.id, {
        ...blueprint,
        connections: blueprint.connections.filter(
          (cs) => !(cs.pluginName === connection.plugin && cs.connectionId === connection.id),
        ),
      }),
    );

    if (success) {
      modal.success({
        closable: true,
        centered: true,
        width: 500,
        title: 'Data Scope Changed',
        content: 'Re-collect the data to get the project metrics updated?',
        footer: (
          <div style={{ marginTop: 20, textAlign: 'center' }}>
            <Button type="primary" loading={operating} onClick={() => handleRun()}>
              Recollect Data
            </Button>
          </div>
        ),
        onCancel: () => {
          navigate(
            pname ? PATHS.PROJECT(pname, { tab: 'configuration' }) : PATHS.BLUEPRINT(blueprint.id, 'configuration'),
          );
        },
      });
    }
  };

  const handleChangeDataScope = async (scopeIds: any) => {
    const [success] = await operator(
      () =>
        API.blueprint.update(blueprint.id, {
          ...blueprint,
          connections: blueprint.connections.map((cs) => {
            if (cs.pluginName === connection.plugin && cs.connectionId === connection.id) {
              return {
                ...cs,
                scopes: scopeIds.map((scopeId: any) => ({ scopeId })),
              };
            }
            return cs;
          }),
        }),
      {
        hideToast: true,
      },
    );

    if (success) {
      handleHideDataScope();
      modal.success({
        closable: true,
        centered: true,
        width: 500,
        title: 'Data Scope Changed',
        content: 'Re-collect the data to get the project metrics updated?',
        footer: (
          <div style={{ marginTop: 20, textAlign: 'center' }}>
            <Button type="primary" loading={operating} onClick={() => handleRun()}>
              Recollect Data
            </Button>
          </div>
        ),
        onCancel: () => {
          setVersion(version + 1);
        },
      });
    }
  };

  const handleChangeScopeConfig = () => {
    modal.success({
      closable: true,
      centered: true,
      width: 550,
      title: 'Scope Config Saved',
      content: 'Please re-transform data to apply the updated scope config.',
      footer: (
        <div style={{ marginTop: 20, textAlign: 'center' }}>
          <Button type="primary" loading={operating} onClick={() => handleRun({ skipCollectors: true })}>
            Re-transform now
          </Button>
        </div>
      ),
      onCancel: () => {
        setVersion(version + 1);
      },
    });
  };

  return (
    <PageHeader
      breadcrumbs={
        pname
          ? [
              { name: 'Projects', path: PATHS.PROJECTS() },
              { name: pname, path: PATHS.PROJECT(pname) },
              { name: `Connection - ${connection.name}`, path: '' },
            ]
          : [
              { name: 'Advanced', path: PATHS.BLUEPRINTS() },
              { name: 'Blueprints', path: PATHS.BLUEPRINTS() },
              { name: bid as any, path: PATHS.BLUEPRINT(bid as any) },
              { name: `Connection - ${connection.name}`, path: '' },
            ]
      }
    >
      <Helmet>
        <title>
          {pname ? pname : blueprint.name} - {connection.name} - {brandName}
        </title>
      </Helmet>
      <S.Top>
        <span>
          To manage the complete data scope and scope config for this connection, please{' '}
          <ExternalLink link={`/connections/${connection.plugin}/${connection.id}`}>
            go to the connection detail page
          </ExternalLink>
          .
        </span>
        <Popconfirm
          placement="top"
          title="Are you sure you want to remove the connection from this project/blueprint?"
          cancelButtonProps={{
            style: {
              display: 'none',
            },
          }}
          okText="Confirm"
          onConfirm={handleRemoveConnection}
        >
          <Button type="primary" danger icon={<DeleteOutlined />}>
            Remove this Connection
          </Button>
        </Popconfirm>
      </S.Top>
      <Flex vertical gap="middle">
        <Flex>
          <Button type="primary" icon={<FormOutlined />} onClick={handleShowDataScope}>
            Manage Data Scope
          </Button>
        </Flex>
        <Table
          rowKey="id"
          size="middle"
          columns={[
            {
              title: 'Data Scope',
              dataIndex: 'name',
              key: 'name',
            },
            {
              title: 'Scope Config',
              key: 'scopeConfig',
              render: (_, { id, name, scopeConfigId, scopeConfigName }) => (
                <ScopeConfig
                  plugin={plugin}
                  connectionId={connectionId}
                  scopeId={id}
                  scopeName={name}
                  id={scopeConfigId}
                  name={scopeConfigName}
                  onSuccess={handleChangeScopeConfig}
                />
              ),
            },
          ]}
          dataSource={scopes}
        />
      </Flex>
      <Modal open={open} width={820} centered title="Manage Data Scope" footer={null} onCancel={handleHideDataScope}>
        <DataScopeSelect
          plugin={connection.plugin}
          connectionId={connection.id}
          showWarning
          initialScope={scopes}
          onCancel={handleHideDataScope}
          onSubmit={handleChangeDataScope}
        />
      </Modal>
      {contextHolder}
    </PageHeader>
  );
};
