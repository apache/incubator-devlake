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
  Icon,
  Spinner,
  Colors,
  Tooltip,
  Position,
  Intent,
} from '@blueprintjs/core'
import { ProviderIcons } from '@/data/Providers'

const StageTaskIndicator = (props) => {
  const { task } = props

  return (
    <div
      className='task-module-status'
      style={{
        display: 'flex',
        justifyContent: 'center',
        padding: '8px',
        width: '32px',
        minWidth: '32px',
        outline: 'none'
      }}
    >
      <span>
        {ProviderIcons[task.plugin] && ProviderIcons[task?.plugin](14, 14)}
      </span>
      {/* {task.status === 'TASK_COMPLETED' && (
        <Tooltip content={`Task Complete [STAGE ${task.pipelineRow}]`} position={Position.TOP} intent={Intent.SUCCESS}>
          <Icon icon='small-tick' size={18} color={Colors.GREEN5} style={{ marginLeft: '0', outline: 'none' }} />
        </Tooltip>
      )}
      {task.status === 'TASK_FAILED' && (
        <Tooltip content={`Task Failed [STAGE ${task.pipelineRow}]`} position={Position.TOP} intent={Intent.PRIMARY}>
          <Spinner
            className='task-module-spinner'
            size={14}
            intent={Intent.WARNING}
            value={task.progress}
          />
        </Tooltip>
      )}
      {task.status === 'TASK_RUNNING' && (
        <Tooltip content={`Task Running [STAGE ${task.pipelineRow}]`} position={Position.TOP}>
          <Spinner
            className='task-module-spinner'
            size={14}
            intent={task.status === 'TASK_COMPLETED' ? 'success' : 'warning'}
            value={task.status === 'TASK_COMPLETED' ? 1 : task.progress}
          />
        </Tooltip>
      )}
      {task.status === 'TASK_CREATED' && (
        <Tooltip content={`Task Created (Pending) [STAGE ${task.pipelineRow}]`} position={Position.TOP}>
          <Spinner
            className='task-module-spinner'
            size={14}
            value={0}
          />
        </Tooltip>
      )} */}
    </div>
  )
}

export default StageTaskIndicator
