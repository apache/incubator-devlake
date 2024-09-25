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

import { useState, useMemo } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { MoreOutlined, DeleteOutlined, WarningOutlined } from '@ant-design/icons';
import { theme, Card, Modal, Switch, Button, Tooltip, Dropdown, Flex, Space } from 'antd';

import API from '@/api';
import { Message } from '@/components';
import { getCron } from '@/config';
import { selectAllConnections } from '@/features/connections';
import { useAppSelector, useRefreshData } from '@/hooks';
import { PipelineInfo, PipelineTasks, PipelineTable } from '@/routes/pipeline';
import { IBlueprint } from '@/types';
import { formatTime, operator } from '@/utils';

import { FromEnum } from '../types';

interface Props {
  from: FromEnum;
  blueprint: IBlueprint;
  pipelineId?: ID;
  onRefresh: () => void;
}

export const StatusPanel = ({ from, blueprint, pipelineId, onRefresh }: Props) => {
  const [type, setType] = useState<'delete' | 'fullSync' | 'checkTokenFailed'>();
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [operating, setOperating] = useState(false);
  const [connectionFailed, setConnectionFailed] = useState<
    Array<{ unique: string; name: string; plugin: string; connectionId: ID }>
  >([]);

  const navigate = useNavigate();

  const {
    token: { orange5 },
  } = theme.useToken();

  const connections = useAppSelector(selectAllConnections);

  const cron = useMemo(() => getCron(blueprint.isManual, blueprint.cronConfig), [blueprint]);

  const { ready, data } = useRefreshData(
    () => API.blueprint.pipelines(blueprint.id, { page, pageSize }),
    [blueprint.id, page, pageSize],
  );

  const handleResetType = () => {
    setType(undefined);
  };

  const handleRun = async ({
    skipCollectors = false,
    fullSync = false,
  }: {
    skipCollectors?: boolean;
    fullSync?: boolean;
  }) => {
    if (!skipCollectors) {
      const [success, res] = await operator(() => API.blueprint.connectionsTokenCheck(blueprint.id), {
        hideToast: true,
        setOperating,
      });

      if (success && res.length) {
        const connectionFailed = res
          .filter((it: any) => !it.success)
          .map((it: any) => {
            const unique = `${it.pluginName}-${it.connectionId}`;
            const connection = connections.find((c) => c.unique === unique);
            return {
              unique,
              name: connection?.name ?? '',
              plugin: it.pluginName,
              connectionId: it.connectionId,
            };
          });

        setType('checkTokenFailed');
        setConnectionFailed(connectionFailed);
        return;
      }
    }

    const [success] = await operator(() => API.blueprint.trigger(blueprint.id, { skipCollectors, fullSync }), {
      setOperating,
      formatMessage: () => 'Trigger blueprint successful.',
    });

    if (success) {
      onRefresh();
    }
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
        formatMessage: () =>
          from === FromEnum.project ? 'Update project successful.' : 'Update blueprint successful.',
      },
    );

    if (success) {
      onRefresh();
    }
  };

  const handleDelete = async () => {
    const [success] = await operator(() => API.blueprint.remove(blueprint.id), {
      setOperating,
      formatMessage: () => 'Delete blueprint successful.',
    });

    if (success) {
      navigate('/advanced/blueprints');
    }
  };

  return (
    <Flex vertical>
      {from === FromEnum.project && (
        <Flex justify="flex-end" align="center">
          <Space>
            <span>
              {cron.label === 'Manual' ? 'Manual' : `Next Run: ${formatTime(cron.nextTime, 'YYYY-MM-DD HH:mm')}`}
            </span>
            <Tooltip
              placement="top"
              title="It is recommended to re-transform your data in this project if you have updated the transformation of the data scope in this project."
            >
              <Button
                type="primary"
                disabled={!blueprint.enable}
                loading={operating}
                onClick={() => handleRun({ skipCollectors: true, fullSync: true })}
              >
                Re-transform Data
              </Button>
            </Tooltip>
            <Button type="primary" disabled={!blueprint.enable} loading={operating} onClick={() => handleRun({})}>
              Collect Data
            </Button>
            <Dropdown
              menu={{
                items: [
                  {
                    key: '1',
                    label: 'Collect Data in Full Refresh Mode',
                    disabled: !blueprint.enable,
                  },
                ],
                onClick: ({ key }) => {
                  if (key === '1') {
                    setType('fullSync');
                  }
                },
              }}
            >
              <Button icon={<MoreOutlined />} />
            </Dropdown>
          </Space>
        </Flex>
      )}

      {from === FromEnum.blueprint && (
        <Flex justify="center" align="center">
          <Space>
            <Button type="primary" disabled={!blueprint.enable} onClick={() => handleRun({})}>
              Run Now
            </Button>
            <Switch
              style={{ marginBottom: 0 }}
              disabled={!!blueprint.projectName}
              checked={blueprint.enable}
              onChange={(enable) => handleUpdate({ enable })}
            />
            Blueprint Enabled
            <Tooltip title="Delete Blueprint">
              <Button
                type="primary"
                loading={operating}
                disabled={!!blueprint.projectName}
                icon={<DeleteOutlined />}
                onClick={() => setType('delete')}
              />
            </Tooltip>
          </Space>
        </Flex>
      )}

      <Space direction="vertical" size="large">
        <h3>Current Pipeline</h3>

        {!pipelineId ? (
          <Card>There is no current run for this blueprint.</Card>
        ) : (
          <>
            <Card>
              <PipelineInfo id={pipelineId} />
            </Card>
            <Card>
              <PipelineTasks id={pipelineId} />
            </Card>
          </>
        )}

        <h3>Historical Pipelines</h3>

        {!data?.count ? (
          <Card>There are no historical runs associated with this blueprint.</Card>
        ) : (
          <PipelineTable
            loading={!ready}
            dataSource={data.pipelines}
            pagination={{
              current: page,
              pageSize,
              total: data.count,
              onChange: setPage,
            }}
          />
        )}
      </Space>

      {type === 'delete' && (
        <Modal
          open
          width={820}
          centered
          title="Are you sure you want to delete this Blueprint?"
          okText="Confirm"
          okButtonProps={{
            loading: operating,
          }}
          onCancel={handleResetType}
          onOk={handleDelete}
        >
          <Message
            content="Please note: deleting the Blueprint will not delete the historical data of the Data Scopes in this
              Blueprint. If you would like to delete the historical data of Data Scopes, please visit the Connection
              page and do so."
          />
        </Modal>
      )}

      {type === 'fullSync' && (
        <Modal
          open
          centered
          okText="Run Now"
          okButtonProps={{
            loading: operating,
          }}
          onCancel={handleResetType}
          onOk={() => handleRun({ fullSync: true })}
        >
          <Message content="This operation may take a long time as it will empty all of your existing data and re-collect it." />
        </Modal>
      )}

      {type === 'checkTokenFailed' && (
        <Modal
          open
          title={
            <>
              <WarningOutlined style={{ marginRight: 8, fontSize: 20, color: orange5 }} />
              <span>Invalid Token(s) Detected</span>
            </>
          }
          width={820}
          footer={null}
          onCancel={() => {
            handleResetType();
            setConnectionFailed([]);
          }}
        >
          <p>There are invalid tokens in the following connections. Please update them before re-syncing the data.</p>
          <ul style={{ paddingLeft: 20 }}>
            {connectionFailed.map((it) => (
              <li key={it.unique} style={{ listStyle: 'initial' }}>
                <Link to={`/connections/${it.plugin}/${it.connectionId}`} target="_blank">
                  {it.name}
                </Link>
              </li>
            ))}
          </ul>
        </Modal>
      )}
    </Flex>
  );
};
