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

import { Dialog, PageHeader, PageLoading } from '@/components';
import { EntitiesLabel } from '@/config';
import { useRefreshData } from '@/hooks';
import { DataScope, getPluginConfig, Transformation } from '@/plugins';

import * as API from './api';
import * as S from './styled';

export const BlueprintConnectionDetailPage = () => {
  const [version, setVersion] = useState(1);
  const [isOpen, setIsOpen] = useState(false);

  const { pname, bid, unique } = useParams<{ pname?: string; bid: string; unique: string }>();
  const history = useHistory();

  const { ready, data } = useRefreshData(async () => {
    const [plugin, connectionId] = unique.split('-');
    const blueprint = await API.getBlueprint(bid);
    const connection = await API.getConnection(plugin, connectionId);
    const scope = blueprint.settings.connections.find(
      (cs: any) => cs.plugin === plugin && cs.connectionId === +connectionId,
    ).scopes;
    const config = getPluginConfig(plugin);
    const origin = await Promise.all(scope.map((sc: any) => API.getDataScope(plugin, connectionId, sc.id)));

    return {
      blueprint,
      bpName: blueprint.name,
      csName: connection.name,
      entities: scope[0].entities,
      connection: {
        unique,
        plugin,
        connectionId: +connectionId,
        name: connection.name,
        icon: config.icon,
        scope,
        origin,
        transformationType: config.transformationType,
      },
    };
  }, [version]);

  if (!ready || !data) {
    return <PageLoading />;
  }

  const { blueprint, bpName, csName, entities, connection } = data;

  const handleShowDataScope = () => setIsOpen(true);
  const handleHideDataScope = () => setIsOpen(false);

  const handleChangeDataScope = async (connections: MixConnection[]) => {
    const [connection] = connections;
    await API.updateBlueprint(blueprint.id, {
      ...blueprint,
      settings: {
        ...blueprint.settings,
        connections: blueprint.settings.connections.map((cs: any) => {
          if (cs.plugin === connection.plugin && cs.connectionId === connection.connectionId) {
            return {
              ...cs,
              scopes: connection.scope.map((sc: any) => ({
                id: `${sc.id}`,
                entities: sc.entities,
              })),
            };
          }
          return cs;
        }),
      },
    });
    setVersion((v) => v + 1);
  };

  const handleChangeTransformation = () => {
    setVersion((v) => v + 1);
  };

  const handleRemoveConnection = async () => {
    await API.updateBlueprint(blueprint.id, {
      ...blueprint,
      settings: {
        ...blueprint.settings,
        connections: blueprint.settings.connections.filter(
          (cs: any) => !(cs.plugin === connection.plugin && cs.connectionId === connection.connectionId),
        ),
      },
    });
    history.push(pname ? `/projects/:${pname}` : `/blueprints/${blueprint.id}`);
  };

  return (
    <PageHeader
      breadcrumbs={[
        ...(pname
          ? [
              {
                name: 'Projects',
                path: '/projects',
              },
              {
                name: pname,
                path: `/projects/${pname}`,
              },
            ]
          : [{ name: bpName, path: `/blueprints/${bid}` }]),
        { name: `Connection - ${csName}`, path: '' },
      ]}
    >
      <S.Action>
        <span>
          <Button intent={Intent.PRIMARY} icon="annotation" onClick={handleShowDataScope}>
            Edit Data Scope
          </Button>
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
      </S.Action>
      <S.Entities>
        <h4>Data Entities</h4>
        <ul>
          {entities.map((it: string) => (
            <li key={it}>{EntitiesLabel[it]}</li>
          ))}
        </ul>
      </S.Entities>
      <Transformation connections={[connection]} noFooter onSubmit={handleChangeTransformation} />
      <Dialog
        isOpen={isOpen}
        title="Change Data Scope"
        footer={null}
        style={{ width: 820 }}
        onCancel={handleHideDataScope}
      >
        <DataScope
          connections={[connection]}
          onCancel={handleHideDataScope}
          onSubmit={handleChangeDataScope}
          onNext={handleHideDataScope}
        />
      </Dialog>
    </PageHeader>
  );
};
