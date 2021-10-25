import React, { useEffect, useState } from 'react'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  useParams,
  Link,
  useHistory
} from 'react-router-dom'
import {
  Alignment,
  Icon,
} from '@blueprintjs/core'
// import { FormGroup, InputGroup, Button, Tooltip, Position, Label } from '@blueprintjs/core'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import Content from '@/components/Content'
import SaveAlert from '@/components/SaveAlert'
import { SERVER_HOST } from '@/utils/config'

import { ReactComponent as GitlabProvider } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProvider } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProvider } from '@/images/integrations/jira.svg'

import '@/styles/integration.scss'

export default function Integration () {
  const history = useHistory()

  const [alertOpen, setAlertOpen] = useState(false)
  const [dbUrl, setDbUrl] = useState()
  const [port, setPort] = useState()
  const [mode, setMode] = useState()

  const [integrations, setIntegrations] = useState([
    {
      id: 'gitlab',
      name: 'Gitlab',
      icon: <GitlabProvider className='providerIconSvg' width='48' height='48' />
    },
    {
      id: 'jenkins',
      name: 'Jenkins',
      icon: <JenkinsProvider className='providerIconSvg' width='48' height='48' />
    },
    {
      id: 'jira',
      name: 'jira',
      icon: <JiraProvider className='providerIconSvg' width='48' height='48' />
    },
  ])

  const [activeProvider, setActiveProvider] = useState(integrations[0])
  const [invalidProvider, setInvalidProvider] = useState(false)

  function updateEnv (key, value) {
    fetch(`${SERVER_HOST}/api/setenv/${key}/${encodeURIComponent(value)}`)
  }

  // function saveAll (e) {
  //   e.preventDefault()
  //   updateEnv('DB_URL', dbUrl)
  //   updateEnv('PORT', port)
  //   updateEnv('MODE', mode)
  //   setAlertOpen(true)
  // }

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
    fetch(`${SERVER_HOST}/api/getenv`)
      .then(response => response.json())
      .then(env => {
        setDbUrl(env.DB_URL)
        setPort(env.PORT)
        setMode(env.MODE)
      })
  }, [])

  return (
    <>
      <div className='container'>
        <Nav />
        <Sidebar />
        <Content>
          <main className='main'>
            <div className='headlineContainer'>
              <h1>Data Integrations</h1>
              <p className='description'>3 sources are available for data collection.</p>
            </div>
            <div className='integrationProviders'>
              {integrations.map((provider) => (
                <div
                  className='iProvider'
                  key={`provider-${provider.id}`}
                  onClick={() => handleProviderClick(provider.id)}
                >
                  <div className='providerIcon'>
                    {provider.icon}
                  </div>
                  <div className='providerName'>
                    {provider.name}
                  </div>
                </div>
              ))}
              {/* <div className='iProvider'>
                <div className='providerIcon'>
                  <JenkinsProvider className='providerIconSvg' width='48' height='48' />
                </div>
                <div className='providerName'>
                  Jenkins
                </div>
              </div>
              <div className='iProvider'>
                <div className='providerIcon'>
                  <JiraProvider className='providerIconSvg' width='48' height='48' />
                </div>
                <div className='providerName'>
                  Jira
                </div>
              </div> */}
            </div>
          </main>
        </Content>
      </div>

      <SaveAlert alertOpen={alertOpen} onClose={() => setAlertOpen(false)} />
    </>
  )
}
