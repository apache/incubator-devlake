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
import { Tabs } from 'antd';
import useUrlState from '@ahooksjs/use-url-state';

import API from '@/api';
import { PageLoading } from '@/components';
import { useRefreshData } from '@/hooks';

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

  const [query, setQuery] = useUrlState({ tab: 'status' });

  const { ready, data } = useRefreshData(async () => {
    const [bpRes, pipelineRes] = await Promise.all([API.blueprint.get(id), API.blueprint.pipelines(id)]);
    return [bpRes, pipelineRes.pipelines[0]];
  }, [version]);

  const handlRefresh = () => {
    setVersion((v) => v + 1);
  };

  const handleChangeTab = (tab: string) => {
    setQuery({ tab });
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
                onChangeTab={handleChangeTab}
              />
            ),
          },
        ]}
        activeKey={query.tab}
        onChange={(tab) => setQuery({ tab })}
      />
    </S.Wrapper>
  );
};
