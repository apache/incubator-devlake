import React, { useEffect, useState } from 'react'
import {
  BrowserRouter as Router,
  useHistory
} from 'react-router-dom'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import { integrationsData } from '@/data/integrations'

import '@/styles/integration.scss'

export default function Integration () {
  const history = useHistory()

  const [integrations, setIntegrations] = useState(integrationsData)

  const [activeProvider, setActiveProvider] = useState(integrations[0])
  const [invalidProvider, setInvalidProvider] = useState(false)

  const handleProviderClick = (providerId) => {
    const theProvider = integrations.find(p => p.id === providerId)
    if (theProvider) {
      setActiveProvider(theProvider)
      history.push(`/integrations/${theProvider.id}`)
    } else {
      setInvalidProvider(true)
      setActiveProvider(integrations[0])
    }
  }

  useEffect(() => {
    // Selected Provider
    console.log(activeProvider)
  }, [activeProvider, history])

  useEffect(() => {

  }, [])

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
                { href: '/integrations', icon: false, text: 'Integrations', current: true }
              ]}
            />
            <div className='headlineContainer'>
              <h1>Data Integrations</h1>
              <p className='page-description'>{integrationsData.length} connections are available for data collection.</p>
            </div>
            <div className='integrationProviders'>
              {integrations.map((provider) => (
                <div
                  className='iProvider'
                  key={`provider-${provider.id}`}
                  onClick={() => handleProviderClick(provider.id)}
                >
                  <div className='providerIcon'>
                    {provider.iconDashboard}
                  </div>
                  <div className='providerName'>
                    {provider.name}
                  </div>
                </div>
              ))}
            </div>
          </main>
        </Content>
      </div>
    </>
  )
}
