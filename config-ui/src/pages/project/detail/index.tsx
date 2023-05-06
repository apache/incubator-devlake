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

import { useState, useEffect } from 'react';
import { useParams, useLocation, useHistory } from 'react-router-dom';
import { Tabs, Tab } from '@blueprintjs/core';

import { PageHeader, PageLoading } from '@/components';

import { useProject } from './use-project';
import { BlueprintPanel, IncomingWebhooksPanel, SettingsPanel } from './panel';
import * as S from './styled';

export const ProjectDetailPage = () => {
  const [tabId, setTabId] = useState('bp');

  const { pname } = useParams<{ pname: string }>();
  const { search } = useLocation();
  const history = useHistory();

  const query = new URLSearchParams(search);
  const urlTabId = query.get('tabId');

  useEffect(() => {
    setTabId(urlTabId ?? 'bp');
  }, [urlTabId]);

  const { loading, project, saving, onUpdate, onSelectWebhook, onCreateWebhook, onDeleteWebhook } = useProject(pname);

  const handleChangeTabId = (tabId: string) => {
    query.delete('tabId');
    query.append('tabId', tabId);
    history.push({ search: query.toString() });
  };

  if (loading || !project) {
    return <PageLoading />;
  }

  return (
    <PageHeader
      breadcrumbs={[
        { name: 'Projects', path: '/projects' },
        { name: project.name, path: `/projects/${pname}` },
      ]}
    >
      <S.Wrapper>
        <Tabs selectedTabId={tabId} onChange={handleChangeTabId}>
          <Tab id="bp" title="Blueprint" panel={<BlueprintPanel project={project} />} />
          <Tab
            id="iw"
            title="Incoming Webhooks"
            disabled={!project.blueprint}
            panel={
              <IncomingWebhooksPanel
                project={project}
                saving={saving}
                onSelectWebhook={onSelectWebhook}
                onCreateWebhook={onCreateWebhook}
                onDeleteWebhook={onDeleteWebhook}
              />
            }
          />
          <Tab id="st" title="Settings" panel={<SettingsPanel project={project} onUpdate={onUpdate} />} />
        </Tabs>
      </S.Wrapper>
    </PageHeader>
  );
};
