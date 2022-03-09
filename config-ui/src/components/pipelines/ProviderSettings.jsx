import React, { useEffect, useState } from 'react'
import {
  Providers,
} from '@/data/Providers'
import {
  Button,
  ButtonGroup,
  FormGroup,
  Popover,
  InputGroup,
  MenuItem,
  Menu,
  Intent,
  TagInput,
  Checkbox,
  Tag,
  Switch,
  Tooltip,
  Colors,
} from '@blueprintjs/core'
import { Select } from '@blueprintjs/select'
import RefDiffTasksMenu from '@/components/menus/RefDiffTasksMenu'

const ProviderSettings = (props) => {
  const {
    providerId,
    projectId = [],
    sourceId,
    selectedSource,
    sources = [],
    boardId = [],
    owner,
    repositoryName,
    gitExtractorUrl,
    gitExtractorRepoId,
    refDiffRepoId,
    refDiffPairs = [],
    refDiffTasks = [],
    setProjectId = () => {},
    setSourceId = () => {},
    setSelectedSource = () => {},
    setBoardId = () => {},
    setOwner = () => {},
    setRepositoryName = () => {},
    setGitExtractorUrl = () => {},
    setGitExtractorRepoId = () => {},
    setRefDiffRepoId = () => {},
    setRefDiffPairs = () => {},
    setRefDiffTasks = () => {},
    isEnabled = () => {},
    isRunning = false,
    onReset = () => {}
  } = props

  let providerSettings = null

  const [refDiffOldTag, setRefDiffOldTag] = useState('')
  const [refDiffNewTag, setRefDiffNewTag] = useState('')

  const handleRefDiffTaskSelect = (e, task) => {
    setRefDiffTasks(t => t.includes(task.task)
      ? t.filter(t2 => t2 !== task.task)
      : [...new Set([...t, task.task])])
  }

  const createRefDiffPair = (oldRef, newRef) => {
    return {
      oldRef,
      newRef
    }
  }

  const addRefDiffPairObject = (oldRef, newRef) => {
    setRefDiffPairs(pairs => (!pairs.some(p => p.oldRef === oldRef && p.newRef === newRef))
      ? [...pairs, { oldRef, newRef }]
      : [...pairs])
    setRefDiffNewTag('')
    setRefDiffOldTag('')
  }

  const removeRefDiffPairObject = (oldRef, newRef) => {
    setRefDiffPairs(pairs => pairs.filter(p => !(p.oldRef === oldRef && p.newRef === newRef)))
  }

  // useEffect(() => {
  //   console.log('>>>> REF DIFF PAIRS ARRAY...', refDiffPairs)
  // }, [refDiffPairs])

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
            labelInfo={<span style={{ display: 'block' }}>Choose Connection Instance ID</span>}
            inline={false}
            labelFor='source-id'
            className=''
            contentClassName=''
            fill
          >
            {/* <InputGroup
              id='source-id'
              disabled={isRunning || !isEnabled(providerId)}
              placeholder='eg. 54'
              value={sourceId}
              onChange={(e) => setSourceId(e.target.value)}
              className='input-source-id'
              autoComplete='off'
              fill={false}
            /> */}
            <ButtonGroup>
              <Select
                disabled={isRunning || !isEnabled(providerId)}
                className='selector-source-id'
                multiple
                inline={true}
                fill={true}
                items={sources}
                activeItem={selectedSource}
                itemPredicate={(query, item) => item?.title?.toLowerCase().indexOf(query.toLowerCase()) >= 0}
                itemRenderer={(item, { handleClick, modifiers }) => (
                  <MenuItem
                    active={modifiers.active}
                    key={item.value}
                    label={item.value}
                    onClick={handleClick}
                    text={item.title}
                  />
                )}
                noResults={<MenuItem disabled={true} text='No Connections.' />}
                onItemSelect={(item) => {
                  setSelectedSource(item)
                }}
              >
                <Button
                  disabled={isRunning || !isEnabled(providerId)}
                  style={{ justifyContent: 'space-between', minWidth: '206px', maxWidth: '290px', whiteSpace: 'nowrap' }}
                  text={selectedSource ? `${selectedSource.title} [${selectedSource.value}]` : 'Select Instance'}
                  rightIcon='double-caret-vertical'
                  fill
                />
              </Select>
              <Button
                icon='eraser'
                intent={Intent.WARNING}
                disabled={isRunning || !isEnabled(providerId)}
                onClick={() => setSelectedSource(null)}
              />
            </ButtonGroup>
          </FormGroup>
          <FormGroup
            disabled={isRunning || !isEnabled(providerId)}
            label={
              <strong>Board IDs<span className='requiredStar'>*</span>
                <span
                  className='badge-count'
                  style={{
                    opacity: isEnabled(providerId) ? 0.5 : 0.1
                  }}
                >{boardId.length}
                </span>
              </strong>
            }
            labelInfo={<span style={{ display: 'block' }}>Enter one or more JIRA Board IDs.</span>}
            inline={false}
            labelFor='board-id'
            className=''
            contentClassName=''
            style={{ marginLeft: '12px' }}
            fill
          >
            {/* (DISABLED) Single Input */}
            {/* <InputGroup
              id='board-id'
              disabled={isRunning || !isEnabled(providerId)}
              placeholder='eg. 8'
              value={boardId}
              onChange={(e) => setBoardId(e.target.value)}
              className='input-board-id'
              autoComplete='off'
              fill={false}
            /> */}
            <div style={{ width: '100%' }}>
              <TagInput
                id='board-id'
                disabled={isRunning || !isEnabled(providerId)}
                placeholder='eg. 8, 100, 200'
                values={boardId || []}
                fill={true}
                onChange={(values) => setBoardId([...new Set(values)])}
                addOnPaste={true}
                addOnBlur={true}
                rightElement={
                  <Button
                    disabled={isRunning || !isEnabled(providerId)}
                    icon='eraser'
                    minimal
                    onClick={() => setBoardId([])}
                  />
                  }
                onKeyDown={e => e.key === 'Enter' && e.preventDefault()}
                className='input-board-id tagInput'
              />
            </div>
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
              placeholder='eg. merico-dev'
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
            label={
              <strong>Project IDs<span className='requiredStar'>*</span>
                <span
                  className='badge-count'
                  style={{
                    opacity: isEnabled(providerId) ? 0.5 : 0.1
                  }}
                >{projectId.length}
                </span>
              </strong>
            }
            labelInfo={<span style={{ display: 'block' }}>Enter one or more GitLab Project IDs.</span>}
            inline={false}
            labelFor='project-id'
            className=''
            contentClassName=''
          >
            {/* (DISABLED) Single Input */}
            {/* <InputGroup
              id='project-id'
              disabled={isRunning || !isEnabled(providerId)}
              placeholder='eg. 937810831'
              value={projectId}
              onChange={(e) => setProjectId(pId => e.target.value)}
              className='input-project-id'
              autoComplete='off'
            /> */}
            <div style={{ width: '100%' }}>
              <TagInput
                id='project-id'
                disabled={isRunning || !isEnabled(providerId)}
                placeholder='eg. 937810831, 95781015'
                values={projectId || []}
                fill={true}
                onChange={(values) => setProjectId([...new Set(values)])}
                addOnPaste={true}
                addOnBlur={true}
                rightElement={
                  <Button
                    disabled={isRunning || !isEnabled(providerId)}
                    icon='eraser'
                    minimal
                    onClick={() => setProjectId([])}
                  />
                  }
                onKeyDown={e => e.key === 'Enter' && e.preventDefault()}
                className='input-project-id tagInput'
              />
            </div>
          </FormGroup>
        </>
      )
      break
    case Providers.GITEXTRACTOR:
      providerSettings = (
        <>
          <FormGroup
            disabled={isRunning || !isEnabled(providerId)}
            label={<strong>Git URL<span className='requiredStar'>*</span></strong>}
            labelInfo={<span style={{ display: 'block' }}>Enter Repository URL</span>}
            inline={false}
            labelFor='git-url'
            className=''
            contentClassName=''
            fill
            style={{ minWidth: '372px' }}
          >
            <InputGroup
              id='gitextractor-url'
              disabled={isRunning || !isEnabled(providerId)}
              placeholder='eg. https://github.com/merico-dev/lake.git'
              value={gitExtractorUrl}
              onChange={(e) => setGitExtractorUrl(e.target.value)}
              className='input-gitextractor-url'
              autoComplete='off'
            />
          </FormGroup>
          <FormGroup
            disabled={isRunning || !isEnabled(providerId)}
            label={<strong>Repository ID<span className='requiredStar'>*</span></strong>}
            labelInfo={<span style={{ display: 'block' }}>Enter Repo Column ID</span>}
            inline={false}
            labelFor='gitextractor-repo-id'
            className=''
            contentClassName=''
            style={{ marginLeft: '12px', minWidth: '280px' }}
            fill
          >
            <InputGroup
              id='gitextractor-repo-id'
              disabled={isRunning || !isEnabled(providerId)}
              placeholder='eg. github:GithubRepo:384111310'
              value={gitExtractorRepoId}
              onChange={(e) => setGitExtractorRepoId(e.target.value)}
              className='input-gitextractor-repo-id'
              autoComplete='off'
              fill={false}
            />
          </FormGroup>
        </>
      )
      break
    case Providers.REFDIFF:
      providerSettings = (
        <>
          <FormGroup
            disabled={isRunning || !isEnabled(providerId)}
            label={<strong>Repository ID<span className='requiredStar'>*</span></strong>}
            labelInfo={<span style={{ display: 'block' }}>Enter Repo Column ID</span>}
            inline={false}
            labelFor='refdiff-repo-id'
            className=''
            contentClassName=''
            style={{ minWidth: '280px', marginBottom: 'auto' }}
            fill
          >
            <InputGroup
              id='refdiff-repo-id'
              disabled={isRunning || !isEnabled(providerId)}
              placeholder='eg. github:GithubRepo:384111310'
              value={refDiffRepoId}
              onChange={(e) => setRefDiffRepoId(e.target.value)}
              className='input-refdiff-repo-id'
              autoComplete='off'
              fill={false}
              style={{ marginBottom: '10px' }}
            />
          </FormGroup>
          <FormGroup
            disabled={isRunning || !isEnabled(providerId)}
            label={(
              <strong>Tasks<span className='requiredStar'>*</span>
                <span
                  className='badge-count'
                  style={{
                    opacity: isEnabled(providerId) ? 0.5 : 0.1
                  }}
                >{refDiffTasks.length}
                </span>
              </strong>
              )}
            labelInfo={<span style={{ display: 'block' }}>Select Tasks to Execute</span>}
            inline={false}
            labelFor='refdiff-tasks'
            className=''
            contentClassName=''
            style={{ marginLeft: '12px', marginRight: '12px', marginBottom: 'auto' }}
            fill
          >

            <Popover
              className='provider-tasks-popover'
              disabled={isRunning || !isEnabled(providerId)}
              content={(
                <RefDiffTasksMenu onSelect={handleRefDiffTaskSelect} selected={refDiffTasks} />
              )}
              placement='top-center'
            >
              <Button icon='menu' text='Choose Tasks' disabled={isRunning || !isEnabled(providerId)} />
            </Popover>
          </FormGroup>
          <div style={{ display: 'flex', flex: 1, width: '100%' }}>
            <FormGroup
              disabled={isRunning || !isEnabled(providerId)}
              label={(
                <strong>Tags<span className='requiredStar'>*</span>
                  <span
                    className='badge-count'
                    style={{
                      opacity: isEnabled(providerId) ? 0.5 : 0.1
                    }}
                  >{refDiffPairs.length}
                  </span>
                </strong>
              )}
              labelInfo={<span style={{ display: 'block' }}>Specify tag Ref Pairs</span>}
              inline={false}
              labelFor='refdiff-pair-newref'
              className=''
              contentClassName=''
              style={{ minWidth: '600px', marginBottom: 'auto' }}
              fill={false}
            >
              <div style={{ display: 'flex' }}>
                <div>
                  <InputGroup
                    id='refdiff-pair-newref'
                    round='true'
                    leftElement={(
                      <Tag
                        intent={Intent.WARNING} style={{
                          opacity: isEnabled(providerId) ? 1 : 0.3
                        }}
                      >New Ref
                      </Tag>
                    )}
                    inline={true}
                    disabled={isRunning || !isEnabled(providerId)}
                    placeholder='eg. refs/tags/v0.6.0'
                    value={refDiffNewTag}
                    onChange={(e) => setRefDiffNewTag(e.target.value)}
                    autoComplete='off'
                    fill={false}
                    // small
                  />
                </div>
                <div style={{ marginLeft: '10px', marginRight: '10px' }}>
                  {/* <label>Old Ref</label> */}
                  <InputGroup
                    id='refdiff-pair-oldref'
                    round='true'
                    leftElement={(
                      <Tag
                        style={{
                          opacity: isEnabled(providerId) ? 1 : 0.3
                        }}
                      >Old Ref
                      </Tag>
                    )}
                    inline={true}
                    disabled={isRunning || !isEnabled(providerId)}
                    placeholder='eg. refs/tags/v0.5.0'
                    value={refDiffOldTag}
                    onChange={(e) => setRefDiffOldTag(e.target.value)}
                    autoComplete='off'
                    fill={false}
                    rightElement={(
                      <Tooltip content='Add Tag Pair'>
                        <Button
                          intent={Intent.PRIMARY}
                          disabled={!refDiffOldTag || !refDiffNewTag || refDiffOldTag === refDiffNewTag}
                          icon='plus'
                          small
                          style={{}}
                          onClick={() => addRefDiffPairObject(refDiffOldTag, refDiffNewTag)}
                        />
                      </Tooltip>

                    )}
                    // small
                  />
                </div>
                <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
                  {/* <Button icon='remove' minimal small style={{ marginTop: 'auto', alignSelf: 'center' }} /> */}
                  {/* <Button
                    intent={Intent.PRIMARY}
                    disabled={!refDiffOldTag || !refDiffNewTag || refDiffOldTag === refDiffNewTag}
                    icon='add'
                    small
                    style={{ marginTop: 'auto', alignSelf: 'center' }}
                    onClick={() => addRefDiffPairObject(refDiffOldTag, refDiffNewTag)}
                  /> */}
                </div>
              </div>
              <div
                style={{
                  borderRadius: '4px',
                  padding: '10px',
                  margin: '5px 0'
                }}
              >
                {refDiffPairs.map((pair, pairIdx) => (
                  <div key={`refdiff-added-pairs-itemkey-$${pairIdx}`} style={{ display: 'flex' }}>
                    <div style={{ flex: 1 }}>
                      <Tag intent={Intent.WARNING} round='true' small>new</Tag> {pair.newRef}
                    </div>
                    <div style={{ flex: 1, marginLeft: '10px', marginRight: '10px' }}>
                      <Tag round='true' small>old</Tag> {pair.oldRef}
                    </div>
                    <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
                      <Button
                        icon='remove'
                        minimal
                        small
                        style={{ marginTop: 'auto', alignSelf: 'center' }}
                        onClick={() => removeRefDiffPairObject(pair.oldRef, pair.newRef)}
                      />
                    </div>
                  </div>
                ))}
                {refDiffPairs.length === 0 && <><span style={{ color: Colors.GRAY3 }}>( No Tag Pairs Added)</span></>}
              </div>
            </FormGroup>
          </div>
        </>

      )
      break
    default:
      break
  }

  return providerSettings
}

export default ProviderSettings
