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
import { Colors, Icon } from '@blueprintjs/core'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
// @todo: replace with Integrations Hook
// import { integrationsData } from '@/data/integrations'
import useIntegrations from '@/hooks/useIntegrations'
import { ReactComponent as WebHookProviderIcon } from '@/images/integrations/incoming-webhook.svg'

import '@/styles/integration.scss'

export default function Integration() {
  const history = useHistory()

  const {
    registry,
    plugins: Plugins,
    integrations: Integrations,
    activeProvider,
    setActiveProvider
  } = useIntegrations()
  // const [activeProvider, setActiveProvider] = useState(integrationsData[0])

  const handleProviderClick = (providerId) => {
    const theProvider = Plugins.find((p) => p.id === providerId)
    if (theProvider) {
      setActiveProvider(theProvider)
      history.push(`/integrations/${theProvider.id}`)
    } else {
      setActiveProvider(Plugins.find[0])
    }
  }

  // useEffect(() => {
  //   // Selected Provider
  //   console.log(activeProvider)
  // }, [activeProvider, history])

  // useEffect(() => {}, [])

  return (
    <>
      <div className='container'>
        <Nav />
        <Sidebar key={Integrations} integrations={Integrations} />
        <Content>
          <main className='main'>
            <AppCrumbs
              items={[
                { href: '/', icon: false, text: 'Dashboard' },
                {
                  href: '/integrations',
                  icon: false,
                  text: 'Connections',
                  current: true
                }
              ]}
            />
            <div className='headlineContainer'>
              <h1>Data Connections</h1>
              <p className='page-description'>
                {Integrations.length} connections are available for data
                collection.
              </p>
            </div>
            <div className='integrationProviders'>
              {Integrations.map((provider) => (
                <div
                  className='iProvider'
                  key={`provider-${provider.id}`}
                  onClick={() => handleProviderClick(provider.id)}
                  style={{ position: 'relative' }}
                >
                  {provider?.private && (
                    <span
                      style={{
                        width: '20px',
                        height: '20px',
                        position: 'absolute',
                        top: '-5px',
                        right: '-5px',
                        textAlign: 'center',
                        lineHeight: '16px',
                        backgroundColor: '#fff',
                        display: 'block',
                        borderRadius: '50%',
                        border: '1px solid #eee'
                      }}
                    >
                      <Icon
                        icon='lock'
                        size={10}
                        style={{ color: Colors.RED5 }}
                      />
                    </span>
                  )}
                  <div className='providerIcon'>
                    <img
                      className='providerIconSvg'
                      src={provider.icon}
                      width={40}
                      height={40}
                      style={{ width: '40px', height: '40px' }}
                    />
                  </div>
                  <div className='providerName'>
                    {provider.name}{' '}
                    {provider.isBeta && (
                      <>
                        <sup>(beta)</sup>
                      </>
                    )}
                  </div>
                </div>
              ))}
            </div>
            <div className='headlineContainer'>
              <h1>Webhooks</h1>
              <p className='page-description'>
                You can use webhooks to import deployments and incidents from
                the unsupported data integrations to calculate DORA metrics,
                etc. Please note: webhooks cannot be created or managed in
                Blueprints.
              </p>
            </div>
            <div className='integrationProviders'>
              <div
                className='iProvider'
                style={{ width: 130 }}
                onClick={() => history.push('/connections/incoming-webhook')}
              >
                <div className='providerIcon'>
                  <WebHookProviderIcon
                    className='providerIconSvg'
                    width='40'
                    height='40'
                  />
                </div>
                <div className='providerName'>
                  Issue/Deployment Incoming Webhook
                </div>
              </div>
            </div>
          </main>
        </Content>
      </div>
    </>
  )
}
