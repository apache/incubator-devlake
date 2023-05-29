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
import type { TabId } from '@blueprintjs/core';
import { Tabs, Tab, Switch } from '@blueprintjs/core';

import { PageLoading } from '@/components';
import { useRefreshData } from '@/hooks';
import { operator } from '@/utils';

import { Configuration } from './panel/configuration';
import { Status } from './panel/status';
import * as API from './api';
import * as S from './styled';

interface Props {
  id: ID;
}

export const BlueprintDetail = ({ id }: Props) => {
  const [activeTab, setActiveTab] = useState<TabId>('configuration');
  const [version, setVersion] = useState(1);
  const [operating, setOperating] = useState(false);

  const { ready, data } = useRefreshData(
    async () => Promise.all([API.getBlueprint(id), API.getBlueprintPipelines(id)]),
    [version],
  );

  if (!ready || !data) {
    return <PageLoading />;
  }

  const [blueprint, pipelines] = data;

  const handleUpdate = async (payload: any, callback?: () => void) => {
    const [success] = await operator(
      () =>
        API.updateBlueprint(id, {
          ...blueprint,
          ...payload,
        }),
      {
        setOperating,
      },
    );

    if (success) {
      setVersion((v) => v + 1);
      callback?.();
    }
  };

  const handleRun = async () => {
    const [success] = await operator(() => API.runBlueprint(id), {
      setOperating,
    });

    if (success) {
      setVersion((v) => v + 1);
    }
  };

  return (
    <S.Wrapper>
      <Tabs selectedTabId={activeTab} onChange={(at) => setActiveTab(at)}>
        <Tab
          id="status"
          title="Status"
          panel={
            <Status blueprint={blueprint} pipelineId={pipelines?.[0]?.id} operating={operating} onRun={handleRun} />
          }
        />
        <Tab
          id="configuration"
          title="Configuration"
          panel={<Configuration blueprint={blueprint} operating={operating} onUpdate={handleUpdate} />}
        />
        <Tabs.Expander />
        <Switch
          style={{ marginBottom: 0 }}
          label="Blueprint Enabled"
          checked={blueprint.enable}
          onChange={(e) => handleUpdate({ enable: (e.target as HTMLInputElement).checked })}
        />
      </Tabs>
    </S.Wrapper>
  );
};
