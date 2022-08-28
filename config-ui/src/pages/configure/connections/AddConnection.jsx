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
import React, { useEffect, useMemo, useState } from 'react'
import { Link, useHistory, useParams } from 'react-router-dom'
import { Icon, } from '@blueprintjs/core'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import ConnectionForm from '@/pages/configure/connections/ConnectionForm'
import { integrationsData } from '@/data/integrations'
import {
  ProviderConfigMap,
  ProviderConnectionLimits,
  ProviderFormLabels,
  ProviderFormPlaceholders,
  ProviderLabels,
  Providers
} from '@/data/Providers'

import useConnectionManager from '@/hooks/useConnectionManager'
import useConnectionValidation from '@/hooks/useConnectionValidation'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

export default function AddConnection () {
  const history = useHistory()
  const { providerId } = useParams()

  const activeProvider = useMemo(() => integrationsData.find(p => p.id === providerId), [providerId])

  const {
    testConnection,
    saveConnection,
    errors,
    isSaving,
    isTesting,
    showError,
    testStatus,
    testResponse,
    allTestResponses,
    fetchAllConnections,
    connectionLimitReached,
    // Providers

    editingConnection,
    setConnectionColumn,
    initialTokenStore,
  } = useConnectionManager({
    provider: activeProvider,
  })

  const {
    validate,
    errors: validationErrors,
    isValid: isValidForm
  } = useConnectionValidation({
    activeProvider,
  })

  const cancel = () => {
    history.push(`/integrations/${activeProvider.id}`)
  }

  useEffect(() => {
    // Selected Provider
    if (activeProvider?.id) {
      fetchAllConnections()
    }
  }, [activeProvider.id, fetchAllConnections])

  return (
    <>
      <div className='container'>
        <Nav />
        <Sidebar />
        <Content>
          <main className='main'>
            <AppCrumbs
              items={[
                { href: '/', icon: false, text: 'Dashboard' },
                { href: '/integrations', icon: false, text: 'Integrations' },
                { href: `/integrations/${activeProvider.id}`, icon: false, text: `${activeProvider.name}` },
                { href: `/connections/add/${activeProvider.id}`, icon: false, text: 'Add Connection', current: true }
              ]}
            />
            <div style={{ width: '100%' }}>
              <Link style={{ float: 'right', marginLeft: '10px', color: '#777777' }} to={`/integrations/${activeProvider.id}`}>
                <Icon icon='undo' size={16} /> Go Back
              </Link>
              <div style={{ display: 'flex' }}>
                <div>
                  <span style={{ marginRight: '10px' }}>{activeProvider.icon}</span>
                </div>
                <div>
                  <h1 style={{ margin: 0 }}>
                    {activeProvider.name} Add Connection
                  </h1>
                  <p className='page-description'>Create a new connection for this provider.</p>
                </div>
              </div>
              <div className='addConnection' style={{ display: 'flex' }}>
                <ConnectionForm
                  isLocked={connectionLimitReached}
                  isValid={isValidForm}
                  validationErrors={validationErrors}
                  activeProvider={activeProvider}
                  editingConnection={editingConnection}
                  onConnectionColumnChange={setConnectionColumn}
                  initialTokenStore={initialTokenStore}
                  onSave={() => saveConnection({})}
                  onTest={testConnection}
                  onCancel={cancel}
                  onValidate={validate}
                  isSaving={isSaving}
                  isTesting={isTesting}
                  testStatus={testStatus}
                  testResponse={testResponse}
                  allTestResponses={allTestResponses}
                  errors={errors}
                  showError={showError}
                  activeProviderConfig={ProviderConfigMap[activeProvider.id]}
                />
              </div>
              {/* {validationErrors.length > 0 && (
                <FormValidationErrors errors={validationErrors} />
              )} */}
            </div>
          </main>
        </Content>
      </div>
    </>
  )
}
