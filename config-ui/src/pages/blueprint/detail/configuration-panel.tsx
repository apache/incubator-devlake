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
import { Link } from 'react-router-dom';
import { Button, Intent } from '@blueprintjs/core';

import { IconButton, Table, NoData, Buttons } from '@/components';
import { getCron } from '@/config';
import { useConnections } from '@/hooks';
import { getPluginConfig } from '@/plugins';
import { formatTime, operator } from '@/utils';

import { encodeName } from '../../project/utils';

import { BlueprintType, FromEnum } from '../types';
import { ModeEnum } from '../types';
import { validRawPlan } from '../utils';

import { AdvancedEditor, UpdateNameDialog, UpdatePolicyDialog, AddConnectionDialog } from './components';
import * as API from './api';
import * as S from './styled';

interface Props {
  from: FromEnum;
  blueprint: BlueprintType;
  onRefresh: () => void;
  onChangeTab: (tab: string) => void;
}

export const ConfigurationPanel = ({ from, blueprint, onRefresh, onChangeTab }: Props) => {
  const [type, setType] = useState<'name' | 'policy' | 'add-connection'>();
  const [rawPlan, setRawPlan] = useState('');
  const [operating, setOperating] = useState(false);

  useEffect(() => {
    setRawPlan(JSON.stringify(blueprint.plan, null, '  '));
  }, [blueprint]);

  const { onGet } = useConnections();

  const connections = useMemo(
    () =>
      blueprint.connections
        .filter((cs) => cs.pluginName !== 'webhook')
        .map((cs: any) => {
          const unique = `${cs.pluginName}-${cs.connectionId}`;
          const plugin = getPluginConfig(cs.pluginName);
          const connection = onGet(unique);
          return {
            unique,
            icon: plugin.icon,
            name: connection?.name,
            scope: cs.scopes,
          };
        })
        .filter(Boolean),
    [blueprint],
  );

  const handleCancel = () => {
    setType(undefined);
  };

  const handleShowNameDialog = () => {
    setType('name');
  };

  const handleShowPolicyDialog = () => {
    setType('policy');
  };

  const handleShowAddConnectionDialog = () => {
    setType('add-connection');
  };

  const handleUpdate = async (payload: any) => {
    const [success] = await operator(
      () =>
        API.updateBlueprint(blueprint.id, {
          ...blueprint,
          ...payload,
        }),
      {
        setOperating,
        formatMessage: () => 'Update blueprint successful.',
      },
    );

    if (success) {
      onRefresh();
      handleCancel();
    }
  };

  const handleRun = async () => {
    const [success] = await operator(() => API.runBlueprint(blueprint.id, false), {
      setOperating,
      formatMessage: () => 'Trigger blueprint successful.',
    });

    if (success) {
      onRefresh();
      onChangeTab('status');
    }
  };

  return (
    <S.ConfigurationPanel>
      <div className="block">
        <h3>Blueprint Name</h3>
        <div>
          <span>{blueprint.name}</span>
          <IconButton icon="annotation" tooltip="Edit" onClick={handleShowNameDialog} />
        </div>
      </div>
      <div className="block">
        <h3>
          <span>Sync Policy</span>
          <IconButton icon="annotation" tooltip="Edit" onClick={handleShowPolicyDialog} />
        </h3>
        <Table
          columns={[
            {
              title: 'Data Time Range',
              dataIndex: 'timeRange',
              key: 'timeRange',
              render: (val) => (blueprint.mode === ModeEnum.normal ? `${formatTime(val)} to Now` : 'N/A'),
            },
            {
              title: 'Sync Frequency',
              dataIndex: 'frequency',
              key: 'frequency',
              align: 'center',
              render: (val, row) => {
                const cron = getCron(row.isManual, val);
                return `${cron.label}${cron.description}`;
              },
            },
            {
              title: 'Skip Failed Tasks',
              dataIndex: 'skipFailed',
              key: 'skipFailed',
              align: 'center',
              render: (val) => (val ? 'Enabled' : 'Disabled'),
            },
          ]}
          dataSource={[
            {
              timeRange: blueprint.timeAfter,
              frequency: blueprint.cronConfig,
              isManual: blueprint.isManual,
              skipFailed: blueprint.skipOnFail,
            },
          ]}
        />
      </div>
      {blueprint.mode === ModeEnum.normal && (
        <div className="block">
          <h3>Data Connections</h3>
          {!connections.length ? (
            <NoData
              text={
                <>
                  If you have not created data connections yet, please <Link to="/connections">create connections</Link>{' '}
                  first and then add them to the project.
                </>
              }
              action={
                <Button
                  intent={Intent.PRIMARY}
                  icon="add"
                  text="Add a Connection"
                  onClick={handleShowAddConnectionDialog}
                />
              }
            />
          ) : (
            <>
              <Buttons>
                <Button
                  intent={Intent.PRIMARY}
                  icon="add"
                  text="Add a Connection"
                  onClick={handleShowAddConnectionDialog}
                />
              </Buttons>
              <S.ConnectionList>
                {connections.map((cs) => (
                  <S.ConnectionItem key={cs.unique}>
                    <div className="title">
                      <img src={cs.icon} alt="" />
                      <span>{cs.name}</span>
                    </div>
                    <div className="count">
                      <span>{cs.scope.length} data scope</span>
                    </div>
                    <div className="link">
                      <Link
                        to={
                          from === FromEnum.blueprint
                            ? `/blueprints/${blueprint.id}/${cs.unique}`
                            : `/projects/${encodeName(blueprint.projectName)}/${cs.unique}`
                        }
                      >
                        Edit Data Scope and Scope Config
                      </Link>
                    </div>
                  </S.ConnectionItem>
                ))}
              </S.ConnectionList>
              <Buttons position="bottom" align="center">
                <Button intent={Intent.PRIMARY} text="Collect Data" onClick={handleRun} />
              </Buttons>
            </>
          )}
        </div>
      )}
      {blueprint.mode === ModeEnum.advanced && (
        <div className="block">
          <h3>JSON Configuration</h3>
          <AdvancedEditor value={rawPlan} onChange={setRawPlan} />
          <div className="btns">
            <Button
              intent={Intent.PRIMARY}
              text="Save"
              onClick={() =>
                handleUpdate({
                  plan: !validRawPlan(rawPlan) ? JSON.parse(rawPlan) : JSON.stringify([[]], null, '  '),
                })
              }
            />
          </div>
        </div>
      )}
      {type === 'name' && (
        <UpdateNameDialog
          name={blueprint.name}
          operating={operating}
          onCancel={handleCancel}
          onSubmit={(name) => handleUpdate({ name })}
        />
      )}
      {type === 'policy' && (
        <UpdatePolicyDialog
          blueprint={blueprint}
          isManual={blueprint.isManual}
          cronConfig={blueprint.cronConfig}
          skipOnFail={blueprint.skipOnFail}
          timeAfter={blueprint.timeAfter}
          operating={operating}
          onCancel={handleCancel}
          onSubmit={(payload) => handleUpdate(payload)}
        />
      )}
      {type === 'add-connection' && (
        <AddConnectionDialog
          disabled={connections.map((cs) => cs.unique)}
          onCancel={handleCancel}
          onSubmit={(connection) =>
            handleUpdate({
              connections: [...blueprint.connections, connection],
            })
          }
        />
      )}
    </S.ConfigurationPanel>
  );
};
