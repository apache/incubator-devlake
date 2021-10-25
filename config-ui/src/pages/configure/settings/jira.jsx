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

import { ReactComponent as GitlabProvider } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProvider } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProvider } from '@/images/integrations/jira.svg'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function JiraSettings (props) {
  const { connection, provider, isSaving } = props
  const history = useHistory()
  const { providerId, connectionId } = useParams()

  const [typeMappingBug, setTypeMappingBug] = useState()
  const [typeMappingIncident, setTypeMappingIncident] = useState()
  const [typeMappingRequirement, setTypeMappingRequirement] = useState()
  const [typeMappingAll, setTypeMappingAll] = useState()

  const [statusMappings, setStatusMappings] = useState()

  function setStatusMapping (key, values, status) {
    setStatusMappings(statusMappings.map(mapping => {
      if (mapping.key === key) {
        mapping.mapping[status] = values
      }
      return mapping
    }))
  }

  const [customStatusOverlay, setCustomStatusOverlay] = useState(false)
  const [customStatusName, setCustomStatusName] = useState('')

  function addStatusMapping (e) {
    const type = customStatusName.trim().toUpperCase()
    if (statusMappings.find(e => e.type === type)) {
      return
    }
    const result = [
      ...statusMappings,
      {
        type,
        key: `JIRA_ISSUE_${type}_STATUS_MAPPING`,
        mapping: {
          Resolved: [],
          Rejected: [],
        }
      }
    ]
    setStatusMappings(result)
    setCustomStatusOverlay(false)
    e.preventDefault()
  }

  useEffect(() => {

  }, [])

  return (
    <>
      <div className='headlineContainer'>
        <h3 className='headline'>Issue Type Mappings</h3>
        <p className='description'>Map your own issue types to Dev Lake's standard types</p>
      </div>

      <MappingTag
        labelName='Bug'
        labelIntent='danger'
        typeOrStatus='type'
        placeholderText='Add Issue Types...'
        values={typeMappingBug}
        helperText='JIRA_ISSUE_TYPE_MAPPING'
        rightElement={<ClearButton onClick={() => setTypeMappingBug([])} />}
        onChange={(values) => setTypeMappingBug(values)}
        disabled={isSaving}
      />

      <MappingTag
        labelName='Incident'
        labelIntent='warning'
        typeOrStatus='type'
        placeholderText='Add Issue Types...'
        values={typeMappingIncident}
        helperText='JIRA_ISSUE_TYPE_MAPPING'
        rightElement={<ClearButton onClick={() => setTypeMappingIncident([])} />}
        onChange={(values) => setTypeMappingIncident(values)}
        disabled={isSaving}
      />

      <MappingTag
        labelName='Requirement'
        labelIntent='primary'
        typeOrStatus='type'
        placeholderText='Add Issue Types...'
        values={typeMappingRequirement}
        helperText='JIRA_ISSUE_TYPE_MAPPING'
        rightElement={<ClearButton onClick={() => setTypeMappingRequirement([])} />}
        onChange={(values) => setTypeMappingRequirement(values)}
        disabled={isSaving}
      />

      <div className='headlineContainer'>
        <h3 className='headline'>Epic Key</h3>
        <p className='description'>Choose the Jira field you’re using to represent the key of an Epic to which an issue belongs to.</p>
      </div>

      <div className='headlineContainer'>
        <h3 className='headline'>Requirement Granularity</h3>
        <p className='description'>Choose the Jira field you’re using to represent the granularity of a requirement-type issue.</p>
      </div>


      <div className='headlineContainer'>
        <h3 className='headline'>Board ID (Optional)</h3>
        <p className='description'>Choose the specific Jira board(s) to collect issues from.</p>
      </div>

    </>
  )
}
