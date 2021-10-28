import React, { useEffect, useState } from 'react'
import {
  useParams,
  useHistory
} from 'react-router-dom'
import {
  FormGroup, InputGroup, Label
} from '@blueprintjs/core'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function GitlabSettings (props) {
  const { connection, provider, isSaving, onSettingsChange } = props
  const history = useHistory()
  const { providerId, connectionId } = useParams()
  const [jiraBoardGitlabeProjects, setJiraBoardGitlabeProjects] = useState()

  useEffect(() => {
    const settings = {
      JIRA_BOARD_GITLAB_PROJECTS: jiraBoardGitlabeProjects
    }
    onSettingsChange(settings)
    console.log('>> GITLAB INSTANCE SETTINGS FIELDS CHANGED!', settings)
  }, [
    jiraBoardGitlabeProjects,
    onSettingsChange
  ])

  useEffect(() => {
    if (connection && connection.id) {
      setJiraBoardGitlabeProjects(connection.JIRA_BOARD_GITLAB_PROJECTS)
    }
  }, [connection])

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
