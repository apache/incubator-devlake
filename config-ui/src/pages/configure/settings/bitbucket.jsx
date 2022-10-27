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
import React, { useState, useEffect, useCallback } from 'react'
import {
  FormGroup,
  Checkbox,
  InputGroup,
  NumericInput
} from '@blueprintjs/core'

import { DataEntityTypes } from '@/data/DataEntities'
import Deployment from '@/components/blueprints/transformations/CICD/Deployment'

import { StatusSelect } from './bitbucket/status-select'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

export default function BitbucketSettings(props) {
  const {
    provider,
    connection,
    entities = [],
    transformation = {},
    isSaving = false,
    isSavingConnection = false,
    onSettingsChange = () => {}
    // configuredProject
    // configuredBoard
  } = props

  const [enableAdditionalCalculations, setEnableAdditionalCalculations] =
    useState(false)

  const [IssueStatusTODO, setIssueStatusTODO] = useState([])
  const [IssueStatusINPROGRESS, setIssueStatusINPROGRESS] = useState([])
  const [IssueStatusDONE, setIssueStatusDONE] = useState([])
  const [IssueStatusOTHER, setIssueStatusOTHER] = useState([])

  useEffect(() => {
    if (transformation?.IssueStatusTODO) {
      setIssueStatusTODO(transformation?.IssueStatusTODO ?? [])
    }

    if (transformation?.IssueStatusINPROGRESS) {
      setIssueStatusINPROGRESS(transformation?.IssueStatusINPROGRESS ?? [])
    }

    if (transformation?.IssueStatusDONE) {
      setIssueStatusDONE(transformation?.IssueStatusDONE ?? [])
    }

    if (transformation?.IssueStatusOTHER) {
      setIssueStatusOTHER(transformation?.IssueStatusOTHER ?? [])
    }
  }, [transformation])

  useEffect(() => {
    console.log(
      '>>>> BITBUCKET: TRANSFORMATION SETTINGS OBJECT....',
      transformation
    )
    setEnableAdditionalCalculations(!!transformation?.refdiff)
  }, [transformation])

  const handleAdditionalEnable = useCallback(
    (enable) => {
      setEnableAdditionalCalculations(enable)
      onSettingsChange({
        refdiff: enable
          ? { tagsOrder: '', tagsPattern: '', tagsLimit: 10 }
          : null
      })
    },
    [setEnableAdditionalCalculations, onSettingsChange]
  )

  const handleItemSelect = (name, selectedItem) => {
    let newValue

    switch (name) {
      case 'TODO':
        newValue = !IssueStatusTODO.includes(selectedItem)
          ? [...IssueStatusTODO, selectedItem]
          : [...IssueStatusTODO]
        setIssueStatusTODO(newValue)
        onSettingsChange({ IssueStatusTODO: newValue })
        break
      case 'IN-PROGRESS':
        newValue = !IssueStatusINPROGRESS.includes(selectedItem)
          ? [...IssueStatusINPROGRESS, selectedItem]
          : [...IssueStatusINPROGRESS]
        setIssueStatusINPROGRESS(newValue)
        onSettingsChange({ IssueStatusINPROGRESS: newValue })
        break
      case 'DONE':
        newValue = !IssueStatusDONE.includes(selectedItem)
          ? [...IssueStatusDONE, selectedItem]
          : [...IssueStatusDONE]
        setIssueStatusDONE(newValue)
        onSettingsChange({ IssueStatusDONE: newValue })
        break
      case 'OTHER':
        newValue = !IssueStatusOTHER.includes(selectedItem)
          ? [...IssueStatusOTHER, selectedItem]
          : [...IssueStatusOTHER]
        setIssueStatusOTHER(newValue)
        onSettingsChange({ IssueStatusOTHER: newValue })
        break
    }
  }

  const handleItemRemove = (name, removedItem) => {
    let newValue

    switch (name) {
      case 'TODO':
        newValue = IssueStatusTODO.filter((it) => it !== removedItem)
        setIssueStatusTODO(newValue)
        onSettingsChange({ IssueStatusTODO: newValue })
        break
      case 'IN-PROGRESS':
        newValue = IssueStatusINPROGRESS.filter((it) => it !== removedItem)
        setIssueStatusINPROGRESS(newValue)
        onSettingsChange({ IssueStatusINPROGRESS: newValue })
        break
      case 'DONE':
        newValue = IssueStatusDONE.filter((it) => it !== removedItem)
        setIssueStatusDONE(newValue)
        onSettingsChange({ IssueStatusDONE: newValue })
        break
      case 'OTHER':
        newValue = IssueStatusOTHER.filter((it) => it !== removedItem)
        setIssueStatusOTHER(newValue)
        onSettingsChange({ IssueStatusOTHER: newValue })
        break
    }
  }

  const handleItemClear = (name) => {
    switch (name) {
      case 'TODO':
        setIssueStatusTODO([])
        onSettingsChange({ IssueStatusTODO: [] })
        break
      case 'IN-PROGRESS':
        setIssueStatusINPROGRESS([])
        onSettingsChange({ IssueStatusINPROGRESS: [] })
        break
      case 'DONE':
        setIssueStatusDONE([])
        onSettingsChange({ IssueStatusDONE: [] })
        break
      case 'OTHER':
        setIssueStatusOTHER([])
        onSettingsChange({ IssueStatusOTHER: [] })
        break
    }
  }

  return (
    <>
      {entities.some((e) => e.value === DataEntityTypes.TICKET) && (
        <>
          <h5>Issue Tracking</h5>
          <h6>Issue Status Mapping</h6>
          <p>
            Standardize your issue statuses to the following issue statuses to
            view metrics such as `Requirement Delivery Rate` in built-in
            dashboards.
          </p>

          <StatusSelect
            style={{ marginBottom: 10 }}
            name='TODO'
            saving={isSaving}
            selectedItems={IssueStatusTODO}
            disabledItems={[
              ...IssueStatusINPROGRESS,
              ...IssueStatusDONE,
              ...IssueStatusOTHER
            ]}
            onItemSelect={handleItemSelect}
            onItemRemove={handleItemRemove}
            onItemClear={handleItemClear}
          />
          <StatusSelect
            style={{ marginBottom: 10 }}
            name='IN-PROGRESS'
            saving={isSaving}
            selectedItems={IssueStatusINPROGRESS}
            disabledItems={[
              ...IssueStatusTODO,
              ...IssueStatusDONE,
              ...IssueStatusOTHER
            ]}
            onItemSelect={handleItemSelect}
            onItemRemove={handleItemRemove}
            onItemClear={handleItemClear}
          />
          <StatusSelect
            name='DONE'
            saving={isSaving}
            selectedItems={IssueStatusDONE}
            disabledItems={[
              ...IssueStatusTODO,
              ...IssueStatusINPROGRESS,
              ...IssueStatusOTHER
            ]}
            onItemSelect={handleItemSelect}
            onItemRemove={handleItemRemove}
            onItemClear={handleItemClear}
          />
          <StatusSelect
            name='OTHER'
            saving={isSaving}
            selectedItems={IssueStatusOTHER}
            disabledItems={[
              ...IssueStatusTODO,
              ...IssueStatusINPROGRESS,
              ...IssueStatusDONE
            ]}
            onItemSelect={handleItemSelect}
            onItemRemove={handleItemRemove}
            onItemClear={handleItemClear}
          />
        </>
      )}

      {entities.some((e) => e.value === DataEntityTypes.DEVOPS) && (
        <Deployment
          provider={provider}
          entities={entities}
          transformation={transformation}
          connection={connection}
          onSettingsChange={onSettingsChange}
          isSaving={isSaving || isSavingConnection}
        />
      )}

      {entities.some((e) => e.value === DataEntityTypes.CODE_REVIEW) && (
        <>
          <h5>Additional Settings</h5>
          <div>
            <Checkbox
              checked={enableAdditionalCalculations}
              label='Enable calculation of commit and issue difference'
              onChange={(e) =>
                handleAdditionalEnable(!enableAdditionalCalculations)
              }
            />
            {enableAdditionalCalculations && (
              <>
                <div className='additional-settings-refdiff'>
                  <FormGroup
                    disabled={isSaving || isSavingConnection}
                    inline={true}
                    label='Tags Limit'
                    className='formGroup'
                    contentClassName='formGroupContent'
                  >
                    <NumericInput
                      id='refdiff-tags-limit'
                      disabled={isSaving || isSavingConnection}
                      fill={true}
                      placeholder='10'
                      allowNumericCharactersOnly={true}
                      onValueChange={(tagsLimitNumeric) =>
                        onSettingsChange({
                          refdiff: {
                            ...transformation?.refdiff,
                            tagsLimit: tagsLimitNumeric
                          }
                        })
                      }
                      value={transformation?.refdiff?.tagsLimit}
                    />
                  </FormGroup>
                  <FormGroup
                    disabled={isSaving || isSavingConnection}
                    inline={true}
                    label='Tags Pattern'
                    className='formGroup'
                    contentClassName='formGroupContent'
                  >
                    <InputGroup
                      id='refdiff-tags-pattern'
                      placeholder='(regex)$'
                      value={transformation?.refdiff?.tagsPattern}
                      onChange={(e) =>
                        onSettingsChange({
                          refdiff: {
                            ...transformation?.refdiff,
                            tagsPattern: e.target.value
                          }
                        })
                      }
                      disabled={isSaving || isSavingConnection}
                      className='input'
                      maxLength={255}
                    />
                  </FormGroup>
                  <FormGroup
                    disabled={isSaving || isSavingConnection}
                    inline={true}
                    label='Tags Order'
                    className='formGroup'
                    contentClassName='formGroupContent'
                  >
                    <InputGroup
                      id='refdiff-tags-order'
                      placeholder='reverse semver'
                      value={transformation?.refdiff?.tagsOrder}
                      onChange={(e) =>
                        onSettingsChange({
                          refdiff: {
                            ...transformation?.refdiff,
                            tagsOrder: e.target.value
                          }
                        })
                      }
                      disabled={isSaving || isSavingConnection}
                      className='input'
                      maxLength={255}
                    />
                  </FormGroup>
                </div>
              </>
            )}
          </div>
        </>
      )}
    </>
  )
}
