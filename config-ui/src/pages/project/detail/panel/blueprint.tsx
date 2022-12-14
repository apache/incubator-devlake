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
import { useHistory } from 'react-router-dom'
import { Button, Intent } from '@blueprintjs/core'

import NoData from '@/images/no-data.svg'
import { Card } from '@/components'
import { BlueprintDetail } from '@/pages'

import type { ProjectType } from '../types'

interface Props {
  project: ProjectType
}

export const BlueprintPanel = ({ project }: Props) => {
  const history = useHistory()

  const handleGoCreateBlueprint = () =>
    history.push(`/projects/${project.name}/create-blueprint`)

  return !project.blueprint ? (
    <Card>
      <div className='blueprint'>
        <div className='logo'>
          <img src={NoData} alt='' />
        </div>
        <div className='desc'>
          <p>Create a blueprint to collect data from data sources.</p>
        </div>
        <div className='action'>
          <Button
            intent={Intent.PRIMARY}
            icon='plus'
            text='Create a Blueprint'
            onClick={handleGoCreateBlueprint}
          />
        </div>
      </div>
    </Card>
  ) : (
    <BlueprintDetail id={project.blueprint.id} />
  )
}
