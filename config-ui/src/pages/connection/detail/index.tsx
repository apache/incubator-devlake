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

import { useEffect, useState } from 'react';
import { useParams, useHistory } from 'react-router-dom';
import { Button, Icon, Intent } from '@blueprintjs/core';

import { PageHeader, Buttons, Dialog, IconButton, Table } from '@/components';
import { useTips, useConnections, useRefreshData } from '@/hooks';
import {
  ConnectionForm,
  ConnectionStatus,
  DataScopeSelectRemote,
  getPluginId,
  ScopeConfigForm,
  ScopeConfigSelect,
} from '@/plugins';
import { operator } from '@/utils';

import * as API from './api';
import * as S from './styled';

export const ConnectionDetailPage = () => {
  const [type, setType] = useState<
    | 'deleteConnection'
    | 'updateConnection'
    | 'createDataScope'
    | 'clearDataScope'
    | 'deleteDataScope'
    | 'associateScopeConfig'
  >();
  const [operating, setOperating] = useState(false);
  const [version, setVersion] = useState(1);
  const [scopeId, setScopeId] = useState<ID>();
  const [scopeIds, setScopeIds] = useState<ID[]>([]);
  const [scopeConfigId, setScopeConfigId] = useState<ID>();

  const { plugin, id } = useParams<{ plugin: string; id: string }>();
  const history = useHistory();
  const { onGet, onTest, onRefresh } = useConnections();
  const { setTips } = useTips();
  const { ready, data } = useRefreshData(() => API.getDataScopes(plugin, id), [version]);

  const { unique, status, name, icon } = onGet(`${plugin}-${id}`) || {};

  useEffect(() => {
    onTest(`${plugin}-${id}`);
  }, [plugin, id]);

  const handleHideDialog = () => {
    setType(undefined);
  };

  const handleShowTips = () => {
    setTips(
      <div>
        <Icon icon="warning-sign" style={{ marginRight: 8 }} color="#F4BE55" />
        <span>
          The transformation of certain data scope has been updated. If you would like to re-transform the data in the
          related project(s), please go to the Project page and do so.
        </span>
      </div>,
    );
  };

  const handleShowDeleteDialog = () => {
    setType('deleteConnection');
  };

  const handleDelete = async () => {
    const [success] = await operator(() => API.deleteConnection(plugin, id), {
      setOperating,
      formatMessage: () => 'Delete Connection Successful.',
    });

    if (success) {
      onRefresh(plugin);
      history.push('/connections');
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
    handleShowTips();
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

    const [success] = await operator(() => API.deleteDataScope(plugin, id, scopeId, onlyData), {
      setOperating,
      formatMessage: () => (onlyData ? 'Clear historical data successful.' : 'Delete Data Scope successful.'),
    });

    if (success) {
      setVersion((v) => v + 1);
      handleShowTips();
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
            const scope = await API.getDataScope(plugin, id, scopeId);
            return API.updateDataScope(plugin, id, scopeId, {
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
      handleShowTips();
      handleHideDialog();
    }
  };

  return (
    <PageHeader
      breadcrumbs={[
        { name: 'Connections', path: '/connections' },
        { name, path: '' },
      ]}
      extra={<Button intent={Intent.DANGER} icon="trash" text="Delete Connection" onClick={handleShowDeleteDialog} />}
    >
      <S.Wrapper>
        <div className="authentication">
          <span style={{ marginRight: 4 }}>Authentication Status:</span>
          <ConnectionStatus status={status} unique={unique} onTest={onTest} />
          <IconButton icon="annotation" tooltip="Edit Connection" onClick={handleShowUpdateDialog} />
        </div>
        <Buttons position="top" align="left">
          <Button intent={Intent.PRIMARY} icon="add" text="Add Data Scope" onClick={handleShowCreateDataScopeDialog} />
          {plugin !== 'tapd' && (
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
              title: 'Scope Config',
              dataIndex: 'scopeConfigName',
              key: 'scopeConfig',
              width: 400,
              render: (val, row) => (
                <>
                  <span>{val ?? 'No Scope Config'}</span>
                  <IconButton
                    icon="link"
                    tooltip="Associate Scope Config"
                    onClick={() => {
                      handleShowScopeConfigSelectDialog([row[getPluginId(plugin)]]);
                      if (plugin === 'tapd') {
                        setScopeConfigId(row.scopeConfigId);
                      }
                    }}
                  />
                </>
              ),
            },
            {
              title: '',
              dataIndex: getPluginId(plugin),
              key: 'id',
              width: 100,
              render: (id) => (
                <>
                  <IconButton
                    icon="unarchive"
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
          dataSource={data}
          noData={{
            text: 'Add data to this connection.',
            btnText: 'Add Data Scope',
            onCreate: handleShowCreateDataScopeDialog,
          }}
          rowSelection={{
            rowKey: getPluginId(plugin),
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
          <S.DialogBody>
            <Icon icon="warning-sign" />
            <span>
              This operation cannot be undone. Deleting a Data Connection will delete all data that have been collected
              in this Connection.
            </span>
          </S.DialogBody>
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
          <ConnectionForm plugin={plugin} connectionId={id} onSuccess={handleUpdate} />
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
          <DataScopeSelectRemote
            plugin={plugin}
            connectionId={id}
            disabledScope={data}
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
          <S.DialogBody>
            <Icon icon="warning-sign" />
            <span>This operation cannot be undone.</span>
          </S.DialogBody>
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
          <S.DialogBody>
            <Icon icon="warning-sign" />
            <span>
              This operation cannot be undone. Deleting Data Scope will delete all data that have been collected in the
              past.
            </span>
          </S.DialogBody>
        </Dialog>
      )}
      {type === 'associateScopeConfig' && (
        <Dialog isOpen style={{ width: 820 }} footer={null} title="Associate Scope Config" onCancel={handleHideDialog}>
          {plugin === 'tapd' ? (
            <ScopeConfigForm
              plugin={plugin}
              connectionId={id}
              scopeId={scopeIds[0]}
              scopeConfigId={scopeConfigId}
              onCancel={handleHideDialog}
              onSubmit={handleAssociateScopeConfig}
            />
          ) : (
            <ScopeConfigSelect
              plugin={plugin}
              connectionId={id}
              onCancel={handleHideDialog}
              onSubmit={handleAssociateScopeConfig}
            />
          )}
        </Dialog>
      )}
    </PageHeader>
  );
};
