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
import { useParams } from 'react-router-dom';
import { Helmet } from 'react-helmet';
import { Tabs } from 'antd';
import useUrlState from '@ahooksjs/use-url-state';

import API from '@/api';
import { PageHeader, PageLoading } from '@/components';
import { PATHS } from '@/config';
import { useRefreshData } from '@/hooks';
import { BlueprintDetail, FromEnum } from '@/routes';

import { WebhooksPanel } from './webhooks-panel';
import { SettingsPanel } from './settings-panel';
import * as S from './styled';

const brandName = import.meta.env.DEVLAKE_BRAND_NAME ?? 'DevLake';

export const ProjectDetailPage = () => {
  const [version, setVersion] = useState(1);

  const { pname } = useParams() as { pname: string };
  const [query, setQuery] = useUrlState({ tabId: 'blueprint' });

  const { ready, data } = useRefreshData(() => API.project.get(pname), [pname, version]);

  const handleChangeTabId = (tabId: string) => {
    setQuery({ tabId });
  };

  const handleRefresh = () => {
    setVersion((v) => v + 1);
  };

  if (!ready || !data) {
    return <PageLoading />;
  }

  return (
    <PageHeader
      breadcrumbs={[
        { name: 'Projects', path: PATHS.PROJECTS() },
        { name: data.name, path: PATHS.PROJECT(pname) },
      ]}
    >
      <Helmet>
        <title>
          {data.name} - {brandName}
        </title>
      </Helmet>
      <S.Wrapper>
        <Tabs
          items={[
            {
              key: 'blueprint',
              label: 'Blueprint',
              children: <BlueprintDetail id={data.blueprint.id} from={FromEnum.project} />,
            },
            {
              key: 'webhook',
              label: 'Webhooks',
              children: <WebhooksPanel project={data} onRefresh={handleRefresh} />,
            },
            {
              key: 'settings',
              label: 'Settings',
              children: <SettingsPanel project={data} onRefresh={handleRefresh} />,
            },
          ]}
          activeKey={query.tabId}
          onChange={handleChangeTabId}
        />
      </S.Wrapper>
    </PageHeader>
  );
};
