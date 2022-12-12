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

import React, { useState, useEffect, useMemo } from 'react'
import {
  RadioGroup,
  Radio,
  Checkbox,
  InputGroup,
  Icon,
  Position
} from '@blueprintjs/core'
import { Popover2 } from '@blueprintjs/popover2'

import { Divider } from '@/components'
import CronHelp from '@/images/cron-help.png'

import { useBlueprint } from '../hooks'

import * as S from './styled'
import StartFromSelector from "@/components/blueprints/StartFromSelector";

const cronPresets = [
  {
    label: 'Hourly',
    config: '59 * * * *',
    desc: 'At minute 59 on every day-of-week from Monday through Sunday.'
  },
  {
    label: 'Daily',
    config: '0 0 * * *',
    desc: 'At 00:00 (Midnight) on every day-of-week from Monday through Sunday.'
  },
  {
    label: 'Weekly',
    config: '0 0 * * 1',
    desc: 'At 00:00 (Midnight) on Monday.'
  },
  {
    label: 'Monthly',
    config: '0 0 1 * *',
    desc: 'At 00:00 (Midnight) on day-of-month 1.'
  }
]

export const StepFour = () => {
  const [frequency, setFrequency] = useState('')

  const {
    cronConfig,
    isManual,
    skipOnFail,
    createdDateAfter,
    onChangeCronConfig,
    onChangeIsManual,
    onChangeSkipOnFail,
    onChangeCreatedDateAfter
  } = useBlueprint()

  const description = useMemo(() => {
    switch (frequency) {
      case 'manual':
        return 'Manual'
      case 'custom':
        return 'Custom'
      default:
        return cronPresets.find((it) => it.config === frequency)?.desc ?? ''
    }
  }, [frequency])

  useEffect(() => {
    if (isManual) {
      setFrequency('manual')
    } else if (frequency !== 'custom') {
      setFrequency(cronConfig)
    }
  }, [cronConfig, isManual])

  const handleChangeFrequency = (e: React.FormEvent<HTMLInputElement>) => {
    const value = (e.target as HTMLInputElement).value
    if (value === 'manual') {
      onChangeIsManual(true)
    } else if (value === 'custom') {
      onChangeIsManual(false)
      setFrequency('custom')
    } else {
      onChangeIsManual(false)
      onChangeCronConfig(value)
      setFrequency(value)
    }
  }

  const handleChangeCronConfig = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChangeCronConfig(e.target.value)
  }

  return (
    <S.Card>
      <h2>Set Sync Policy</h2>
      <Divider />
      <div className='block'>
          <h4>Time Filter *</h4>
          <p>Select the data range you wish to collect. DevLake will collect the last six months of data by default.</p>
          <StartFromSelector
            date={createdDateAfter}
            onSave={onChangeCreatedDateAfter}
          />
      </div>
      <div className='block'>
        <h3>Frequency</h3>
        <p>Blueprints will run recurringly based on the sync frequency.</p>
        <p style={{ margin: '10px 0' }}>{description}</p>
        <RadioGroup selectedValue={frequency} onChange={handleChangeFrequency}>
          <Radio label='Manual' value='manual' />
          {cronPresets.map((cron) => (
            <Radio key={cron.label} label={cron.label} value={cron.config} />
          ))}
          <Radio label='Custom' value='custom' />
        </RadioGroup>
        {frequency === 'custom' && (
          <S.Input>
            <InputGroup value={cronConfig} onChange={handleChangeCronConfig} />
            <Popover2
              position={Position.RIGHT}
              content={
                <S.Help>
                  <div className='title'>
                    <Icon icon='help' />
                    <span>Cron Expression Format</span>
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
                  <img src={CronHelp} alt='' />
                </S.Help>
              }
            >
              <Icon
                icon='help'
                size={14}
                style={{ marginLeft: '10px', transition: 'none' }}
              />
            </Popover2>
          </S.Input>
        )}
      </div>
      <div className='block'>
        <h3>Running Policy</h3>
        <Checkbox
          label='Skip failed tasks (Recommended when collecting large volume of data, eg. 10+ GitHub repos/Jira boards)'
          checked={skipOnFail}
          onChange={(e) =>
            onChangeSkipOnFail((e.target as HTMLInputElement).checked)
          }
        />
        <p>
          A task is a unit of a pipeline. A pipeline is an execution of a
          blueprint. By default, when a task is failed, the whole pipeline will
          fail and all the data that has been collected will be discarded. By
          skipping failed tasks, the pipeline will continue to run, and the data
          collected by other tasks will not be affected. After the pipeline is
          finished, you can rerun these failed tasks.
        </p>
      </div>
    </S.Card>
  )
}
