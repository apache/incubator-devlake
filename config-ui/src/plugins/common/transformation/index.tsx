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
import { GitHubTransformation } from '@/plugins/github'
import { JIRATransformation } from '@/plugins/jira'
import { GitLabTransformation } from '@/plugins/gitlab'
import { JenkinsTransformation } from '@/plugins/jenkins'

import type { UseTransformationProps } from './use-transformation'
import { useTransformation } from './use-transformation'
import * as S from './styled'

interface Props extends UseTransformationProps {
  connectionId: ID
  onCancel?: () => void
}

export const Transformation = ({
  plugin,
  connectionId,
  name,
  initialValues,
  onCancel,
  onSaveAfter
}: Props) => {
  const { saving, transformation, setTransformation, onSave } =
    useTransformation({
      plugin,
      name,
      initialValues,
      onSaveAfter
    })

  return (
    <S.Container>
      {plugin === Plugins.GitHub && (
        <GitHubTransformation
          transformation={transformation}
          setTransformation={setTransformation}
        />
      )}

      {plugin === Plugins.JIRA && (
        <JIRATransformation
          connectionId={connectionId}
          transformation={transformation}
          setTransformation={setTransformation}
        />
      )}

      {plugin === Plugins.GitLab && (
        <GitLabTransformation
          transformation={transformation}
          setTransformation={setTransformation}
        />
      )}

      {plugin === Plugins.Jenkins && (
        <JenkinsTransformation
          transformation={transformation}
          setTransformation={setTransformation}
        />
      )}
      <ButtonGroup>
        <Button
          outlined
          intent={Intent.PRIMARY}
          text='Cancel and Go Back'
          onClick={onCancel}
        />
        <Button
          outlined
          disabled={!name}
          loading={saving}
          intent={Intent.PRIMARY}
          text='Save'
          onClick={onSave}
        />
      </ButtonGroup>
    </S.Container>
  )
}
