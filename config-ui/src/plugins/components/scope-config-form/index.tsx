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
import { InputGroup, Button, Intent } from '@blueprintjs/core';

import { Alert, ExternalLink, Card, FormItem, MultiSelector, Buttons, Divider } from '@/components';
import { transformEntities, EntitiesLabel } from '@/config';
import { getPluginConfig } from '@/plugins';
import { GitHubTransformation } from '@/plugins/register/github';
import { JiraTransformation } from '@/plugins/register/jira';
import { GitLabTransformation } from '@/plugins/register/gitlab';
import { JenkinsTransformation } from '@/plugins/register/jenkins';
import { BitbucketTransformation } from '@/plugins/register/bitbucket';
import { AzureTransformation } from '@/plugins/register/azure';
import { TapdTransformation } from '@/plugins/register/tapd';
import { operator } from '@/utils';

import { AdditionalSettings } from './fields';
import { TIPS_MAP } from './misc';
import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  connectionId: ID;
  scopeId?: ID;
  scopeConfigId?: ID;
  onCancel?: () => void;
  onSubmit?: (trId: string) => void;
}

export const ScopeConfigForm = ({ plugin, connectionId, scopeId, scopeConfigId, onCancel, onSubmit }: Props) => {
  const [step, setStep] = useState(1);
  const [name, setName] = useState('');
  const [entities, setEntities] = useState<string[]>([]);
  const [transformation, setTransformation] = useState<any>({});
  const [hasRefDiff, setHasRefDiff] = useState(false);
  const [operating, setOperating] = useState(false);

  const config = useMemo(() => getPluginConfig(plugin), []);

  useEffect(() => {
    setHasRefDiff(!!config.transformation.refdiff);
  }, [config.transformation]);

  useEffect(() => {
    if (!scopeConfigId) return;

    (async () => {
      try {
        const res = await API.getScopeConfig(plugin, connectionId, scopeConfigId);
        setName(res.name);
        setEntities(res.entities);
        setTransformation(omit(res, ['id', 'connectionId', 'name', 'entities', 'createdAt', 'updatedAt']));
      } catch {}
    })();
  }, [scopeConfigId]);

  const handleNextStep = () => {
    setStep(2);
  };

  const handleSubmit = async () => {
    const [success, res] = await operator(
      () =>
        !scopeConfigId
          ? API.createScopeConfig(plugin, connectionId, { name, entities, ...transformation })
          : API.updateScopeConfig(plugin, connectionId, scopeConfigId, { name, entities, ...transformation }),
      {
        setOperating,
        formatMessage: () => (!scopeConfigId ? 'Create scope config successful.' : 'Update scope config successful'),
      },
    );

    if (success) {
      onCancel?.();
      onSubmit?.(res.id);
    }
  };

  return (
    <S.Wrapper>
      {TIPS_MAP[plugin] && (
        <Alert style={{ marginBottom: 24 }}>
          To learn about how {TIPS_MAP[plugin].name} transformation is used in DevLake,{' '}
          <ExternalLink link={TIPS_MAP[plugin].link}>check out this doc</ExternalLink>.
        </Alert>
      )}
      {step === 1 && (
        <>
          <Card>
            <FormItem
              label="Scope Config Name"
              subLabel="Give this Scope Config a unique name so that you can identify it in the future."
              required
            >
              <InputGroup placeholder="My Scope Config 1" value={name} onChange={(e) => setName(e.target.value)} />
            </FormItem>
          </Card>
          <Card>
            <FormItem
              label="Data Entities"
              subLabel={
                <>
                  Select the data entities you wish to collect for the Data Scope.
                  <ExternalLink link="">Learn about data entities</ExternalLink>
                </>
              }
              required
            >
              <MultiSelector
                items={transformEntities(config.entities)}
                getKey={(it) => it.value}
                getName={(it) => it.label}
                selectedItems={entities.map((it) => ({ label: EntitiesLabel[it], value: it }))}
                onChangeItems={(its) => setEntities(its.map((it) => it.value))}
              />
            </FormItem>
          </Card>
          <Buttons>
            <Button outlined intent={Intent.PRIMARY} text="Cancel" onClick={onCancel} />
            <Button disabled={!name || !entities.length} intent={Intent.PRIMARY} text="Next" onClick={handleNextStep} />
          </Buttons>
        </>
      )}
      {step === 2 && (
        <>
          <Card>
            {plugin === 'github' && (
              <GitHubTransformation transformation={transformation} setTransformation={setTransformation} />
            )}

            {plugin === 'jira' && (
              <JiraTransformation
                connectionId={connectionId}
                transformation={transformation}
                setTransformation={setTransformation}
              />
            )}

            {plugin === 'gitlab' && (
              <GitLabTransformation transformation={transformation} setTransformation={setTransformation} />
            )}

            {plugin === 'jenkins' && (
              <JenkinsTransformation transformation={transformation} setTransformation={setTransformation} />
            )}

            {plugin === 'bitbucket' && (
              <BitbucketTransformation transformation={transformation} setTransformation={setTransformation} />
            )}

            {plugin === 'azuredevops' && (
              <AzureTransformation transformation={transformation} setTransformation={setTransformation} />
            )}

            {plugin === 'tapd' && scopeId && (
              <TapdTransformation
                connectionId={connectionId}
                scopeId={scopeId}
                transformation={transformation}
                setTransformation={setTransformation}
              />
            )}

            {hasRefDiff && (
              <>
                <Divider />
                <AdditionalSettings transformation={transformation} setTransformation={setTransformation} />
              </>
            )}
          </Card>
          <Buttons>
            <Button outlined intent={Intent.PRIMARY} text="Cancel" onClick={onCancel} />
            <Button loading={operating} intent={Intent.PRIMARY} text="Save" onClick={handleSubmit} />
          </Buttons>
        </>
      )}
    </S.Wrapper>
  );
};
