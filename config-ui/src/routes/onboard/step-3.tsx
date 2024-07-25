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

import { useState, useContext, useEffect, useMemo } from 'react';
import { Flex, Button } from 'antd';
import dayjs from 'dayjs';

import API from '@/api';
import { Markdown } from '@/components';
import { DataScopeRemote, getPluginScopeId } from '@/plugins';
import { operator, formatTime } from '@/utils';

import { Context } from './context';
import * as S from './styled';

export const Step3 = () => {
  const [QA, setQA] = useState('');
  const [operating, setOperating] = useState(false);
  const [scopes, setScopes] = useState<any[]>([]);

  const { step, records, done, projectName, plugin, setStep, setRecords } = useContext(Context);

  useEffect(() => {
    fetch(`/onboard/step-3/${plugin}.md`)
      .then((res) => res.text())
      .then((text) => setQA(text));
  }, [plugin]);

  const connectionId = useMemo(() => {
    const record = records.find((it) => it.plugin === plugin);
    return record?.connectionId ?? null;
  }, [plugin, records]);

  const handleSubmit = async () => {
    if (!projectName || !plugin || !connectionId) {
      return;
    }

    const [success] = await operator(
      async () => {
        // 1. create a new project
        const { blueprint } = await API.project.create({
          name: projectName,
          description: '',
          metrics: [
            {
              pluginName: 'dora',
              pluginOption: '',
              enable: true,
            },
          ],
        });

        // 2. add data scopes to the connection
        await API.scope.batch(plugin, connectionId, { data: scopes.map((it) => it.data) });

        // 3. add data scopes to the blueprint
        await API.blueprint.update(blueprint.id, {
          timeAfter: formatTime(dayjs().subtract(14, 'day').startOf('day').toDate(), 'YYYY-MM-DD[T]HH:mm:ssZ'),
          connections: [
            {
              pluginName: plugin,
              connectionId,
              scopes: scopes.map((it) => ({
                scopeId: getPluginScopeId(plugin, it.data),
              })),
            },
          ],
        });

        // 4. trigger this blueprint
        await API.blueprint.trigger(blueprint.id, { skipCollectors: false, fullSync: false });

        // 5. get current run pipeline
        const pipeline = await API.blueprint.pipelines(blueprint.id);

        const newRecords = records.map((it) =>
          it.plugin !== plugin
            ? it
            : {
                ...it,
                blueprintId: blueprint.id,
                pipelineId: pipeline.pipelines[0].id,
                scopeName: scopes[0]?.fullName ?? scopes[0].name,
              },
        );

        setRecords(newRecords);

        // 6. update store
        await API.store.set('onboard', {
          step: 4,
          records: newRecords,
          done,
          projectName,
          plugin,
        });
      },
      {
        setOperating,
        formatMessage: () => 'Congratulations！You have successfully connected to your first repository!',
      },
    );

    if (success) {
      setStep(step + 1);
    }
  };

  if (!plugin || !connectionId) {
    return null;
  }

  return (
    <>
      <S.StepContent>
        <div className="content">
          <DataScopeRemote
            mode="single"
            plugin={plugin}
            connectionId={connectionId}
            selectedScope={scopes}
            onChangeSelectedScope={setScopes}
            footer={null}
          />
        </div>
        <Markdown className="qa">{QA}</Markdown>
      </S.StepContent>
      <Flex style={{ marginTop: 36 }} justify="space-between">
        <Button ghost type="primary" loading={operating} onClick={() => setStep(step - 1)}>
          Previous Step
        </Button>
        <Button type="primary" loading={operating} disabled={!scopes.length} onClick={handleSubmit}>
          Next Step
        </Button>
      </Flex>
    </>
  );
};
