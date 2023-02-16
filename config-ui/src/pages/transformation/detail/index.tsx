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

import React from 'react';
import { useParams, useHistory } from 'react-router-dom';
import { InputGroup, ButtonGroup, Button, Intent } from '@blueprintjs/core';

import { PageLoading, PageHeader, Card } from '@/components';
import { GitHubTransformation } from '@/plugins/register/github';
import { GitLabTransformation } from '@/plugins/register/gitlab';
import { JenkinsTransformation } from '@/plugins/register/jenkins';
import { BitbucketTransformation } from '@/plugins/register/bitbucket';

import { useDetail } from './use-detail';
import * as S from './styled';

export const TransformationDetailPage = () => {
  const { plugin, tid } = useParams<{ plugin: string; tid?: string }>();
  const history = useHistory();

  const { loading, operating, name, transformation, onChangeName, onChangeTransformation, onSave } = useDetail({
    plugin,
    id: tid,
  });

  if (loading) {
    return <PageLoading />;
  }

  return (
    <PageHeader
      breadcrumbs={[
        { name: 'Transformations', path: '/transformations' },
        {
          name: plugin,
          path: '/transformations',
        },
        {
          name: 'Create',
          path: `/transformations/${plugin}/${tid ? tid : 'Create'}`,
        },
      ]}
    >
      <S.Wrapper>
        <Card className="name card">
          <h3>Transformation Name *</h3>
          <p>Give this set of transformation rules a unique name so that you can identify it in the future.</p>
          <InputGroup
            placeholder="Enter Transformation Name"
            value={name}
            onChange={(e) => onChangeName(e.target.value)}
          />
        </Card>
        <Card className="card">
          {plugin === 'github' && (
            <GitHubTransformation transformation={transformation} setTransformation={onChangeTransformation} />
          )}

          {plugin === 'gitlab' && (
            <GitLabTransformation transformation={transformation} setTransformation={onChangeTransformation} />
          )}

          {plugin === 'jenkins' && (
            <JenkinsTransformation transformation={transformation} setTransformation={onChangeTransformation} />
          )}

          {plugin === 'bitbucket' && (
            <BitbucketTransformation transformation={transformation} setTransformation={onChangeTransformation} />
          )}

          <div className="action">
            <ButtonGroup>
              <Button disabled={operating} outlined text="Cancel" onClick={() => history.push('/transformations')} />
              <Button
                disabled={!name}
                loading={operating}
                outlined
                intent={Intent.PRIMARY}
                text="Save"
                onClick={onSave}
              />
            </ButtonGroup>
          </div>
        </Card>
      </S.Wrapper>
    </PageHeader>
  );
};
