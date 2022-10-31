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
import React, { useEffect, useMemo, useContext } from 'react'
import {
  Button,
  Card,
  Divider,
  Elevation,
  Intent,
  TagInput
} from '@blueprintjs/core'
import IntegrationsContext from '@/store/integrations-context'
// import { ProviderIcons, Providers } from '@/data/Providers'
import ConnectionTabs from '@/components/blueprints/ConnectionTabs'
import BoardsSelector from '@/components/blueprints/BoardsSelector'
import DataEntitiesSelector from '@/components/blueprints/DataEntitiesSelector'
import NoData from '@/components/NoData'
import GitlabProjectsSelector from '@/components/blueprints/GitlabProjectsSelector'
import GitHubProject from '@/models/GithubProject'
import JenkinsJobsSelector from '@/components/blueprints/JenkinsJobsSelector'

const DataScopes = (props) => {
  const {
    activeStep,
    activeConnectionTab,
    blueprintConnections = [],
    dataEntitiesList = [],
    boardsList = [],
    fetchGitlabProjects = () => [],
    isFetchingGitlab = false,
    gitlabProjects = [],
    fetchJenkinsJobs = () => [],
    isFetchingJenkins = false,
    jenkinsJobs = [],
    dataEntities = [],
    projects = [],
    boards = [],
    validationErrors = [],
    configuredConnection,
    handleConnectionTabChange = () => {},
    setDataEntities = () => {},
    setProjects = () => {},
    setBoards = () => {},
    setBoardSearch = () => {},
    prevStep = () => {},
    fieldHasError = () => {},
    getFieldError = () => {},
    isSaving = false,
    isRunning = false,
    isFetching = false,
    enableConnectionTabs = true,
    elevation = Elevation.TWO,
    cardStyle = {}
  } = props

  const { Providers, ProviderIcons } = useContext(IntegrationsContext)

  const selectedBoards = useMemo(
    () => boards[configuredConnection.id],
    [boards, configuredConnection?.id]
  )
  const selectedProjects = useMemo(
    () => projects[configuredConnection.id],
    [projects, configuredConnection?.id]
  )

  useEffect(() => {
    console.log('>> OVER HERE!!!', selectedBoards)
  }, [selectedBoards])

  useEffect(() => {
    console.log('>> OVER HERE FOR Projects!!!', selectedProjects)
  }, [selectedProjects])

  return (
    <div
      className='workflow-step workflow-step-data-scope'
      data-step={activeStep?.id}
    >
      {blueprintConnections.length > 0 && (
        <div style={{ display: 'flex' }}>
          {enableConnectionTabs && (
            <div
              className='connection-tab-selector'
              style={{ minWidth: '200px' }}
            >
              <Card
                className='workflow-card connection-tabs-card'
                elevation={Elevation.TWO}
                style={{ padding: '10px' }}
              >
                <ConnectionTabs
                  connections={blueprintConnections}
                  onChange={handleConnectionTabChange}
                  selectedTabId={activeConnectionTab}
                  errors={validationErrors}
                />
              </Card>
            </div>
          )}
          <div
            className='connection-scope'
            style={{ marginLeft: '10px', width: '100%' }}
          >
            <Card
              className='workflow-card worfklow-panel-card'
              elevation={elevation}
              style={{ ...cardStyle }}
            >
              {configuredConnection && (
                <>
                  <h3>
                    <span style={{ float: 'left', marginRight: '8px' }}>
                      {ProviderIcons[configuredConnection.provider] ? (
                        ProviderIcons[configuredConnection.provider](24, 24)
                      ) : (
                        <></>
                      )}
                    </span>{' '}
                    {configuredConnection.title}
                  </h3>
                  <Divider className='section-divider' />

                  {[Providers.GITHUB].includes(
                    configuredConnection.provider
                  ) && (
                    <>
                      <h4>Projects *</h4>
                      <p>Enter the project names you would like to sync.</p>
                      <TagInput
                        id='project-id'
                        disabled={isRunning}
                        placeholder='username/repo, username/another-repo'
                        values={
                          projects[configuredConnection.id]?.map(
                            (p) => p.value
                          ) || []
                        }
                        fill={true}
                        onChange={(values) =>
                          setProjects((p) => ({
                            ...p,
                            [configuredConnection.id]: [
                              ...values.map(
                                (v, vIdx) =>
                                  new GitHubProject({
                                    id: v,
                                    key: v,
                                    title: v,
                                    value: v,
                                    type: 'string'
                                  })
                              )
                            ]
                          }))
                        }
                        addOnPaste={true}
                        addOnBlur={true}
                        rightElement={
                          <Button
                            disabled={isRunning}
                            icon='eraser'
                            minimal
                            onClick={() =>
                              setProjects((p) => ({
                                ...p,
                                [configuredConnection.id]: []
                              }))
                            }
                          />
                        }
                        onKeyDown={(e) =>
                          e.key === 'Enter' && e.preventDefault()
                        }
                        tagProps={{
                          intent: validationErrors.some((e) =>
                            e.startsWith('Projects:')
                          )
                            ? Intent.WARNING
                            : Intent.PRIMARY,
                          minimal: true
                        }}
                        className='input-project-id tagInput'
                      />
                    </>
                  )}

                  {[Providers.JIRA].includes(configuredConnection.provider) && (
                    <>
                      <h4>Boards *</h4>
                      <p>Select the boards you would like to sync.</p>
                      <BoardsSelector
                        items={boardsList}
                        selectedItems={selectedBoards}
                        onQueryChange={setBoardSearch}
                        onItemSelect={setBoards}
                        onClear={setBoards}
                        onRemove={setBoards}
                        disabled={isSaving}
                        configuredConnection={configuredConnection}
                        isLoading={isFetching}
                      />
                    </>
                  )}

                  {[Providers.GITLAB].includes(
                    configuredConnection.provider
                  ) && (
                    <>
                      <h4>Projects *</h4>
                      <p>Select the project you would like to sync.</p>
                      <GitlabProjectsSelector
                        onFetch={fetchGitlabProjects}
                        isFetching={isFetchingGitlab}
                        items={gitlabProjects}
                        selectedItems={selectedProjects}
                        onItemSelect={setProjects}
                        onClear={setProjects}
                        onRemove={setProjects}
                        disabled={isSaving}
                        configuredConnection={configuredConnection}
                        isLoading={isFetching}
                      />
                    </>
                  )}

                  {[Providers.JENKINS].includes(
                    configuredConnection.provider
                  ) && (
                    <>
                      <h4>Jobs *</h4>
                      <p>Select the job you would like to sync.</p>
                      <JenkinsJobsSelector
                        onFetch={fetchJenkinsJobs}
                        isFetching={isFetchingJenkins}
                        items={jenkinsJobs}
                        selectedItems={selectedProjects}
                        onItemSelect={setProjects}
                        onClear={setProjects}
                        onRemove={setProjects}
                        disabled={isSaving}
                        configuredConnection={configuredConnection}
                        isLoading={isFetching}
                      />
                    </>
                  )}

                  <h4>Data Entities</h4>
                  <p>
                    Select the data entities you wish to collect for the
                    projects.{' '}
                    <a
                      href='https://devlake.apache.org/docs/DataModels/DevLakeDomainLayerSchema'
                      target='_blank'
                      rel='noreferrer'
                    >
                      Learn about data entities
                    </a>
                  </p>
                  <DataEntitiesSelector
                    items={dataEntitiesList}
                    selectedItems={dataEntities[configuredConnection.id] || []}
                    // restrictedItems={getRestrictedDataEntities()}
                    onItemSelect={setDataEntities}
                    onClear={setDataEntities}
                    fieldHasError={fieldHasError}
                    getFieldError={getFieldError}
                    onRemove={setDataEntities}
                    disabled={isSaving}
                    configuredConnection={configuredConnection}
                    isSaving={isSaving}
                  />
                </>
              )}
            </Card>
          </div>
        </div>
      )}
      {blueprintConnections.length === 0 && (
        <NoData
          title='No Data Connections'
          message='Please select at least one connection source.'
          onClick={prevStep}
        />
      )}
    </div>
  )
}

export default DataScopes
