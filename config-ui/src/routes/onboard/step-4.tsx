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

import { useState, useContext, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons';
import { theme, Progress, Space, Button, Modal } from 'antd';
import styled from 'styled-components';

import API from '@/api';
import { useAutoRefresh } from '@/hooks';
import { ConnectionName, ConnectionForm } from '@/plugins';
import { operator } from '@/utils';

import { Logs } from './components';
import { Context } from './context';

const Wrapper = styled.div`
  margin-top: 150px;
  padding: 24px;
  background-color: #fff;
  box-shadow: 0px 2.4px 4.8px -0.8px rgba(0, 0, 0, 0.1), 0px 1.6px 8px 0px rgba(0, 0, 0, 0.07);

  .top {
    margin-bottom: 42px;
    text-align: center;

    .info {
      margin-bottom: 34px;
      font-size: 16px;
      font-weight: 600;
    }

    .tip {
      margin-bottom: 42px;
      font-size: 12px;
      color: #818388;
    }

    .action {
      margin-top: 30px;
    }
  }

  .logs {
    .tip {
      font-size: 12px;
      font-weight: 600;
      color: #70727f;
    }

    .detail {
      display: flex;
      margin-top: 12px;

      & > div {
        flex: 1;
      }
    }
  }
`;

const DashboardURLMap: Record<string, string> = {
  github: '/grafana/d/KXWvOFQnz/github?orgId=1&var-repo_id=All&var-interval=WEEKDAY',
  gitlab: '/grafana/d/msSjEq97z/gitlab?orgId=1&var-repo_id=All&var-interval=WEEKDAY',
  bitbucket: '/grafana/d/4LzQHZa4k/bitbucket?orgId=1&var-repo_id=All&var-interval=WEEKDAY',
  azuredevops:
    '/grafana/d/ba7e3a95-80ed-4067-a54b-2a82758eb3dd/azure-devops?orgId=1&var-repo_id=All&var-interval=WEEKDAY',
};

const getStatus = (data: any) => {
  if (!data) {
    return 'running';
  }

  switch (data.status) {
    case 'TASK_COMPLETED':
      return 'success';
    case 'TASK_PARTIAL':
      return 'partial';
    case 'TASK_FAILED':
      return 'failed';
    case 'TASK_RUNNING':
    default:
      return 'running';
  }
};

export const Step4 = () => {
  const [operating, setOperating] = useState(false);
  const [open, setOpen] = useState(false);

  const navigate = useNavigate();

  const { step, records, projectName, plugin } = useContext(Context);

  const record = useMemo(() => records.find((it) => it.plugin === plugin), [plugin, records]);

  const { data } = useAutoRefresh(
    async () => {
      return await API.pipeline.subTasks(record?.pipelineId as string);
    },
    [],
    {
      cancel: (data) => {
        return !!(data && ['TASK_COMPLETED', 'TASK_PARTIAL', 'TASK_FAILED'].includes(data.status));
      },
    },
  );

  const [status, percent, collector, extractor] = useMemo(() => {
    const status = getStatus(data);
    const percent = Math.floor((data?.completionRate ?? 0) * 100);

    const collectorTask = (data?.subtasks ?? [])[0] ?? {};
    const extractorTask = (data?.subtasks ?? [])[1] ?? {};

    const collector = {
      plugin: collectorTask.plugin,
      scopeName: collectorTask.options?.name,
      status: collectorTask.status,
      tasks: (collectorTask.subtaskDetails ?? [])
        .filter((it) => it.isCollector)
        .map((it) => ({
          step: it.sequence,
          name: it.name,
          status: it.isFailed ? 'failed' : !it.beganAt ? 'pending' : it.finishedAt ? 'success' : 'running',
          finishedRecords: it.finishedRecords,
        })),
    };

    const extractor = {
      plugin: extractorTask.plugin,
      scopeName: extractorTask.options?.name,
      status: extractorTask.status,
      tasks: (extractorTask.subtaskDetails ?? [])
        .filter((it) => it.isCollector)
        .map((it) => ({
          step: it.sequence,
          name: it.name,
          status: it.isFailed ? 'failed' : !it.beganAt ? 'pending' : it.finishedAt ? 'success' : 'running',
          finishedRecords: it.finishedRecords,
        })),
    };

    return [status, percent, collector, extractor];
  }, [data]);

  const {
    token: { green5, orange5, red5, colorPrimary },
  } = theme.useToken();

  const handleFinish = async () => {
    const [success] = await operator(
      () =>
        API.store.set('onboard', {
          step,
          records,
          done: true,
          projectName,
          plugin,
        }),
      {
        setOperating,
      },
    );

    if (success) {
      navigate('/');
    }
  };

  if (!plugin || !record) {
    return null;
  }

  const { connectionId, scopeName } = record;

  return (
    <Wrapper>
      {status === 'running' && (
        <div className="top">
          <div className="info">Syncing up data from {scopeName}...</div>
          <div className="tip">
            This may take a few minutes to hours, depending on the size of your data and rate limits of the tool you
            choose.
          </div>
          <Progress type="circle" size={120} percent={percent} />
        </div>
      )}
      {status === 'success' && (
        <div className="top">
          <div className="info">{scopeName} is successfully collected !</div>
          <CheckCircleOutlined style={{ fontSize: 120, color: green5 }} />
          <div className="action">
            <Space direction="vertical">
              <Button type="primary" onClick={() => window.open(DashboardURLMap[plugin])}>
                Check Dashboard
              </Button>
              <Button loading={operating} onClick={handleFinish}>
                Finish and Exit
              </Button>
            </Space>
          </div>
        </div>
      )}
      {status === 'partial' && (
        <div className="top">
          <div className="info">{scopeName} is parted collected！</div>
          <CheckCircleOutlined style={{ fontSize: 120, color: orange5 }} />
          <div className="action">
            <Space>
              <Button type="primary">Re-collect Data</Button>
              <Button type="primary" onClick={() => window.open(DashboardURLMap[plugin])}>
                Check Dashboard
              </Button>
            </Space>
          </div>
        </div>
      )}
      {status === 'failed' && (
        <div className="top">
          <div className="info" style={{ marginBottom: 10 }}>
            Something went wrong with the collection process.
          </div>
          <div className="info">
            Please check out the{' '}
            <span style={{ color: colorPrimary, cursor: 'pointer' }} onClick={() => setOpen(true)}>
              network and token permission
            </span>{' '}
            and retry data collection
          </div>
          <CloseCircleOutlined style={{ fontSize: 120, color: red5 }} />
          <div className="action">
            <Space direction="vertical">
              <Button type="primary">Re-collect Data</Button>
            </Space>
          </div>
        </div>
      )}
      <div className="logs">
        <div className="tip">Data synchronization progress:</div>
        <div className="detail">
          <Logs log={collector} />
          <Logs log={extractor} style={{ marginLeft: 16 }} />
        </div>
      </div>
      <Modal
        open={open}
        width={820}
        centered
        title={<ConnectionName plugin={plugin} />}
        footer={null}
        onCancel={() => setOpen(false)}
      >
        <ConnectionForm plugin={plugin} connectionId={connectionId} onSuccess={() => setOpen(false)} />
      </Modal>
    </Wrapper>
  );
};
