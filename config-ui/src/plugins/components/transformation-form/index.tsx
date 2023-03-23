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
import { InputGroup, Button, Intent } from '@blueprintjs/core';

import { PageLoading, ExternalLink, Card } from '@/components';
import { useRefreshData, useOperator } from '@/hooks';
import type { PluginConfigType } from '@/plugins';
import { PluginConfig } from '@/plugins';
import { GitHubTransformation } from '@/plugins/register/github';
import { JiraTransformation } from '@/plugins/register/jira';
import { GitLabTransformation } from '@/plugins/register/gitlab';
import { JenkinsTransformation } from '@/plugins/register/jenkins';
import { BitbucketTransformation } from '@/plugins/register/bitbucket';

import { TIPS_MAP } from './misc';
import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  connectionId: ID;
  id?: ID;
  onCancel?: () => void;
}

export const TransformationForm = ({ plugin, connectionId, id, onCancel }: Props) => {
  const [name, setName] = useState('');
  const [transformation, setTransformation] = useState({});

  const config = useMemo(() => PluginConfig.find((p) => p.plugin === plugin) as PluginConfigType, []);

  const { ready, data } = useRefreshData(async () => {
    if (!id) return null;
    return API.getTransformation(plugin, connectionId, id);
  }, [id]);

  const { operating, onSubmit } = useOperator(
    (payload) =>
      id
        ? API.updateTransformation(plugin, connectionId, id, payload)
        : API.createTransformation(plugin, connectionId, payload),
    {
      callback: onCancel,
      formatMessage: () => 'Transformation created successfully',
    },
  );

  useEffect(() => {
    setTransformation(data ?? config.transformation);
  }, [data, config.transformation]);

  if (!ready) {
    return <PageLoading />;
  }

  return (
    <S.Wrapper>
      {TIPS_MAP[plugin] && (
        <S.Tips>
          To learn about how {TIPS_MAP[plugin].name} transformation is used in DevLake,{' '}
          <ExternalLink link={TIPS_MAP[plugin].link}>check out this doc</ExternalLink>.
        </S.Tips>
      )}

      <Card style={{ marginTop: 24 }}>
        <h3>Transformation Name *</h3>
        <p>Give this set of transformation rules a unique name so that you can identify it in the future.</p>
        <InputGroup placeholder="Enter Transformation Name" value={name} onChange={(e) => setName(e.target.value)} />
      </Card>

      <Card style={{ marginTop: 24 }}>
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
      </Card>

      <S.Btns>
        <Button outlined intent={Intent.PRIMARY} text="Cancel" onClick={onCancel} />
        <Button
          intent={Intent.PRIMARY}
          disabled={!name}
          loading={operating}
          text="Save"
          onClick={() => onSubmit({ ...transformation, name })}
        />
      </S.Btns>
    </S.Wrapper>
  );
};
