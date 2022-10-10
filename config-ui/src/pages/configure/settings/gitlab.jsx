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
import React, { useEffect } from 'react'
import { useHistory } from 'react-router-dom'
import { DataEntityTypes } from '@/data/DataEntities'
import Deployment from '@/components/blueprints/transformations/CICD/Deployment'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

export default function GitlabSettings(props) {
  const {
    connection,
    entities = [],
    transformation = {},
    entityIdKey,
    provider,
    configuredProject,
    isSaving = false,
    isSavingConnection = false,
    onSettingsChange = () => {}
  } = props

  useEffect(() => {
    console.log(
      '>>>> GITLAB: TRANSFORMATION SETTINGS OBJECT....',
      transformation
    )
  }, [transformation])

  return (
    <>
      {entities.some((e) => e.value === DataEntityTypes.DEVOPS) &&
      configuredProject ? (
        <Deployment
          provider={provider}
          entities={entities}
          entityIdKey={entityIdKey}
          transformation={transformation}
          connection={connection}
          onSettingsChange={onSettingsChange}
          isSaving={isSaving || isSavingConnection}
        />
      ) : (
        <></>
      )}
    </>
  )
}
