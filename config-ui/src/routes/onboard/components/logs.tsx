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

import { LoadingOutlined, CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons';
import { theme, Tooltip, Progress } from 'antd';
import styled from 'styled-components';

const Wrapper = styled.div`
  padding: 10px 20px;
  font-size: 12px;
  color: #70727f;
  background: #f6f6f8;

  .title {
    display: flex;
    align-items: center;
    justify-content: space-between;

    & > span.name {
      width: 220px;
      font-weight: 600;
      white-space: nowrap;
      text-overflow: ellipsis;
      overflow: hidden;
    }

    & > span.progress {
      margin-left: 12px;
      flex: auto;
    }
  }

  ul {
    margin-top: 12px;
  }

  li {
    display: flex;
    margin-top: 6px;
    position: relative;

    &:first-child {
      margin-top: 0;
    }
  }

  span.name {
    flex: auto;
  }

  span.status {
    flex: 0 0 150px;
  }

  span.anticon {
    position: absolute;
    right: -15px;
  }
`;

const getStatus = (task: { step: number; name: string; status: string; finishedRecords: number }) => {
  if (task.status === 'pending') {
    return 'Pending';
  }

  if (task.status === 'running' && task.name === 'Clone Git Repo') {
    return 'N/A';
  }

  if (task.status === 'success' && task.name === 'Clone Git Repo') {
    return 'Completed';
  }

  if (task.status === 'failed' && task.name === 'Clone Git Repo') {
    return 'Failed';
  }

  if (['running', 'success', 'failed'].includes(task.status)) {
    return `Records collected: ${task.finishedRecords}`;
  }
};

interface LogsProps {
  style?: React.CSSProperties;
  log: {
    plugin: string;
    name: string;
    percent: number;
    tasks: Array<{
      step: number;
      name: string;
      status: string;
      finishedRecords: number;
    }>;
  };
}

export const Logs = ({ style, log: { plugin, name, percent, tasks } }: LogsProps) => {
  const {
    token: { green5, red5, colorPrimary },
  } = theme.useToken();

  if (!plugin) {
    return null;
  }

  return (
    <Wrapper style={style}>
      <div className="title">
        <Tooltip title={name}>
          <span className="name">{name}</span>
        </Tooltip>
        <span className="progress">
          <Progress size="small" percent={percent} showInfo={false} />
        </span>
      </div>
      <ul>
        {tasks.map((task) => (
          <li>
            <span className="name">
              Step {task.step} - {task.name}
            </span>
            <span className="status">{getStatus(task)}</span>
            {task.status === 'running' && <LoadingOutlined style={{ color: colorPrimary }} />}
            {task.status === 'success' && <CheckCircleOutlined style={{ color: green5 }} />}
            {task.status === 'failed' && <CloseCircleOutlined style={{ color: red5 }} />}
          </li>
        ))}
      </ul>
    </Wrapper>
  );
};
