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
  Icon,
  RadioGroup,
  Radio,
  InputGroup,
  ButtonGroup,
  Button,
  Intent
} from '@blueprintjs/core'

import { Table, Divider } from '@/components'
import {
  DataScopeSelector,
  RuleSelector,
  Transformation,
  Plugins
} from '@/plugins'

import type { BPConnectionItemType } from '../types'
import { useBlueprint } from '../hooks'
import * as API from '../api'

import { useColumns } from './use-columns'
import * as S from './styled'

type TransformationType = 'create' | 'createByExist' | 'selectExist'

export const StepThree = () => {
  const [connection, setConnection] = useState<BPConnectionItemType>()
  const [type, setType] = useState<TransformationType>()
  const [name, setName] = useState('')
  const [selectedRule, setSelectedRule] = useState<any>()
  const [selectedScope, setSelectedScope] = useState<any>([])

  const { connections, onChangeShowDetail } = useBlueprint()

  const handleGoDetail = (c: BPConnectionItemType) => {
    setConnection(c)
    onChangeShowDetail(true)
  }

  const handleBack = () => {
    setConnection(undefined)
    setType(undefined)
    setSelectedRule(undefined)
    onChangeShowDetail(false)
  }

  const handleChangeType = (e: React.FormEvent<HTMLInputElement>) => {
    setType((e.target as HTMLInputElement).value as TransformationType)
    setSelectedRule(undefined)
  }

  const handleUpdateScopeRule = async (
    plugin: Plugins,
    connectionId: ID,
    tsId: ID
  ) => {
    const paylod = selectedScope.map((sc: any) => ({
      ...sc,
      transformationRuleId: tsId
    }))

    if (paylod.length) {
      await API.updateDataScope(plugin, connectionId, {
        data: paylod
      })
    }

    handleBack()
  }

  const columns = useColumns({ onDetail: handleGoDetail })

  return !connection ? (
    <S.Card style={{ padding: 0 }}>
      <Table columns={columns} dataSource={connections} />
    </S.Card>
  ) : (
    <>
      <S.Card>
        <div className='back' onClick={handleBack}>
          <Icon icon='arrow-left' size={14} />
          <span>Cancel and Go Back</span>
        </div>
        <h2>Create/Select a Transformation</h2>
        <Divider />
        <h3 style={{ marginBottom: 12 }}>
          Create a New or Select an Existing Transformation *
        </h3>
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
      </S.Card>
      {type && (
        <S.Card>
          {(type === 'createByExist' || type === 'selectExist') && (
            <div className='block'>
              <h3>Select a Transformation to Duplicate *</h3>
              <p>
                The selected transformation will be used as template that you
                can tweak and save as a new transformation.
              </p>
              <RuleSelector
                plugin={connection.plugin}
                selectedRule={selectedRule}
                onChangeRule={setSelectedRule}
              />
            </div>
          )}
          {type !== 'selectExist' && (
            <div className='block'>
              <h3>Transformation Name *</h3>
              <p>
                Give this set of transformation rules a unique name so that you
                can identify it in the future.
              </p>
              <InputGroup
                placeholder='Enter Transformation Name'
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </div>
          )}
          <div className='block'>
            <h3>Applied Data Scope</h3>
            <p>
              Select the data scope for which you want to apply this
              transformation for.
            </p>
            <DataScopeSelector
              plugin={connection.plugin}
              connectionId={connection.id}
              scopeIds={connection.scope.map((sc) => sc.id)}
              selectedItems={selectedScope}
              onChangeItems={setSelectedScope}
            />
            {type === 'selectExist' && (
              <ButtonGroup>
                <Button
                  outlined
                  intent={Intent.PRIMARY}
                  text='Cancel and Go Back'
                  onClick={handleBack}
                />
                <Button
                  outlined
                  disabled={!selectedRule}
                  intent={Intent.PRIMARY}
                  text='Save'
                  onClick={() =>
                    handleUpdateScopeRule(
                      connection.plugin,
                      connection.id,
                      selectedRule.id
                    )
                  }
                />
              </ButtonGroup>
            )}
          </div>
        </S.Card>
      )}
      {(type === 'create' || (type === 'createByExist' && selectedRule)) && (
        <Transformation
          plugin={connection.plugin}
          connectionId={connection.id}
          name={name}
          initialValues={selectedRule}
          onSaveAfter={(tsid) =>
            handleUpdateScopeRule(connection.plugin, connection.id, tsid)
          }
          onCancel={handleBack}
        />
      )}
    </>
  )
}
