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

import { useState, useEffect, useMemo } from 'react'
import { useHistory } from 'react-router-dom'

import { operator } from '@/utils'

import type { ProjectType } from './types'
import * as API from './api'

export const useProject = (name: string) => {
  const [loading, setLoading] = useState(false)
  const [project, setProject] = useState<ProjectType>(null)

  const history = useHistory()

  const getProject = async (name: string) => {
    setLoading(true)
    try {
      const res = await API.getProject(name)
      const doraMetrics = res.metrics.find(
        (ms: any) => ms.pluginName === 'dora'
      )

      setProject({
        name: res.name,
        description: res.description,
        blueprint: res.blueprint,
        enableDora: doraMetrics.enable
      })
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    getProject(name)
  }, [name])

  const updateProject = async (newName: string, enableDora: boolean) => {
    const payload = {
      name: newName,
      description: '',
      metrics: [
        {
          pluginName: 'dora',
          pluginOption: '',
          enable: enableDora
        }
      ]
    }

    const [success] = await operator(() => API.updateProject(name, payload))

    if (success) {
      history.push(`/projects/${newName}`)
    }
  }

  return useMemo(
    () => ({
      loading,
      project,
      updateProject
    }),
    [loading, project]
  )
}
