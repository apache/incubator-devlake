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
import { SmileFilled, CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons';
import { theme, Progress, Space, Button, Modal } from 'antd';
import styled from 'styled-components';

import API from '@/api';
import { useAutoRefresh } from '@/hooks';
import { ConnectionName, ConnectionForm } from '@/plugins';

import { Logs } from './components';
import { Context } from './context';

const Top = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 100px;
  margin-bottom: 24px;
  height: 70px;

  span.text {
    margin-left: 8px;
    font-size: 20px;
  }
`;

const Content = styled.div`
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
`;

const LogsWrapper = styled.div`
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
`;

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
  const [open, setOpen] = useState(false);

  const { records, plugin } = useContext(Context);

  const record = useMemo(() => records.find((it) => it.plugin === plugin), [plugin, records]);

  const { data } = useAutoRefresh(
    async () => {
      const taskRes = await API.pipeline.subTasks(record?.pipelineId as string);
      return taskRes.tasks;
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
    const percent = (data?.completionRate ?? 0) * 100;

    const collectorTask = (data?.subtasks ?? [])[0] ?? {};
    const extractorTask = (data?.subtasks ?? [])[1] ?? {};

    const collector = {
      plugin: collectorTask.plugin,
      scopeName: collectorTask.option?.name,
      status: collectorTask.status,
      tasks: (collectorTask.subtaskDetails ?? [])
        .filter((it: any) => it.is_collector === '1')
        .map((it: any) => ({
          step: it.sequence,
          name: it.name,
          status: it.is_failed === '1' ? 'failed' : !it.began_at ? 'pending' : it.finished_at ? 'success' : 'running',
          finishedRecords: it.finished_records,
        })),
    };

    const extractor = {
      plugin: extractorTask.plugin,
      scopeName: extractorTask.option?.name,
      status: extractorTask.status,
      tasks: (extractorTask.subtaskDetails ?? [])
        .filter((it: any) => it.is_collector === '1')
        .map((it: any) => ({
          step: it.sequence,
          name: it.name,
          status: it.is_failed === '1' ? 'failed' : !it.began_at ? 'pending' : it.finished_at ? 'success' : 'running',
          finishedRecords: it.finished_records,
        })),
    };

    return [status, percent, collector, extractor];
  }, [data]);

  const {
    token: { green5, orange5, red5 },
  } = theme.useToken();

  if (!plugin || !record) {
    return null;
  }

  const { connectionId, scopeName } = record;

  return (
    <>
      <Top>
        <SmileFilled style={{ fontSize: 36, color: green5 }} />
        <span className="text">Congratulations！You have successfully connected to your first repository!</span>
      </Top>
      <Content>
        {status === 'running' && (
          <div className="top">
            <div className="info">syncing up data from {scopeName}...</div>
            <div className="tip">
              This may take a few minutes to hours, depending on the size of your data and rate limits of the tool you
              choose. Exit
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
                <Button type="primary">Check Dashboard</Button>
                <Button type="link">finish</Button>
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
                <Button type="primary">Check Dashboard</Button>
              </Space>
            </div>
          </div>
        )}
        {status === 'failed' && (
          <div className="top">
            <div className="info">Something went wrong with the collection process. </div>
            <div className="info">
              Please check out the
              <Button type="link" onClick={() => setOpen(true)}>
                network and token permission
              </Button>
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
        <LogsWrapper>
          <div className="tip">Sync progress details</div>
          <div className="detail">
            <Logs log={collector} />
            <Logs log={extractor} style={{ marginLeft: 16 }} />
          </div>
        </LogsWrapper>
      </Content>
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
    </>
  );
};
