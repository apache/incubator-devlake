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
import { useParams } from 'react-router-dom'
import { FormGroup, InputGroup, Label, Tag } from '@blueprintjs/core'

// import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function GitlabSettings (props) {
  const { connection, isSaving, isSavingConnection, onSettingsChange } = props

  // const { providerId, connectionId } = useParams()
  const [jiraBoardGitlabProjects, setJiraBoardGitlabProjects] = useState('')
  const { providerId, connectionId } = useParams()
  const [mrType, setMrType] = useState('')
  const [mrComponent, setMrComponent] = useState('')
  const [issueSeverity, setIssueSeverity] = useState('')
  const [issueComponent, setIssueComponent] = useState('')
  const [issuePriority, setIssuePriority] = useState('')
  const [issueTypeBug, setIssueTypeBug] = useState('')
  const [issueTypeRequirement, setIssueTypeRequirement] = useState('')
  const [issueTypeIncident, setIssueTypeIncident] = useState('')
  const [errors, setErrors] = useState([])
  // const [showBoardMappingDialog, setShowBoardMappingDialog] = useState(false)
  //
  // const [selectedBoard, setSelectedBoard] = useState({ id: 1, title: 'Open', value: 1 })
  // const [boards, setBoards] = useState(boardsData)
  //
  // const [selectedProject, setSelectedProject] = useState({ id: 0, title: 'GL PRJ 3E4', value: 938191 })
  // const [projects, setProjects] = useState(projectsData)
  //
  // const [boardMappings, setBoardMappings] = useState([
  //   {
  //     id: 0,
  //     boards: [8],
  //     projects: [8967944, 8967945],
  //     deleted: false
  //   },
  //   {
  //     id: 1,
  //     boards: [9],
  //     projects: [8967946, 8967947],
  //     deleted: false
  //   },
  //   {
  //     id: 2,
  //     boards: [10],
  //     projects: [8967946, 8967947, 1967900],
  //     deleted: false
  //   }
  // ])
  // const [deletedBoards, setDeletedBoards] = useState([])
  //
  // const boardMapExists = (boardMap) => {
  //   return boardMappings.some(b => b.boards.includes(boardMap.boards[0]) && b.projects.includes(boardMap.projects[0]))
  // }

  useEffect(() => {
    const settings = {
      JIRA_BOARD_GITLAB_PROJECTS: jiraBoardGitlabProjects
    }
    onSettingsChange(settings)
    console.log('>> GITLAB INSTANCE SETTINGS FIELDS CHANGED!', settings)
  }, [
    jiraBoardGitlabProjects,
    onSettingsChange
  ])
  useEffect(() => {
    onSettingsChange({
      errors,
      providerId,
      connectionId
    })
  }, [errors, onSettingsChange, connectionId, providerId])
  // const createBoardMapping = () => {
  //   setShowBoardMappingDialog(true)
  // }
  //
  // const addBoardMapping = (boardMap) => {
  //   setBoardMappings([
  //     ...boardMappings,
  //     boardMap
  //   ])
  // }
  //
  // const editBoardMapping = () => {
  //
  // }
  //
  // const deleteBoardMapping = (boardMap) => {
  //   const newBoards = boardMappings.filter(board => board.id !== boardMap.id)
  //   setBoardMappings(newBoards)
  //   console.log('>>> EDIT BOARD MAPPINGS', newBoards)
  // }
  //
  // const linkProject = () => {
  //   setShowBoardMappingDialog(false)
  //   // @todo MULTI-PROJECT MAPPING!
  //   const newBoardMap = {
  //     id: boardMappings.length,
  //     boards: [selectedBoard.value],
  //     projects: [selectedProject.value],
  //     deleted: false
  //   }
  //   if (!boardMapExists(newBoardMap)) {
  //     addBoardMapping(newBoardMap)
  //   } else {
  //     ToastNotification.show({ message: 'Gitlab Project Link already ', intent: 'danger', icon: 'error' })
  //   }
  //   console.log('>> LINKING JIRA BOARD', selectedBoard, ' TO GITLAB PROJECT ', selectedProject)
  // }
  useEffect(() => {
    setErrors(['This integration doesn’t require any configuration.'])
  }, [])
  useEffect(() => {
    // @todo FETCH BOARDS FROM BE API
    // setBoards([])

    // @todo FETCH PROJECTS FROM BE API
    // setProjects([])
    if (connection && connection.ID) {
      console.log('>> GITLAB CONNECTION OBJECT RECEIVED...', connection, connection.JIRA_BOARD_GITLAB_PROJECTS)
      setJiraBoardGitlabProjects(connection.JIRA_BOARD_GITLAB_PROJECTS)
    } else {
      console.log('>>>> WARNING!! NO CONNECTION OBJECT', connection)
    }
  }, [connection])
  useEffect(() => {
    onSettingsChange({
      errors,
      providerId,
      connectionId
    })
  }, [errors, onSettingsChange, connectionId, providerId])

  useEffect(() => {
    setMrType(connection.mrType)
    setMrComponent(connection.mrComponent)
    setIssueSeverity(connection.issueSeverity)
    setIssuePriority(connection.issuePriority)
    setIssueComponent(connection.issueComponent)
    setIssueTypeBug(connection.issueTypeBug)
    setIssueTypeRequirement(connection.issueTypeRequirement)
    setIssueTypeIncident(connection.issueTypeIncident)
  }, [connection])

  useEffect(() => {
    const settings = {
      mrType: mrType,
      mrComponent: mrComponent,
      issueSeverity: issueSeverity,
      issueComponent: issueComponent,
      issuePriority: issuePriority,
      issueTypeRequirement: issueTypeRequirement,
      issueTypeBug: issueTypeBug,
      issueTypeIncident: issueTypeIncident,
    }
    console.log('>> GITLAB INSTANCE SETTINGS FIELDS CHANGED!', settings)
    onSettingsChange(settings)
  }, [
    mrType,
    mrComponent,
    issueSeverity,
    issueComponent,
    issuePriority,
    issueTypeRequirement,
    issueTypeBug,
    issueTypeIncident,
    onSettingsChange
  ])
  return (
    <>
      <h3 className='headline'>Issue Enrichment Options <Tag className='bp3-form-helper-text'>RegExp</Tag></h3>
      <p className=''>Enrich Gitlab Issues using Label data.</p>
      <div style={{ maxWidth: '60%' }}>
        <div className='formContainer'>
          <FormGroup
            disabled={isSaving || isSavingConnection}
            labelFor='gitlab-issue-severity'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              Severity
            </Label>
            <InputGroup
              id='gitlab-issue-severity'
              placeholder='severity/(.*)$'
              defaultValue={issueSeverity}
              onChange={(e) => setIssueSeverity(e.target.value)}
              onKeyUp={(e) => e.target.value.length === 0 ? setIssueSeverity('') : null}
              disabled={isSaving || isSavingConnection}
              className='input'
              maxLength={255}
            />
          </FormGroup>
        </div>
        <div className='formContainer'>
          <FormGroup
            disabled={isSaving || isSavingConnection}
            labelFor='gitlab-issue-component'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              Component
            </Label>
            <InputGroup
              id='gitlab-issue-component'
              placeholder='component/(.*)$'
              defaultValue={issueComponent}
              onChange={(e) => setIssueComponent(e.target.value)}
              onKeyUp={(e) => e.target.value.length === 0 ? setIssueComponent('') : null}
              disabled={isSaving || isSavingConnection}
              className='input'
              maxLength={255}
            />
          </FormGroup>
        </div>
        <div className='formContainer'>
          <FormGroup
            disabled={isSaving || isSavingConnection}
            labelFor='gitlab-issue-priority'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              Priority
            </Label>
            <InputGroup
              id='gitlab-issue-priority'
              placeholder='(highest|high|medium|low)$'
              defaultValue={issuePriority}
              onChange={(e) => setIssuePriority(e.target.value)}
              onKeyUp={(e) => e.target.value.length === 0 ? setIssuePriority('') : null}
              disabled={isSaving || isSavingConnection}
              className='input'
              maxLength={255}
            />
          </FormGroup>
        </div>
        <div className='formContainer'>
          <FormGroup
            disabled={isSaving || isSavingConnection}
            labelFor='gitlab-issue-requirement'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              <span className='bp3-tag tag-requirement'>Type - Requirement</span>
            </Label>
            <InputGroup
              id='gitlab-issue-requirement'
              placeholder='(feat|feature|proposal|requirement)$'
              defaultValue={issueTypeRequirement}
              onChange={(e) => setIssueTypeRequirement(e.target.value)}
              onKeyUp={(e) => e.target.value.length === 0 ? setIssueTypeRequirement('') : null}
              disabled={isSaving || isSavingConnection}
              className='input'
              maxLength={255}
            />
          </FormGroup>
        </div>
        <div className='formContainer'>
          <FormGroup
            disabled={isSaving || isSavingConnection}
            labelFor='gitlab-issue-bug'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              <span className='bp3-tag tag-bug'>Type - Bug</span>
            </Label>
            <InputGroup
              id='gitlab-issue-bug'
              placeholder='(bug|broken)$'
              defaultValue={issueTypeBug}
              onChange={(e) => setIssueTypeBug(e.target.value)}
              onKeyUp={(e) => e.target.value.length === 0 ? setIssueTypeBug('') : null}
              disabled={isSaving || isSavingConnection}
              className='input'
              maxLength={255}
            />
          </FormGroup>
        </div>
        <div className='formContainer'>
          <FormGroup
            disabled={isSaving || isSavingConnection}
            labelFor='gitlab-issue-bug'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              <span className='bp3-tag tag-incident'>Type - Incident</span>
            </Label>
            <InputGroup
              id='gitlab-issue-incident'
              placeholder='(incident|p0|p1|p2)$'
              defaultValue={issueTypeIncident}
              onChange={(e) => setIssueTypeIncident(e.target.value)}
              onKeyUp={(e) => e.target.value.length === 0 ? setIssueTypeIncident('') : null}
              disabled={isSaving || isSavingConnection}
              className='input'
              maxLength={255}
            />
          </FormGroup>
        </div>
      </div>

      <h3 className='headline'>Merge Request Enrichment Options <Tag className='bp3-form-helper-text'>RegExp</Tag></h3>
      <p className=''>Enrich Gitlab MRs using Label data.</p>

      <div style={{ maxWidth: '60%' }}>
        <div className='formContainer'>
          <FormGroup
            disabled={isSaving || isSavingConnection}
            labelFor='gitlab-mr-type'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              Type
            </Label>
            <InputGroup
              id='gitlab-mr-type'
              placeholder='type/(.*)$'
              defaultValue={mrType}
              onChange={(e) => setMrType(e.target.value)}
              onKeyUp={(e) => e.target.value.length === 0 ? setMrType('') : null}
              disabled={isSaving || isSavingConnection}
              className='input'
              maxLength={255}
            />
          </FormGroup>
        </div>
        <div className='formContainer'>
          <FormGroup
            disabled={isSaving || isSavingConnection}
            labelFor='gitlab-mr-component'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              Component
            </Label>
            <InputGroup
              id='gitlab-mr-type'
              placeholder='component/(.*)$'
              defaultValue={mrComponent}
              onChange={(e) => setMrComponent(e.target.value)}
              onKeyUp={(e) => e.target.value.length === 0 ? setMrComponent('') : null}
              disabled={isSaving || isSavingConnection}
              className='input'
              maxLength={255}
            />
          </FormGroup>
        </div>
      </div>
      {/* <h3 className='headline'>Github Proxy</h3>
      <p className=''>Optional</p>
      <div className='formContainer'>
        <FormGroup
          disabled={isSaving || isSavingConnection}
          labelFor='gitlab-mroxy'
          helperText='PROXY'
          className='formGroup'
          contentClassName='formGroupContent'
        >
          <Label>
            Proxy URL
          </Label>
          <InputGroup
            id='gitlab-mroxy'
            placeholder='http://your-proxy-server.com:1080'
            defaultValue={gitlabProxy}
            onChange={(e) => setGithubProxy(e.target.value)}
            disabled={isSaving || isSavingConnection}
            className='input'
          />
        </FormGroup>
      </div> */}
    </>
  )
}
