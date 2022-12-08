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
import { Button, Intent } from '@blueprintjs/core'

import NoData from '@/images/no-webhook.svg'

import type { ProjectType } from '../types'
import * as S from '../styled'

interface Props {
  project?: ProjectType
}

export const IncomingWebhooksPanel = ({ project }: Props) => {
  return (
    <S.Panel>
      <div className='webhook'>
        <div className='logo'>
          <img src={NoData} alt='' />
        </div>
        <div className='desc'>
          <p>
            Push `incidents` or `deployments` from your tools by incoming
            webhooks.
          </p>
        </div>
        <div className='action'>
          <Button intent={Intent.PRIMARY} icon='plus' text='Add a Webhook' />
          <span className='or'>or</span>
          <Button
            outlined
            intent={Intent.PRIMARY}
            text='Select Existing Webhooks'
          />
        </div>
      </div>
    </S.Panel>
  )
}
