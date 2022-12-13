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
import { useParams } from 'react-router-dom'

import { PageHeader, PageLoading } from '@/components'

import { useDetail } from './use-detail'
import { BlueprintDetail } from './blueprint-detail'

export const BlueprintDetailPage = () => {
  const { id } = useParams<{ id: string }>()

  const { loading, blueprint } = useDetail({ id })

  if (loading || !blueprint) {
    return <PageLoading />
  }

  return (
    <PageHeader
      breadcrumbs={[
        { name: 'Blueprints', path: '/blueprints' },
        { name: blueprint.name, path: `/blueprints/${blueprint.id}` }
      ]}
    >
      <BlueprintDetail id={blueprint.id} />
    </PageHeader>
  )
}
