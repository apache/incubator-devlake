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
import {
  Checkbox,
  Icon,
  Intent,
  Tooltip,
  FormGroup,
  InputGroup,
  Position,
  RadioGroup,
  Popover,
  Radio,
  Divider,
  Elevation,
  Card,
  Colors
} from '@blueprintjs/core'

import InputValidationError from '@/components/validation/InputValidationError'

import CronHelp from '@/images/cron-help.png'
import StartFromSelector from '@/components/blueprints/StartFromSelector'

const DataSync = (props) => {
  const {
    skipOnFail,
    setSkipOnFail,
    createdDateAfter,
    setCreatedDateAfter,
    activeStep,
    cronConfig,
    customCronConfig,
    setCronConfig = () => {},
    getCronPreset = () => {},
    fieldHasError = () => {},
    getFieldError = () => {},
    createCron = () => {},
    setCustomCronConfig = () => {},
    getCronPresetByConfig = () => {},
    advancedMode = false,
    enableHeader = true,
    elevation = Elevation.TWO,
    cardStyle = {}
  } = props

  return (
    <div
      className='workflow-step workflow-step-set-sync-frequency'
      data-step={activeStep?.id}
    >
      <Card
        className='workflow-card'
        elevation={elevation}
        style={{ ...cardStyle }}
      >
        {enableHeader && (
          <>
            <h3 style={{ marginBottom: '8px' }}>Set Sync Frequency</h3>
            {getCronPresetByConfig(cronConfig) ? (
              <p
                style={{
                  display: 'block'
                }}
              >
                <strong>Automated</strong> &mdash;{' '}
                {getCronPresetByConfig(cronConfig).description}
              </p>
            ) : (
              <small
                style={{
                  fontSize: '10px',
                  color: Colors.GRAY2,
                  textTransform: 'uppercase'
                }}
              >
                {cronConfig}
              </small>
            )}
            <Divider className='section-divider' />
          </>
        )}

        <h4>Time Filter *</h4>
        <p>
          Select the data range you wish to collect. DevLake will collect the
          last six months of data by default.
        </p>
        <StartFromSelector
          date={createdDateAfter}
          onSave={setCreatedDateAfter}
        />

        <h4>Frequency</h4>
        <p>Blueprints will run recurringly based on the sync frequency.</p>

        <RadioGroup
          inline={false}
          label={false}
          name='blueprint-frequency'
          onChange={(e) => setCronConfig(e.target.value)}
          selectedValue={cronConfig}
          required
        >
          <Radio
            label='Manual'
            value='manual'
            style={{
              fontWeight: cronConfig === 'manual' ? 'bold' : 'normal'
            }}
          />
          {/* Dynamic Presets from Connection Manager */}
          {[
            getCronPreset('hourly'),
            getCronPreset('daily'),
            getCronPreset('weekly'),
            getCronPreset('monthly')
          ].map((preset, prIdx) => (
            <Radio
              key={`cron-preset-tooltip-key${prIdx}`}
              label={
                <>
                  <Tooltip
                    position={Position.RIGHT}
                    intent={Intent.PRIMARY}
                    content={preset.description}
                  >
                    {preset.label}
                  </Tooltip>
                </>
              }
              value={preset.cronConfig}
              style={{
                fontWeight:
                  cronConfig === preset.cronConfig ? 'bold' : 'normal',
                outline: 'none !important'
              }}
            />
          ))}
          <Radio
            label='Custom'
            value='custom'
            style={{
              fontWeight: cronConfig === 'custom' ? 'bold' : 'normal'
            }}
          />
        </RadioGroup>
        <div
          style={{
            display: cronConfig === 'custom' ? 'flex' : 'none'
          }}
        >
          <FormGroup
            disabled={cronConfig !== 'custom'}
            label=''
            inline={true}
            labelFor='cron-custom'
            className='formGroup-inline'
            contentClassName='formGroupContent'
            style={{ marginBottom: '5px' }}
            fill={false}
          >
            <InputGroup
              id='cron-custom'
              inline={true}
              fill={false}
              readOnly={cronConfig !== 'custom'}
              leftElement={
                cronConfig !== 'custom' ? (
                  <Icon
                    icon='lock'
                    size={11}
                    style={{
                      alignSelf: 'center',
                      margin: '4px 10px -2px 6px'
                    }}
                  />
                ) : null
              }
              rightElement={
                <>
                  <InputValidationError
                    error={getFieldError('Blueprint Cron')}
                  />
                </>
              }
              placeholder='Enter Crontab Syntax'
              value={cronConfig !== 'custom' ? cronConfig : customCronConfig}
              onChange={(e) => setCustomCronConfig(e.target.value)}
              className={`cron-custom-input ${
                fieldHasError('Blueprint Cron') ? 'invalid-field' : ''
              }`}
              style={{ transition: 'none' }}
            />
          </FormGroup>
          <div
            style={{
              display: 'inline',
              marginTop: 'auto',
              paddingBottom: '15px'
            }}
          >
            <Popover
              className='trigger-crontab-help'
              popoverClassName='popover-help-crontab'
              position={Position.RIGHT}
              autoFocus={false}
              enforceFocus={false}
              usePortal={false}
            >
              <a rel='noreferrer'>
                <Icon
                  icon='help'
                  size={14}
                  style={{ marginLeft: '10px', transition: 'none' }}
                />
              </a>
              <>
                <div
                  style={{
                    textShadow: 'none',
                    fontSize: '12px',
                    padding: '12px',
                    maxWidth: '300px'
                  }}
                >
                  <div
                    style={{
                      marginBottom: '10px',
                      fontWeight: 700,
                      fontSize: '14px'
                    }}
                  >
                    <Icon
                      icon='help'
                      size={16}
                      style={{ marginRight: '5px' }}
                    />{' '}
                    Cron Expression Format
                  </div>
                  <p>
                    Need Help? &mdash; For additional information on{' '}
                    <strong>Crontab</strong>, please reference the{' '}
                    <a
                      href='https://man7.org/linux/man-pages/man5/crontab.5.html'
                      rel='noreferrer'
                      target='_blank'
                      style={{ textDecoration: 'underline' }}
                    >
                      Crontab Linux manual
                    </a>
                    .
                  </p>
                  <img
                    src={CronHelp}
                    style={{
                      border: 0,
                      margin: 0,
                      maxWidth: '100%'
                    }}
                  />
                </div>
              </>
            </Popover>
          </div>
        </div>

        <h4>Running Policy</h4>
        <div style={{ marginTop: 20 }}>
          <Checkbox
            label='Skip failed tasks (Recommended when collecting large volume of data, eg. 10+ GitHub repos/Jira boards)'
            checked={skipOnFail}
            onChange={(e) => setSkipOnFail(e.target.checked)}
          />
          <p>
            A task is a unit of a pipeline. A pipeline is an execution of a
            blueprint. By default, when a task is failed, the whole pipeline
            will fail and all the data that has been collected will be
            discarded. By skipping failed tasks, the pipeline will continue to
            run, and the data collected by other tasks will not be affected.
            After the pipeline is finished, you can rerun these failed tasks.
          </p>
        </div>
      </Card>
    </div>
  )
}

export default DataSync
