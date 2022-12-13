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

import React, { useState } from 'react'
import type { TabId } from '@blueprintjs/core'
import { Tabs, Tab } from '@blueprintjs/core'

// TO-DO: use new panel to replace it
import Status from '@/pages/blueprints/blueprint-detail'
import Configuration from '@/pages/blueprints/blueprint-settings'

import { PageLoading } from '@/components'

import type { UseDetailProps } from './use-detail'
import { useDetail } from './use-detail'

interface Props extends UseDetailProps {}

export const BlueprintDetail = ({ id }: Props) => {
  const [activeTab, setActiveTab] = useState<TabId>('configuration')

  const { loading, blueprint, saving, onUpdate } = useDetail({ id })

  if (loading || !blueprint) {
    return <PageLoading />
  }

  return (
    <Tabs selectedTabId={activeTab} onChange={(at) => setActiveTab(at)}>
      <Tab id='status' title='Status' panel={<Status id={id} />} />
      <Tab
        id='configuration'
        title='Configuration'
        panel={<Configuration id={id} />}
      />
    </Tabs>
  )
}
