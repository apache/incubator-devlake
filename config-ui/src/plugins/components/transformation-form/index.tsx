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
import { Button, InputGroup, Intent } from '@blueprintjs/core';

import { Card, ExternalLink, PageLoading, Divider } from '@/components';
import { useRefreshData } from '@/hooks';
import { operator } from '@/utils';
import { getPluginConfig } from '@/plugins';
import { GitHubTransformation } from '@/plugins/register/github';
import { JiraTransformation } from '@/plugins/register/jira';
import { GitLabTransformation } from '@/plugins/register/gitlab';
import { JenkinsTransformation } from '@/plugins/register/jenkins';
import { BitbucketTransformation } from '@/plugins/register/bitbucket';
import { AzureTransformation } from '@/plugins/register/azure';
import { TapdTransformation } from '@/plugins/register/tapd';
import { KubeDeploymentTransformation } from '@/plugins/register/myplug';

import { TIPS_MAP } from './misc';
import { AdditionalSettings } from './fields';
import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  connectionId: ID;
  scopeId: ID;
  id?: ID;
  onCancel?: (transformationRule?: any) => void;
}

export const TransformationForm = ({ plugin, connectionId, scopeId, id, onCancel }: Props) => {
  const [saving, setSaving] = useState(false);
  const [name, setName] = useState('');
  const [transformation, setTransformation] = useState({});
  const [hasRefDiff, setHasRefDiff] = useState(false);

  const config = useMemo(() => getPluginConfig(plugin), []);

  const { ready, data } = useRefreshData(async () => {
    if (!id) return null;
    return API.getTransformation(plugin, connectionId, id);
  }, [id]);

  useEffect(() => {
    setTransformation(data ?? config.transformation);
    setHasRefDiff(!!config.transformation.refdiff);
    setName(data?.name ?? '');
  }, [data, config.transformation]);

  const handleSubmit = async () => {
    const [success, res] = await operator(
      () =>
        id
          ? API.updateTransformation(plugin, connectionId, id, { ...transformation, name })
          : API.createTransformation(plugin, connectionId, { ...transformation, name }),
      {
        setOperating: setSaving,
        formatMessage: () => 'Transformation created successfully',
      },
    );

    if (success) {
      onCancel?.(res);
    }
  };

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

        {plugin === 'azuredevops' && (
          <AzureTransformation transformation={transformation} setTransformation={setTransformation} />
        )}

        {plugin === 'tapd' && (
          <TapdTransformation
            connectionId={connectionId}
            scopeId={scopeId}
            transformation={transformation}
            setTransformation={setTransformation}
          />
        )}

        {plugin === 'kube_deployment' && (
          <KubeDeploymentTransformation
            connectionId={connectionId}
            transformation={transformation}
            setTransformation={setTransformation} />
        )}


        {hasRefDiff && (
          <>
            <Divider />
            <AdditionalSettings transformation={transformation} setTransformation={setTransformation} />
          </>
        )}
      </Card>

      <S.Btns>
        <Button outlined intent={Intent.PRIMARY} text="Cancel" onClick={() => onCancel?.(undefined)} />
        <Button intent={Intent.PRIMARY} disabled={!name} loading={saving} text="Save" onClick={handleSubmit} />
      </S.Btns>
    </S.Wrapper>
  );
};
