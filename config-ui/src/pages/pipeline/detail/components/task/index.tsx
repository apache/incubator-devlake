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
import { Icon, Intent } from '@blueprintjs/core'
import { Tooltip2 } from '@blueprintjs/popover2'

import type { PluginConfigType } from '@/plugins'
import { Plugins, PluginConfig } from '@/plugins'
import { duration } from '@/utils'

import { StatusEnum, TaskType } from '../../../types'
import { STATUS_CLS } from '../../../misc'

import * as S from './styled'

interface Props {
  task: TaskType
  operating: boolean
  onRerun: (id: ID) => void
}

export const Task = ({ task, operating, onRerun }: Props) => {
  const { beganAt, finishedAt, status, message, progressDetail } = task

  const [icon, name] = useMemo(() => {
    const config = PluginConfig.find((p) => p.plugin === task.plugin) as PluginConfigType
    const options = JSON.parse(task.options)

    let name = config.name

    switch (true) {
      case [Plugins.GitHub, Plugins.GitHubGraphql].includes(config.plugin):
        name = `${name}:${options.name}`
        break
      case Plugins.GitExtractor === config.plugin:
        name = `${name}:${options.repoId}`
        break
      case [Plugins.DORA, Plugins.RefDiff].includes(config.plugin):
        name = `${name}:${options.projectName}`
        break
      case Plugins.GitLab === config.plugin:
        name = `${name}:projectId:${options.projectId}`
        break
    }

    return [config.icon, name]
  }, [task])

  const statusCls = STATUS_CLS(status)

  const handleRerun = () => {
    if (operating) return
    onRerun(task.id)
  }

  return (
    <S.Wrapper>
      <S.Info>
        <div className='title'>
          <img src={icon} alt='' />
          <strong>Task{task.id}</strong>
          <span title={name}>{name}</span>
        </div>
        {[status === StatusEnum.CREATED, StatusEnum.PENDING].includes(
          status
        ) && <p className={statusCls}>Subtasks pending</p>}

        {[StatusEnum.ACTIVE, StatusEnum.RUNNING].includes(status) && (
          <p className={statusCls}>
            Subtasks running
            <strong style={{ marginLeft: 8 }}>
              {progressDetail?.finishedSubTasks}/{progressDetail?.totalSubTasks}
            </strong>
          </p>
        )}

        {status === StatusEnum.COMPLETED && (
          <p className={statusCls}>All Subtasks completed</p>
        )}

        {status === StatusEnum.FAILED && (
          <Tooltip2 content={message} intent={Intent.DANGER}>
            <p className={statusCls}>Task failed: hover to view the reason</p>
          </Tooltip2>
        )}

        {status === StatusEnum.CANCELLED && (
          <p className={statusCls}>Subtasks canceled</p>
        )}
      </S.Info>
      <S.Duration>
        {[
          StatusEnum.COMPLETED,
          StatusEnum.FAILED,
          StatusEnum.CANCELLED
        ].includes(status) && <Icon icon='repeat' onClick={handleRerun} />}
        <span>{duration(beganAt, finishedAt)}</span>
      </S.Duration>
    </S.Wrapper>
  )
}
