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

import React, { useMemo } from 'react'
import {
  RadioGroup,
  Radio,
  InputGroup,
  Icon,
  Checkbox,
  Position
} from '@blueprintjs/core'
import { Popover2 } from '@blueprintjs/popover2'

import StartFromSelector from '@/components/blueprints/StartFromSelector'
import { getCron, getCronOptions } from '@/config'
import CronHelp from '@/images/cron-help.png'

import * as S from './styled'

interface Props {
  isManual: boolean
  cronConfig: string
  skipOnFail: boolean
  createdDateAfter: string | null
  onChangeIsManual: (val: boolean) => void
  onChangeCronConfig: (val: string) => void
  onChangeSkipOnFail: (val: boolean) => void
  onChangeCreatedDateAfter: (val: string | null) => void
}

export const SyncPolicy = ({
  isManual,
  cronConfig,
  skipOnFail,
  createdDateAfter,
  onChangeIsManual,
  onChangeCronConfig,
  onChangeSkipOnFail,
  onChangeCreatedDateAfter
}: Props) => {
  const cron = useMemo(
    () => getCron(isManual, cronConfig),
    [isManual, cronConfig]
  )

  const options = useMemo(() => getCronOptions(), [])

  const handleChangeFrequency = (e: React.FormEvent<HTMLInputElement>) => {
    const value = (e.target as HTMLInputElement).value
    if (value === 'manual') {
      onChangeIsManual(true)
    } else if (value === 'custom') {
      onChangeIsManual(false)
      onChangeCronConfig('* * * * *')
    } else {
      onChangeIsManual(false)
      onChangeCronConfig(value)
    }
  }

  return (
    <S.Wrapper>
      <div className='block'>
        <h3>Time Filter *</h3>
        <p>
          Select the data range you wish to collect. DevLake will collect the
          last six months of data by default.
        </p>
        <StartFromSelector
          autoFillDefault={true}
          date={createdDateAfter}
          onSave={onChangeCreatedDateAfter}
        />
      </div>
      <div className='block'>
        <h3>Frequency</h3>
        <p>Blueprints will run recurringly based on the sync frequency.</p>
        <p style={{ margin: '10px 0' }}>{cron.description}</p>
        <RadioGroup selectedValue={cron.value} onChange={handleChangeFrequency}>
          {options.map(({ label, value }) => (
            <Radio label={label} value={value} />
          ))}
        </RadioGroup>
        {cron.value === 'custom' && (
          <S.Input>
            <InputGroup
              value={cronConfig}
              onChange={(e) => onChangeCronConfig(e.target.value)}
            />
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
    </S.Wrapper>
  )
}
