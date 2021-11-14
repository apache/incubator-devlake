import React, { useEffect, useState } from 'react'
import {
  useParams,
  useHistory
} from 'react-router-dom'
import {
  FormGroup,
  InputGroup,
  Label,
} from '@blueprintjs/core'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

export default function GithubSettings (props) {
  const { connection, provider, isSaving, onSettingsChange } = props
  const history = useHistory()
  const { providerId, connectionId } = useParams()
  const [githubProxy, setGithubProxy] = useState(null)

  const [errors, setErrors] = useState([])

  const cancel = () => {
    history.push(`/integrations/${provider.id}`)
  }

  useEffect(() => {
    setErrors(['This integration doesn’t require any configuration.'])
  }, [])

  useEffect(() => {
    onSettingsChange({
      errors,
      providerId,
      connectionId
    })
  }, [errors, onSettingsChange, connectionId, providerId])

  useEffect(() => {
    setGithubProxy(connection.proxy)
  }, [connection])

  useEffect(() => {
    const settings = {
      GITHUB_PROXY: githubProxy
    }
    console.log('>> GITHUB INSTANCE SETTINGS FIELDS CHANGED!', settings)
    onSettingsChange(settings)
  }, [
    githubProxy,
    onSettingsChange
  ])

  return (
    <>
      <h3 className='headline'>Github Proxy</h3>
      <p className=''>Optional</p>
      <div className='formContainer'>
        <FormGroup
          disabled={isSaving}
          labelFor='github-proxy'
          helperText='GITHUB_PROXY'
          className='formGroup'
          contentClassName='formGroupContent'
        >
          <Label>
            Proxy URL
          </Label>
          <InputGroup
            id='github-proxy'
            placeholder='http://your-proxy-server.com:1080'
            defaultValue={githubProxy}
            onChange={(e) => setGithubProxy(e.target.value)}
            disabled={isSaving}
            className='input'
          />
        </FormGroup>
      </div>
    </>
  )
}
