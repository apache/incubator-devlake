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
import { Tabs, Tab } from '@blueprintjs/core';
import useUrlState from '@ahooksjs/use-url-state';

import { PageHeader, PageLoading } from '@/components';
import { useRefreshData } from '@/hooks';
import { BlueprintDetail, FromEnum } from '@/pages';

import { WebhooksPanel } from './webhooks-panel';
import { SettingsPanel } from './settings-panel';
import * as API from './api';
import * as S from './styled';

export const ProjectDetailPage = () => {
  const [version, setVersion] = useState(1);

  const { pname } = useParams() as { pname: string };
  const [query, setQuery] = useUrlState({ tabId: 'blueprint' });

  const { ready, data } = useRefreshData(() => Promise.all([API.getProject(pname)]), [pname, version]);

  const handleChangeTabId = (tabId: string) => {
    setQuery({ tabId });
  };

  const handleRefresh = () => {
    setVersion((v) => v + 1);
  };

  if (!ready || !data) {
    return <PageLoading />;
  }

  const [project] = data;

  return (
    <PageHeader
      breadcrumbs={[
        { name: 'Projects', path: '/projects' },
        { name: project.name, path: `/projects/${pname}` },
      ]}
    >
      <S.Wrapper>
        <Tabs selectedTabId={query.tabId} onChange={handleChangeTabId}>
          <Tab
            id="blueprint"
            title="Blueprint"
            panel={<BlueprintDetail id={project.blueprint.id} from={FromEnum.project} />}
          />
          <Tab id="webhook" title="Webhooks" panel={<WebhooksPanel project={project} onRefresh={handleRefresh} />} />
          <Tab id="settings" title="Settings" panel={<SettingsPanel project={project} onRefresh={handleRefresh} />} />
        </Tabs>
      </S.Wrapper>
    </PageHeader>
  );
};
