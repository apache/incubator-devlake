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

import { useState, useEffect, useMemo } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { Button, Intent } from '@blueprintjs/core';

import { PageHeader, Buttons, Dialog, IconButton, Table, Message, toast } from '@/components';
import { useTips, useConnections, useRefreshData } from '@/hooks';
import ClearImg from '@/images/icons/clear.svg';
import {
  ConnectionForm,
  ConnectionStatus,
  DataScopeRemote,
  getPluginConfig,
  getPluginScopeId,
  ScopeConfigForm,
  ScopeConfigSelect,
} from '@/plugins';
import { operator } from '@/utils';

import * as API from './api';
import * as S from './styled';

export const ConnectionDetailPage = () => {
  const { plugin, id } = useParams() as { plugin: string; id: string };
  return <ConnectionDetail plugin={plugin} connectionId={+id} />;
};

interface Props {
  plugin: string;
  connectionId: ID;
}

const ConnectionDetail = ({ plugin, connectionId }: Props) => {
  const [type, setType] = useState<
    | 'deleteConnection'
    | 'updateConnection'
    | 'createDataScope'
    | 'clearDataScope'
    | 'deleteDataScope'
    | 'associateScopeConfig'
    | 'deleteConnectionFailed'
    | 'deleteDataScopeFailed'
  >();
  const [operating, setOperating] = useState(false);
  const [version, setVersion] = useState(1);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [scopeId, setScopeId] = useState<ID>();
  const [scopeIds, setScopeIds] = useState<ID[]>([]);
  const [scopeConfigId, setScopeConfigId] = useState<ID>();
  const [conflict, setConflict] = useState<string[]>([]);
  const [errorMsg, setErrorMsg] = useState('');

  const navigate = useNavigate();
  const { onGet, onTest, onRefresh } = useConnections();
  const { setTips } = useTips();
  const { ready, data } = useRefreshData(
    () => API.getDataScopes(plugin, connectionId, { page, pageSize }),
    [version, page, pageSize],
  );

  const { unique, status, name, icon } = onGet(`${plugin}-${connectionId}`) || {};

  const pluginConfig = useMemo(() => getPluginConfig(plugin), [plugin]);

  const [dataSource, total] = useMemo(
    () => [
      data?.scopes.map((it: any) => ({
        id: getPluginScopeId(plugin, it.scope),
        name: it.scope.fullName ?? it.scope.name,
        projects: it.blueprints?.map((bp: any) => bp.projectName) ?? [],
        configId: it.scopeConfig?.id,
        configName: it.scopeConfig?.name,
      })) ?? [],
      data?.count ?? 0,
    ],
    [data],
  );

  useEffect(() => {
    onTest(`${plugin}-${connectionId}`);
  }, [plugin, connectionId]);

  const handleHideDialog = () => {
    setType(undefined);
  };

  const handleShowDeleteDialog = () => {
    setType('deleteConnection');
  };

  const handleDelete = async () => {
    const [, res] = await operator(
      async () => {
        try {
          await API.deleteConnection(plugin, connectionId);
          return { status: 'success' };
        } catch (err: any) {
          const { status, data } = err.response;
          return {
            status: status === 409 ? 'conflict' : 'error',
            conflict: data.data ? [...data.data.projects, ...data.data.blueprints] : [],
            message: data.message,
          };
        }
      },
      {
        setOperating,
        hideToast: true,
      },
    );

    if (res.status === 'success') {
      toast.success('Delete Connection Successful.');
      onRefresh(plugin);
      navigate('/connections');
    } else if (res.status === 'conflict') {
      setType('deleteConnectionFailed');
      setConflict(res.conflict);
      setErrorMsg(res.message);
    } else {
      toast.error('Operation failed.');
      handleHideDialog();
    }
  };

  const handleShowUpdateDialog = () => {
    setType('updateConnection');
  };

  const handleUpdate = () => {
    onRefresh(plugin);
    handleHideDialog();
  };

  const handleShowCreateDataScopeDialog = () => {
    setType('createDataScope');
  };

  const handleCreateDataScope = () => {
    setVersion((v) => v + 1);
    handleHideDialog();
  };

  const handleShowClearDataScopeDialog = (scopeId: ID) => {
    setType('clearDataScope');
    setScopeId(scopeId);
  };

  const handleShowDeleteDataScopeDialog = (scopeId: ID) => {
    setType('deleteDataScope');
    setScopeId(scopeId);
  };

  const handleDeleteDataScope = async (onlyData: boolean) => {
    if (!scopeId) return;

    const [, res] = await operator(
      async () => {
        try {
          await API.deleteDataScope(plugin, connectionId, scopeId, onlyData);
          return { status: 'success' };
        } catch (err: any) {
          const { status, data } = err.response;
          return {
            status: status === 409 ? 'conflict' : 'error',
            conflict: data.data ? [...data.data.projects, ...data.data.blueprints] : [],
            message: data.message,
          };
        }
      },
      {
        setOperating,
        hideToast: true,
      },
    );

    if (res.status === 'success') {
      setVersion((v) => v + 1);
      toast.success(onlyData ? 'Clear historical data successful.' : 'Delete Data Scope successful.');
      handleHideDialog();
    } else if (res.status === 'conflict') {
      setType('deleteDataScopeFailed');
      setConflict(res.conflict);
      setErrorMsg(res.message);
    } else {
      toast.error('Operation failed.');
      handleHideDialog();
    }
  };

  const handleShowScopeConfigSelectDialog = (scopeIds: ID[]) => {
    setType('associateScopeConfig');
    setScopeIds(scopeIds);
  };

  const handleAssociateScopeConfig = async (trId: ID) => {
    const [success] = await operator(
      () =>
        Promise.all(
          scopeIds.map(async (scopeId) => {
            const scope = await API.getDataScope(plugin, connectionId, scopeId);
            return API.updateDataScope(plugin, connectionId, scopeId, {
              ...scope,
              scopeConfigId: +trId,
            });
          }),
        ),
      {
        setOperating,
        formatMessage: () => `Associate scope config successful.`,
      },
    );

    if (success) {
      setVersion((v) => v + 1);
      setTips(
        <Message
          content="Scope Config(s) have been updated. If you would like to re-transform or re-collect the data in the related
        project(s), please go to the Project page and do so."
        />,
      );
      handleHideDialog();
    }
  };

  return (
    <PageHeader
      breadcrumbs={[
        { name: 'Connections', path: '/connections' },
        { name, path: '' },
      ]}
      extra={
        <S.PageHeaderExtra>
          <span style={{ marginRight: 4 }}>Status:</span>
          <ConnectionStatus status={status} unique={unique} onTest={onTest} />
          <Buttons style={{ marginLeft: 8 }}>
            <Button outlined intent={Intent.PRIMARY} icon="annotation" text="Edit" onClick={handleShowUpdateDialog} />
            <Button intent={Intent.DANGER} icon="trash" text="Delete" onClick={handleShowDeleteDialog} />
          </Buttons>
        </S.PageHeaderExtra>
      }
    >
      <S.Wrapper>
        <div className="top">Please note: In order to view DORA metrics, you will need to add Scope Configs.</div>
        <Buttons position="top">
          <Button intent={Intent.PRIMARY} icon="add" text="Add Data Scope" onClick={handleShowCreateDataScopeDialog} />
          {plugin !== 'tapd' && pluginConfig.scopeConfig && (
            <Button
              disabled={!scopeIds.length}
              intent={Intent.PRIMARY}
              icon="many-to-one"
              text="Associate Scope Config"
              onClick={() => handleShowScopeConfigSelectDialog(scopeIds)}
            />
          )}
        </Buttons>
        <Table
          loading={!ready}
          columns={[
            {
              title: 'Data Scope',
              dataIndex: 'name',
              key: 'name',
            },
            {
              title: 'Project',
              dataIndex: 'projects',
              key: 'projects',
              render: (projects) => (
                <>
                  {projects.length ? (
                    <ul>
                      {projects.map((it: string) => (
                        <li key={it}>
                          <Link to={`/projects/${it}`}>{it}</Link>
                        </li>
                      ))}
                    </ul>
                  ) : (
                    '-'
                  )}
                </>
              ),
            },
            {
              title: 'Scope Config',
              dataIndex: ['id', 'configId', 'configName'],
              key: 'scopeConfig',
              width: 400,
              render: ({ id, configId, configName }) => (
                <>
                  <span>{configId ? configName : 'N/A'}</span>
                  {pluginConfig.scopeConfig && (
                    <IconButton
                      icon="link"
                      tooltip="Associate Scope Config"
                      onClick={() => {
                        handleShowScopeConfigSelectDialog([id]);
                        setScopeConfigId(configId);
                      }}
                    />
                  )}
                </>
              ),
            },
            {
              title: '',
              dataIndex: 'id',
              key: 'id',
              width: 100,
              render: (id) => (
                <>
                  <IconButton
                    image={<img src={ClearImg} alt="clear" />}
                    tooltip="Clear historical data"
                    onClick={() => handleShowClearDataScopeDialog(id)}
                  />
                  <IconButton
                    icon="trash"
                    tooltip="Delete Data Scope"
                    onClick={() => handleShowDeleteDataScopeDialog(id)}
                  />
                </>
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
            text: 'Add data to this connection.',
            btnText: 'Add Data Scope',
            onCreate: handleShowCreateDataScopeDialog,
          }}
          rowSelection={{
            getRowKey: (row) => row.id,
            selectedRowKeys: scopeIds,
            onChange: (selectedRowKeys) => setScopeIds(selectedRowKeys),
          }}
        />
      </S.Wrapper>
      {type === 'deleteConnection' && (
        <Dialog
          isOpen
          style={{ width: 820 }}
          title="Would you like to delete this Data Connection?"
          okText="Confirm"
          okLoading={operating}
          onCancel={handleHideDialog}
          onOk={handleDelete}
        >
          <Message
            content=" This operation cannot be undone. Deleting a Data Connection will delete all data that have been collected
              in this Connection."
          />
        </Dialog>
      )}
      {type === 'updateConnection' && (
        <Dialog
          isOpen
          style={{ width: 820 }}
          footer={null}
          title={
            <S.DialogTitle>
              <img src={icon} alt="" />
              <span>Authentication</span>
            </S.DialogTitle>
          }
          onCancel={handleHideDialog}
        >
          <ConnectionForm plugin={plugin} connectionId={connectionId} onSuccess={handleUpdate} />
        </Dialog>
      )}
      {type === 'createDataScope' && (
        <Dialog
          isOpen
          style={{ width: 820 }}
          footer={null}
          title={
            <S.DialogTitle>
              <img src={icon} alt="" />
              <span>Add Data Scope: {name}</span>
            </S.DialogTitle>
          }
          onCancel={handleHideDialog}
        >
          <DataScopeRemote
            plugin={plugin}
            connectionId={connectionId}
            disabledScope={dataSource}
            onCancel={handleHideDialog}
            onSubmit={handleCreateDataScope}
          />
        </Dialog>
      )}
      {type === 'clearDataScope' && (
        <Dialog
          isOpen
          style={{ width: 820 }}
          title="Would you like to clear the historical data of the selected Data Scope?"
          okText="Confirm"
          okLoading={operating}
          onCancel={handleHideDialog}
          onOk={() => handleDeleteDataScope(true)}
        >
          <Message content="This operation cannot be undone." />
        </Dialog>
      )}
      {type === 'deleteDataScope' && (
        <Dialog
          isOpen
          style={{ width: 820 }}
          title="Would you like to delete the selected Data Scope?"
          okText="Confirm"
          okLoading={operating}
          onCancel={handleHideDialog}
          onOk={() => handleDeleteDataScope(false)}
        >
          <Message
            content="This operation cannot be undone. Deleting Data Scope will delete all data that have been collected in the
              past."
          />
        </Dialog>
      )}
      {type === 'associateScopeConfig' && (
        <Dialog isOpen style={{ width: 960 }} footer={null} title="Associate Scope Config" onCancel={handleHideDialog}>
          {plugin === 'tapd' ? (
            <ScopeConfigForm
              plugin={plugin}
              connectionId={connectionId}
              scopeId={scopeIds[0]}
              scopeConfigId={scopeConfigId}
              onCancel={handleHideDialog}
              onSubmit={handleAssociateScopeConfig}
            />
          ) : (
            <ScopeConfigSelect
              plugin={plugin}
              connectionId={connectionId}
              scopeConfigId={scopeConfigId}
              onCancel={handleHideDialog}
              onSubmit={handleAssociateScopeConfig}
            />
          )}
        </Dialog>
      )}
      {type === 'deleteConnectionFailed' && (
        <Dialog
          isOpen
          style={{ width: 820 }}
          footer={null}
          title={`This Data Connection can not be deleted.`}
          onCancel={handleHideDialog}
        >
          <S.DialogBody>
            {!conflict.length ? (
              <Message content={errorMsg} />
            ) : (
              <>
                <Message
                  content={`This Data Connection can not be deleted because it has been used in the following projects/blueprints:`}
                />
                <ul>
                  {conflict.map((it) => (
                    <li key={it}>{it}</li>
                  ))}
                </ul>
              </>
            )}
            <Buttons position="bottom" align="right">
              <Button intent={Intent.PRIMARY} text="OK" onClick={handleHideDialog} />
            </Buttons>
          </S.DialogBody>
        </Dialog>
      )}
      {type === 'deleteDataScopeFailed' && (
        <Dialog
          isOpen
          style={{ width: 820 }}
          footer={null}
          title={`This Data Scope can not be deleted.`}
          onCancel={handleHideDialog}
        >
          <S.DialogBody>
            {!conflict.length ? (
              <Message content={errorMsg} />
            ) : (
              <>
                <Message content="This Data Scope can not be deleted because it has been used in the following projects/blueprints:" />
                <ul>
                  {conflict.map((it) => (
                    <li key={it}>{it}</li>
                  ))}
                </ul>
              </>
            )}
            <Buttons position="bottom" align="right">
              <Button intent={Intent.PRIMARY} text="OK" onClick={handleHideDialog} />
            </Buttons>
          </S.DialogBody>
        </Dialog>
      )}
    </PageHeader>
  );
};
