import React, { useEffect, useState } from 'react'
import {
  useParams,
  useHistory
} from 'react-router-dom'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

export default function JenkinsSettings (props) {
  const { connection, provider, isSaving, isSavingConnection, onSettingsChange } = props
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
        <h3 className='headline'>No Additional Settings</h3>
        <p className='description'>
          This integration doesn’t require any configuration.
          You can continue to&nbsp;
          <a href='#' style={{ textDecoration: 'underline' }} onClick={cancel}>add other data sources</a>&nbsp;
          or trigger collection at the <a href='#' style={{ textDecoration: 'underline' }} onClick={cancel}>previous page</a>.
        </p>
      </div>
    </>
  )
}
