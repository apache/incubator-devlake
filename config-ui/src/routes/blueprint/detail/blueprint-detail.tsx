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

import { useState, useEffect, useReducer } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { WarningOutlined } from '@ant-design/icons';
import { theme, Tabs, Modal } from 'antd';

import API from '@/api';
import { PageLoading } from '@/components';
import { useRefreshData } from '@/hooks';
import { operator } from '@/utils';

import { FromEnum } from '../types';

import { ConnectionCheck } from './components';
import { ConfigurationPanel } from './configuration-panel';
import { StatusPanel } from './status-panel';
import * as S from './styled';

type ConnectionFailed = {
  open: boolean;
  failedList?: Array<{ plugin: string; connectionId: ID }>;
};

function reducer(state: ConnectionFailed, action: { type: string; failedList?: ConnectionFailed['failedList'] }) {
  switch (action.type) {
    case 'open':
      return { open: true, failedList: action.failedList };
    case 'close':
      return { open: false, failedList: [] };
    default:
      return state;
  }
}

interface Props {
  id: ID;
  from: FromEnum;
}

export const BlueprintDetail = ({ id, from }: Props) => {
  const [version, setVersion] = useState(1);
  const [activeKey, setActiveKey] = useState('status');
  const [operating, setOperating] = useState(false);

  const [{ open, failedList }, dispatch] = useReducer(reducer, {
    open: false,
  });

  const { state } = useLocation();
  const navigate = useNavigate();

  const {
    token: { orange5 },
  } = theme.useToken();

  useEffect(() => {
    setActiveKey(state?.activeKey ?? 'status');
  }, [state]);

  const { ready, data } = useRefreshData(async () => {
    const [bpRes, pipelineRes] = await Promise.all([API.blueprint.get(id), API.blueprint.pipelines(id)]);
    return [bpRes, pipelineRes.pipelines[0]];
  }, [version]);

  const handleDelete = async () => {
    const [success] = await operator(() => API.blueprint.remove(blueprint.id), {
      setOperating,
      formatMessage: () => 'Delete blueprint successful.',
    });

    if (success) {
      navigate('/advanced/blueprints');
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
      setVersion((v) => v + 1);
    }
  };

  const handleTrigger = async (payload?: { skipCollectors?: boolean; fullSync?: boolean }) => {
    const { skipCollectors, fullSync } = payload ?? { skipCollectors: false, fullSync: false };

    if (!skipCollectors) {
      const [success, res] = await operator(() => API.blueprint.connectionsTokenCheck(blueprint.id), {
        hideToast: true,
        setOperating,
      });

      if (success && res.length) {
        const connectionFailed = res
          .filter((it: any) => !it.success)
          .map((it: any) => {
            return {
              plugin: it.pluginName,
              connectionId: it.connectionId,
            };
          });

        if (connectionFailed.length) {
          dispatch({ type: 'open', failedList: connectionFailed });
          return;
        }
      }
    }

    const [success] = await operator(() => API.blueprint.trigger(blueprint.id, { skipCollectors, fullSync }), {
      setOperating,
      formatMessage: () => 'Trigger blueprint successful.',
    });

    if (success) {
      setVersion((v) => v + 1);
      setActiveKey('status');
    }
  };

  if (!ready || !data) {
    return <PageLoading />;
  }

  const [blueprint, lastPipeline] = data;

  return (
    <S.Wrapper>
      <Tabs
        items={[
          {
            key: 'status',
            label: 'Status',
            children: (
              <StatusPanel
                from={from}
                blueprint={blueprint}
                pipelineId={lastPipeline?.id}
                operating={operating}
                onDelete={handleDelete}
                onUpdate={handleUpdate}
                onTrigger={handleTrigger}
              />
            ),
          },
          {
            key: 'configuration',
            label: 'General Settings',
            children: (
              <ConfigurationPanel
                from={from}
                blueprint={blueprint}
                operating={operating}
                onUpdate={handleUpdate}
                onTrigger={handleTrigger}
              />
            ),
          },
        ]}
        activeKey={activeKey}
        onChange={setActiveKey}
      />
      <Modal
        open={open}
        title={
          <>
            <WarningOutlined style={{ marginRight: 8, fontSize: 20, color: orange5 }} />
            <span>Invalid Token(s) Detected</span>
          </>
        }
        width={820}
        footer={null}
        onCancel={() => dispatch({ type: 'close' })}
      >
        <p>There are invalid tokens in the following connections. Please update them before re-syncing the data.</p>
        <ul>
          {(failedList ?? []).map(({ plugin, connectionId }) => (
            <li key={`${plugin}-${connectionId}`}>
              <ConnectionCheck plugin={plugin} connectionId={connectionId} />
            </li>
          ))}
        </ul>
      </Modal>
    </S.Wrapper>
  );
};
