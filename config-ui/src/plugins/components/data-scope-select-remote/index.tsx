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

import { useMemo, useState } from 'react';
import { Button, Intent } from '@blueprintjs/core';

import { getPluginId } from '@/plugins';

import { GitHubDataScope } from '@/plugins/register/github';
import { JiraDataScope } from '@/plugins/register/jira';
import { GitLabDataScope } from '@/plugins/register/gitlab';
import { JenkinsDataScope } from '@/plugins/register/jenkins';
import { BitbucketDataScope } from '@/plugins/register/bitbucket';
import { AzureDataScope } from '@/plugins/register/azure';
import { SonarQubeDataScope } from '@/plugins/register/sonarqube';
import { PagerDutyDataScope } from '@/plugins/register/pagerduty';
import { TapdDataScope } from '@/plugins/register/tapd';
import { ZentaoDataScope } from '@/plugins/register/zentao';

import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  connectionId: ID;
  disabledScope?: any[];
  onCancel: () => void;
  onSubmit: (origin: any) => void;
}

export const DataScopeSelectRemote = ({ plugin, connectionId, disabledScope, onSubmit, onCancel }: Props) => {
  const [operating, setOperating] = useState(false);
  const [scope, setScope] = useState<any>([]);

  const error = useMemo(() => (!scope.length ? 'No Data Scope is Selected' : ''), [scope]);

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

      onSubmit(res);
    } finally {
      setOperating(false);
      onCancel();
    }
  };

  return (
    <S.Wrapper>
      {plugin === 'github' && (
        <GitHubDataScope
          connectionId={connectionId}
          disabledItems={disabledScope}
          selectedItems={scope}
          onChangeItems={setScope}
        />
      )}

      {plugin === 'jira' && (
        <JiraDataScope
          connectionId={connectionId}
          disabledItems={disabledScope}
          selectedItems={scope}
          onChangeItems={setScope}
        />
      )}

      {plugin === 'gitlab' && (
        <GitLabDataScope
          connectionId={connectionId}
          disabledItems={disabledScope}
          selectedItems={scope}
          onChangeItems={setScope}
        />
      )}

      {plugin === 'jenkins' && (
        <JenkinsDataScope
          connectionId={connectionId}
          disabledItems={disabledScope}
          selectedItems={scope}
          onChangeItems={setScope}
        />
      )}

      {plugin === 'bitbucket' && (
        <BitbucketDataScope
          disabledItems={disabledScope}
          connectionId={connectionId}
          selectedItems={scope}
          onChangeItems={setScope}
        />
      )}

      {plugin === 'azuredevops' && (
        <AzureDataScope
          disabledItems={disabledScope}
          connectionId={connectionId}
          selectedItems={scope}
          onChangeItems={setScope}
        />
      )}

      {plugin === 'sonarqube' && (
        <SonarQubeDataScope
          disabledItems={disabledScope}
          connectionId={connectionId}
          selectedItems={scope}
          onChangeItems={setScope}
        />
      )}

      {plugin === 'pagerduty' && (
        <PagerDutyDataScope
          connectionId={connectionId}
          disabledItems={disabledScope}
          selectedItems={scope}
          onChangeItems={setScope}
        />
      )}

      {plugin === 'tapd' && (
        <TapdDataScope
          connectionId={connectionId}
          disabledItems={disabledScope}
          selectedItems={scope}
          onChangeItems={setScope}
        />
      )}

      {plugin === 'zentao' && (
        <ZentaoDataScope
          connectionId={connectionId}
          disabledItems={disabledScope}
          selectedItems={scope}
          onChangeItems={setScope}
        />
      )}

      <div className="action">
        <Button outlined intent={Intent.PRIMARY} text="Cancel" disabled={operating} onClick={onCancel} />
        <Button
          outlined
          intent={Intent.PRIMARY}
          text="Save"
          loading={operating}
          disabled={!!error}
          onClick={handleSubmit}
        />
      </div>
    </S.Wrapper>
  );
};
