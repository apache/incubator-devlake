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

import { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Tabs, message } from 'antd';
import axios from 'axios';
import API from '@/api';
import { PageLoading } from '@/components';
import { useRefreshData } from '@/hooks';
import { PATHS } from '@/config';

import { FromEnum } from '../types';

import { ConfigurationPanel } from './configuration-panel';
import { StatusPanel } from './status-panel';
import * as S from './styled';

interface Props {
  id: ID;
  from: FromEnum;
}

export const BlueprintDetail = ({ id, from }: Props) => {
  const [version, setVersion] = useState(1);
  const [activeKey, setActiveKey] = useState('status');

  const { state } = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    setActiveKey(state?.activeKey ?? 'status');
  }, [state]);

  const { ready, data, error } = useRefreshData(async () => {
    const [bpRes, pipelineRes] = await Promise.all([
      API.blueprint.get(id),
      API.blueprint.pipelines(id),
    ]);
    return [bpRes, pipelineRes.pipelines[0]];
  }, [version]);

  useEffect(() => {
     if (axios.isAxiosError(error) && error.response?.status === 404) {
      message.error(`Blueprint not found with id: ${id}`);
      navigate(PATHS.BLUEPRINTS(), { replace: true });
    }
  }, [error, navigate, id]);

  const handlRefresh = () => {
    setVersion((v) => v + 1);
  };

  const handleChangeActiveKey = (activeKey: string) => {
    setActiveKey(activeKey);
  };

  if (!ready || !data) {
    return <PageLoading />;
  }

  const [blueprint, lastPipeline] = data;

  return (
    <S.Wrapper>
      <Tabs
        centered
        items={[
          {
            key: 'status',
            label: 'Status',
            children: (
              <StatusPanel from={from} blueprint={blueprint} pipelineId={lastPipeline?.id} onRefresh={handlRefresh} />
            ),
          },
          {
            key: 'configuration',
            label: 'Configuration',
            children: (
              <ConfigurationPanel
                from={from}
                blueprint={blueprint}
                onRefresh={handlRefresh}
                onChangeTab={handleChangeActiveKey}
              />
            ),
          },
        ]}
        activeKey={activeKey}
        onChange={handleChangeActiveKey}
      />
    </S.Wrapper>
  );
};
