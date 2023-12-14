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

import { useState, useEffect, useMemo } from 'react';
import { omit } from 'lodash';
import { Flex, Form, Input, Card, Alert, Divider, Select } from 'antd';
import { Button, Intent } from '@blueprintjs/core';

import API from '@/api';
import { ExternalLink, Block, Message } from '@/components';
import { transformEntities, EntitiesLabel } from '@/config';
import { getPluginConfig } from '@/plugins';
import { GitHubTransformation } from '@/plugins/register/github';
import { JiraTransformation } from '@/plugins/register/jira';
import { GitLabTransformation } from '@/plugins/register/gitlab';
import { JenkinsTransformation } from '@/plugins/register/jenkins';
import { BitbucketTransformation } from '@/plugins/register/bitbucket';
import { AzureTransformation } from '@/plugins/register/azure';
import { TapdTransformation } from '@/plugins/register/tapd';
import { BambooTransformation } from '@/plugins/register/bamboo';
import { operator } from '@/utils';

import { TIPS_MAP } from './misc';

interface Props {
  plugin: string;
  connectionId: ID;
  showWarning?: boolean;
  scopeId?: ID;
  scopeConfigId?: ID;
  onCancel: () => void;
  onSubmit: (trId: string) => void;
}

export const ScopeConfigForm = ({
  plugin,
  connectionId,
  showWarning = false,
  scopeId,
  scopeConfigId,
  onCancel,
  onSubmit,
}: Props) => {
  const [step, setStep] = useState(1);
  const [name, setName] = useState('');
  const [entities, setEntities] = useState<string[]>([]);
  const [transformation, setTransformation] = useState<any>({});
  const [hasError, setHasError] = useState(false);
  const [operating, setOperating] = useState(false);

  const config = useMemo(() => getPluginConfig(plugin), []);

  useEffect(() => {
    setTransformation(config.scopeConfig?.transformation ?? {});
  }, [config.scopeConfig?.transformation]);

  useEffect(() => {
    setEntities(config.scopeConfig?.entities ?? []);
  }, [config.scopeConfig?.entities]);

  useEffect(() => {
    if (!scopeConfigId) return;

    (async () => {
      try {
        const res = await API.scopeConfig.get(plugin, connectionId, scopeConfigId);
        setName(res.name);
        setEntities(res.entities ?? []);
        setTransformation(omit(res, ['id', 'connectionId', 'name', 'entities', 'createdAt', 'updatedAt']));
      } catch {}
    })();
  }, [scopeConfigId]);

  const handlePrevStep = () => {
    setStep(1);
  };

  const handleNextStep = () => {
    setStep(2);
  };

  const handleSubmit = async () => {
    const [success, res] = await operator(
      () =>
        !scopeConfigId
          ? API.scopeConfig.create(plugin, connectionId, { name, entities, ...transformation })
          : API.scopeConfig.update(plugin, connectionId, scopeConfigId, { name, entities, ...transformation }),
      {
        setOperating,
        formatMessage: () => (!scopeConfigId ? 'Create scope config successful.' : 'Update scope config successful'),
      },
    );

    if (success) {
      onSubmit(res.id);
    }
  };

  return (
    <Flex vertical gap="middle">
      {TIPS_MAP[plugin] && (
        <Alert
          message={
            <>
              To learn about how {TIPS_MAP[plugin].name} transformation is used in DevLake,{' '}
              <ExternalLink link={TIPS_MAP[plugin].link}>check out this doc</ExternalLink>.
            </>
          }
        />
      )}
      {step === 1 && (
        <>
          <Card>
            <Block
              title="Scope Config Name"
              description="Give this Scope Config a unique name so that you can identify it in the future."
              required
            >
              <Input placeholder="My Scope Config 1" value={name} onChange={(e) => setName(e.target.value)} />
            </Block>
          </Card>
          <Card>
            <Block
              title="Data Entities"
              description={
                <>
                  Select the data entities you wish to collect for the Data Scope.
                  <ExternalLink link="">Learn about data entities</ExternalLink>
                </>
              }
              required
            >
              <Select
                style={{ width: '100%' }}
                mode="multiple"
                options={transformEntities(config.scopeConfig?.entities ?? [])}
                value={entities}
                onChange={(value) => setEntities(value)}
              />
            </Block>
            {showWarning && (
              <Message
                content="Please note: if you edit Data Entities and expect to see the Dashboards updated, you will need to visit
                  the Project page of the Data Scope that has been associated with this Scope Config and click on “Collect
                  All Data”."
              />
            )}
          </Card>
          <Flex justify="flex-end" gap="small">
            <Button outlined intent={Intent.PRIMARY} text="Cancel" onClick={onCancel} />
            <Button disabled={!name || !entities.length} intent={Intent.PRIMARY} text="Next" onClick={handleNextStep} />
          </Flex>
        </>
      )}
      {step === 2 && (
        <>
          <Card>
            <h1 style={{ marginBottom: 16 }}>Transformations</h1>
            <Divider />
            {showWarning && (
              <Message
                style={{ marginBottom: 16 }}
                content="Please note: if you only edit the following Scope Configs without editing Data Entities in the previous step, you will only need to re-transform data on the Project page to see the Dashboard updated."
              />
            )}

            <Form labelCol={{ span: 4 }} wrapperCol={{ span: 16 }}>
              {plugin === 'azuredevops' && (
                <AzureTransformation
                  entities={entities}
                  transformation={transformation}
                  setTransformation={setTransformation}
                />
              )}

              {plugin === 'bamboo' && (
                <BambooTransformation
                  entities={entities}
                  transformation={transformation}
                  setTransformation={setTransformation}
                />
              )}

              {plugin === 'bitbucket' && (
                <BitbucketTransformation
                  entities={entities}
                  transformation={transformation}
                  setTransformation={setTransformation}
                />
              )}

              {plugin === 'github' && (
                <GitHubTransformation
                  entities={entities}
                  transformation={transformation}
                  setTransformation={setTransformation}
                  setHasError={setHasError}
                />
              )}

              {plugin === 'gitlab' && (
                <GitLabTransformation
                  entities={entities}
                  transformation={transformation}
                  setTransformation={setTransformation}
                  setHasError={setHasError}
                />
              )}

              {plugin === 'jenkins' && (
                <JenkinsTransformation
                  entities={entities}
                  transformation={transformation}
                  setTransformation={setTransformation}
                />
              )}

              {plugin === 'jira' && (
                <JiraTransformation
                  entities={entities}
                  connectionId={connectionId}
                  transformation={transformation}
                  setTransformation={setTransformation}
                />
              )}

              {plugin === 'tapd' && scopeId && (
                <TapdTransformation
                  entities={entities}
                  connectionId={connectionId}
                  scopeId={scopeId}
                  transformation={transformation}
                  setTransformation={setTransformation}
                />
              )}
            </Form>
          </Card>
          <Flex justify="flex-end" gap="small">
            <Button outlined intent={Intent.PRIMARY} text="Prev" onClick={handlePrevStep} />
            <Button
              loading={operating}
              disabled={hasError}
              intent={Intent.PRIMARY}
              text="Save"
              onClick={handleSubmit}
            />
          </Flex>
        </>
      )}
    </Flex>
  );
};
