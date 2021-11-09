import React, { useEffect, useState } from 'react'
import {
  useParams,
  useHistory
} from 'react-router-dom'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function GithubSettings (props) {
  const { connection, provider, isSaving, onSettingsChange } = props
  const history = useHistory()
  const { providerId, connectionId } = useParams()

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

  return (
    <>
      <div className='headlineContainer'>
        <h3 className='headline'>GitHub Settings</h3>
        <p className='description'>
          This integration doesn’t require any configuration.
        </p>
      </div>
    </>
  )
}
