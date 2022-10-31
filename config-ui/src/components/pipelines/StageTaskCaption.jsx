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
import { Colors } from '@blueprintjs/core'
import dayjs from '@/utils/time'
// import { Providers } from '@/data/Providers'

const StageTaskCaption = (props) => {
  const { task, options } = props

  return (
    <span
      className='task-module-caption'
      style={{
        opacity: 0.4,
        display: 'block',
        width: '100%',
        fontSize: '9px',
        overflow: 'hidden',
        whiteSpace: 'nowrap',
        textOverflow: 'ellipsis'
      }}
    >
      {(task.status === 'TASK_RUNNING' || task.status === 'TASK_COMPLETED') && (
        <span style={{ float: 'right' }}>
          {task.status === 'TASK_RUNNING'
            ? dayjs(task.beganAt).toNow(true)
            : dayjs(task.beganAt).from(task.finishedAt || task.updatedAt, true)}
        </span>
      )}
      {task.status === 'TASK_RUNNING' && (
        <span>
          Subtask {task?.progressDetail?.finishedSubTasks} /{' '}
          {task?.progressDetail?.totalSubTasks}
        </span>
      )}
      {task.status === 'TASK_COMPLETED' && (
        <span>
          {task?.progressDetail?.finishedSubTasks || 'All'} Subtasks completed
        </span>
      )}
      {task.status === 'TASK_COMPLETED' && (
        <span>{task?.progressDetail?.finishedRecords}</span>
      )}
      {task.status === 'TASK_CREATED' && <span>Records Pending</span>}
      {task.status === 'TASK_FAILED' && (
        <span style={{ color: Colors.RED3 }}>
          Task failed &mdash;{' '}
          <strong>{task?.failedSubTask || task?.message}</strong>
        </span>
      )}
    </span>
  )
}

export default StageTaskCaption
