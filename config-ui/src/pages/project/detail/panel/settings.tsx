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

import React, { useEffect, useState } from 'react'
import {
  InputGroup,
  Checkbox,
  ButtonGroup,
  Button,
  Intent
} from '@blueprintjs/core'

import type { ProjectType } from '../types'
import * as S from '../styled'

interface Props {
  project: ProjectType
  onUpdate: (name: string, enableDora: boolean) => void
}

export const SettingsPanel = ({ project, onUpdate }: Props) => {
  const [name, setName] = useState('')
  const [enableDora, setEnableDora] = useState(false)

  useEffect(() => {
    if (project) {
      setName(project.name)
      setEnableDora(project.enableDora)
    }
  }, [project])

  const handleSave = () => onUpdate(name, enableDora)

  return (
    <S.Panel>
      <div className='settings'>
        <div className='block'>
          <h3>Project Name *</h3>
          <p>Edit your project name.</p>
          <InputGroup value={name} onChange={(e) => setName(e.target.value)} />
        </div>
        <div className='block'>
          <Checkbox
            label='Enable DORA Metrics'
            checked={enableDora}
            onChange={(e) =>
              setEnableDora((e.target as HTMLInputElement).checked)
            }
          />
          <p>
            DORA metrics are four widely-adopted metrics for measuring software
            delivery performance.
          </p>
        </div>
        <ButtonGroup>
          <Button text='Save' intent={Intent.PRIMARY} onClick={handleSave} />
        </ButtonGroup>
      </div>
    </S.Panel>
  )
}
