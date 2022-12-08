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

import React from 'react'
import {
  BrowserRouter as Router,
  Route,
  Switch,
  Redirect
} from 'react-router-dom'

import 'normalize.css'
import '@/styles/app.scss'
// import 'typeface-montserrat'
// import 'jetbrains-mono'
import '@fontsource/inter/400.css'
import '@fontsource/inter/600.css'
import '@fontsource/inter/variable-full.css'
// Theme variables (@styles/theme.scss) injected via Webpack w/ @sass-loader additionalData option!
// import '@/styles/theme.scss'
import useDatabaseMigrations from '@/hooks/useDatabaseMigrations'

import { BaseLayout } from '@/layouts'
import {
  ProjectHomePage,
  CreateBlueprintPage,
  WebHookConnectionPage
} from '@/pages'
import ErrorBoundary from '@/components/ErrorBoundary'
import Integration from '@/pages/configure/integration/index'
import ManageIntegration from '@/pages/configure/integration/manage'
import AddConnection from '@/pages/configure/connections/AddConnection'
import ConfigureConnection from '@/pages/configure/connections/ConfigureConnection'
import Offline from '@/pages/offline/index'
import Blueprints from '@/pages/blueprints/index'
import CreateBlueprint from '@/pages/blueprints/create-blueprint'
import BlueprintDetail from '@/pages/blueprints/blueprint-detail'
import BlueprintSettings from '@/pages/blueprints/blueprint-settings'
import { IncomingWebhook as IncomingWebhookConnection } from '@/pages/connections/incoming-webhook'
import MigrationAlertDialog from '@/components/MigrationAlertDialog'

function App(props) {
  const {
    isProcessing,
    migrationWarning,
    migrationAlertOpened,
    wasMigrationSuccessful,
    hasMigrationFailed,
    handleConfirmMigration,
    handleCancelMigration,
    handleMigrationDialogClose
  } = useDatabaseMigrations()

  return (
    <Router>
      <Switch>
        <Route exact path='/offline' component={() => <Offline />} />

        <Route>
          <BaseLayout>
            <Switch>
              <Route
                path='/'
                exact
                component={() => <Redirect to='/projects' />}
              />
              <Route
                exact
                path='/projects'
                component={() => <ProjectHomePage />}
              />
              <Route
                exact
                path='/integrations'
                component={() => (
                  <ErrorBoundary>
                    <Integration />
                  </ErrorBoundary>
                )}
              />
              <Route
                path='/integrations/:providerId'
                component={() => (
                  <ErrorBoundary>
                    <ManageIntegration />
                  </ErrorBoundary>
                )}
              />
              <Route
                path='/connections/add/:providerId'
                component={() => (
                  <ErrorBoundary>
                    <AddConnection />
                  </ErrorBoundary>
                )}
              />
              <Route
                path='/connections/configure/:providerId/:connectionId'
                component={() => (
                  <ErrorBoundary>
                    <ConfigureConnection />
                  </ErrorBoundary>
                )}
              />
              <Route
                exact
                path='/connections/incoming-webhook'
                component={() => <WebHookConnectionPage />}
              />
              <Route
                exact
                path='/blueprints'
                component={() => (
                  <ErrorBoundary>
                    <Blueprints />
                  </ErrorBoundary>
                )}
              />
              <Route
                exact
                path='/blueprints/create'
                component={() => <CreateBlueprintPage from='blueprint' />}
              />

              <Route
                exact
                path='/blueprints/detail/:bId'
                component={() => (
                  <ErrorBoundary>
                    <BlueprintDetail />
                  </ErrorBoundary>
                )}
              />
              <Route
                exact
                path='/blueprints/settings/:bId'
                component={() => (
                  <ErrorBoundary>
                    <BlueprintSettings />
                  </ErrorBoundary>
                )}
              />
            </Switch>
            <MigrationAlertDialog
              isOpen={migrationAlertOpened}
              onClose={handleMigrationDialogClose}
              onCancel={handleCancelMigration}
              onConfirm={handleConfirmMigration}
              isMigrating={isProcessing}
              wasSuccessful={wasMigrationSuccessful}
              hasFailed={hasMigrationFailed}
            />
          </BaseLayout>
        </Route>
      </Switch>
    </Router>
  )
}

export default App
