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
import { useParams, useHistory } from 'react-router-dom';
import { Button, Icon, Intent } from '@blueprintjs/core';

import { PageHeader, Dialog, IconButton, Table } from '@/components';
import { useTips, useConnections, useRefreshData } from '@/hooks';
import { ConnectionForm, ConnectionStatus, DataScopeSelectRemote, getPluginId } from '@/plugins';
import { operator } from '@/utils';

import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  id: ID;
}

const ConnectionDetail = ({ plugin, id }: Props) => {
  const [type, setType] = useState<
    'deleteConnection' | 'updateConnection' | 'createDataScope' | 'clearDataScope' | 'deleteDataScope'
  >();
  const [operating, setOperating] = useState(false);
  const [version, setVersion] = useState(1);
  const [scopeId, setScopeId] = useState<ID>();

  const history = useHistory();
  const { onGet, onTest, onRefresh } = useConnections();
  const { setTips } = useTips();
  const { ready, data } = useRefreshData(() => API.getDataScope(plugin, id), [version]);

  const { unique, status, name, icon } = onGet(`${plugin}-${id}`);

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
        <div className="action">
          <Button intent={Intent.PRIMARY} icon="add" text="Add Data Scope" onClick={handleShowCreateDataScopeDialog} />
        </div>
        <Table
          loading={!ready}
          columns={[
            {
              title: 'Data Scope',
              dataIndex: 'name',
              key: 'name',
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
    </PageHeader>
  );
};

export const ConnectionDetailPage = () => {
  const { plugin, id } = useParams<{ plugin: string; id: string }>();

  return <ConnectionDetail plugin={plugin} id={+id} />;
};
