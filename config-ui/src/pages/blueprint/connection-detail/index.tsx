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
import { useHistory, useParams } from 'react-router-dom';
import { Button, Intent, Position } from '@blueprintjs/core';
import { Popover2 } from '@blueprintjs/popover2';

import { PageLoading, PageHeader, ExternalLink, Buttons, Table, Dialog } from '@/components';
import { useRefreshData } from '@/hooks';
import { DataScopeSelect, getPluginId } from '@/plugins';
import { operator } from '@/utils';

import * as API from './api';
import * as S from './styled';

export const BlueprintConnectionDetailPage = () => {
  const [version, setVersion] = useState(1);
  const [isOpen, setIsOpen] = useState(false);

  const { bid, unique } = useParams<{ bid: string; unique: string }>();
  const history = useHistory();

  const { ready, data } = useRefreshData(async () => {
    const [plugin, connectionId] = unique.split('-');
    const [blueprint, connection, scopes] = await Promise.all([
      API.getBlueprint(bid),
      API.getConnection(plugin, connectionId),
      API.getDataScopes(plugin, connectionId),
    ]);

    const scopeIds = blueprint.settings.connections
      .find((cs: any) => cs.plugin === plugin && cs.connectionId === +connectionId)
      .scopes.map((sc: any) => +sc.id);

    return {
      blueprint,
      connection: {
        unique,
        plugin,
        id: +connectionId,
        name: connection.name,
      },
      scopes: scopes.filter((sc: any) => scopeIds.includes(sc[getPluginId(plugin)])),
    };
  }, [version]);

  if (!ready || !data) {
    return <PageLoading />;
  }

  const { blueprint, connection, scopes } = data;

  const handleShowDataScope = () => setIsOpen(true);
  const handleHideDataScope = () => setIsOpen(false);

  const handleRemoveConnection = async () => {
    const [success] = await operator(() =>
      API.updateBlueprint(blueprint.id, {
        ...blueprint,
        settings: {
          ...blueprint.settings,
          connections: blueprint.settings.connections.filter(
            (cs: any) => !(cs.plugin === connection.plugin && cs.connectionId === connection.id),
          ),
        },
      }),
    );

    if (success) {
      history.push(`/blueprints/${blueprint.id}`);
    }
  };

  const handleChangeDataScope = async (scope: any) => {
    const [success] = await operator(
      () =>
        API.updateBlueprint(blueprint.id, {
          ...blueprint,
          settings: {
            ...blueprint.settings,
            connections: blueprint.settings.connections.map((cs: any) => {
              if (cs.plugin === connection.plugin && cs.connectionId === connection.id) {
                return {
                  ...cs,
                  scopes: scope.map((sc: any) => ({ id: `${sc[getPluginId(connection.plugin)]}` })),
                };
              }
              return cs;
            }),
          },
        }),
      {
        formatMessage: () => 'Update data scope successful.',
      },
    );

    if (success) {
      handleHideDataScope();
      setVersion((v) => v + 1);
    }
  };

  return (
    <PageHeader
      breadcrumbs={[
        { name: blueprint.name, path: `/blueprints/${bid}` },
        { name: `Connection - ${connection.name}`, path: '' },
      ]}
    >
      <S.Top>
        <span>
          If you would like to manage Data Entities and Data Scope of this Connection, please{' '}
          <ExternalLink link={`/connections/${connection.plugin}/${connection.id}`}>
            go to the Connection detail page
          </ExternalLink>
          .
        </span>
        <Popover2
          position={Position.BOTTOM}
          content={
            <S.ActionDelete>
              <div className="content">Are you sure you want to delete this connection?</div>
              <div className="btns" onClick={handleRemoveConnection}>
                <Button intent={Intent.PRIMARY} text="Confirm" />
              </div>
            </S.ActionDelete>
          }
        >
          <Button intent={Intent.DANGER} icon="trash">
            Remove this Connection
          </Button>
        </Popover2>
      </S.Top>
      <Buttons position="top" align="left">
        <Button intent={Intent.PRIMARY} icon="annotation" text="Manage Data Scope" onClick={handleShowDataScope} />
      </Buttons>
      <Table columns={[{ title: 'Data Scope', dataIndex: 'name', key: 'name' }]} dataSource={scopes} />
      <Dialog
        isOpen={isOpen}
        title="Change Data Scope"
        footer={null}
        style={{ width: 820 }}
        onCancel={handleHideDataScope}
      >
        <DataScopeSelect
          plugin={connection.plugin}
          connectionId={connection.id}
          initialScope={scopes}
          onCancel={handleHideDataScope}
          onSubmit={handleChangeDataScope}
        />
      </Dialog>
    </PageHeader>
  );
};
