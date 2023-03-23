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

import { useMemo } from 'react';
import { Button, Intent } from '@blueprintjs/core';

import { Card, MultiSelector } from '@/components';
import { transformEntities } from '@/config';
import type { PluginConfigType } from '@/plugins';
import { PluginConfig } from '@/plugins';
import { GitHubDataScope } from '@/plugins/register/github';
import { JiraDataScope } from '@/plugins/register/jira';
import { GitLabDataScope } from '@/plugins/register/gitlab';
import { JenkinsDataScope } from '@/plugins/register/jenkins';
import { BitbucketDataScope } from '@/plugins/register/bitbucket';
import { SonarQubeDataScope } from '@/plugins/register/sonarqube';
import { ZentaoDataScope } from '@/plugins/register/zentao';

import type { UseDataScope } from './use-data-scope';
import { useDataScope } from './use-data-scope';
import * as S from './styled';

interface Props extends Pick<UseDataScope, 'plugin' | 'connectionId' | 'onSubmit'> {
  cancelBtnProps?: {
    text?: string;
  };
  submitBtnProps?: {
    text?: string;
  };
  initialScope?: any[];
  initialEntities?: string[];
  onCancel?: () => void;
}

export const DataScopeForm = ({
  plugin,
  connectionId,
  initialScope,
  initialEntities,
  onCancel,
  cancelBtnProps,
  submitBtnProps,
  ...props
}: Props) => {
  const config = useMemo(() => PluginConfig.find((p) => p.plugin === plugin) as PluginConfigType, []);

  const { saving, scope, setScope, entities, setEntites, onSave } = useDataScope({
    ...props,
    plugin,
    connectionId,
    initialScope: initialScope ?? [],
    initialEntities: initialEntities ?? config.entities,
  });

  return (
    <S.Wrapper>
      <Card>
        <div className="block">
          {plugin === 'github' && (
            <GitHubDataScope connectionId={connectionId} selectedItems={scope} onChangeItems={setScope} />
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

          {plugin === 'sonarqube' && (
            <SonarQubeDataScope connectionId={connectionId} selectedItems={scope} onChangeItems={setScope} />
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
          disabled={saving}
          text="Cancel"
          onClick={onCancel}
          {...cancelBtnProps}
        />
        <Button
          outlined
          intent={Intent.PRIMARY}
          loading={saving}
          disabled={!scope.length || !entities.length}
          text="Save"
          onClick={onSave}
          {...submitBtnProps}
        />
      </div>
    </S.Wrapper>
  );
};
