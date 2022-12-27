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

import React, { useMemo } from 'react';
import { useParams } from 'react-router-dom';

import { PageHeader, PageLoading } from '@/components';
import { Plugins } from '@/plugins';

import { useDetail } from './use-detail';
import { BlueprintDetail } from './blueprint-detail';

import * as S from './styled';

export const BlueprintDetailPage = () => {
  const { id } = useParams<{ id: string }>();

  const { loading, blueprint } = useDetail({ id });

  const showJenkinsTips = useMemo(() => {
    const jenkins = blueprint && blueprint.settings.connections.find((cs) => cs.plugin === Plugins.Jenkins);

    return !jenkins?.scopes.length;
  }, [blueprint]);

  if (loading || !blueprint) {
    return <PageLoading />;
  }

  return (
    <PageHeader
      breadcrumbs={[
        { name: 'Blueprints', path: '/blueprints' },
        { name: blueprint.name, path: `/blueprints/${blueprint.id}` },
      ]}
    >
      <BlueprintDetail id={blueprint.id} />
      {showJenkinsTips && (
        <S.JenkinsTips>
          <p>Please add the "Jenkins jobs" to collect before this Blueprint can run again.</p>
        </S.JenkinsTips>
      )}
    </PageHeader>
  );
};
