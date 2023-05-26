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
import { transformEntities } from '@/config';
import { useTips, useConnections, useRefreshData } from '@/hooks';
import { ConnectionForm, ConnectionStatus, DataScopeForm2 } from '@/plugins';
import { operator } from '@/utils';

import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  id: ID;
}

const ConnectionDetail = ({ plugin, id }: Props) => {
  const [type, setType] = useState<'deleteConnection' | 'updateConnection' | 'createDataScope'>();
  const [operating, setOperating] = useState(false);
  const [version, setVersion] = useState(1);

  const history = useHistory();
  const { onGet, onTest, onRefresh } = useConnections();
  const { setTips } = useTips();
  const { ready, data } = useRefreshData(() => API.getDataScope(plugin, id), [version]);

  const { unique, status, name, icon, entities } = onGet(`${plugin}-${id}`);

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

  return (
    <PageHeader
      breadcrumbs={[
        { name: 'Connections', path: '/connections' },
        { name, path: '' },
      ]}
      extra={<Button intent={Intent.DANGER} icon="trash" text="Delete Connection" onClick={handleShowDeleteDialog} />}
    >
      <S.Wrapper>
        <div className="top">
          <div className="entities">
            <h3>Data Entities</h3>
            <span>
              {transformEntities(entities)
                .map((it) => it.label)
                .join(',')}
            </span>
          </div>
          <div className="authentication">
            <h3>
              <span>Authentication</span>
              <IconButton icon="annotation" tooltip="Edit Connection" onClick={handleShowUpdateDialog} />
            </h3>
            <span>Status: </span>
            <span>
              Status: <ConnectionStatus status={status} unique={unique} onTest={onTest} />
            </span>
          </div>
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
          <DataScopeForm2
            plugin={plugin}
            connectionId={id}
            disabledScope={data}
            onCancel={handleHideDialog}
            onSubmit={handleCreateDataScope}
          />
        </Dialog>
      )}
    </PageHeader>
  );
};

export const ConnectionDetailPage = () => {
  const { plugin, id } = useParams<{ plugin: string; id: string }>();

  return <ConnectionDetail plugin={plugin} id={+id} />;
};
