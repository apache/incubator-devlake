import React from 'react'
import {
  Providers,
} from '@/data/Providers'
import {
  FormGroup,
  InputGroup
} from '@blueprintjs/core'

const ProviderSettings = (props) => {
  const {
    providerId,
    projectId,
    sourceId,
    boardId,
    owner,
    repositoryName,
    setProjectId = () => {},
    setSourceId = () => {},
    setBoardId = () => {},
    setOwner = () => {},
    setRepositoryName = () => {},
    isEnabled = () => {},
    isRunning = false,
  } = props

  let providerSettings = null

  switch (providerId) {
    case Providers.JENKINS:
      providerSettings = <p><strong style={{ fontWeight: 900 }}>AUTO-CONFIGURED</strong><br />No Additional Settings</p>
      break
    case Providers.JIRA:
      providerSettings = (
        <>
          <FormGroup
            disabled={isRunning || !isEnabled(providerId)}
            label={<strong>Source ID<span className='requiredStar'>*</span></strong>}
            labelInfo={<span style={{ display: 'block' }}>Enter Connection Instance ID</span>}
            inline={false}
            labelFor='source-id'
            className=''
            contentClassName=''
            fill
          >
            <InputGroup
              id='source-id'
              disabled={isRunning || !isEnabled(providerId)}
              placeholder='eg. 54'
              value={sourceId}
              onChange={(e) => setSourceId(e.target.value)}
              className='input-source-id'
              autoComplete='off'
              fill={false}
            />
          </FormGroup>
          <FormGroup
            disabled={isRunning || !isEnabled(providerId)}
            label={<strong>Board ID<span className='requiredStar'>*</span></strong>}
            labelInfo={<span style={{ display: 'block' }}>Enter JIRA Board No.</span>}
            inline={false}
            labelFor='board-id'
            className=''
            contentClassName=''
            style={{ marginLeft: '12px' }}
            fill
          >
            <InputGroup
              id='board-id'
              disabled={isRunning || !isEnabled(providerId)}
              placeholder='eg. 8'
              value={boardId}
              onChange={(e) => setBoardId(e.target.value)}
              className='input-board-id'
              autoComplete='off'
              fill={false}
            />
          </FormGroup>
        </>
      )
      break
    case Providers.GITHUB:
      providerSettings = (
        <>
          <FormGroup
            disabled={isRunning || !isEnabled(providerId)}
            label={<strong>Owner<span className='requiredStar'>*</span></strong>}
            labelInfo={<span style={{ display: 'block' }}>Enter Project Owner</span>}
            inline={false}
            labelFor='owner'
            className=''
            contentClassName=''
            fill
          >
            <InputGroup
              id='owner'
              disabled={isRunning || !isEnabled(providerId)}
              placeholder='eg. merio-dev'
              value={owner}
              onChange={(e) => setOwner(e.target.value)}
              className='input-owner'
              autoComplete='off'
            />
          </FormGroup>
          <FormGroup
            disabled={isRunning || !isEnabled(providerId)}
            label={<strong>Repository Name<span className='requiredStar'>*</span></strong>}
            labelInfo={<span style={{ display: 'block' }}>Enter Git repository</span>}
            inline={false}
            labelFor='repository-name'
            className=''
            contentClassName=''
            style={{ marginLeft: '12px' }}
            fill
          >
            <InputGroup
              id='repository-name'
              disabled={isRunning || !isEnabled(providerId)}
              placeholder='eg. lake'
              value={repositoryName}
              onChange={(e) => setRepositoryName(e.target.value)}
              className='input-repository-name'
              autoComplete='off'
              fill={false}
            />
          </FormGroup>
        </>
      )
      break
    case Providers.GITLAB:
      providerSettings = (
        <>
          <FormGroup
            disabled={isRunning || !isEnabled(providerId)}
            label={<strong>Project ID<span className='requiredStar'>*</span></strong>}
            labelInfo={<span style={{ display: 'block' }}>Enter the GitLab Project ID No.</span>}
            inline={false}
            labelFor='project-id'
            className=''
            contentClassName=''
          >
            <InputGroup
              id='project-id'
              disabled={isRunning || !isEnabled(providerId)}
              placeholder='eg. 937810831'
              value={projectId}
              onChange={(e) => setProjectId(pId => e.target.value)}
              className='input-project-id'
              autoComplete='off'
            />
          </FormGroup>
        </>
      )
      break
    default:
      break
  }

  return providerSettings
}

export default ProviderSettings
