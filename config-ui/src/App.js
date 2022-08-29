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

import React, { useEffect, useState, useCallback } from 'react'
import {
  BrowserRouter as Router,
  Route
} from 'react-router-dom'
import {
  Intent
} from '@blueprintjs/core'

import 'normalize.css'
import '@/styles/app.scss'
// import 'typeface-montserrat'
// import 'jetbrains-mono'
import '@fontsource/inter/400.css'
import '@fontsource/inter/600.css'
import '@fontsource/inter/variable-full.css'
// Theme variables (@styles/theme.scss) injected via Webpack w/ @sass-loader additionalData option!
// import '@/styles/theme.scss'

import { MigrationOptions } from '@/config/migration'
import request from '@/utils/request'

import ErrorBoundary from '@/components/ErrorBoundary'
import { ToastNotification } from '@/components/Toast'
// import Configure from './pages/configure/index'
import Integration from '@/pages/configure/integration/index'
import ManageIntegration from '@/pages/configure/integration/manage'
import AddConnection from '@/pages/configure/connections/AddConnection'
// import EditConnection from '@/pages/configure/connections/EditConnection'
import ConfigureConnection from '@/pages/configure/connections/ConfigureConnection'
// import Triggers from '@/pages/triggers/index'
import Offline from '@/pages/offline/index'
// import Pipelines from '@/pages/pipelines/index'
// import CreatePipeline from '@/pages/pipelines/create'
// import PipelineActivity from '@/pages/pipelines/activity'
import Blueprints from '@/pages/blueprints/index'
import CreateBlueprint from '@/pages/blueprints/create-blueprint'
import BlueprintDetail from '@/pages/blueprints/blueprint-detail'
import BlueprintSettings from '@/pages/blueprints/blueprint-settings'
import Connections from '@/pages/connections/index'
import MigrationAlertDialog from '@/components/MigrationAlertDialog'

function App (props) {
  const [isProcessing, setIsProcessing] = useState(false)

  const [migrationWarning, setMigrationWarning] = useState(localStorage.getItem(MigrationOptions.warningId))
  const [migrationAlertOpened, setMigrationAlertOpened] = useState(false)
  const [wasMigrationSuccessful, setWasMigrationSuccessful] = useState(false)
  const [hasMigrationFailed, setHasMigrationFailed] = useState(false)

  const handleConfirmMigration = useCallback(() => {
    setIsProcessing(true)
    const m = request.get(MigrationOptions.apiProceedEndpoint)
    setWasMigrationSuccessful(m?.status === 200 && m?.success === true)
    setTimeout(() => {
      setIsProcessing(false)
      setHasMigrationFailed(m?.status !== 200)
    }, 3000)
  }, [])

  const handleCancelMigration = useCallback(() => {
    setIsProcessing(true)
    localStorage.removeItem(MigrationOptions.warningId)
    setMigrationAlertOpened(false)
    setIsProcessing(false)
    ToastNotification.clear()
    ToastNotification.show({
      // eslint-disable-next-line max-len
      message: MigrationOptions.cancelToastMessage,
      intent: Intent.NONE,
      icon: 'warning-sign'
    })
  }, [])

  const handleMigrationDialogClose = useCallback(() => {
    setMigrationAlertOpened(false)
  }, [setMigrationAlertOpened])

  useEffect(() => {
    setMigrationAlertOpened(migrationWarning !== null)
  }, [migrationWarning, setMigrationAlertOpened])

  useEffect(() => {
    if (wasMigrationSuccessful) {
      localStorage.removeItem(MigrationOptions.warningId)
    }
  }, [wasMigrationSuccessful])

  useEffect(() => {
    if (hasMigrationFailed) {
      ToastNotification.clear()
      ToastNotification.show({
        // eslint-disable-next-line max-len
        message: MigrationOptions.failedToastMessage,
        intent: Intent.DANGER,
        icon: 'error'
      })
    }
  }, [hasMigrationFailed])

  useEffect(() => {
    if (migrationWarning) {
      // eslint-disable-next-line max-len
      console.log(`>>> MIGRATION WARNING DETECTED !! Local Storage Key = [${MigrationOptions.warningId}]:`, migrationWarning)
    }
  }, [migrationWarning])

  return (
    <Router>
      <Route exact path='/'>
        <ErrorBoundary>
          <Integration />
        </ErrorBoundary>
      </Route>
      <Route path='/integrations/:providerId'>
        <ErrorBoundary>
          <ManageIntegration />
        </ErrorBoundary>
      </Route>
      <Route path='/connections/add/:providerId'>
        <ErrorBoundary>
          <AddConnection />
        </ErrorBoundary>
      </Route>
      <Route path='/connections/configure/:providerId/:connectionId'>
        <ErrorBoundary>
          <ConfigureConnection />
        </ErrorBoundary>
      </Route>
      <Route exact path='/integrations'>
        <ErrorBoundary>
          <Integration />
        </ErrorBoundary>
      </Route>
      <Route exact path='/blueprints/create'>
        <ErrorBoundary>
          <CreateBlueprint />
        </ErrorBoundary>
      </Route>
      <Route exact path='/blueprints/detail/:bId'>
        <ErrorBoundary>
          <BlueprintDetail />
        </ErrorBoundary>
      </Route>
      <Route exact path='/blueprints/settings/:bId'>
        <ErrorBoundary>
          <BlueprintSettings />
        </ErrorBoundary>
      </Route>
      <Route exact path='/blueprints'>
        <ErrorBoundary>
          <Blueprints />
        </ErrorBoundary>
      </Route>
      <Route exact path='/connections'>
        <Connections />
      </Route>
      <Route exact path='/offline'>
        <Offline />
      </Route>
      <MigrationAlertDialog
        isOpen={migrationAlertOpened}
        onClose={handleMigrationDialogClose}
        onCancel={handleCancelMigration}
        onConfirm={handleConfirmMigration}
        isMigrating={isProcessing}
        wasSuccesful={wasMigrationSuccessful}
        hasFailed={hasMigrationFailed}
      />
    </Router>

  )
}

export default App
