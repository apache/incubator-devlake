import React, { useEffect, useState } from 'react'
import {
  useParams,
  useHistory
} from 'react-router-dom'
import {
  FormGroup,
  InputGroup,
  Button,
  Intent,
  Label,
  Tag
} from '@blueprintjs/core'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

export default function GithubSettings (props) {
  const { connection, provider, isSaving, isSavingConnection, onSettingsChange } = props
  const history = useHistory()
  const { providerId, connectionId } = useParams()
  const [prType, setPrType] = useState('')
  const [prComponent, setPrComponent] = useState('')
  const [issueSeverity, setIssueSeverity] = useState('')
  const [issueComponent, setIssueComponent] = useState('')
  const [issuePriority, setIssuePriority] = useState('')
  const [issueTypeBug, setIssueTypeBug] = useState('')
  const [issueTypeRequirement, setIssueTypeRequirement] = useState('')
  const [issueTypeIncident, setIssueTypeIncident] = useState('')

  const [errors, setErrors] = useState([])

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
    setPrType(connection.PrType)
    setPrComponent(connection.PrComponent)
    setIssueSeverity(connection.IssueSeverity)
    setIssuePriority(connection.IssuePriority)
    setIssueComponent(connection.IssueComponent)
    setIssueTypeBug(connection.IssueTypeBug)
    setIssueTypeRequirement(connection.IssueTypeRequirement)
    setIssueTypeIncident(connection.IssueTypeIncident)
  }, [connection])

  useEffect(() => {
    const settings = {
      GITHUB_PR_TYPE: prType,
      GITHUB_PR_COMPONENT: prComponent,
      GITHUB_ISSUE_SEVERITY: issueSeverity,
      GITHUB_ISSUE_COMPONENT: issueComponent,
      GITHUB_ISSUE_PRIORITY: issuePriority,
      GITHUB_ISSUE_TYPE_REQUIREMENT: issueTypeRequirement,
      GITHUB_ISSUE_TYPE_BUG: issueTypeBug,
      GITHUB_ISSUE_TYPE_INCIDENT: issueTypeIncident,
    }
    console.log('>> GITHUB INSTANCE SETTINGS FIELDS CHANGED!', settings)
    onSettingsChange(settings)
  }, [
    prType,
    prComponent,
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
      <h3 className='headline'>Pull Request Enrichment Options <Tag className='bp3-form-helper-text'>RegExp</Tag></h3>
      <p className=''>Enrich GitHub PRs using Label data.</p>

      <div style={{ maxWidth: '60%' }}>
        <div className='formContainer'>
          <FormGroup
            disabled={isSaving || isSavingConnection}
            labelFor='github-pr-type'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              Type
            </Label>
            <InputGroup
              id='github-pr-type'
              placeholder='type/(.*)$'
              defaultValue={prType}
              onChange={(e) => setPrType(e.target.value)}
              onKeyUp={(e) => e.target.value.length === 0 ? setPrType('') : null}
              disabled={isSaving || isSavingConnection}
              className='input'
              maxLength={255}
            />
          </FormGroup>
        </div>
        <div className='formContainer'>
          <FormGroup
            disabled={isSaving || isSavingConnection}
            labelFor='github-pr-component'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              Component
            </Label>
            <InputGroup
              id='github-pr-type'
              placeholder='component/(.*)$'
              defaultValue={prComponent}
              onChange={(e) => setPrComponent(e.target.value)}
              onKeyUp={(e) => e.target.value.length === 0 ? setPrComponent('') : null}
              disabled={isSaving || isSavingConnection}
              className='input'
              maxLength={255}
            />
          </FormGroup>
        </div>
      </div>

      <h3 className='headline'>Issue Type Enrichment Options <Tag className='bp3-form-helper-text'>RegExp</Tag></h3>
      <p className=''>Enrich GitHub Issues using Label data.</p>
      <div style={{ maxWidth: '60%' }}>
        <div className='formContainer'>
          <FormGroup
            disabled={isSaving || isSavingConnection}
            labelFor='github-issue-severity'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              Severity
            </Label>
            <InputGroup
              id='github-issue-severity'
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
            labelFor='github-issue-component'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              Component
            </Label>
            <InputGroup
              id='github-issue-component'
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
            labelFor='github-issue-priority'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              Priority
            </Label>
            <InputGroup
              id='github-issue-priority'
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
            labelFor='github-issue-requirement'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              <span className='bp3-tag tag-requirement'>Requirement</span>
            </Label>
            <InputGroup
              id='github-issue-requirement'
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
            labelFor='github-issue-bug'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              <span className='bp3-tag tag-bug'>Bug</span>
            </Label>
            <InputGroup
              id='github-issue-bug'
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
            labelFor='github-issue-bug'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              <span className='bp3-tag tag-incident'>Incident</span>
            </Label>
            <InputGroup
              id='github-issue-incident'
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
      {/* <h3 className='headline'>Github Proxy</h3>
      <p className=''>Optional</p>
      <div className='formContainer'>
        <FormGroup
          disabled={isSaving || isSavingConnection}
          labelFor='github-proxy'
          helperText='PROXY'
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
            disabled={isSaving || isSavingConnection}
            className='input'
          />
        </FormGroup>
      </div> */}
    </>
  )
}
