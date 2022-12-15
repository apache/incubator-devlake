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

import React, { useState, useEffect } from 'react'

import { Dialog } from '@/components'
import { SyncPolicy } from '@/plugins'

interface Props {
  isManual: boolean
  cronConfig: string
  skipOnFail: boolean
  createdDateAfter: string | null
  saving: boolean
  onCancel: () => void
  onSubmit: (params: any) => Promise<void>
}

export const UpdatePolicyDialog = ({
  saving,
  onCancel,
  onSubmit,
  ...props
}: Props) => {
  const [isManual, setIsManual] = useState(false)
  const [cronConfig, setCronConfig] = useState('')
  const [skipOnFail, setSkipOnFail] = useState(false)
  const [createdDateAfter, setCreatedDateAfter] = useState<string | null>(null)

  useEffect(() => {
    setIsManual(props.isManual)
    setCronConfig(props.cronConfig)
    setSkipOnFail(props.skipOnFail)
    setCreatedDateAfter(props.createdDateAfter)
  }, [props])

  const handleSubmit = () => {
    onSubmit({
      isManual,
      cronConfig,
      skipOnFail,
      createdDateAfter
    })
  }

  return (
    <Dialog
      isOpen
      title='Change Sync Policy'
      style={{
        width: 720
      }}
      okText='Save'
      okLoading={saving}
      onCancel={onCancel}
      onOk={handleSubmit}
    >
      <SyncPolicy
        isManual={isManual}
        cronConfig={cronConfig}
        skipOnFail={skipOnFail}
        createdDateAfter={createdDateAfter}
        onChangeIsManual={setIsManual}
        onChangeCronConfig={setCronConfig}
        onChangeSkipOnFail={setSkipOnFail}
        onChangeCreatedDateAfter={setCreatedDateAfter}
      />
    </Dialog>
  )
}
