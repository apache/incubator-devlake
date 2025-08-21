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
import { FormOutlined, PlusOutlined } from '@ant-design/icons';
import { Flex, Table, Button } from 'antd';

import API from '@/api';
import { NoData } from '@/components';
import { getCron, PATHS } from '@/config';
import { ConnectionName } from '@/features';
import { getPluginConfig } from '@/plugins';
import { IBlueprint, IBPMode } from '@/types';
import { formatTime, operator } from '@/utils';

import { FromEnum } from '../types';
import { validRawPlan } from '../utils';

import { AdvancedEditor, UpdateNameDialog, UpdatePolicyDialog, AddConnectionDialog } from './components';
import * as S from './styled';

interface Props {
  from: FromEnum;
  blueprint: IBlueprint;
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

  const connections = useMemo(
    () =>
      blueprint.connections
        .filter((cs) => cs.pluginName !== 'webhook')
        .map((cs: any) => {
          const plugin = getPluginConfig(cs.pluginName);
          return {
            plugin: plugin.plugin,
            connectionId: cs.connectionId,
            icon: plugin.icon,
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
        API.blueprint.update(blueprint.id, {
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
    const [success] = await operator(
      () => API.blueprint.trigger(blueprint.id, { skipCollectors: false, fullSync: false }),
      {
        setOperating,
        formatMessage: () => 'Trigger blueprint successful.',
      },
    );

    if (success) {
      onRefresh();
      onChangeTab('status');
    }
  };

  return (
    <S.ConfigurationPanel>
      <div className="block">
        <h3>Blueprint Name</h3>
        <span>{blueprint.name}</span>
        <Button type="link" icon={<FormOutlined />} onClick={handleShowNameDialog} />
      </div>
      <div className="block">
        <h3>
          <span>Sync Policy</span>
          <Button type="link" icon={<FormOutlined />} onClick={handleShowPolicyDialog} />
        </h3>
        <Table
          rowKey="id"
          size="middle"
          columns={[
            {
              title: 'Data Time Range',
              dataIndex: 'timeRange',
              key: 'timeRange',
              render: (val) => (blueprint.mode === IBPMode.NORMAL ? `${formatTime(val)} to Now` : 'N/A'),
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
              id: blueprint.id,
              timeRange: blueprint.timeAfter,
              frequency: blueprint.cronConfig,
              isManual: blueprint.isManual,
              skipFailed: blueprint.skipOnFail,
            },
          ]}
          pagination={false}
        />
      </div>
      {blueprint.mode === IBPMode.NORMAL && (
        <div className="block">
          <h3>Data Connections</h3>
          {!connections.length ? (
            <NoData
              text={
                <>
                  If you have not created data connections yet, please{' '}
                  <Link to={PATHS.CONNECTIONS()}>create connections</Link> first and then add them to the project.
                </>
              }
              action={
                <Button type="primary" icon={<PlusOutlined />} onClick={handleShowAddConnectionDialog}>
                  Add a Connection
                </Button>
              }
            />
          ) : (
            <Flex vertical gap="middle">
              <Flex>
                <Button type="primary" icon={<PlusOutlined />} onClick={handleShowAddConnectionDialog}>
                  Add a Connection
                </Button>
              </Flex>
              <S.ConnectionList>
                {connections.map((cs) => (
                  <S.ConnectionItem key={`${cs.plugin}-${cs.connectionId}`}>
                    <ConnectionName plugin={cs.plugin} connectionId={cs.connectionId} />
                    <div className="count">
                      <span>{cs.scope.length} data scope</span>
                    </div>
                    <div className="link">
                      <Link
                        to={
                          from === FromEnum.blueprint
                            ? PATHS.BLUEPRINT_CONNECTION(blueprint.id, cs.plugin, cs.connectionId)
                            : PATHS.PROJECT_CONNECTION(blueprint.projectName, cs.plugin, cs.connectionId)
                        }
                      >
                        Edit Data Scope and Scope Config
                      </Link>
                    </div>
                  </S.ConnectionItem>
                ))}
              </S.ConnectionList>
              <Flex justify="center">
                <Button type="primary" disabled={!blueprint.enable} onClick={handleRun}>
                  Collect Data
                </Button>
              </Flex>
            </Flex>
          )}
        </div>
      )}
      {blueprint.mode === IBPMode.ADVANCED && (
        <div className="block">
          <h3>JSON Configuration</h3>
          <AdvancedEditor value={rawPlan} onChange={setRawPlan} />
          <div className="btns">
            <Button
              type="primary"
              onClick={() =>
                handleUpdate({
                  plan: !validRawPlan(rawPlan) ? JSON.parse(rawPlan) : JSON.stringify([[]], null, '  '),
                })
              }
            >
              Save
            </Button>
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
          disabled={connections.map((cs) => `${cs.plugin}-${cs.connectionId}`)}
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
