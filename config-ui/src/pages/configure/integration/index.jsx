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
import React, { useEffect, useState } from 'react'
import { useHistory } from 'react-router-dom'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import { integrationsData } from '@/data/integrations'
import { ReactComponent as WebHookProviderIcon } from '@/images/integrations/webhook.svg'

import '@/styles/integration.scss'

export default function Integration() {
  const history = useHistory()

  const [activeProvider, setActiveProvider] = useState(integrationsData[0])

  const handleProviderClick = (providerId) => {
    const theProvider = integrationsData.find((p) => p.id === providerId)
    if (theProvider) {
      setActiveProvider(theProvider)
      history.push(`/integrations/${theProvider.id}`)
    } else {
      setActiveProvider(integrationsData[0])
    }
  }

  useEffect(() => {
    // Selected Provider
    console.log(activeProvider)
  }, [activeProvider, history])

  useEffect(() => {}, [])

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
                { href: '/integrations', icon: false, text: 'Integrations', current: true },
              ]}
            />
            <div className='headlineContainer'>
              <h1>Data Integrations</h1>
              <p className='page-description'>{integrationsData.length} connections are available for data collection.</p>
            </div>
            <div className='integrationProviders'>
              {integrationsData.map((provider) => (
                <div
                  className='iProvider'
                  key={`provider-${provider.id}`}
                  onClick={() => handleProviderClick(provider.id)}
                >
                  <div className='providerIcon'>
                    {provider.iconDashboard}
                  </div>
                  <div className='providerName'>
                    {provider.name} {provider.isBeta && <><sup>(beta)</sup></>}
                  </div>
                </div>
              ))}
            </div>
            <div className='headlineContainer'>
              <h1>Webhooks</h1>
              <p className='page-description'>
                You can use Webhooks to define Issues and Deployments to be used in calculating DORA metrics. Please note: Webhooks cannot
                be created or managed in Blueprints.
              </p>
            </div>
            <div className='integrationProviders'>
              <div className='iProvider' style={{ width: 130 }} onClick={() => history.push('/connections/webhook')}>
                <div className='providerIcon'>
                  <WebHookProviderIcon className='providerIconSvg' width='40' height='40' />
                </div>
                <div className='providerName'>Issue/Deployment Webhook</div>
              </div>
            </div>
          </main>
        </Content>
      </div>
    </>
  )
}
