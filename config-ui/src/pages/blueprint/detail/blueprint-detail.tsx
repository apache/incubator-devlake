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
import { useHistory } from 'react-router-dom';
import type { TabId } from '@blueprintjs/core';
import { Tabs, Tab, Switch, Button, Icon, Intent } from '@blueprintjs/core';

import { PageLoading, Dialog } from '@/components';
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
  const [isOpen, setIsOpen] = useState(false);

  const history = useHistory();

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

  const handleShowDeleteDialog = () => {
    setIsOpen(true);
  };

  const handleHideDeleteDialog = () => {
    setIsOpen(false);
  };

  const handleDelete = async () => {
    const [success] = await operator(() => API.deleteBluprint(id), {
      setOperating,
      formatMessage: () => 'Delete blueprint successful.',
    });

    if (success) {
      history.push('/blueprints');
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
        <Button intent={Intent.DANGER} text="Delete Blueprint" onClick={handleShowDeleteDialog} />
      </Tabs>
      <Dialog
        isOpen={isOpen}
        style={{ width: 820 }}
        title="Are you sure you want to delete this Blueprint?"
        okText="Confirm"
        okLoading={operating}
        onCancel={handleHideDeleteDialog}
        onOk={handleDelete}
      >
        <S.DialogBody>
          <Icon icon="warning-sign" />
          <span>
            Please note: deleting the Blueprint will not delete the historical data of the Data Scopes in this
            Blueprint. If you would like to delete the historical data of Data Scopes, please visit the Connection page
            and do so.
          </span>
        </S.DialogBody>
      </Dialog>
    </S.Wrapper>
  );
};
