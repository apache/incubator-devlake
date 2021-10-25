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

export default function GitlabSettings (props) {
  const { connection, provider, isSaving } = props
  const history = useHistory()
  const { providerId, connectionId } = useParams()
  const [jiraBoardGitlabeProjects, setJiraBoardGitlabeProjects] = useState()

  useEffect(() => {

  }, [])

  return (
    <>
      <div className='formContainer'>
        <FormGroup
          disabled={isSaving}
          labelFor='jira-board-projects'
          helperText='JIRA_BOARD_GITLAB_PROJECTS'
          className='formGroup'
          contentClassName='formGroupContent'
        >
          <Label>
            Map JIRA Boards to <strong>GitLab</strong>
          </Label>
          <InputGroup
            id='jira-storypoint-field'
            placeholder='<JIRA_BOARD>:<GITLAB_PROJECT_ID>,...; eg. 8:8967944,8967945;9:8967946,8967947'
            defaultValue={jiraBoardGitlabeProjects}
            onChange={(e) => setJiraBoardGitlabeProjects(e.target.value)}
            disabled={isSaving}
            className='input'
          />
        </FormGroup>
      </div>
    </>
  )
}
