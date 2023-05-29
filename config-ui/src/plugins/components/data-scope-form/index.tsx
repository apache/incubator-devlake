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

import { useEffect, useMemo, useState } from 'react';
import { Button, Intent } from '@blueprintjs/core';

import { Card, MultiSelector } from '@/components';
import { transformEntities } from '@/config';
import { getPluginConfig, getPluginId } from '@/plugins';

import { GitHubDataScope } from '@/plugins/register/github';
import { JiraDataScope } from '@/plugins/register/jira';
import { GitLabDataScope } from '@/plugins/register/gitlab';
import { JenkinsDataScope } from '@/plugins/register/jenkins';
import { BitbucketDataScope } from '@/plugins/register/bitbucket';
import { AzureDataScope } from '@/plugins/register/azure';
import { SonarQubeDataScope } from '@/plugins/register/sonarqube';
import { PagerDutyDataScope } from '@/plugins/register/pagerduty';
import { ZentaoDataScope } from '@/plugins/register/zentao';
import { KubeDeploymentDataScope } from '@/plugins/register/myplug';

import * as API from './api';
import * as S from './styled';
import { TapdDataScope } from '@/plugins/register/tapd';

interface Props {
  plugin: string;
  connectionId: ID;
  initialScope?: any[];
  initialEntities?: string[];
  cancelBtnProps?: {
    text?: string;
  };
  submitBtnProps?: {
    text?: string;
  };
  onCancel?: () => void;
  onSubmit?: (scope: Array<{ id: string; entities: string[] }>, origin: any) => void;
}

export const DataScopeForm = ({
  plugin,
  connectionId,
  initialScope,
  initialEntities,
  onSubmit,
  onCancel,
  cancelBtnProps,
  submitBtnProps,
}: Props) => {
  const [operating, setOperating] = useState(false);
  const [scope, setScope] = useState<any>([]);
  const [entities, setEntites] = useState<string[]>([]);

  const config = useMemo(() => getPluginConfig(plugin), []);

  const error = useMemo(
    () => (!scope.length || !entities.length ? 'No Data Scope is Selected' : ''),
    [scope, entities],
  );

  useEffect(() => {
    setScope(initialScope ?? []);
    setEntites(initialEntities ?? config.entities);
  }, []);

  const getDataScope = async (scope: any) => {
    try {
      const res = await API.getDataScope(plugin, connectionId, scope[getPluginId(plugin)]);
      return {
        ...scope,
        transformationRuleId: res.transformationRuleId,
      };
    } catch {
      return scope;
    }
  };

  const handleSubmit = async () => {
    setOperating(true);
    try {
      const data = await Promise.all(scope.map((sc: any) => getDataScope(sc)));
      const res =
        plugin === 'zentao'
          ? [
              ...(await API.updateDataScopeWithType(plugin, connectionId, 'product', {
                data: data.filter((s) => s.type !== 'project'),
              })),
              ...(await API.updateDataScopeWithType(plugin, connectionId, 'project', {
                data: data.filter((s) => s.type === 'project'),
              })),
            ]
          : await API.updateDataScope(plugin, connectionId, {
              data,
            });

      onSubmit?.(
        res.map((it: any) => ({
          id: plugin === 'zentao' ? `${it.type}/${it.id}` : it[getPluginId(plugin)],
          entities,
        })),
        res,
      );
    } finally {
      setOperating(false);
    }
  };

  return (
    <S.Wrapper>
      <Card>
        <div className="block">
          {plugin === 'github' && (
            <GitHubDataScope connectionId={connectionId} selectedItems={scope} onChangeItems={setScope} />
          )}

          {plugin === 'kube_deployment' && (
            <KubeDeploymentDataScope connectionId={connectionId} selectedItems={scope} onChangeItems={setScope} />
          )}

          {plugin === 'jira' && (
            <JiraDataScope connectionId={connectionId} selectedItems={scope} onChangeItems={setScope} />
          )}

          {plugin === 'gitlab' && (
            <GitLabDataScope connectionId={connectionId} selectedItems={scope} onChangeItems={setScope} />
          )}

          {plugin === 'jenkins' && (
            <JenkinsDataScope connectionId={connectionId} selectedItems={scope} onChangeItems={setScope} />
          )}

          {plugin === 'bitbucket' && (
            <BitbucketDataScope connectionId={connectionId} selectedItems={scope} onChangeItems={setScope} />
          )}

          {plugin === 'azuredevops' && (
            <AzureDataScope connectionId={connectionId} selectedItems={scope} onChangeItems={setScope} />
          )}

          {plugin === 'sonarqube' && (
            <SonarQubeDataScope connectionId={connectionId} selectedItems={scope} onChangeItems={setScope} />
          )}

          {plugin === 'pagerduty' && (
            <PagerDutyDataScope connectionId={connectionId} selectedItems={scope} onChangeItems={setScope} />
          )}

          {plugin === 'tapd' && (
            <TapdDataScope connectionId={connectionId} selectedItems={scope} onChangeItems={setScope} />
          )}

          {plugin === 'zentao' && (
            <ZentaoDataScope connectionId={connectionId} selectedItems={scope} onChangeItems={setScope} />
          )}
        </div>

        <div className="block">
          <h3>Data Entities</h3>
          <p>
            <span>Select the data entities you wish to collect for the projects.</span>{' '}
            <a
              href="https://devlake.apache.org/docs/DataModels/DevLakeDomainLayerSchema/#data-models"
              target="_blank"
              rel="noreferrer"
            >
              Learn about data entities
            </a>
          </p>
          <MultiSelector
            items={transformEntities(config.entities)}
            getKey={(item) => item.value}
            getName={(item) => item.label}
            selectedItems={transformEntities(entities)}
            onChangeItems={(items) => setEntites(items.map((it) => it.value))}
          />
        </div>
      </Card>

      <div className="action">
        <Button
          outlined
          intent={Intent.PRIMARY}
          text="Cancel"
          {...cancelBtnProps}
          disabled={operating}
          onClick={onCancel}
        />
        <Button
          intent={Intent.PRIMARY}
          text="Save"
          {...submitBtnProps}
          loading={operating}
          disabled={!!error}
          onClick={handleSubmit}
        />
      </div>
    </S.Wrapper>
  );
};
