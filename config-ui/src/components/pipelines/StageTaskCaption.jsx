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
import { Providers } from '@/data/Providers'

const StageTaskCaption = (props) => {
  const { task, options } = props

  return (
    <span
      className='task-module-caption'
      style={{
        opacity: 0.4,
        display: 'block',
        width: '90%',
        fontSize: '9px',
        overflow: 'hidden',
        whiteSpace: 'nowrap',
        textOverflow: 'ellipsis'
      }}
    >
      {(task.plugin === Providers.GITLAB || task.plugin === Providers.JIRA) && (<>ID {options.projectId || options.boardId}</>)}
      {task.plugin === Providers.GITHUB && (<>@{options.owner}/{options.repositoryName}</>)}
      {task.plugin === Providers.JENKINS && (<>Task #{task.ID}</>)}
      {(task.plugin === Providers.GITEXTRACTOR || task.plugin === Providers.REFDIFF) && (<>{options.repoId || `ID ${task.ID}`}</>)}
    </span>
  )
}

export default StageTaskCaption
