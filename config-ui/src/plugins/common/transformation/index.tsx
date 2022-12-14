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

import React, { useState } from 'react'
import {
  RadioGroup,
  Radio,
  InputGroup,
  ButtonGroup,
  Button,
  Intent
} from '@blueprintjs/core'

import { Divider, Selector, MultiSelector } from '@/components'

import { Plugins } from '@/plugins'
import { GitHubTransformation } from '@/plugins/github'
import { JIRATransformation } from '@/plugins/jira'
import { GitLabTransformation } from '@/plugins/gitlab'
import { JenkinsTransformation } from '@/plugins/jenkins'

import type { TransformationType, RuleItem } from './types'
import type { UseTransformationProps } from './use-transformation'
import { useTransformation } from './use-transformation'
import * as S from './styled'

interface Props
  extends Omit<
    UseTransformationProps,
    'name' | 'selectedRule' | 'setSelectedScope'
  > {
  onCancel?: () => void
}

export const Transformation = ({
  plugin,
  connectionId,
  onCancel,
  ...props
}: Props) => {
  const [type, setType] = useState<TransformationType>()
  const [name, setName] = useState('')
  const [selectedRule, setSelectedRule] = useState<RuleItem>()
  const [selectedScope, setSelectedScope] = useState<any>([])

  const {
    loading,
    rules,
    scope,
    saving,
    transformation,
    getScopeKey,
    onSave,
    onUpdateScope,
    onChangeTransformation
  } = useTransformation({
    ...props,
    plugin,
    connectionId,
    name,
    selectedRule,
    selectedScope
  })

  const handleChangeType = (e: React.FormEvent<HTMLInputElement>) => {
    setType((e.target as HTMLInputElement).value as TransformationType)
    setSelectedRule(undefined)
  }

  return (
    <S.Wrapper>
      <div className='block'>
        <h3>Create a New or Select an Existing Transformation *</h3>
        <RadioGroup selectedValue={type} onChange={handleChangeType}>
          <Radio
            label='Create a new transformation from a blank template'
            value='create'
          />
          <Radio
            label='Create a new transformation by duplicating an exisitng transformation as the template'
            value='createByExist'
          />
          <Radio
            label='Select an existing transformation'
            value='selectExist'
          />
        </RadioGroup>
      </div>

      {type && (
        <>
          <Divider />
          <div className='block'>
            {(type === 'createByExist' || type === 'selectExist') && (
              <div className='item'>
                <h3>Select a Transformation to Duplicate *</h3>
                <p>
                  The selected transformation will be used as template that you
                  can tweak and save as a new transformation.
                </p>
                <Selector
                  items={rules}
                  getKey={(it) => it.id}
                  getName={(it) => it.name}
                  selectedItem={selectedRule}
                  onChangeItem={setSelectedRule}
                />
              </div>
            )}
            {type !== 'selectExist' && (
              <div className='item'>
                <h3>Transformation Name *</h3>
                <p>
                  Give this set of transformation rules a unique name so that
                  you can identify it in the future.
                </p>
                <InputGroup
                  placeholder='Enter Transformation Name'
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                />
              </div>
            )}
            <div className='item'>
              <h3>Applied Data Scope</h3>
              <p>
                Select the data scope for which you want to apply this
                transformation for.
              </p>
              <MultiSelector
                loading={loading}
                items={scope}
                getKey={getScopeKey}
                getName={(sc) => sc.name}
                selectedItems={selectedScope}
                onChangeItems={setSelectedScope}
              />
              {type === 'selectExist' && (
                <ButtonGroup>
                  <Button
                    outlined
                    intent={Intent.PRIMARY}
                    text='Cancel and Go Back'
                    onClick={onCancel}
                  />
                  <Button
                    outlined
                    disabled={!selectedRule || !selectedScope.length}
                    intent={Intent.PRIMARY}
                    text='Save'
                    onClick={() => onUpdateScope(selectedRule?.id)}
                  />
                </ButtonGroup>
              )}
            </div>
          </div>
        </>
      )}

      {(type === 'create' || (type === 'createByExist' && selectedRule)) && (
        <>
          <Divider />
          <div className='block'>
            {plugin === Plugins.GitHub && (
              <GitHubTransformation
                transformation={transformation}
                setTransformation={onChangeTransformation}
              />
            )}

            {plugin === Plugins.JIRA && (
              <JIRATransformation
                connectionId={connectionId}
                transformation={transformation}
                setTransformation={onChangeTransformation}
              />
            )}

            {plugin === Plugins.GitLab && (
              <GitLabTransformation
                transformation={transformation}
                setTransformation={onChangeTransformation}
              />
            )}

            {plugin === Plugins.Jenkins && (
              <JenkinsTransformation
                transformation={transformation}
                setTransformation={onChangeTransformation}
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
          </div>
        </>
      )}
    </S.Wrapper>
  )
}
