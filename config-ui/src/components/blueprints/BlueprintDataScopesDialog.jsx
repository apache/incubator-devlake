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
import classNames from 'classnames'
import React, { useEffect, useState, useRef, useCallback } from 'react'
import {
  Button,
  Classes,
  Colors,
  Dialog,
  DialogStep,
  MultistepDialog,
  MultistepDialogNavPosition,
  Elevation,
  FormGroup,
  Icon,
  Intent,
  Label,
  MenuItem,
  Position,
} from '@blueprintjs/core'

import { NullBlueprint } from '@/data/NullBlueprint'
import DataScopes from '@/components/blueprints/create-workflow/DataScopes'
import DataTransformations from '@/components/blueprints/create-workflow/DataTransformations'

const Modes = {
  CREATE: 'create',
  EDIT: 'edit',
}

const DialogPanel = (props) => {
  const {
    children
  } = props

  return (
    <div className={classNames(Classes.DIALOG_BODY)} style={{ minHeight: '300px' }}>
      {children}
    </div>
  )
}

const BlueprintDataScopesDialog = (props) => {
  const {
    isOpen = false,
    title = 'Change Data Scope',
    blueprintConnections = [],
    blueprint = NullBlueprint,
    provider,
    activeTransformation,
    configuredConnection,
    configuredProject,
    configuredBoard,
    scopeConnection,
    dataEntitiesList = [],
    boardsList = [],
    issueTypesList = [],
    fieldsList = [],
    boards = {},
    // dataEntities = [],
    projects = {},
    mode = Modes.EDIT,
    canOutsideClickClose = false,
    showCloseButtonInFooter = true,
    resetOnClose = true,
    isCloseButtonShown = true,
    initialStepIndex = 0,
    navPosition = 'top',
    usePortal = true,
    hasTitle = true,
    activeStep = null,
    onStepChange = () => {},
    onClose = () => {},
    onCancel = () => {},
    onSave = () => {},
    setDataEntities = () => {},
    setProjects = () => {},
    setBoards = () => {},
    setConfiguredProject = () => {},
    setConfiguredBoard = () => {},
    fieldHasError = () => {},
    getFieldError = () => {},
    isSaving = false,
    isValid = true,
    isTesting = false,
    isFetchingJIRA = false,
    jiraProxyError,
    content = null,
    backButtonProps = {
      // disabled:
      intent: Intent.PRIMARY,
      text: 'Previous Step',
      outlined: true
    },
    nextButtonProps = {
      disabled: !isValid,
      intent: Intent.PRIMARY,
      text: 'Next Step',
      outlined: true
    },
    finalButtonProps = {
      disabled: !isValid,
      intent: Intent.PRIMARY,
      onClick: onSave,
      text: 'Save Changes',
    },
    closeButtonProps = {
      // disabled:
      intent: Intent.PRIMARY,
      text: 'Cancel',
      outlined: true
    }
  } = props

  // const [boards, setBoards] = useState({ [configuredConnection?.id]: [] })
  // const [projects, setProjects] = useState({ [configuredConnection?.id]: [] })

  useEffect(() => {
    console.log('>>> MY BOARDS LIST!!!!', boardsList)
  }, [boardsList])

  useEffect(() => {
    console.log('>>> MY SELECTED BOARDS!!!!', boards)
  }, [boards])

  // useEffect(() => {
  //   console.log('>>> MY CONNECTION SCOPE!!!!', scopeConnection)
  //   setBoards({ [configuredConnection?.id]: scopeConnection?.boardsList })
  //   setProjects({ [configuredConnection?.id]: scopeConnection?.projects })
  // }, [scopeConnection, configuredConnection])

  return (
    <>
      <MultistepDialog
        disabled
        isOpen={isOpen}
        className='blueprint-data-scopes-dialog'
        // icon='info-sign'
        navigationPosition={Position.BOTTOM}
        closeButtonProps={closeButtonProps}
        backButtonProps={backButtonProps}
        nextButtonProps={nextButtonProps}
        finalButtonProps={finalButtonProps}
        title={title}
        hasTitle={hasTitle}
        initialStepIndex={initialStepIndex}
        showCloseButtonInFooter={showCloseButtonInFooter}
        isCloseButtonShown={isCloseButtonShown}
        canOutsideClickClose={canOutsideClickClose}
        resetOnClose={resetOnClose}
        onClose={onClose}
        onClosed={() => {}}
        onChange={onStepChange}
      >
        <DialogStep
          id='scopes'
          panel={
            <DialogPanel>
              <DataScopes
                provider={provider}
                activeStep={activeStep}
                // advancedMode={false}
                // activeConnectionTab={activeConnectionTab}
                blueprintConnections={blueprintConnections}
                dataEntitiesList={dataEntitiesList}
                boardsList={boardsList}
                // boards={{
                //   [configuredConnection?.id]: scopeConnection?.boardsList
                // }}
                boards={boards}
                dataEntities={{
                  [configuredConnection?.id]: scopeConnection?.entityList
                  // [configuredConnection?.id]: []
                }}
                // projects={{
                //   [configuredConnection?.id]: scopeConnection?.projects
                // }}
                projects={projects}
                configuredConnection={configuredConnection}
                // handleConnectionTabChange={handleConnectionTabChange}
                setDataEntities={setDataEntities}
                setProjects={setProjects}
                setBoards={setBoards}
                // prevStep={prevStep}
                isSaving={isSaving}
                // isRunning={isRunning}
                validationErrors={[]}
                enableConnectionTabs={false}
                elevation={Elevation.ZERO}
                cardStyle={{ padding: 0, backgroundColor: 'transparent' }}
              />
            </DialogPanel>
          }
          title='Data Scopes'
        />
        <DialogStep
          id='transformations'
          panel={
            <DialogPanel>
              <DataTransformations
                provider={provider}
                activeTransformation={activeTransformation}
                blueprintConnections={blueprintConnections}
                dataEntities={{
                  [configuredConnection?.id]: scopeConnection?.entityList
                }}
                projects={{
                  [configuredConnection?.id]: scopeConnection?.projects
                }}
                boards={{
                  [configuredConnection?.id]: scopeConnection?.boardsList
                }}
                // transformations={[projects].map(entity => ({[entity]: {}}).reduce((pV, cV) => ({...pV}), {}))}
                boardsList={boardsList}
                issueTypes={issueTypesList}
                fields={fieldsList}
                configuredConnection={configuredConnection}
                configuredProject={configuredProject}
                configuredBoard={configuredBoard}
                // handleConnectionTabChange={handleConnectionTabChange}
                // prevStep={prevStep}
                addBoardTransformation={(board) => setConfiguredBoard(board)}
                addProjectTransformation={(project) => setConfiguredProject(project)}
                // transformations={transformations}
                // activeTransformation={activeTransformation}
                // setTransformations={setTransformations}
                // setTransformationSettings={setTransformationSettings}
                isSaving={isSaving}
                // isSavingConnection={isSavingConnection}
                // isRunning={isRunning}
                // onSave={handleTransformationSave}
                // onCancel={handleTransformationCancel}
                // onClear={handleTransformationClear}
                jiraProxyError={jiraProxyError}
                isFetchingJIRA={isFetchingJIRA}
                fieldHasError={fieldHasError}
                getFieldError={getFieldError}
                enableConnectionTabs={false}
                enableNoticeAlert={false}
                useDropdownSelector={true}
                enableGoBack={false}
                elevation={Elevation.ZERO}
                cardStyle={{ padding: 0, backgroundColor: 'transparent' }}
              />
            </DialogPanel>
          }
          title='Transformations'
        />
      </MultistepDialog>
    </>
  )
}

export default BlueprintDataScopesDialog
