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
import React, { useCallback, useContext, useMemo } from 'react'
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
import DataDomainsSelector from '@/components/blueprints/DataDomainsSelector'
import NoData from '@/components/NoData'
import { GitLabMillerColumns, GitLabProjectSelector } from '@/components/gitlab'
import GitlabProject from '@/models/GitlabProject'
import { JIRAMillerColumns } from '@/components/jira'
import JiraBoard from '@/models/JiraBoard'
import { GitHubMillerColumns } from '@/components/github'
import GitHubProject from '@/models/GithubProject'
import { JenkinsMillerColumns } from '@/components/jenkins'
import JenkinsJobsSelector from '@/components/blueprints/JenkinsJobsSelector'
import JenkinsJob from '@/models/JenkinsJob'

const DataScopes = (props) => {
  const {
    activeStep,
    activeConnectionTab,
    blueprintConnections = [],
    jiraBoards = [],
    fetchJenkinsJobs = () => [],
    isFetchingJenkins = false,
    jenkinsJobs = [],
    dataDomainsGroup = [],
    scopeEntitiesGroup = [],
    validationErrors = [],
    configuredConnection,
    handleConnectionTabChange = () => {},
    setDataDomainsGroup = () => {},
    setScopeEntitiesGroup = () => {},
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

  const { Integrations, Providers, ProviderIcons } =
    useContext(IntegrationsContext)

  const selectedScopeEntities = useMemo(
    () => scopeEntitiesGroup[configuredConnection.id],
    [scopeEntitiesGroup, configuredConnection?.id]
  )

  const setScopeEntities = useCallback(
    (scopeEntities) => {
      setScopeEntitiesGroup((g) => ({
        ...g,
        [configuredConnection.id]: scopeEntities
      }))
    },
    [setScopeEntitiesGroup, configuredConnection?.id]
  )

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
                      <h4>Repositories *</h4>
                      {!!activeStep && (
                        <>
                          <p>Select the repositories you would like to sync.</p>
                          <GitHubMillerColumns
                            connectionId={configuredConnection.connectionId}
                            onChangeItems={(items) =>
                              setScopeEntities([
                                ...scopeEntitiesGroup[
                                  configuredConnection.id
                                ].filter((it) => it.type !== 'miller-columns'),
                                ...items.map((it) => new GitHubProject(it))
                              ])
                            }
                          />
                        </>
                      )}
                      <div style={{ margin: '16px 0 8px' }}>
                        Add repositories outside of your organizations
                      </div>
                      <p>
                        Enter the repositories using the format “owner/repo” and
                        separate multiple repos with a comma.
                      </p>
                      <TagInput
                        id='project-id'
                        disabled={isRunning}
                        placeholder='username/repo, username/another-repo'
                        values={
                          selectedScopeEntities
                            ?.filter((it) => it.type !== 'miller-columns')
                            .map((p) => p.value) || []
                        }
                        fill={true}
                        onChange={(values) =>
                          setScopeEntities([
                            ...scopeEntitiesGroup[
                              configuredConnection.id
                            ].filter((it) => it.type === 'miller-columns'),
                            ...values.map(
                              (v, vIdx) =>
                                new GitHubProject({
                                  id: v,
                                  key: v,
                                  owner: v.includes('/') ? v.split('/')[0] : '',
                                  repo: v.includes('/') ? v.split('/')[1] : '',
                                  title: v,
                                  value: v,
                                  type: 'string'
                                })
                            )
                          ])
                        }
                        addOnPaste={true}
                        addOnBlur={true}
                        rightElement={
                          <Button
                            disabled={isRunning}
                            icon='eraser'
                            minimal
                            onClick={() => setScopeEntities([])}
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
                      {!activeStep ? (
                        <BoardsSelector
                          items={jiraBoards}
                          selectedItems={selectedScopeEntities}
                          onQueryChange={setBoardSearch}
                          onItemSelect={setScopeEntities}
                          onClear={setScopeEntities}
                          onRemove={setScopeEntities}
                          disabled={isSaving}
                          configuredConnection={configuredConnection}
                          isLoading={isFetching}
                        />
                      ) : (
                        <JIRAMillerColumns
                          connectionId={configuredConnection.connectionId}
                          onChangeItems={(items) =>
                            setScopeEntities(
                              items.map((it) => new JiraBoard(it))
                            )
                          }
                        />
                      )}
                    </>
                  )}

                  {[Providers.GITLAB].includes(
                    configuredConnection.provider
                  ) && (
                    <>
                      <h4>Projects *</h4>
                      {!!activeStep && (
                        <>
                          <p>Select the project you would like to sync.</p>
                          <GitLabMillerColumns
                            connectionId={configuredConnection.connectionId}
                            disabledItemIds={selectedScopeEntities
                              .filter((it) => it.type !== 'miller-columns')
                              .map((it) => it.id)}
                            onChangeItems={(items) =>
                              setScopeEntities([
                                ...scopeEntitiesGroup[
                                  configuredConnection.id
                                ].filter((it) => it.type !== 'miller-columns'),
                                ...items.map(
                                  (it) =>
                                    new GitlabProject({
                                      ...it,
                                      type: 'miller-columns'
                                    })
                                )
                              ])
                            }
                          />
                        </>
                      )}
                      <div style={{ margin: '16px 0 8px' }}>
                        Add repositories outside of your projects
                      </div>
                      <p>
                        Enter the repositories using the format “owner/repo” and
                        separate multiple repos with a comma.
                      </p>
                      <GitLabProjectSelector
                        connectionId={configuredConnection.connectionId}
                        disabledItemIds={selectedScopeEntities
                          .filter((it) => it.type !== 'project-selector')
                          .map((it) => it.id)}
                        selectedItems={selectedScopeEntities
                          .filter((it) => it.type === 'project-selector')
                          .map((it) => ({
                            id: it.id,
                            key: it.id,
                            title: it.title,
                            shortTitle: it.shortTitle,
                            value: it.id
                          }))}
                        onChangeItems={(items) =>
                          setScopeEntities([
                            ...scopeEntitiesGroup[
                              configuredConnection.id
                            ].filter((it) => it.type !== 'project-selector'),
                            ...items.map(
                              (it) =>
                                new GitlabProject({
                                  ...it,
                                  type: 'project-selector'
                                })
                            )
                          ])
                        }
                      />
                    </>
                  )}

                  {[Providers.JENKINS].includes(
                    configuredConnection.provider
                  ) && (
                    <>
                      <h4>Jobs *</h4>
                      <p>Select the jobs you would like to sync.</p>
                      {!activeStep ? (
                        <JenkinsJobsSelector
                          onFetch={fetchJenkinsJobs}
                          isFetching={isFetchingJenkins}
                          items={jenkinsJobs}
                          selectedItems={selectedScopeEntities}
                          onItemSelect={setScopeEntities}
                          onClear={setScopeEntities}
                          onRemove={setScopeEntities}
                          disabled={isSaving}
                          configuredConnection={configuredConnection}
                          isLoading={isFetching}
                        />
                      ) : (
                        <JenkinsMillerColumns
                          connectionId={configuredConnection.connectionId}
                          onChangeItems={(items) =>
                            setScopeEntities(
                              items.map((it) => new JenkinsJob(it))
                            )
                          }
                        />
                      )}
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
                  <DataDomainsSelector
                    items={
                      Integrations.find(
                        (p) => p.id === configuredConnection.provider
                      )?.getAvailableDataDomains() || []
                    }
                    selectedItems={
                      dataDomainsGroup[configuredConnection.id] || []
                    }
                    onItemSelect={setDataDomainsGroup}
                    onClear={setDataDomainsGroup}
                    fieldHasError={fieldHasError}
                    getFieldError={getFieldError}
                    onRemove={setDataDomainsGroup}
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
