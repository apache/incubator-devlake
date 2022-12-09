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
import { useParams } from 'react-router-dom'

import { PageHeader, PageLoading } from '@/components'

import { useProject } from './use-project'
import { BlueprintPanel, IncomingWebhooksPanel, SettingsPanel } from './panel'
import * as S from './styled'

export const ProjectDetailPage = () => {
  const { pname } = useParams<{ pname: string }>()

  const [activeTab, setActiveTab] = useState<'bp' | 'iw' | 'st'>('bp')

  const { loading, project, updateProject } = useProject(pname)

  if (loading) {
    return <PageLoading />
  }

  return (
    <PageHeader
      breadcrumbs={[
        { name: 'Projects', path: '/projects' },
        { name: pname, path: `/projects/${pname}` }
      ]}
    >
      <S.Tabs>
        <S.Tab active={activeTab === 'bp'} onClick={() => setActiveTab('bp')}>
          Blueprint
        </S.Tab>
        <S.Tab
          active={activeTab === 'iw'}
          disabled={!project?.blueprint}
          onClick={() => {
            if (!project?.blueprint) {
              return
            }
            setActiveTab('iw')
          }}
        >
          Incoming Webhooks
        </S.Tab>
        <S.Tab active={activeTab === 'st'} onClick={() => setActiveTab('st')}>
          Settings
        </S.Tab>
      </S.Tabs>
      {activeTab === 'bp' && <BlueprintPanel name={pname} project={project} />}
      {activeTab === 'iw' && <IncomingWebhooksPanel project={project} />}
      {activeTab === 'st' && (
        <SettingsPanel project={project} onUpdate={updateProject} />
      )}
    </PageHeader>
  )
}
