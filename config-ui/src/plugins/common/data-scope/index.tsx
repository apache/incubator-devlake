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

import React from 'react'
import { ButtonGroup, Button, Intent } from '@blueprintjs/core'

import { Plugins } from '@/plugins'
import { GitHubDataScope } from '@/plugins/github'
import { JIRADataScope } from '@/plugins/jira'
import { GitLabDataScope } from '@/plugins/gitlab'
import { JenkinsDataScope } from '@/plugins/jenkins'
import { MultiSelector } from '@/components'

import type { UseDataScope } from './use-data-scope'
import { useDataScope } from './use-data-scope'
import * as S from './styled'

interface Props extends UseDataScope {
  onCancel?: () => void
}

export const DataScope = ({
  plugin,
  connectionId,
  entities,
  onCancel,
  ...props
}: Props) => {
  const {
    saving,
    selectedScope,
    selectedEntities,
    onChangeScope,
    onChangeEntites,
    onSave
  } = useDataScope({
    ...props,
    plugin,
    connectionId,
    entities
  })

  return (
    <S.Wrapper>
      <div className='block'>
        {plugin === Plugins.GitHub && (
          <GitHubDataScope
            connectionId={connectionId}
            selectedItems={selectedScope}
            onChangeItems={onChangeScope}
          />
        )}

        {plugin === Plugins.JIRA && (
          <JIRADataScope
            connectionId={connectionId}
            selectedItems={selectedScope}
            onChangeItems={onChangeScope}
          />
        )}

        {plugin === Plugins.GitLab && (
          <GitLabDataScope
            connectionId={connectionId}
            selectedItems={selectedScope}
            onChangeItems={onChangeScope}
          />
        )}

        {plugin === Plugins.Jenkins && (
          <JenkinsDataScope
            connectionId={connectionId}
            selectedItems={selectedScope}
            onChangeItems={onChangeScope}
          />
        )}
      </div>

      <div className='block'>
        <h3>Data Entities</h3>
        <p>
          <span>
            Select the data entities you wish to collect for the projects.
          </span>
          <a
            href='https://devlake.apache.org/docs/DataModels/DevLakeDomainLayerSchema/'
            target='_blank'
          >
            Learn about data entities
          </a>
        </p>
        <MultiSelector
          items={entities}
          selectedItems={selectedEntities}
          onChangeItems={onChangeEntites}
        />
      </div>

      <ButtonGroup>
        <Button outlined disabled={saving} text='Cancel' onClick={onCancel} />
        <Button
          outlined
          intent={Intent.PRIMARY}
          loading={saving}
          disabled={!selectedScope.length || !selectedEntities.length}
          text='Save'
          onClick={onSave}
        />
      </ButtonGroup>
    </S.Wrapper>
  )
}
