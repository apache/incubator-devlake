import React, { useEffect, useState } from 'react'
import axios from 'axios'
import {
  // BrowserRouter as Router,
  // Switch,
  // Route,
  useParams,
  Link,
  useHistory
} from 'react-router-dom'
import {
  Button, mergeRefs, Card, Elevation, Colors,
  FormGroup, InputGroup, Tooltip, Label,
  Position,
  Alignment,
  Icon,
  Toaster,
  ToasterPosition,
  IToasterProps,
  IToastProps,
} from '@blueprintjs/core'
// import { Column, Table } from '@blueprintjs/table'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import Content from '@/components/Content'
import { ToastNotification } from '@/components/Toast'
import MappingTag from '@/pages/plugins/jira/MappingTag'
import MappingTagStatus from '@/pages/plugins/jira//MappingTagStatus'
import ClearButton from '@/pages/plugins/jira//ClearButton'
import { SERVER_HOST, DEVLAKE_ENDPOINT } from '@/utils/config'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function JenkinsSettings (props) {
  const { connection, provider, isSaving } = props
  const history = useHistory()
  const { providerId, connectionId } = useParams()

  const cancel = () => {
    history.push(`/integrations/${provider.id}`)
  }

  useEffect(() => {

  }, [])

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
