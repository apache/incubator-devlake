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

import React, { useState, useMemo, useCallback } from 'react'
import { InputGroup, Icon } from '@blueprintjs/core'
import { CopyToClipboard } from 'react-copy-to-clipboard'

import { Dialog, toast } from '@/components'

import * as S from './styled'

const prefix = `${window.location.origin}/api`

interface Props {
  isOpen: boolean
  saving: boolean
  onSubmit: (name: string) => Promise<any>
  onCancel: () => void
}

export const CreateDialog = ({ isOpen, saving, onSubmit, onCancel }: Props) => {
  const [step, setStep] = useState(1)
  const [name, setName] = useState('')
  const [record, setRecord] = useState({
    postIssuesEndpoint: '',
    closeIssuesEndpoint: '',
    postDeploymentsCurl: ''
  })

  const [okText, okDisabled] = useMemo(
    () => [step === 1 ? 'Generate POST URL' : 'Done', step === 1 && !name],
    [step, name]
  )

  const handleSubmit = useCallback(async () => {
    if (step === 1) {
      const res = await onSubmit(name)
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
    } else {
      onCancel()
    }
  }, [step, name])

  return (
    <Dialog
      isOpen={isOpen}
      title='Add a New Incoming Webhook'
      style={{ width: 600, top: -100 }}
      okText={okText}
      okDisabled={okDisabled}
      okLoading={saving}
      onCancel={onCancel}
      onOk={handleSubmit}
    >
      {step === 1 && (
        <S.Wrapper>
          <h3>Incoming Webhook Name *</h3>
          <p>
            Give your Incoming Webhook a unique name to help you identify it in
            the future.
          </p>
          <InputGroup value={name} onChange={(e) => setName(e.target.value)} />
        </S.Wrapper>
      )}
      {step === 2 && (
        <S.Wrapper>
          <h2>
            <Icon icon='endorsed' size={30} />
            <span>POST URL Generated!</span>
          </h2>
          <h3>POST URL</h3>
          <p>
            Copy the following URLs to your issue tracking tool for Incidents
            and CI tool for Deployments by making a POST to DevLake.
          </p>
          <h3>Incident</h3>
          <p>POST to register an incident</p>
          <div className='block'>
            <span>{record.postIssuesEndpoint}</span>
            <CopyToClipboard
              text={record.postIssuesEndpoint}
              onCopy={() => toast.success('Copy successfully.')}
            >
              <Icon icon='clipboard' />
            </CopyToClipboard>
          </div>
          <p>POST to close a registered incident</p>
          <div className='block'>
            <span>{record.closeIssuesEndpoint}</span>
            <CopyToClipboard
              text={record.closeIssuesEndpoint}
              onCopy={() => toast.success('Copy successfully.')}
            >
              <Icon icon='clipboard' />
            </CopyToClipboard>
          </div>
          <h3>Deployment</h3>
          <p>POST to register a deployment</p>
          <div className='block'>
            <span>{record.postDeploymentsCurl}</span>
            <CopyToClipboard
              text={record.postDeploymentsCurl}
              onCopy={() => toast.success('Copy successfully.')}
            >
              <Icon icon='clipboard' />
            </CopyToClipboard>
          </div>
        </S.Wrapper>
      )}
    </Dialog>
  )
}
