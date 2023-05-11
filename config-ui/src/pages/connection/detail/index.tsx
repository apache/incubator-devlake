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

import { PageHeader, Dialog, IconButton } from '@/components';
import { transformEntities } from '@/config';
import { ConnectionForm } from '@/plugins';
import type { ConnectionItemType } from '@/store';
import { ConnectionContextProvider, useConnection, ConnectionStatus } from '@/store';
import { operator } from '@/utils';

import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  id: ID;
}

const ConnectionDetail = ({ plugin, id }: Props) => {
  const [type, setType] = useState<'deleteConnection' | 'updateConnection'>();
  const [operating, setOperating] = useState(false);

  const history = useHistory();
  const { connections, onRefresh, onTest } = useConnection();
  const { unique, status, name, icon, entities } = connections.find(
    (cs) => cs.unique === `${plugin}-${id}`,
  ) as ConnectionItemType;

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
    setType(undefined);
    onRefresh(plugin);
  };

  const handleHideDialog = () => {
    setType(undefined);
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
          style={{ width: 820 }}
          footer={null}
          isOpen
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
    </PageHeader>
  );
};

export const ConnectionDetailPage = () => {
  const { plugin, id } = useParams<{ plugin: string; id: string }>();

  return (
    <ConnectionContextProvider plugin={plugin}>
      <ConnectionDetail plugin={plugin} id={id} />
    </ConnectionContextProvider>
  );
};
