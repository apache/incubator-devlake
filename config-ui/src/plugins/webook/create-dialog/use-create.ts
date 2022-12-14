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

import { useState, useMemo } from 'react'

import { operator } from '@/utils'

import * as API from '../api'

export interface UseCreateProps {
  onSubmitAfter?: (id: ID) => void
}

export const useCreate = ({ onSubmitAfter }: UseCreateProps) => {
  const [saving, setSaving] = useState(false)
  const [step, setStep] = useState(1)
  const [name, setName] = useState('')
  const [record, setRecord] = useState({
    postIssuesEndpoint: '',
    closeIssuesEndpoint: '',
    postDeploymentsCurl: ''
  })

  const prefix = useMemo(() => `${window.location.origin}/api`, [])

  const handleCreate = async () => {
    const [success, res] = await operator(
      async () => {
        const res = await API.createConnection({ name })
        return API.getConnection(res.id)
      },
      {
        setOperating: setSaving
      }
    )

    if (success) {
      setStep(2)
      setRecord({
        postIssuesEndpoint: `${prefix}${res.postIssuesEndpoint}`,
        closeIssuesEndpoint: `${prefix}${res.closeIssuesEndpoint}`,
        postDeploymentsCurl: `curl ${prefix}${res.postPipelineDeployTaskEndpoint} -X 'POST' -d "{
        \\"commit_sha\\":\\"the sha of deployment commit\\",
        \\"repo_url\\":\\"the repo URL of the deployment commit\\",
        \\"start_time\\":\\"Optional, eg. 2020-01-01T12:00:00+00:00\\"
      }"`
      })
      onSubmitAfter?.(res.id)
    }
  }

  return useMemo(
    () => ({
      saving,
      step,
      name,
      record,
      onChangeName: setName,
      onSubmit: handleCreate
    }),
    [saving, step, name, record]
  )
}
